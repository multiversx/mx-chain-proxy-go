package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("api/middleware")

const (
	prefixDurationTooLong      = "[too long]"
	prefixBadRequest           = "[bad request]"
	prefixInternalError        = "[internal error]"
	maxLengthRequestOrResponse = 400
)

// TODO: remove this file and use the same middleware from mx-chain-go after it is merged

type responseLoggerMiddleware struct {
	thresholdDurationForLoggingRequest time.Duration
	printRequestFunc                   func(title string, path string, duration time.Duration, status int, clientIP string, request string, response string)
}

// NewResponseLoggerMiddleware returns a new instance of responseLoggerMiddleware
func NewResponseLoggerMiddleware(thresholdDurationForLoggingRequest time.Duration) *responseLoggerMiddleware {
	rlm := &responseLoggerMiddleware{
		thresholdDurationForLoggingRequest: thresholdDurationForLoggingRequest,
	}

	rlm.printRequestFunc = rlm.printRequest

	return rlm
}

// MiddlewareHandlerFunc logs detail about a request if it is not successful or it's duration is higher than a threshold
func (rlm *responseLoggerMiddleware) MiddlewareHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// read the body for logging purposes and restore it into the context
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		requestBodyString := string(bodyBytes)

		bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bw

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()

		shouldLogRequest := latency > rlm.thresholdDurationForLoggingRequest || c.Writer.Status() != http.StatusOK
		if shouldLogRequest {
			requestBodyString = prepareLog(requestBodyString)
			responseBodyString := prepareLog(bw.body.String())
			rlm.logRequestAndResponse(c, latency, status, requestBodyString, responseBodyString)
		}
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (rlm *responseLoggerMiddleware) IsInterfaceNil() bool {
	return rlm == nil
}

func (rlm *responseLoggerMiddleware) logRequestAndResponse(c *gin.Context, duration time.Duration, status int, request string, response string) {
	title := rlm.computeLogTitle(status)

	rlm.printRequestFunc(title, c.Request.RequestURI, duration, status, c.ClientIP(), request, response)
}

func (rlm *responseLoggerMiddleware) computeLogTitle(status int) string {
	logPrefix := prefixDurationTooLong
	if status == http.StatusBadRequest {
		logPrefix = prefixBadRequest
	} else if status == http.StatusInternalServerError {
		logPrefix = prefixInternalError
	} else if status != http.StatusOK {
		logPrefix = fmt.Sprintf("http code %d", status)
	}

	return fmt.Sprintf("%s api request", logPrefix)
}

func (rlm *responseLoggerMiddleware) printRequest(title string, path string, duration time.Duration, status int, clientIP string, request string, response string) {
	log.Warn(title,
		"path", path,
		"duration", duration,
		"status", status,
		"client IP", clientIP,
		"request", request,
		"response", response,
	)
}

func prepareLog(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}

	result := b.String()
	if len(result) > maxLengthRequestOrResponse {
		return result[:maxLengthRequestOrResponse] + "..."
	}
	return b.String()
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
