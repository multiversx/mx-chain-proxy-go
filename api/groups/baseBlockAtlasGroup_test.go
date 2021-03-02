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

const blockAtlasPath = "/blockatlas"

type blockResponseData struct {
	Block data.AtlasBlock `json:"block"`
}

type blockResponse struct {
	Data  blockResponseData `json:"data"`
	Error string            `json:"error"`
	Code  string            `json:"code"`
}

func TestNewBlockAtlasGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewBlockAtlasGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetBlockByShardIDAndNonceFromElastic_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}

	baseBlockAtlasGroup, err := groups.NewBlockAtlasGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(baseBlockAtlasGroup, blockAtlasPath)

	req, _ := http.NewRequest("GET", "/blockatlas/invalid_shard_id/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseShardID.Error(), apiResp.Error)
}

func TestGetBlockByShardIDAndNonceFromElastic_FailWhenNonceParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{}
	baseBlockAtlasGroup, err := groups.NewBlockAtlasGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(baseBlockAtlasGroup, blockAtlasPath)

	req, _ := http.NewRequest("GET", "/blockatlas/0/invalid_nonce", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, apiErrors.ErrCannotParseNonce.Error(), apiResp.Error)
}

func TestGetBlockByShardIDAndNonceFromElastic_FailWhenFacadeGetAccountFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := &mock.Facade{
		GetBlockByShardIDAndNonceHandler: func(_ uint32, _ uint64) (data.AtlasBlock, error) {
			return data.AtlasBlock{}, returnedError
		},
	}
	baseBlockAtlasGroup, err := groups.NewBlockAtlasGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(baseBlockAtlasGroup, blockAtlasPath)

	req, _ := http.NewRequest("GET", "/blockatlas/0/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetBlockByShardIDAndNonceFromElastic_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	hash := "hashhh"
	facade := &mock.Facade{
		GetBlockByShardIDAndNonceHandler: func(_ uint32, _ uint64) (data.AtlasBlock, error) {
			return data.AtlasBlock{
				Nonce: nonce,
				Hash:  hash,
			}, nil
		},
	}

	baseBlockAtlasGroup, err := groups.NewBlockAtlasGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(baseBlockAtlasGroup, blockAtlasPath)

	req, _ := http.NewRequest("GET", "/blockatlas/0/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := blockResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, apiResp.Data.Block.Nonce, nonce)
	assert.Equal(t, apiResp.Data.Block.Hash, hash)
	assert.Empty(t, apiResp.Error)
	assert.Equal(t, string(data.ReturnCodeSuccess), apiResp.Code)
}
