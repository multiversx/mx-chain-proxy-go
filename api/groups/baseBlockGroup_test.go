package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const blockPath = "/block"

func TestNewBlockGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewBlockGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetBlockByNonce_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/invalid_shard_id/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetBlockByNonce_FailWhenNonceParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-nonce/invalid_nonce", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseNonce.Error(), apiResp.Error)
}

func TestGetBlockByNonce_FailWhenWithTxsParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-nonce/5?withTxs=not_a_bool", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.NotEmpty(t, apiResp.Error)
}

func TestGetBlockByNonce_FailWhenFacadeGetBlockByNonceFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetBlockByNonceCalled: func(_ uint32, _ uint64, _ common.BlockQueryOptions) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{}, returnedError
		},
	}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.BlockApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetBlockByNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	hash := "hashhh"
	facade := &mock.Facade{
		GetBlockByNonceCalled: func(_ uint32, _ uint64, _ common.BlockQueryOptions) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{
				Data: data.BlockApiResponsePayload{Block: api.Block{Nonce: nonce, Hash: hash}},
			}, nil
		},
	}

	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := blockResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, apiResp.Data.Block.Nonce, nonce)
	assert.Equal(t, apiResp.Data.Block.Hash, hash)
	assert.Empty(t, apiResp.Error)
}

func TestGetBlockByHash_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/invalid_shard_id/by-hash/aaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetBlockByHash_FailWhenHashParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-hash/invalid-hash", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrInvalidBlockHashParam.Error(), apiResp.Error)
}

func TestGetBlockByHash_FailWhenWithTxsParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-hash/aaaa?withTxs=not_a_bool", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.NotEmpty(t, apiResp.Error)
}

func TestGetBlockByHash_FailWhenFacadeGetBlockByHashFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetBlockByHashCalled: func(_ uint32, _ string, _ common.BlockQueryOptions) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{}, returnedError
		},
	}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetBlockByHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	hash := "hashhh"
	facade := &mock.Facade{
		GetBlockByHashCalled: func(_ uint32, _ string, _ common.BlockQueryOptions) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{
				Data: data.BlockApiResponsePayload{Block: api.Block{Nonce: nonce, Hash: hash}},
			}, nil
		},
	}

	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	req, _ := http.NewRequest("GET", "/block/0/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := blockResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, apiResp.Data.Block.Nonce, nonce)
	assert.Equal(t, apiResp.Data.Block.Hash, hash)
	assert.Empty(t, apiResp.Error)
}

func getAlteredAccounts(t *testing.T, ws *gin.Engine, url string, expectedRespCode int) *data.AlteredAccountsApiResponse {
	req, _ := http.NewRequest("GET", url, nil)
	resp := httptest.NewRecorder()

	ws.ServeHTTP(resp, req)
	require.Equal(t, expectedRespCode, resp.Code)

	apiResp := data.AlteredAccountsApiResponse{}
	loadResponse(resp.Body, &apiResp)
	return &apiResp
}

func TestGetAlteredAccountsByNonce_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	blockGroup, err := groups.NewBlockGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(blockGroup, blockPath)

	t.Run("invalid shard id, should return error", func(t *testing.T) {
		apiResp := getAlteredAccounts(t, ws, "/block/invalid_shard_id/altered-accounts/by-nonce/1", http.StatusBadRequest)

		require.Equal(t, data.ReturnCodeRequestError, apiResp.Code)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
	})

	t.Run("invalid nonce, should return error", func(t *testing.T) {
		apiResp := getAlteredAccounts(t, ws, "/block/0/altered-accounts/by-nonce/invalid_nonce", http.StatusBadRequest)

		require.Equal(t, data.ReturnCodeRequestError, apiResp.Code)
		require.Empty(t, apiResp.Data)
		require.Equal(t, apiErrors.ErrCannotParseNonce.Error(), apiResp.Error)
	})

	t.Run("invalid options, should return error", func(t *testing.T) {
		apiResp := getAlteredAccounts(t, ws, "/block/0/altered-accounts/by-nonce/4?withMetadata=invalid", http.StatusBadRequest)

		require.Equal(t, data.ReturnCodeRequestError, apiResp.Code)
		require.Empty(t, apiResp.Data)
		require.True(t, strings.Contains(apiResp.Error, apiErrors.ErrBadUrlParams.Error()))
	})

	t.Run("could not get response from facade, should return error", func(t *testing.T) {
		expectedError := errors.New("err getting altered accounts")
		invalidFacade := &mock.Facade{
			GetAlteredAccountsByNonceCalled: func(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
				return nil, expectedError
			},
		}
		blockGroupInvalidFacade, err := groups.NewBlockGroup(invalidFacade)
		require.NoError(t, err)

		wsInvalidFacade := startProxyServer(blockGroupInvalidFacade, blockPath)
		apiResp := getAlteredAccounts(t, wsInvalidFacade, "/block/0/altered-accounts/by-nonce/4", http.StatusInternalServerError)

		assert.Equal(t, data.ReturnCodeInternalError, apiResp.Code)
		assert.Empty(t, apiResp.Data)
		assert.Equal(t, expectedError.Error(), apiResp.Error)
	})

	t.Run("should work", func(t *testing.T) {
		alteredAcc1 := &outport.AlteredAccount{
			Address: "addr1",
			Balance: "1000",
			Nonce:   4,
			Tokens: []*outport.AccountTokenData{
				{
					Identifier: "token1",
					Balance:    "10000",
					Nonce:      5,
					Properties: "properties",
				},
			},
		}
		alteredAcc2 := &outport.AlteredAccount{
			Address: "addr2",
			Balance: "4444",
			Nonce:   3333,
			Tokens:  nil,
		}
		expectedApiResponse := &data.AlteredAccountsApiResponse{
			Data: data.AlteredAccountsPayload{
				Accounts: []*outport.AlteredAccount{alteredAcc1, alteredAcc2},
			},
			Error: "",
			Code:  "success",
		}
		facadeValid := &mock.Facade{
			GetAlteredAccountsByNonceCalled: func(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
				return expectedApiResponse, nil
			},
		}
		blockGroupValid, err := groups.NewBlockGroup(facadeValid)
		require.NoError(t, err)

		wsValid := startProxyServer(blockGroupValid, blockPath)

		apiResp := getAlteredAccounts(t, wsValid, "/block/0/altered-accounts/by-nonce/4", http.StatusOK)
		require.Equal(t, expectedApiResponse, apiResp)
	})
}
