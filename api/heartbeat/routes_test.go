package heartbeat_test

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
	"github.com/ElrondNetwork/elrond-proxy-go/api/heartbeat"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
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
	heartbeatRoute := ws.Group("/heartbeat")
	heartbeat.Routes(heartbeatRoute)
	return ws
}

func startNodeServer(handler address.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	heartbeatRoutes := ws.Group("/heartbeat")
	if handler != nil {
		heartbeatRoutes.Use(api.WithElrondProxyFacade(handler))
	}
	heartbeat.Routes(heartbeatRoutes)
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
	req, _ := http.NewRequest("GET", "/heartbeat/", nil)
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

	req, _ := http.NewRequest("GET", "/heartbeat/", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHeartbeat_GetHeartbeatDataReturnsOkResults(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return &data.HeartbeatResponse{Heartbeats: []data.PubKeyHeartbeat{
				{
					NodeDisplayName: "name1",
				},
				{
					NodeDisplayName: "name2",
				},
			}}, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/heartbeat/", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result data.HeartbeatResponse
	loadResponse(resp.Body, &result)
	assert.Equal(t, "name1", result.Heartbeats[0].NodeDisplayName)
	assert.Equal(t, "name2", result.Heartbeats[1].NodeDisplayName)
}

func TestHeartbeat_GetHeartbeatBadRequestShouldErr(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return nil, errors.New("bad request")
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/heartbeat/", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
