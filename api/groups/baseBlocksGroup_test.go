package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/require"
)

const blocksPath = "/blocks"

func TestNewBlocksGroup_WrongFacade_ExpectError(t *testing.T) {
	t.Parallel()
	bg, err := groups.NewBlocksGroup(&mock.WrongFacade{})

	require.Nil(t, bg)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestGetBlocksByRound_InvalidRound_ExpectFail(t *testing.T) {
	t.Parallel()

	bg, _ := groups.NewBlocksGroup(&mock.Facade{})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/invalid_round", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.GenericAPIResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Empty(t, apiResp.Data)
	require.Equal(t, apiErrors.ErrCannotParseRound.Error(), apiResp.Error)
}
func TestGetBlocksByRound_InvalidWithTxs_ExpectFail(t *testing.T) {
	t.Parallel()

	bg, _ := groups.NewBlocksGroup(&mock.Facade{})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/0?withTxs=invalid_bool", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.GenericAPIResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Empty(t, apiResp.Data)
	require.NotEmpty(t, apiResp.Error)
}

func TestGetBlocksByRound_InvalidFacadeGetBlocksByRound_ExpectFail(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("local error")
	bg, _ := groups.NewBlocksGroup(&mock.Facade{
		GetBlocksByRoundCalled: func(round uint64, withTxs bool) (*data.BlocksApiResponse, error) {
			return &data.BlocksApiResponse{}, expectedErr
		},
	})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/0?withTxs=true", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.GenericAPIResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Empty(t, apiResp.Data)
	require.Equal(t, expectedErr.Error(), apiResp.Error)
}

func TestGetBlocksByRound_ExpectSuccessful(t *testing.T) {
	t.Parallel()

	blocks := []*data.Block{
		{
			Round: 1,
		},
		{
			Round: 2,
		},
	}

	bg, _ := groups.NewBlocksGroup(&mock.Facade{
		GetBlocksByRoundCalled: func(_ uint64, _ bool) (*data.BlocksApiResponse, error) {
			return &data.BlocksApiResponse{
				Data: data.BlocksApiResponsePayload{
					Blocks: blocks,
				},
			}, nil
		},
	})

	proxyServer := startProxyServer(bg, blocksPath)

	request, _ := http.NewRequest("GET", "/blocks/by-round/0?withTxs=true", nil)
	response := httptest.NewRecorder()
	proxyServer.ServeHTTP(response, request)

	apiResp := data.BlocksApiResponse{}
	loadResponse(response.Body, &apiResp)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, apiResp.Data.Blocks, blocks)
	require.Empty(t, apiResp.Error)
}
