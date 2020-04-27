package node_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/api/address"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/node"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
}

//heartbeatResponse structure
type heartbeatResponse struct {
	GeneralResponse
	heartbeats data.HeartbeatResponse
}

type nodeMetricsResponse struct {
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
	heartbeatRoute := ws.Group("/node")
	node.Routes(heartbeatRoute)
	return ws
}

func startNodeServer(handler address.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	heartbeatRoutes := ws.Group("/node")
	if handler != nil {
		heartbeatRoutes.Use(api.WithElrondProxyFacade(handler))
	}
	node.Routes(heartbeatRoutes)
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

func TestGetAccount_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := heartbeatResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestHeartbeat_GetHeartbeatDataReturnsStatusOk(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return &data.HeartbeatResponse{Heartbeats: []data.PubKeyHeartbeat{}}, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHeartbeat_GetHeartbeatDataReturnsOkResults(t *testing.T) {
	t.Parallel()

	name1, identity1 := "name1", "identity1"
	name2, identity2 := "name2", "identity2"

	facade := mock.Facade{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return &data.HeartbeatResponse{Heartbeats: []data.PubKeyHeartbeat{
				{
					NodeDisplayName: name1,
					Identity:        identity1,
				},
				{
					NodeDisplayName: name2,
					Identity:        identity2,
				},
			}}, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result data.HeartbeatResponse
	loadResponse(resp.Body, &result)
	assert.Equal(t, name1, result.Heartbeats[0].NodeDisplayName)
	assert.Equal(t, name2, result.Heartbeats[1].NodeDisplayName)
	assert.Equal(t, identity1, result.Heartbeats[0].Identity)
	assert.Equal(t, identity2, result.Heartbeats[1].Identity)
}

func TestHeartbeat_GetHeartbeatBadRequestShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return nil, errors.New("bad request")
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestEpochMetrics_GetEpochDataBadRequestShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetEpochMetricsHandler: func(shardID uint32) (map[string]interface{}, error) {
			return nil, errors.New("bad request")
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/epoch/0", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestEpochMetrics_GetEpochDataNoShardProvidedShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/epoch", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestEpochMetrics_GetEpochDataFacadeErrShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("expected error")
	facade := mock.Facade{
		GetEpochMetricsHandler: func(shardID uint32) (map[string]interface{}, error) {
			return nil, expectedErr
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/epoch/0", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var result nodeMetricsResponse
	loadResponse(resp.Body, &result)

	assert.Equal(t, expectedErr.Error(), result.Error)
}

func TestEpochMetrics_GetEpochDataOkRequestShouldWork(t *testing.T) {
	t.Parallel()

	key := "erd_current_round"
	value := float64(37)
	facade := mock.Facade{
		GetEpochMetricsHandler: func(shardID uint32) (map[string]interface{}, error) {
			return map[string]interface{}{
				key: value,
			}, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/epoch/0", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result nodeMetricsResponse
	loadResponse(resp.Body, &result)

	res, ok := result.Metrics[key]
	assert.True(t, ok)
	assert.Equal(t, value, res)
}
