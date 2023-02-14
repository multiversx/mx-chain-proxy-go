package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
)

func startApiServerResponseLogger(handler groups.AccountsFacadeHandler, respLogMiddleware *responseLoggerMiddleware) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(respLogMiddleware.MiddlewareHandlerFunc())
	accGr, _ := groups.NewAccountsGroup(handler)

	group := ws.Group("/address")
	accGr.RegisterRoutes(group, data.ApiRoutesConfig{}, emptyGinHandler, emptyGinHandler, emptyGinHandler)
	return ws
}

type responseLogFields struct {
	title    string
	path     string
	request  string
	duration time.Duration
	status   int
	clientIP string
	response string
}

func TestNewResponseLoggerMiddleware(t *testing.T) {
	t.Parallel()

	rlm := NewResponseLoggerMiddleware(10)

	assert.False(t, check.IfNil(rlm))
}

func TestResponseLoggerMiddleware_DurationExceedsTimeout(t *testing.T) {
	t.Parallel()

	thresholdDuration := 10 * time.Millisecond
	addr := "testAddress"
	facade := mock.FacadeStub{
		GetAccountHandler: func(s string, _ common.AccountQueryOptions) (i *data.AccountModel, e error) {
			time.Sleep(thresholdDuration + 1*time.Millisecond)
			return &data.AccountModel{
				Account: data.Account{
					Balance: "37777",
				},
			}, nil
		},
	}

	rlf := responseLogFields{}
	printHandler := func(title string, path string, duration time.Duration, status int, clientIP string, request string, response string) {
		rlf.title = title
		rlf.path = path
		rlf.duration = duration
		rlf.status = status
		rlf.clientIP = clientIP
		rlf.response = response
		rlf.request = request
	}

	rlm := NewResponseLoggerMiddleware(thresholdDuration)
	rlm.printRequestFunc = printHandler

	ws := startApiServerResponseLogger(&facade, rlm)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/balance", addr), nil)
	req.RemoteAddr = "bad address"
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, strings.Contains(rlf.title, prefixDurationTooLong))
	assert.True(t, rlf.duration > thresholdDuration)
	assert.Equal(t, http.StatusOK, rlf.status)
	assert.True(t, strings.Contains(rlf.response, "37777"))
}

func TestResponseLoggerMiddleware_InternalError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	thresholdDuration := 10000 * time.Millisecond
	facade := mock.FacadeStub{
		GetAccountHandler: func(_ string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return nil, expectedErr
		},
	}

	rlf := responseLogFields{}
	printHandler := func(title string, path string, duration time.Duration, status int, clientIP string, request string, response string) {
		rlf.title = title
		rlf.path = path
		rlf.duration = duration
		rlf.status = status
		rlf.response = response
		rlf.clientIP = clientIP
		rlf.request = request
	}

	rlm := NewResponseLoggerMiddleware(thresholdDuration)
	rlm.printRequestFunc = printHandler

	ws := startApiServerResponseLogger(&facade, rlm)

	req, _ := http.NewRequest("GET", "/address/addr/balance", nil)
	req.RemoteAddr = "bad address"
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(rlf.title, prefixInternalError))
	assert.True(t, rlf.duration < thresholdDuration)
	assert.Equal(t, http.StatusInternalServerError, rlf.status)
	assert.True(t, strings.Contains(rlf.response, prepareLog(expectedErr.Error())))
}

func TestResponseLoggerMiddleware_ShouldNotCallHandler(t *testing.T) {
	t.Parallel()

	thresholdDuration := 10000 * time.Millisecond
	facade := mock.FacadeStub{
		GetAccountHandler: func(s string, _ common.AccountQueryOptions) (i *data.AccountModel, e error) {
			return &data.AccountModel{
				Account: data.Account{
					Balance: "5555",
				},
			}, nil
		},
	}

	handlerWasCalled := false
	printHandler := func(title string, path string, duration time.Duration, status int, clientIP string, request string, response string) {
		handlerWasCalled = true
	}

	rlm := NewResponseLoggerMiddleware(thresholdDuration)
	rlm.printRequestFunc = printHandler

	ws := startApiServerResponseLogger(&facade, rlm)

	req, _ := http.NewRequest("GET", "/address/addr/balance", nil)
	req.RemoteAddr = "bad address"
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.False(t, handlerWasCalled)
}
