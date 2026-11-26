package process

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func pageTransport(client *http.Client) *BaseProcessor {
	return &BaseProcessor{httpClient: client, chanTriggerNodesState: make(chan struct{}, 1)}
}

func TestTrieWalkTransportPreservesErrorSignal(t *testing.T) {
	t.Parallel()
	var firstCalls, secondCalls atomic.Int32
	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		firstCalls.Add(1)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(400)
		_, _ = w.Write([]byte(pageUnavailable))
	}))
	defer first.Close()
	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { secondCalls.Add(1); _, _ = w.Write([]byte(pageOK)) }))
	defer second.Close()
	proc := pageTransport(first.Client())
	response, err := fetchPage(context.Background(), proc, continuationRequest(first.URL, second.URL))
	require.NoError(t, err)
	require.Equal(t, data.ReturnCodeSuccess, response.Code)
	require.EqualValues(t, 1, firstCalls.Load())
	require.EqualValues(t, 1, secondCalls.Load())
	status, raw, err := proc.CallTrieWalkEndpoint(context.Background(), http.MethodPost, first.URL, "/page", []byte(`{}`))
	require.NoError(t, err)
	require.Equal(t, 400, status)
	require.JSONEq(t, pageUnavailable, string(raw))
}

func TestTrieWalkTransportDoesNotFollowPostRedirects(t *testing.T) {
	t.Parallel()
	for _, status := range []int{301, 302, 303, 307, 308} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var redirects, targets atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/target" {
					targets.Add(1)
					_, _ = w.Write([]byte(pageOK))
					return
				}
				redirects.Add(1)
				http.Redirect(w, r, "/target", status)
			}))
			defer server.Close()
			client := server.Client()
			proc := pageTransport(client)
			_, err := fetchPage(context.Background(), proc, continuationRequest(server.URL))
			require.ErrorIs(t, err, ErrSendingRequest)
			require.EqualValues(t, 1, redirects.Load())
			require.Zero(t, targets.Load())
			require.Nil(t, client.CheckRedirect)
			// The generic client's existing GET policy is not changed by a paging POST.
			code, _, err := proc.CallTrieWalkEndpoint(context.Background(), http.MethodGet, server.URL, "/page", nil)
			require.NoError(t, err)
			require.Equal(t, 200, code)
			require.EqualValues(t, 1, targets.Load())
		})
	}
}

func TestTrieWalkTransportCancellationDuringFailover(t *testing.T) {
	t.Parallel()
	started, cancelled := make(chan struct{}), make(chan struct{})
	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(pageUnavailable))
	}))
	defer first.Close()
	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		close(started)
		<-r.Context().Done()
		close(cancelled)
	}))
	defer second.Close()
	var thirdCalls atomic.Int32
	third := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { thirdCalls.Add(1) }))
	defer third.Close()
	proc := pageTransport(first.Client())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	finished := make(chan error, 1)
	go func() {
		_, err := fetchPage(ctx, proc, continuationRequest(first.URL, second.URL, third.URL))
		finished <- err
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("second observer did not start")
	}
	cancel()
	select {
	case err := <-finished:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(3 * time.Second):
		t.Fatal("request did not stop")
	}
	select {
	case <-cancelled:
	case <-time.After(3 * time.Second):
		t.Fatal("upstream request was not cancelled")
	}
	require.Zero(t, thirdCalls.Load())
	require.Empty(t, proc.chanTriggerNodesState)
}

func TestTrieWalkTransportTimeoutCanFailOver(t *testing.T) {
	t.Parallel()
	first := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.Copy(io.Discard, r.Body); <-r.Context().Done() }))
	defer first.Close()
	second := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(pageOK)) }))
	defer second.Close()
	client := first.Client()
	client.Timeout = 50 * time.Millisecond
	proc := pageTransport(client)
	response, err := fetchPage(context.Background(), proc, continuationRequest(first.URL, second.URL))
	require.NoError(t, err)
	require.Equal(t, data.ReturnCodeSuccess, response.Code)
	require.Len(t, proc.chanTriggerNodesState, 1)
}

func TestTrieWalkTransportRawErrorAndCancelledContext(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	proc := pageTransport(server.Client())
	status, raw, err := proc.CallTrieWalkEndpoint(context.Background(), http.MethodPost, server.URL, "/missing", nil)
	require.NoError(t, err)
	require.Equal(t, 404, status)
	require.Contains(t, string(raw), "404")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, _, err = proc.CallTrieWalkEndpoint(ctx, http.MethodPost, server.URL, "/missing", nil)
	require.True(t, errors.Is(err, context.Canceled))
	require.Empty(t, proc.chanTriggerNodesState)
}
