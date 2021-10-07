package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func startApiServerResponseLogger(handler groups.AccountsFacadeHandler, respLogMiddleware *responseLoggerMiddleware) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(respLogMiddleware.MiddlewareHandlerFunc())
	accGr, _ := groups.NewAccountsGroup(handler)

	group := ws.Group("/address")
	accGr.RegisterRoutes(group, data.ApiRoutesConfig{}, func(_ *gin.Context) {}, func(_ *gin.Context) {})
	return ws
}

type responseLogFields struct {
	title    string
	path     string
	request  string
	duration time.Duration
	status   int
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
		GetAccountHandler: func(s string) (i *data.Account, e error) {
			time.Sleep(thresholdDuration + 1*time.Millisecond)
			return &data.Account{Balance: "37777"}, nil
		},
	}

	rlf := responseLogFields{}
	printHandler := func(title string, path string, duration time.Duration, status int, request string, response string) {
		rlf.title = title
		rlf.path = path
		rlf.duration = duration
		rlf.status = status
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
		GetAccountHandler: func(_ string) (*data.Account, error) {
			return nil, expectedErr
		},
	}

	rlf := responseLogFields{}
	printHandler := func(title string, path string, duration time.Duration, status int, request string, response string) {
		rlf.title = title
		rlf.path = path
		rlf.duration = duration
		rlf.status = status
		rlf.response = response
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
	assert.True(t, strings.Contains(rlf.response, removeWhitespacesFromString(expectedErr.Error())))
}

func TestResponseLoggerMiddleware_ShouldNotCallHandler(t *testing.T) {
	t.Parallel()

	thresholdDuration := 10000 * time.Millisecond
	facade := mock.FacadeStub{
		GetAccountHandler: func(s string) (i *data.Account, e error) {
			return &data.Account{Balance: "5555"}, nil
		},
	}

	handlerWasCalled := false
	printHandler := func(title string, path string, duration time.Duration, status int, request string, response string) {
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
