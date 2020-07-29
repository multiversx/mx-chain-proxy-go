package block_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-go/api/block"
	"github.com/ElrondNetwork/elrond-proxy-go/api"
	apiBlock "github.com/ElrondNetwork/elrond-proxy-go/api/block"
	"github.com/ElrondNetwork/elrond-proxy-go/api/blockatlas"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type blockResponseData struct {
	Block block.APIBlock `json:"block"`
}

type blockResponse struct {
	Data  blockResponseData `json:"data"`
	Error string            `json:"error"`
	Code  string            `json:"code"`
}

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	blockRoutes := ws.Group("/block")
	apiBlock.Routes(blockRoutes)
	return ws
}

func startNodeServer(handler blockatlas.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	blockRoutes := ws.Group("/block")
	if handler != nil {
		blockRoutes.Use(api.WithElrondProxyFacade(handler))
	}
	apiBlock.Routes(blockRoutes)
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGetBlockByNonce_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/block/0/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, apiErrors.ErrInvalidAppContext.Error(), apiResp.Error)
}

func TestGetBlockByNonce_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

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

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

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

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/block/0/by-nonce/5?withTxs=not_a_bool", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.NotEmpty(t, apiResp.Error)
}

func TestGetBlockByNonce_FailWhenFacadeGetAccountFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := mock.Facade{
		GetBlockByNonceCalled: func(_ uint32, _ uint64, _ bool) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{}, returnedError
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/block/0/by-nonce/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.Equal(t, returnedError.Error(), apiResp.Error)
}

func TestGetBlockByNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	hash := "hashhh"
	facade := mock.Facade{
		GetBlockByNonceCalled: func(_ uint32, _ uint64, _ bool) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: gin.H{"block": block.APIBlock{Nonce: nonce, Hash: hash}},
			}, nil
		},
	}

	ws := startNodeServer(&facade)

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

func TestGetBlockByHash_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/block/0/by-hash/aaa", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, apiErrors.ErrInvalidAppContext.Error(), apiResp.Error)
}

func TestGetBlockByHash_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

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

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

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

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/block/0/by-hash/aaaa?withTxs=not_a_bool", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, apiResp.Data)
	assert.NotEmpty(t, apiResp.Error)
}

func TestGetBlockByHash_FailWhenFacadeGetAccountFails(t *testing.T) {
	t.Parallel()

	returnedError := errors.New("i am an error")
	facade := mock.Facade{
		GetBlockByHashCalled: func(_ uint32, _ string, _ bool) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{}, returnedError
		},
	}
	ws := startNodeServer(&facade)

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
	facade := mock.Facade{
		GetBlockByHashCalled: func(_ uint32, _ string, _ bool) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: gin.H{"block": block.APIBlock{Nonce: nonce, Hash: hash}},
			}, nil
		},
	}

	ws := startNodeServer(&facade)

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
