package groups_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type internalBlockResponseData struct {
	Block testStruct `json:"block"`
}

type internalBlockResponse struct {
	Data  internalBlockResponseData `json:"data"`
	Error string                    `json:"error"`
	Code  string                    `json:"code"`
}

type internalMiniBlockResponseData struct {
	Block testStruct `json:"miniblock"`
}

type internalMiniBlockResponse struct {
	Data  internalMiniBlockResponseData `json:"data"`
	Error string                        `json:"error"`
	Code  string                        `json:"code"`
}

type rawBlockResponseData struct {
	Block []byte `json:"block"`
}

type rawBlockResponse struct {
	Data  rawBlockResponseData `json:"data"`
	Error string               `json:"error"`
	Code  string               `json:"code"`
}

type rawMiniBlockResponseData struct {
	Block []byte `json:"miniblock"`
}

type rawMiniBlockResponse struct {
	Data  rawMiniBlockResponseData `json:"data"`
	Error string                   `json:"error"`
	Code  string                   `json:"code"`
}

type internalValidatorsInfoResponse struct {
	Data struct {
		ValidatorsInfo testStruct `json:"validators"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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
	facade := &mock.FacadeStub{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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
	facade := &mock.FacadeStub{
		GetInternalBlockByHashCalled: func(_ uint32, _ string, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{
		GetInternalBlockByHashCalled: func(_ uint32, _ string, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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
	facade := &mock.FacadeStub{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{
		GetInternalBlockByNonceCalled: func(_ uint32, _ uint64, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{}
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

	facade := &mock.FacadeStub{}
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
	facade := &mock.FacadeStub{
		GetInternalBlockByHashCalled: func(_ uint32, _ string, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

	facade := &mock.FacadeStub{
		GetInternalBlockByHashCalled: func(_ uint32, _ string, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
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

// ---- InternalMiniBlockByHash

func TestGetInternalMiniBlockByHash_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/invalid_shard_id/json/miniblock/by-hash/1/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetInternalMiniBlockByHash_FailWhenHashParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/miniblock/by-hash/invalid-hash/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrInvalidBlockHashParam.Error(), apiResp.Error)
}

func TestGetInternalMiniBlockByHash_FailWhenFacadeGetBlockByHashFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.FacadeStub{
		GetInternalMiniBlockByHashCalled: func(_ uint32, _ string, epoch uint32, _ common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
			return &data.InternalMiniBlockApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/miniblock/by-hash/aaaa/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetInternalMiniBlockByHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "aaaa"

	ts := &testStruct{
		Nonce: nonce,
		Hash:  hash,
	}

	expectedData := &data.InternalMiniBlockApiResponse{
		Data: data.InternalMiniBlockApiResponsePayload{MiniBlock: ts},
	}

	facade := &mock.FacadeStub{
		GetInternalMiniBlockByHashCalled: func(_ uint32, _ string, epoch uint32, _ common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
			return expectedData, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/json/miniblock/by-hash/aaaa/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalMiniBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, nonce, apiResp.Data.Block.Nonce)
	assert.Equal(t, hash, apiResp.Data.Block.Hash)
	assert.Empty(t, apiResp.Error)
}

// ---- RawMiniBlockByHash

func TestGetRawMiniBlockByHash_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/invalid_shard_id/raw/miniblock/by-hash/1/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetRawMiniBlockByHash_FailWhenHashParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/miniblock/by-hash/invalid-hash/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrInvalidBlockHashParam.Error(), apiResp.Error)
}

func TestGetRawMiniBlockByHash_FailWhenEpochParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/miniblock/by-hash/aaaa/epoch/a", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseEpoch.Error(), apiResp.Error)
}

func TestGetRawMiniBlockByHash_FailWhenFacadeGetBlockByHashFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.FacadeStub{
		GetInternalMiniBlockByHashCalled: func(_ uint32, _ string, epoch uint32, _ common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
			return &data.InternalMiniBlockApiResponse{}, returnedError
		},
	}
	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/miniblock/by-hash/aaaa/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &internalBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetRawMiniBlockByHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "aaaa"

	ts := &testStruct{
		Nonce: nonce,
		Hash:  hash,
	}
	tsBytes, err := json.Marshal(ts)
	require.NoError(t, err)

	facade := &mock.FacadeStub{
		GetInternalMiniBlockByHashCalled: func(_ uint32, _ string, epoch uint32, _ common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
			return &data.InternalMiniBlockApiResponse{
				Data: data.InternalMiniBlockApiResponsePayload{MiniBlock: tsBytes},
			}, nil
		},
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(internalGroup, internalPath)

	req, _ := http.NewRequest("GET", "/internal/0/raw/miniblock/by-hash/aaaa/epoch/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &rawMiniBlockResponse{}
	loadResponse(resp.Body, apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, tsBytes, apiResp.Data.Block)
	assert.Empty(t, apiResp.Error)
}

// ---- InternalStartOfEpochMetaBlock

func TestGetInternalStartOfEpochMetaBlock(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	hash := "aaaa"

	t.Run("facade fail to get internal meta block", func(t *testing.T) {
		returnedError := errors.New("i am an error")
		facade := &mock.FacadeStub{
			GetInternalStartOfEpochMetaBlockCalled: func(epoch uint32, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
				return &data.InternalBlockApiResponse{}, returnedError
			},
		}
		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/json/startofepoch/metablock/by-epoch/1", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &internalBlockResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Empty(t, apiResp.Data)
		assert.Equal(t, returnedError.Error(), apiResp.Error)
	})

	t.Run("fail when epoch param is invalid", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{}
		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/raw/startofepoch/metablock/by-epoch/a", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &internalBlockResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Empty(t, apiResp.Data)
		assert.Equal(t, apiErrors.ErrCannotParseEpoch.Error(), apiResp.Error)
	})

	t.Run("internal start of epoch metablock, should work", func(t *testing.T) {
		t.Parallel()

		ts := &testStruct{
			Nonce: nonce,
			Hash:  hash,
		}

		expectedData := &data.InternalBlockApiResponse{
			Data: data.InternalBlockApiResponsePayload{Block: ts},
		}

		facade := &mock.FacadeStub{
			GetInternalStartOfEpochMetaBlockCalled: func(epoch uint32, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
				return expectedData, nil
			},
		}

		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/json/startofepoch/metablock/by-epoch/1", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &internalBlockResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, nonce, apiResp.Data.Block.Nonce)
		assert.Equal(t, hash, apiResp.Data.Block.Hash)
		assert.Empty(t, apiResp.Error)
	})

	t.Run("raw start of epoch metablock, should work", func(t *testing.T) {
		t.Parallel()

		ts := &testStruct{
			Nonce: nonce,
			Hash:  hash,
		}
		tsBytes, err := json.Marshal(ts)
		require.NoError(t, err)

		facade := &mock.FacadeStub{
			GetInternalStartOfEpochMetaBlockCalled: func(epoch uint32, _ common.OutputFormat) (*data.InternalBlockApiResponse, error) {
				return &data.InternalBlockApiResponse{
					Data: data.InternalBlockApiResponsePayload{Block: tsBytes},
				}, nil
			},
		}

		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/raw/startofepoch/metablock/by-epoch/1", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &rawBlockResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, tsBytes, apiResp.Data.Block)
		assert.Empty(t, apiResp.Error)
	})
}

func TestGetInternalStartOfEpochValidatorsInfo(t *testing.T) {
	t.Parallel()

	t.Run("failed when epoch param is invalid", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{}
		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/json/startofepoch/validators/by-epoch/aaa", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &internalValidatorsInfoResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Empty(t, apiResp.Data)
		assert.Equal(t, apiErrors.ErrCannotParseEpoch.Error(), apiResp.Error)
	})

	t.Run("facade error should fail", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		facade := &mock.FacadeStub{
			GetInternalStartOfEpochValidatorsInfoCalled: func(epoch uint32) (*data.ValidatorsInfoApiResponse, error) {
				return nil, expectedErr
			},
		}

		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/json/startofepoch/validators/by-epoch/1", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &internalValidatorsInfoResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, expectedErr.Error(), apiResp.Error)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		ts := &testStruct{
			Nonce: uint64(1),
		}

		facade := &mock.FacadeStub{
			GetInternalStartOfEpochValidatorsInfoCalled: func(epoch uint32) (*data.ValidatorsInfoApiResponse, error) {
				return &data.ValidatorsInfoApiResponse{
					Data: data.InternalStartOfEpochValidators{
						ValidatorsInfo: ts,
					},
				}, nil
			},
		}

		internalGroup, err := groups.NewInternalGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(internalGroup, internalPath)

		req, _ := http.NewRequest("GET", "/internal/json/startofepoch/validators/by-epoch/1", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := &internalValidatorsInfoResponse{}
		loadResponse(resp.Body, apiResp)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Empty(t, apiResp.Error)
	})
}
