package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/require"
)

type statusMetricsResponse struct {
	Data struct {
		Metrics map[string]*data.EndpointMetrics `json:"metrics"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}

const statusPath = "/status"

func TestNewStatusGroup_WrongFacadeShouldErr(t *testing.T) {
	t.Parallel()

	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewStatusGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetMetrics_ShouldErrorIfFacadeReturnsError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("error")
	facade := &mock.Facade{
		GetMetricsCalled: func() (map[string]*data.EndpointMetrics, error) {
			return nil, expectedErr
		},
	}

	statusGroup, err := groups.NewStatusGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(statusGroup, statusPath)

	req, _ := http.NewRequest("GET", "/status/metrics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	var apiResp data.GenericAPIResponse
	loadResponse(resp.Body, &apiResp)
	require.Equal(t, http.StatusInternalServerError, resp.Code)
	require.Equal(t, expectedErr.Error(), apiResp.Error)
}

func TestGetMetrics_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedMetrics := map[string]*data.EndpointMetrics{
		"/network/config": {
			NumRequests:         5,
			NumErrors:           3,
			TotalResponseTime:   100,
			LowestResponseTime:  20,
			HighestResponseTime: 50,
		},
	}
	facade := &mock.Facade{
		GetMetricsCalled: func() (map[string]*data.EndpointMetrics, error) {
			return expectedMetrics, nil
		},
	}

	statusGroup, err := groups.NewStatusGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(statusGroup, statusPath)

	req, _ := http.NewRequest("GET", "/status/metrics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	var apiResp statusMetricsResponse
	loadResponse(resp.Body, &apiResp)
	require.Equal(t, http.StatusOK, resp.Code)

	require.Equal(t, expectedMetrics, apiResp.Data.Metrics)
}
