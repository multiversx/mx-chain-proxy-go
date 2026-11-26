package process

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestTrieWalkProcessorRouting(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name     string
		network  bool
		meta     bool
		withKeys bool
		forced   bool
		block    bool
		paged    bool
		method   string
		path     string
	}{
		{"address legacy", false, false, false, false, false, false, "GET", "/address/deadbeef/keys"},
		{"address opt in", false, false, false, true, false, true, "POST", "/address/deadbeef/keys"},
		{"account with keys", false, false, true, false, false, true, "POST", "/address/deadbeef"},
		{"registry", false, true, false, false, false, true, "POST", "/address/deadbeef/registered-nfts"},
		{"network legacy", true, true, false, false, false, false, "GET", "/network/esdts"},
		{"network historical all or error", true, true, false, false, true, false, "POST", "/network/esdts"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := common.AccountQueryOptions{}
			page := common.TrieWalkPage{Paginated: tc.paged}
			if tc.forced {
				options.ForcedShardID = core.OptionalUint32{HasValue: true, Value: 7}
			}
			if tc.block {
				options.BlockNonce = core.OptionalUint64{HasValue: true, Value: 5}
			}
			if tc.paged {
				page.NumKeys = core.OptionalUint64{HasValue: true, Value: 0}
			}
			calls := 0
			proc := &mock.ProcessorStub{ComputeShardIdCalled: func([]byte) (uint32, error) { return 1, nil }, GetObserversCalled: func(shard uint32, availability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				want := uint32(1)
				if tc.meta {
					want = core.MetachainShardId
				}
				if tc.forced {
					want = 7
				}
				require.Equal(t, want, shard)
				if tc.block {
					require.Equal(t, data.AvailabilityAll, availability)
				} else {
					require.Equal(t, data.AvailabilityRecent, availability)
				}
				return []*data.NodeData{{Address: "observer"}}, nil
			}, CallTrieWalkEndpointCalled: func(ctx context.Context, method, address, path string, body []byte) (int, []byte, error) {
				calls++
				require.Equal(t, tc.method, method)
				u, err := url.Parse(path)
				require.NoError(t, err)
				require.Equal(t, tc.path, u.Path)
				require.Empty(t, u.Query().Get("forced-shard-id"))
				if tc.block {
					require.Equal(t, "5", u.Query().Get("blockNonce"))
				}
				if tc.withKeys {
					require.Equal(t, "true", u.Query().Get("withKeys"))
				}
				if tc.paged {
					require.JSONEq(t, `{"numKeys":0}`, string(body))
				}
				return 200, []byte(`{"data":{"account":{"address":"deadbeef"},"blockInfo":{"rootHash":"aa"}},"code":"successful"}`), nil
			}}
			var err error
			if tc.network {
				nsp := &NodeStatusProcessor{proc: proc}
				_, err = nsp.GetAllIssuedESDTs(context.Background(), "", options, page)
			} else {
				ap, _ := NewAccountProcessor(proc, &mock.PubKeyConverterMock{})
				if tc.meta {
					_, err = ap.GetNFTTokenIDsRegisteredByAddress(context.Background(), "deadbeef", options, page)
				} else if tc.withKeys {
					var account *data.AccountModel
					account, err = ap.GetAccountWithKeys(context.Background(), "deadbeef", options, page)
					if err == nil {
						require.Equal(t, "aa", account.BlockInfo.RootHash)
					}
				} else {
					_, err = ap.GetKeyValuePairs(context.Background(), "deadbeef", options, page)
				}
			}
			require.NoError(t, err)
			require.Equal(t, 1, calls)
		})
	}
}

func TestTrieWalkNetworkHistoricalRequestStillAllOrError(t *testing.T) {
	t.Parallel()
	proc := &mock.ProcessorStub{GetObserversCalled: func(uint32, data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
		return []*data.NodeData{{Address: "one"}, {Address: "two"}}, nil
	}, CallTrieWalkEndpointCalled: func(_ context.Context, method, _, _ string, body []byte) (int, []byte, error) {
		require.Equal(t, "POST", method)
		require.JSONEq(t, `{}`, string(body))
		return 200, []byte(`{"data":{"tokens":[],"iteratorState":{"blockNonce":1}},"code":"successful"}`), nil
	}}
	nsp := &NodeStatusProcessor{proc: proc}
	_, err := nsp.GetAllIssuedESDTs(context.Background(), "", common.AccountQueryOptions{BlockNonce: core.OptionalUint64{HasValue: true, Value: 1}}, common.TrieWalkPage{})
	requirePageStatus(t, err, 500)
	// The checkpoint stays raw when decoded into the account model.
	account, err := decodeAccountPage(&data.GenericAPIResponse{Data: json.RawMessage(`{"account":{},"iteratorState":{"blockNonce":18446744073709551615}}`)})
	require.NoError(t, err)
	require.Contains(t, string(account.IteratorState), "18446744073709551615")
}
