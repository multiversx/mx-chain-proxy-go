package groups_test

import (
	"encoding/json"
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

// TODO: add unit tests for raw data

type internalBlockResponseData struct {
	Block testStruct `json:"block"`
}

type internalBlockResponse struct {
	Data  internalBlockResponseData `json:"data"`
	Error string                    `json:"error"`
	Code  string                    `json:"code"`
}

type rawBlockResponseData struct {
	Block []byte `json:"block"`
}

type rawBlockResponse struct {
	Data  rawBlockResponseData `json:"data"`
	Error string               `json:"error"`
	Code  string               `json:"code"`
}

type testStruct struct {
	Nonce uint64
	Hash  string
}

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

	apiResp := &internalBlockResponse{}
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

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseNonce.Error(), apiResp.Error)
}

func TestGetInternalBlockByNonce_FailWhenFacadeGetBlockByNonceFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetInternalBlockByNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "dummyhash"

	ts := &testStruct{
		Nonce: nonce,
		Hash:  hash,
	}

	facade := &mock.Facade{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{
				Data: data.InternalBlockApiResponsePayload{Block: ts},
			}, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, nonce, apiResp.Data.Block.Nonce)
	assert.Equal(t, hash, apiResp.Data.Block.Hash)
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

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

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

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrInvalidBlockHashParam.Error(), apiResp.Error)
}

func TestGetInternalBlockByHash_FailWhenFacadeGetBlockByHashFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetInternalBlockByHashCalled: func(_ uint32, _ string) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetInternalBlockByHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "aaaa"

	ts := &testStruct{
		Nonce: nonce,
		Hash:  hash,
	}

	expectedData := &data.InternalBlockApiResponse{
		Data: data.InternalBlockApiResponsePayload{Block: ts},
	}

	facade := &mock.Facade{
		GetInternalBlockByHashCalled: func(_ uint32, _ string) (*data.InternalBlockApiResponse, error) {
			return expectedData, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/block/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, nonce, apiResp.Data.Block.Nonce)
	assert.Equal(t, hash, apiResp.Data.Block.Hash)
	assert.Empty(t, apiResp.Error)
}

// ---- RawBlockByNonce

func TestGetRawBlockByNonce_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/invalid_shard_id/raw/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetRawBlockByNonce_FailWhenNonceParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/block/by-nonce/invalid_nonce", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseNonce.Error(), apiResp.Error)
}

func TestGetRawBlockByNonce_FailWhenFacadeGetBlockByNonceFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetRawBlockByNonceCalled: func(_ uint32, _ uint64) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetRawBlockByNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "dummyhash"

	ts := &testStruct{
		Nonce: nonce,
		Hash:  hash,
	}
	tsBytes, err := json.Marshal(ts)
	require.NoError(t, err)

	facade := &mock.Facade{
		GetRawBlockByNonceCalled: func(_ uint32, _ uint64) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{
				Data: data.InternalBlockApiResponsePayload{Block: tsBytes},
			}, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/block/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &rawBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, tsBytes, apiResp.Data.Block)
	assert.Empty(t, apiResp.Error)
}

// ---- RawBlockByHash

func TestGetRawBlockByHash_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/invalid_shard_id/raw/block/by-hash/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetRawBlockByHash_FailWhenHashParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/block/by-hash/invalid-hash", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrInvalidBlockHashParam.Error(), apiResp.Error)
}

func TestGetRawBlockByHash_FailWhenFacadeGetBlockByHashFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetRawBlockByHashCalled: func(_ uint32, _ string) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/block/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetRawBlockByHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "aaaa"

	ts := &testStruct{
		Nonce: nonce,
		Hash:  hash,
	}
	tsBytes, err := json.Marshal(ts)
	require.NoError(t, err)

	facade := &mock.Facade{
		GetRawBlockByHashCalled: func(_ uint32, _ string) (*data.InternalBlockApiResponse, error) {
			return &data.InternalBlockApiResponse{
				Data: data.InternalBlockApiResponsePayload{Block: tsBytes},
			}, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/block/by-hash/aaaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &rawBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, tsBytes, apiResp.Data.Block)
	assert.Empty(t, apiResp.Error)
}
