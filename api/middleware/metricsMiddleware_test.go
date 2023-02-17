package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	apiMock "github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

var emptyGinHandler = func(_ *gin.Context) {}

func startApiServerMetrics(handler groups.AccountsFacadeHandler, metricsMiddleware *metricsMiddleware) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(metricsMiddleware.MiddlewareHandlerFunc())
	accGr, _ := groups.NewAccountsGroup(handler)

	group := ws.Group("/address")
	accGr.RegisterRoutes(group, data.ApiRoutesConfig{}, emptyGinHandler, emptyGinHandler, emptyGinHandler)
	return ws
}

func TestNewMetricsMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("nil status metrics exporter - should err", func(t *testing.T) {
		t.Parallel()

		mm, err := NewMetricsMiddleware(nil)
		require.Nil(t, mm)
		require.Equal(t, ErrNilStatusMetricsExtractor, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		mm, err := NewMetricsMiddleware(&apiMock.StatusMetricsExporterStub{})
		require.NoError(t, err)
		require.NotNil(t, mm)
	})
}

func TestMetricsMiddleware_MiddlewareHandlerFunc(t *testing.T) {
	t.Parallel()

	type receivedRequestData struct {
		path      string
		withError bool
		duration  time.Duration
	}
	receivedData := make([]*receivedRequestData, 0)
	mm, err := NewMetricsMiddleware(&apiMock.StatusMetricsExporterStub{
		AddRequestDataCalled: func(path string, withError bool, duration time.Duration) {
			receivedData = append(receivedData, &receivedRequestData{
				path:      path,
				withError: withError,
				duration:  duration,
			})
		},
	})
	require.NoError(t, err)

	facade := &apiMock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}

	ws := startApiServerMetrics(facade, mm)

	resp := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(resp)
	req, _ := http.NewRequestWithContext(context, "GET", "/address/test", nil)
	ws.ServeHTTP(resp, req)

	require.Len(t, receivedData, 1)
	require.Equal(t, "/address/:address", receivedData[0].path)
	require.False(t, receivedData[0].withError)
}
