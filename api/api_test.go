package api

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestMaxRequestBodySizeMiddleware_OversizedBodyShouldErr(t *testing.T) {
	t.Parallel()

	ws := gin.New()
	ws.Use(maxRequestBodySizeMiddleware(1 << 20))
	ws.POST("/test", func(c *gin.Context) {
		_, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusOK)
	})

	hugeBody := make([]byte, 2*1024*1024)
	for i := range hugeBody {
		hugeBody[i] = 'a'
	}

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(hugeBody))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	require.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)
}

func TestMaxRequestBodySizeMiddleware_NormalBodyShouldPass(t *testing.T) {
	t.Parallel()

	ws := gin.New()
	ws.Use(maxRequestBodySizeMiddleware(1 << 20))
	ws.POST("/test", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if string(body) != `{"key":"value"}` {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	body := []byte(`{"key":"value"}`)
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}
