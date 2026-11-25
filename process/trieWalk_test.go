package process

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

const pageOK = `{"data":{"pairs":{},"blockInfo":{"nonce":9007199254740993}},"error":"","code":"successful"}`
const pageUnavailable = `{"data":null,"error":"different wording","code":"iterator_state_unavailable"}`

func continuationRequest(nodes ...string) trieWalkRequest {
	return trieWalkRequest{
		path: "/address/alice/keys",
		page: common.TrieWalkPage{Paginated: true, IteratorState: json.RawMessage(`{"blockNonce":9007199254740993,"blockRootHash":"aa"}`)},
		getObservers: func(_ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			result := make([]*data.NodeData, 0, len(nodes))
			for _, address := range nodes {
				result = append(result, &data.NodeData{Address: address})
			}
			return result, nil
		},
	}
}

func requirePageStatus(t *testing.T, err error, expected int) {
	t.Helper()
	var apiErr *data.ObserverAPIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, expected, apiErr.StatusCode)
}

func TestTrieWalkContinuationGuard(t *testing.T) {
	t.Parallel()
	type answer struct {
		status int
		body   string
		err    error
	}
	for _, tc := range []struct {
		name    string
		answers []answer
		status  int
		calls   int
	}{
		{"unavailable bounded", []answer{{400, pageUnavailable, nil}, {400, pageUnavailable, nil}, {400, pageUnavailable, nil}}, 400, 3},
		{"mixed bounded", []answer{{400, pageUnavailable, nil}, {503, "", nil}, {404, "not found", nil}}, 400, 3},
		{"all busy", []answer{{429, "", nil}, {429, "", nil}, {429, "", nil}}, 429, 3},
		{"no signal", []answer{{503, "", nil}, {404, "", nil}, {408, "", nil}}, 500, 3},
		{"ordinary 400 text is not signal", []answer{{400, `{"error":"iterator node unavailable","code":"bad_request"}`, nil}}, 400, 1},
		{"terminal after signal", []answer{{400, pageUnavailable, nil}, {500, `{"error":"storage error","code":"internal_issue"}`, nil}}, 500, 2},
		{"413 terminal", []answer{{413, "", nil}}, 413, 1},
		{"success second", []answer{{400, pageUnavailable, nil}, {200, pageOK, nil}}, 200, 2},
		{"success third", []answer{{404, "", nil}, {400, pageUnavailable, nil}, {200, pageOK, nil}}, 200, 3},
		{"transport failure counts", []answer{{0, "", errors.New("offline")}, {400, pageUnavailable, nil}, {503, "", nil}}, 400, 3},
		{"bad 200 counts", []answer{{200, "null", nil}, {200, "{}", nil}, {400, pageUnavailable, nil}}, 400, 3},
		{"redirect counts", []answer{{307, "", nil}, {302, "", nil}, {400, pageUnavailable, nil}}, 400, 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			proc := &mock.ProcessorStub{CallTrieWalkEndpointCalled: func(ctx context.Context, method, address, path string, body []byte) (int, []byte, error) {
				require.Less(t, calls, len(tc.answers), "unexpected extra observer call")
				require.Equal(t, http.MethodPost, method)
				require.Contains(t, string(body), "9007199254740993")
				response := tc.answers[calls]
				calls++
				return response.status, []byte(response.body), response.err
			}}
			response, err := fetchPage(context.Background(), proc, continuationRequest("one", "two", "three", "four"))
			require.Equal(t, tc.calls, calls)
			if tc.status == 200 {
				require.NoError(t, err)
				require.Equal(t, data.ReturnCodeSuccess, response.Code)
				require.Contains(t, string(response.Data.(json.RawMessage)), "9007199254740993")
			} else if tc.name == "no signal" {
				require.ErrorIs(t, err, ErrSendingRequest)
			} else {
				requirePageStatus(t, err, tc.status)
			}
		})
	}
}

func TestTrieWalkPoolsAndDeduplication(t *testing.T) {
	t.Parallel()
	request := continuationRequest()
	var pools []data.ObserverDataAvailabilityType
	request.getObservers = func(pool data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
		pools = append(pools, pool)
		if pool == data.AvailabilityAll {
			return []*data.NodeData{{Address: "one"}, {Address: "one"}}, nil
		}
		return []*data.NodeData{{Address: "one"}, {Address: "two"}, {Address: "two"}, {Address: "three"}, {Address: "four"}}, nil
	}
	var calls []string
	proc := &mock.ProcessorStub{CallTrieWalkEndpointCalled: func(_ context.Context, _, address, _ string, _ []byte) (int, []byte, error) {
		calls = append(calls, address)
		return 400, []byte(pageUnavailable), nil
	}}
	for range 2 { // each client request has a fresh budget
		calls = nil
		pools = nil
		_, err := fetchPage(context.Background(), proc, request)
		requirePageStatus(t, err, 400)
		require.Equal(t, []string{"one", "two", "three"}, calls)
		require.Equal(t, []data.ObserverDataAvailabilityType{data.AvailabilityAll, data.AvailabilityRecent}, pools)
	}
}

func TestTrieWalkLookupAndCandidateExhaustion(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		allError    bool
		recentError bool
		nodes       []string
		status      int
		calls       int
	}{
		{"empty", false, false, nil, 500, 0},
		{"one unavailable", false, false, []string{"one"}, 400, 1},
		{"recent lookup fails after signal", false, true, []string{"one"}, 400, 1},
		{"all lookup fails then fallback", true, false, []string{"one"}, 400, 1},
		{"both lookup fail", true, true, nil, 500, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := continuationRequest(tc.nodes...)
			source := request.getObservers
			request.getObservers = func(pool data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				if (pool == data.AvailabilityAll && tc.allError) || (pool == data.AvailabilityRecent && tc.recentError) {
					return nil, errors.New("lookup failed")
				}
				return source(pool)
			}
			calls := 0
			proc := &mock.ProcessorStub{CallTrieWalkEndpointCalled: func(context.Context, string, string, string, []byte) (int, []byte, error) {
				calls++
				return 400, []byte(pageUnavailable), nil
			}}
			_, err := fetchPage(context.Background(), proc, request)
			require.Equal(t, tc.calls, calls)
			if tc.status == 400 {
				requirePageStatus(t, err, 400)
			} else {
				require.ErrorIs(t, err, ErrSendingRequest)
			}
		})
	}
}

func TestTrieWalkCancellationDuringObserverLookup(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	request := continuationRequest("one")
	request.getObservers = func(data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
		cancel()
		return nil, errors.New("lookup interrupted")
	}
	_, err := fetchPage(ctx, &mock.ProcessorStub{}, request)
	require.ErrorIs(t, err, context.Canceled)
}

func TestTrieWalkLegacyAndInvalidResponses(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		paginated   bool
		retryLegacy bool
		status      int
		body        string
		wantCalls   int
		wantStatus  int
	}{
		{"legacy partial", false, true, 200, `{"data":{"pairs":{},"iteratorState":{"blockNonce":1}},"code":"successful"}`, 1, 500},
		{"opt in partial", true, false, 200, `{"data":{"pairs":{},"iteratorState":{"blockNonce":1}},"code":"successful"}`, 1, 200},
		{"legacy complete", false, false, 200, pageOK, 1, 200},
		{"legacy listing 500", false, false, 500, `{"error":"node failure"}`, 1, 500},
		{"legacy account 500 retries", false, true, 500, `{"error":"node failure"}`, 2, 200},
		{"legacy text404 retries", false, false, 404, "not found", 2, 200},
		{"first post unavailable", true, false, 400, pageUnavailable, 1, 400},
		{"legacy null stops", false, false, 200, "null", 1, 500},
		{"legacy empty object stops", false, false, 200, "{}", 1, 500},
		{"legacy account bad200 retries", false, true, 200, "bad json", 2, 200},
		{"post empty200 retries", true, false, 200, "", 2, 200},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := continuationRequest("one", "two")
			request.page = common.TrieWalkPage{Paginated: tc.paginated}
			request.legacyUpstream = !tc.paginated
			request.retryOnNodeError = tc.retryLegacy
			calls := 0
			proc := &mock.ProcessorStub{CallTrieWalkEndpointCalled: func(_ context.Context, method, _, _ string, body []byte) (int, []byte, error) {
				calls++
				if tc.paginated {
					require.Equal(t, "POST", method)
				} else {
					require.Equal(t, "GET", method)
					require.Empty(t, body)
				}
				if calls > 1 {
					return 200, []byte(pageOK), nil
				}
				return tc.status, []byte(tc.body), nil
			}}
			response, err := fetchPage(context.Background(), proc, request)
			require.Equal(t, tc.wantCalls, calls)
			if tc.wantStatus == 200 {
				require.NoError(t, err)
				require.Equal(t, data.ReturnCodeSuccess, response.Code)
			} else {
				requirePageStatus(t, err, tc.wantStatus)
			}
		})
	}
}
