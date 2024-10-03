package process

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func createMockArgNumShardsProcessor() ArgNumShardsProcessor {
	return ArgNumShardsProcessor{
		HttpClient:                    &mock.HttpClientMock{},
		Observers:                     []string{"obs1, obs2"},
		TimeBetweenNodesRequestsInSec: 2,
		NumShardsTimeoutInSec:         10,
		RequestTimeoutInSec:           5,
	}
}

func TestNewNumShardsProcessor(t *testing.T) {
	t.Parallel()

	t.Run("nil HttpClient should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.HttpClient = nil

		proc, err := NewNumShardsProcessor(args)
		require.Equal(t, ErrNilHttpClient, err)
		require.Nil(t, proc)
	})
	t.Run("empty observers list should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.Observers = []string{}

		proc, err := NewNumShardsProcessor(args)
		require.True(t, errors.Is(err, core.ErrInvalidValue))
		require.True(t, strings.Contains(err.Error(), "Observers"))
		require.Nil(t, proc)
	})
	t.Run("invalid TimeBetweenNodesRequestsInSec should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.TimeBetweenNodesRequestsInSec = 0

		proc, err := NewNumShardsProcessor(args)
		require.True(t, errors.Is(err, core.ErrInvalidValue))
		require.True(t, strings.Contains(err.Error(), "TimeBetweenNodesRequestsInSec"))
		require.Nil(t, proc)
	})
	t.Run("invalid NumShardsTimeoutInSec should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.NumShardsTimeoutInSec = 0

		proc, err := NewNumShardsProcessor(args)
		require.True(t, errors.Is(err, core.ErrInvalidValue))
		require.True(t, strings.Contains(err.Error(), "NumShardsTimeoutInSec"))
		require.Nil(t, proc)
	})
	t.Run("invalid RequestTimeoutInSec should error", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.RequestTimeoutInSec = 0

		proc, err := NewNumShardsProcessor(args)
		require.True(t, errors.Is(err, core.ErrInvalidValue))
		require.True(t, strings.Contains(err.Error(), "RequestTimeoutInSec"))
		require.Nil(t, proc)
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		proc, err := NewNumShardsProcessor(createMockArgNumShardsProcessor())
		require.NoError(t, err)
		require.NotNil(t, proc)
	})
}

func TestNumShardsProcessor_GetNetworkNumShards(t *testing.T) {
	t.Parallel()

	t.Run("context done should exit with timeout", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.TimeBetweenNodesRequestsInSec = 30
		args.NumShardsTimeoutInSec = 30

		proc, err := NewNumShardsProcessor(args)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(time.Millisecond * 200)
			cancel()
		}()
		numShards, err := proc.GetNetworkNumShards(ctx)
		require.Equal(t, errTimeIsOut, err)
		require.Zero(t, numShards)
	})
	t.Run("timeout should exit with timeout", func(t *testing.T) {
		t.Parallel()

		args := createMockArgNumShardsProcessor()
		args.TimeBetweenNodesRequestsInSec = 30
		args.NumShardsTimeoutInSec = 1

		proc, err := NewNumShardsProcessor(args)
		require.NoError(t, err)
		numShards, err := proc.GetNetworkNumShards(context.Background())
		require.True(t, errors.Is(err, errTimeIsOut))
		require.Zero(t, numShards)
	})
	t.Run("should work on 4th observer", func(t *testing.T) {
		t.Parallel()

		providedBody := &networkConfigResponse{
			Data: networkConfigResponseData{
				Config: struct {
					NumShards uint32 `json:"erd_num_shards_without_meta"`
				}(struct{ NumShards uint32 }{NumShards: 2}),
			},
		}
		providedBodyBuff, _ := json.Marshal(providedBody)

		args := createMockArgNumShardsProcessor()
		args.TimeBetweenNodesRequestsInSec = 1
		args.NumShardsTimeoutInSec = 15
		cnt := 0
		args.HttpClient = &mock.HttpClientMock{
			DoCalled: func(req *http.Request) (*http.Response, error) {
				cnt++
				switch cnt {
				case 1: // error on Do
					return nil, errors.New("observer offline")
				case 2: // status code not 200
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, nil
				case 3: // status code ok, but invalid response
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("not the expected response")),
					}, nil
				default: // response ok
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(providedBodyBuff)),
					}, nil
				}
			},
		}

		proc, err := NewNumShardsProcessor(args)
		require.NoError(t, err)
		numShards, err := proc.GetNetworkNumShards(context.Background())
		require.NoError(t, err)
		require.Equal(t, uint32(2), numShards)
	})
}
