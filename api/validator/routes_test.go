package validator_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/validator"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
}

type valStatsResponseData struct {
	Statistics map[string]*data.ValidatorApiResponse `json:"statistics"`
}

// ValStatsResponse structure
type ValStatsResponse struct {
	Error string               `json:"error"`
	Data  valStatsResponseData `json:"data"`
}

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	validatorRoute := ws.Group("/validator")
	validator.Routes(validatorRoute)
	return ws
}

func startNodeServer(handler validator.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	transactionRoute := ws.Group("/validator")
	if handler != nil {
		transactionRoute.Use(api.WithElrondProxyFacade(handler))
	}
	validator.Routes(transactionRoute)
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		logError(err)
	}
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func TestValidatorStatistics_ErrorWithWrongFacade(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/validator/statistics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}

func TestValidatorStatistics_ShouldErrWhenFacadeFails(t *testing.T) {
	t.Parallel()

	errStr := "expected err"
	facade := mock.Facade{
		ValidatorStatisticsHandler: func() (map[string]*data.ValidatorApiResponse, error) {
			return nil, errors.New(errStr)
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/validator/statistics", nil)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.True(t, strings.Contains(response.Error, errStr))
}

func TestValidatorStatistics_ShouldWork(t *testing.T) {
	t.Parallel()

	valStatsMap := make(map[string]*data.ValidatorApiResponse)
	valStatsMap["statistics"] = &data.ValidatorApiResponse{
		NumLeaderSuccess:         4,
		NumLeaderFailure:         5,
		NumValidatorSuccess:      6,
		NumValidatorFailure:      7,
		Rating:                   0.5,
		TempRating:               0.51,
		TotalNumLeaderSuccess:    4,
		TotalNumLeaderFailure:    5,
		TotalNumValidatorSuccess: 6,
		TotalNumValidatorFailure: 7,
	}
	facade := mock.Facade{
		ValidatorStatisticsHandler: func() (map[string]*data.ValidatorApiResponse, error) {
			return valStatsMap, nil
		},
	}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/validator/statistics", nil)

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := ValStatsResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.Statistics["statistics"], valStatsMap["statistics"])
}
