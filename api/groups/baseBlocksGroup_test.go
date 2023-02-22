package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

const blocksPath = "/blocks"

func TestNewBlocksGroup_WrongFacade_ExpectError(t *testing.T) {
	t.Parallel()
	bg, err := groups.NewBlocksGroup(&mock.WrongFacade{})

	require.Nil(t, bg)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetBlocksByRound_InvalidRound_ExpectFail(t *testing.T) {
	t.Parallel()

	bg, _ := groups.NewBlocksGroup(&mock.FacadeStub{})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/invalid_round", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.GenericAPIResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Empty(t, apiResp.Data)
	require.Equal(t, apiErrors.ErrCannotParseRound.Error(), apiResp.Error)
}

func TestGetBlocksByRound_InvalidWithTxs_ExpectFail(t *testing.T) {
	t.Parallel()

	bg, _ := groups.NewBlocksGroup(&mock.FacadeStub{})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/0?withTxs=invalid_bool", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.GenericAPIResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Empty(t, apiResp.Data)
	require.NotEmpty(t, apiResp.Error)
}

func TestGetBlocksByRound_InvalidFacadeGetBlocksByRound_ExpectFail(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("local error")
	bg, _ := groups.NewBlocksGroup(&mock.FacadeStub{
		GetBlocksByRoundCalled: func(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
			return &data.BlocksApiResponse{}, expectedErr
		},
	})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/0?withTxs=true", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.GenericAPIResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Empty(t, apiResp.Data)
	require.Equal(t, expectedErr.Error(), apiResp.Error)
}

func TestGetBlocksByRound_ExpectSuccessful(t *testing.T) {
	t.Parallel()

	tx1 := transaction.ApiTransactionResult{
		Receiver: "receiver1",
		Sender:   "sender1",
	}
	tx2 := transaction.ApiTransactionResult{
		Receiver: "receiver2",
		Sender:   "sender2",
	}
	tx3 := transaction.ApiTransactionResult{
		Receiver: "receiver3",
		Sender:   "sender3",
	}

	mb1 := api.MiniBlock{
		Hash:         "hash1",
		Transactions: []*transaction.ApiTransactionResult{&tx1, &tx2},
	}
	mb2 := api.MiniBlock{
		Hash:         "hash2",
		Transactions: []*transaction.ApiTransactionResult{&tx3},
	}

	block1 := api.Block{
		Round:      4,
		Hash:       "blockHash1",
		MiniBlocks: []*api.MiniBlock{&mb1, &mb2},
	}
	block2 := api.Block{
		Round: 4,
		Hash:  "blockHash2",
	}

	blocks := []*api.Block{&block1, &block2}

	errGetBlockByRound := errors.New("could not get block by round")
	bg, _ := groups.NewBlocksGroup(&mock.FacadeStub{
		GetBlocksByRoundCalled: func(round uint64, _ common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
			if round == 4 {
				return &data.BlocksApiResponse{
					Data: data.BlocksApiResponsePayload{
						Blocks: blocks,
					},
				}, nil
			}
			return nil, errGetBlockByRound
		},
	})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/4?withTxs=true", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.BlocksApiResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, apiResp.Data.Blocks, blocks)
	require.Empty(t, apiResp.Error)

	request, _ = http.NewRequest("GET", "/blocks/by-round/3?withTxs=true", nil)
	response = httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp2 := data.BlocksApiResponse{}
	loadResponse(response.Body, &apiResp2)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Empty(t, apiResp2.Data)
	require.Equal(t, errGetBlockByRound.Error(), apiResp2.Error)
}

func TestGetBlocksByRound_DifferentWithTxsQueryParams_ExpectWithTxsFlagIsSetCorrectlyInFacade(t *testing.T) {
	t.Parallel()

	tests := []struct {
		URL     string
		withTxs bool
	}{
		{
			URL:     "/blocks/by-round/0",
			withTxs: false,
		},
		{
			URL:     "/blocks/by-round/0?withTxs=false",
			withTxs: false,
		},
		{
			URL:     "/blocks/by-round/0?withTxs=true",
			withTxs: true,
		},
	}

	for _, currTest := range tests {
		bg, _ := groups.NewBlocksGroup(&mock.FacadeStub{
			GetBlocksByRoundCalled: func(_ uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
				require.Equal(t, options.WithTransactions, currTest.withTxs)
				return &data.BlocksApiResponse{}, nil
			},
		})

		proxyServer := startProxyServer(bg, blocksPath)

		request, _ := http.NewRequest("GET", currTest.URL, nil)
		response := httptest.NewRecorder()
		proxyServer.ServeHTTP(response, request)

		apiResp := data.BlocksApiResponse{}
		loadResponse(response.Body, &apiResp)

		require.Equal(t, http.StatusOK, response.Code)
		require.Equal(t, apiResp, data.BlocksApiResponse{})
		require.Empty(t, apiResp.Error)
	}
}
