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

// TODO: handler data types better

// type internalBlockApiResponse struct {
// 	Data  internalBlockApiResponsePayload `json:"data"`
// 	Error string                          `json:"error"`
// 	Code  ReturnCode                      `json:"code"`
// }

// type internalBlockApiResponsePayload struct {
// 	Block data.Block `json:"block"`
// }

const internalPath = "/internal"

func TestNewInternalGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewInternalGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

// ---- InternalBlockByNonce

func TestGetInternalBlockByNonce_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/invalid_shard_id/json/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericInternalApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetInternalBlockByNonce_FailWhenNonceParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-nonce/invalid_nonce", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericInternalApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseNonce.Error(), apiResp.Error)
}

func TestGetInternalBlockByNonce_FailWhenFacadeGetBlockByNonceFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64) (*data.GenericInternalApiResponse, error) {
			return &data.GenericInternalApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericInternalApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetInternalBlockByNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "dummyhash"
	facade := &mock.Facade{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64) (*data.GenericInternalApiResponse, error) {
			return &data.GenericInternalApiResponse{
				Data: &data.Block{Nonce: nonce, Hash: hash},
			}, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := blockResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	// assert.Equal(t, apiResp.Data.Block.Nonce, nonce)
	// assert.Equal(t, apiResp.Data.Block.Hash, hash)
	assert.Empty(t, apiResp.Error)
}

// ---- InternalBlockByHash

func TestGetInternalBlockByHash_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/invalid_shard_id/json/block/by-hash/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericInternalApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetInternalBlockByHash_FailWhenHashParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-hash/invalid-hash", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericInternalApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrInvalidBlockHashParam.Error(), apiResp.Error)
}

func TestGetInternalBlockByHash_FailWhenFacadeGetBlockByHashFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetInternalBlockByHashCalled: func(_ uint32, _ string) (*data.GenericInternalApiResponse, error) {
			return &data.GenericInternalApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericInternalApiResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetInternalBlockByHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "aaaa"

	expectedData := &data.GenericInternalApiResponse{
		Data: data.BlockApiResponsePayload{Block: data.Block{Nonce: nonce, Hash: hash}},
	}

	facade := &mock.Facade{
		GetInternalBlockByHashCalled: func(_ uint32, _ string) (*data.GenericInternalApiResponse, error) {
			return expectedData, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &data.GenericInternalApiResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	//assert.Equal(t, expectedData, apiResp)
	assert.Empty(t, apiResp.Error)
}
