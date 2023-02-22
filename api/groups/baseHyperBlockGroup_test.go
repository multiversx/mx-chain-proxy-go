package groups_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

const hyperBlockPath = "/hyperblock"

func TestNewHyperBlockGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewHyperBlockGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetHyperblockByHash(t *testing.T) {
	facade := &mock.FacadeStub{
		GetHyperBlockByHashCalled: func(hash string, _ common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
			if hash == "abcd" {
				return data.NewHyperblockApiResponse(api.Hyperblock{
					Nonce: 42,
				}), nil
			}

			return nil, fmt.Errorf("fooError")
		},
	}

	// Get with success
	response := data.HyperblockApiResponse{}
	statusCode := doGet(t, facade, "/hyperblock/by-hash/abcd", &response)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "successful", string(response.Code))
	require.Equal(t, "", response.Error)
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))

	// Block missing
	response = data.HyperblockApiResponse{}
	statusCode = doGet(t, facade, "/hyperblock/by-hash/dbca", &response)
	require.Equal(t, http.StatusInternalServerError, statusCode)
	require.Equal(t, "internal_issue", string(response.Code))
	require.Equal(t, "fooError", response.Error)

	// Bad hash
	response = data.HyperblockApiResponse{}
	statusCode = doGet(t, facade, "/hyperblock/by-hash/badhash", &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Equal(t, "bad_request", string(response.Code))
	require.Equal(t, "invalid block hash parameter", response.Error)
}

func TestGetHyperblockByNonce(t *testing.T) {
	facade := &mock.FacadeStub{
		GetHyperBlockByNonceCalled: func(nonce uint64, _ common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
			if nonce == 42 {
				return data.NewHyperblockApiResponse(api.Hyperblock{
					Nonce: 42,
				}), nil
			}

			return nil, fmt.Errorf("fooError")
		},
	}

	// Get with success
	response := data.HyperblockApiResponse{}
	statusCode := doGet(t, facade, "/hyperblock/by-nonce/42", &response)
	require.Equal(t, http.StatusOK, statusCode)
	require.Equal(t, "successful", string(response.Code))
	require.Equal(t, "", response.Error)
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))

	// Block missing
	response = data.HyperblockApiResponse{}
	statusCode = doGet(t, facade, "/hyperblock/by-nonce/43", &response)
	require.Equal(t, http.StatusInternalServerError, statusCode)
	require.Equal(t, "internal_issue", string(response.Code))
	require.Equal(t, "fooError", response.Error)

	// Bad nonce
	response = data.HyperblockApiResponse{}
	statusCode = doGet(t, facade, "/hyperblock/by-hash/badnonce", &response)
	require.Equal(t, http.StatusBadRequest, statusCode)
	require.Equal(t, "bad_request", string(response.Code))
	require.Equal(t, "invalid block hash parameter", response.Error)
}

func doGet(t *testing.T, facade interface{}, url string, response interface{}) int {
	hyperBlockGroup, err := groups.NewHyperBlockGroup(facade)
	require.NoError(t, err)

	server := startProxyServer(hyperBlockGroup, hyperBlockPath)
	httpRequest, _ := http.NewRequest("GET", url, nil)

	responseRecorder := httptest.NewRecorder()
	server.ServeHTTP(responseRecorder, httpRequest)

	loadResponse(responseRecorder.Body, &response)
	return responseRecorder.Code
}
