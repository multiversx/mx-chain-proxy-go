package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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
	facade := &mock.FacadeStub{
		GetBlockByNonceCalled: func(_ uint32, _ uint64, _ bool) (*data.BlockApiResponse, error) {
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
	facade := &mock.FacadeStub{
		GetBlockByNonceCalled: func(_ uint32, _ uint64, _ bool) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{
				Data: data.BlockApiResponsePayload{Block: data.Block{Nonce: nonce, Hash: hash}},
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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
	facade := &mock.FacadeStub{
		GetBlockByHashCalled: func(_ uint32, _ string, _ bool) (*data.BlockApiResponse, error) {
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
	facade := &mock.FacadeStub{
		GetBlockByHashCalled: func(_ uint32, _ string, _ bool) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{
				Data: data.BlockApiResponsePayload{Block: data.Block{Nonce: nonce, Hash: hash}},
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
