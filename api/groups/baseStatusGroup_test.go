package groups_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
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
	facade := &mock.FacadeStub{
		GetMetricsCalled: func() map[string]*data.EndpointMetrics {
			return expectedMetrics
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

func TestGetPrometheusMetrics_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedMetrics := `num_requests{endpoint="/network/config"} 37`
	facade := &mock.FacadeStub{
		GetPrometheusMetricsCalled: func() string {
			return expectedMetrics
		},
	}

	statusGroup, err := groups.NewStatusGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(statusGroup, statusPath)

	req, _ := http.NewRequest("GET", "/status/prometheus-metrics", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, expectedMetrics, string(bodyBytes))
}
