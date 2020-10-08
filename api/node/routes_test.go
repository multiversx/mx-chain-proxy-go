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

func TestHeartbeat_FailsWithWrongFacadeTypeConversion(t *testing.T) {
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
			return &data.HeartbeatResponse{
				Heartbeats: []data.PubKeyHeartbeat{
					{
						NodeDisplayName: name1,
						Identity:        identity1,
					},
					{
						NodeDisplayName: name2,
						Identity:        identity2,
					},
				},
			}, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result data.HeartbeatApiResponse
	loadResponse(resp.Body, &result)
	assert.Equal(t, name1, result.Data.Heartbeats[0].NodeDisplayName)
	assert.Equal(t, name2, result.Data.Heartbeats[1].NodeDisplayName)
	assert.Equal(t, identity1, result.Data.Heartbeats[0].Identity)
	assert.Equal(t, identity2, result.Data.Heartbeats[1].Identity)
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

func TestGetEconomicsData_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/node/economics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	ecResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &ecResp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, ecResp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetEconomicsData_ShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal error")
	facade := mock.Facade{
		GetEconomicsDataMetricsHandler: func() (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/economics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	ecDataResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &ecDataResp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, expectedErr.Error(), ecDataResp.Error)
}

func TestGetEconomicsData_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedResp := data.GenericAPIResponse{Data: map[string]interface{}{"erd_total_supply": "12345"}}
	facade := mock.Facade{
		GetEconomicsDataMetricsHandler: func() (*data.GenericAPIResponse, error) {
			return &expectedResp, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/node/economics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	ecDataResp := data.GenericAPIResponse{}
	loadResponse(resp.Body, &ecDataResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedResp, ecDataResp)
	assert.Equal(t, expectedResp.Data, ecDataResp.Data) //extra safe
}
