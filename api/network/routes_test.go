package network_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/network"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
}

type networkResponse struct {
	GeneralResponse
	Metrics map[string]interface{} `json:"message"`
}

func init() {
	gin.SetMode(gin.TestMode)
}

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	networkRoute := ws.Group("/network")
	network.Routes(networkRoute)
	return ws
}

func startNodeServer(handler network.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	networkRoutes := ws.Group("/network")
	if handler != nil {
		networkRoutes.Use(api.WithElrondProxyFacade(handler))
	}
	network.Routes(networkRoutes)
	return ws
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	logError(err)
}

func TestGetNetworkData_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/network/status/0", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := networkResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetNetworkData_NoShardProvidedShouldErr(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/network/status", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := networkResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestGetNetworkData_FacadeFailsShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetNetworkMetricsHandler: func(_ uint32) (map[string]interface{}, error) {
			return nil, errors.New("bad request")
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/network/status/0", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestGetNetworkData_ShouldWork(t *testing.T) {
	t.Parallel()

	respMap := make(map[string]interface{})
	respMap["1"] = "2"
	respMap["2"] = "3"
	facade := mock.Facade{
		GetNetworkMetricsHandler: func(_ uint32) (map[string]interface{}, error) {
			return respMap, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/network/status/0", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result networkResponse
	loadResponse(resp.Body, &result)

	assert.Equal(t, respMap, result.Metrics)
}

func TestEpochMetrics_GetConfigDataBadRequestShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetConfigMetricsHandler: func() (map[string]interface{}, error) {
			return nil, errors.New("bad request")
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/network/config", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestEpochMetrics_GetConfigDataFacadeErrShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	facade := mock.Facade{
		GetConfigMetricsHandler: func() (map[string]interface{}, error) {
			return nil, expectedErr
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/network/config", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var result networkResponse
	loadResponse(resp.Body, &result)

	assert.Equal(t, expectedErr.Error(), result.Error)
}

func TestEpochMetrics_GetConfigDataOkRequestShouldWork(t *testing.T) {
	t.Parallel()

	key := "erd_min_gas_limit"
	value := float64(37)
	facade := mock.Facade{
		GetConfigMetricsHandler: func() (map[string]interface{}, error) {
			return map[string]interface{}{
				key: value,
			}, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/network/config", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result networkResponse
	loadResponse(resp.Body, &result)

	res, ok := result.Metrics[key]
	assert.True(t, ok)
	assert.Equal(t, value, res)
}
