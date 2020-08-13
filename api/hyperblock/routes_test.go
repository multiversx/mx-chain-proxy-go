package hyperblock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetHyperblockByHash(t *testing.T) {
	facade := mock.Facade{
		GetHyperBlockByHashCalled: func(hash string) (*data.HyperblockApiResponse, error) {
			if hash == "abcd" {
				return data.NewHyperblockApiResponse(data.Hyperblock{
					Nonce: 42,
				}), nil
			}

			return nil, fmt.Errorf("fooError")
		},
	}

	// Get with success
	response := data.HyperblockApiResponse{}
	statusCode := doGet(&facade, "/hyperblock/by-hash/abcd", &response)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "successful", string(response.Code))
	require.Equal(t, "", response.Error)
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))

	// Block missing
	response = data.HyperblockApiResponse{}
	statusCode = doGet(&facade, "/hyperblock/by-hash/dbca", &response)
	require.Equal(t, http.StatusInternalServerError, statusCode)
	require.Equal(t, "internal_issue", string(response.Code))
	require.Equal(t, "fooError", response.Error)

	// Bad hash
	response = data.HyperblockApiResponse{}
	statusCode = doGet(&facade, "/hyperblock/by-hash/badhash", &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Equal(t, "bad_request", string(response.Code))
	require.Equal(t, "invalid block hash parameter", response.Error)
}

func TestGetHyperblockByNonce(t *testing.T) {
	facade := mock.Facade{
		GetHyperBlockByNonceCalled: func(nonce uint64) (*data.HyperblockApiResponse, error) {
			if nonce == 42 {
				return data.NewHyperblockApiResponse(data.Hyperblock{
					Nonce: 42,
				}), nil
			}

			return nil, fmt.Errorf("fooError")
		},
	}

	// Get with success
	response := data.HyperblockApiResponse{}
	statusCode := doGet(&facade, "/hyperblock/by-nonce/42", &response)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "successful", string(response.Code))
	require.Equal(t, "", response.Error)
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))

	// Block missing
	response = data.HyperblockApiResponse{}
	statusCode = doGet(&facade, "/hyperblock/by-nonce/43", &response)
	require.Equal(t, http.StatusInternalServerError, statusCode)
	require.Equal(t, "internal_issue", string(response.Code))
	require.Equal(t, "fooError", response.Error)

	// Bad nonce
	response = data.HyperblockApiResponse{}
	statusCode = doGet(&facade, "/hyperblock/by-hash/badnonce", &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Equal(t, "bad_request", string(response.Code))
	require.Equal(t, "invalid block hash parameter", response.Error)
}

func doGet(facade interface{}, url string, response interface{}) int {
	server := startNodeServer(facade)
	httpRequest, _ := http.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	server.ServeHTTP(responseRecorder, httpRequest)

	parseResponse(responseRecorder.Body, &response)
	return responseRecorder.Code
}

func startNodeServer(handler interface{}) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	route := ws.Group("/hyperblock")
	route.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", handler)
		c.Next()
	})
	Routes(route)

	return ws
}

func parseResponse(responseBody io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(responseBody)

	err := jsonParser.Decode(destination)
	if err != nil {
		fmt.Println(err)
	}
}
