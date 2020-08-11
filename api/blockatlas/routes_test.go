package blockatlas_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/api/blockatlas"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type blockResponseData struct {
	Block data.AtlasBlock `json:"block"`
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
	blockAtlasRoutes := ws.Group("/blockatlas")
	blockatlas.Routes(blockAtlasRoutes)
	return ws
}

func startNodeServer(handler blockatlas.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	addressRoutes := ws.Group("/blockatlas")
	if handler != nil {
		addressRoutes.Use(api.WithElrondProxyFacade(handler))
	}
	blockatlas.Routes(addressRoutes)
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGetBlockByShardIDAndNonceFromElastic_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/blockatlas/0/1", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, apiErrors.ErrInvalidAppContext.Error(), apiResp.Error)
}

func TestGetBlockByShardIDAndNonceFromElastic_FailWhenShardParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

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

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

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
	facade := mock.Facade{
		GetBlockByShardIDAndNonceHandler: func(_ uint32, _ uint64) (data.AtlasBlock, error) {
			return data.AtlasBlock{}, returnedError
		},
	}
	ws := startNodeServer(&facade)

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
	facade := mock.Facade{
		GetBlockByShardIDAndNonceHandler: func(_ uint32, _ uint64) (data.AtlasBlock, error) {
			return data.AtlasBlock{
				Nonce: nonce,
				Hash:  hash,
			}, nil
		},
	}

	ws := startNodeServer(&facade)

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
