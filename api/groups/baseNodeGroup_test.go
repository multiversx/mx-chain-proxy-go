package groups_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const nodePath = "/node"

func TestNewNodeGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewNodeGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestHeartbeat_GetHeartbeatDataReturnsStatusOk(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return &data.HeartbeatResponse{Heartbeats: []data.PubKeyHeartbeat{}}, nil
		},
	}

	nodeGroup, err := groups.NewNodeGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(nodeGroup, nodePath)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHeartbeat_GetHeartbeatDataReturnsOkResults(t *testing.T) {
	t.Parallel()

	name1, identity1 := "name1", "identity1"
	name2, identity2 := "name2", "identity2"

	facade := &mock.FacadeStub{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return &data.HeartbeatResponse{
				Heartbeats: []data.PubKeyHeartbeat{
					{
						NodeDisplayName: name1,
						Identity:        identity1,
					},
					{
						NodeDisplayName: name2,
						Identity:        identity2,
					},
				},
			}, nil
		},
	}
	nodeGroup, err := groups.NewNodeGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(nodeGroup, nodePath)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var result data.HeartbeatApiResponse
	loadResponse(resp.Body, &result)
	assert.Equal(t, name1, result.Data.Heartbeats[0].NodeDisplayName)
	assert.Equal(t, name2, result.Data.Heartbeats[1].NodeDisplayName)
	assert.Equal(t, identity1, result.Data.Heartbeats[0].Identity)
	assert.Equal(t, identity2, result.Data.Heartbeats[1].Identity)
}

func TestHeartbeat_GetHeartbeatBadRequestShouldErr(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetHeartbeatDataHandler: func() (*data.HeartbeatResponse, error) {
			return nil, errors.New("bad request")
		},
	}
	nodeGroup, err := groups.NewNodeGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(nodeGroup, nodePath)

	req, _ := http.NewRequest("GET", "/node/heartbeatstatus", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestNodeGroup_IsOldStorageToken(t *testing.T) {
	t.Parallel()

	t.Run("should error due to facade error", func(t *testing.T) {
		t.Parallel()

		expectedError := errors.New("expected error")
		facade := &mock.FacadeStub{
			IsOldStorageForTokenCalled: func(_ string, _ uint64) (bool, error) {
				return true, expectedError
			},
		}
		nodeGroup, _ := groups.NewNodeGroup(facade)

		ws := startProxyServer(nodeGroup, nodePath)

		req, _ := http.NewRequest("GET", "/node/old-storage-token/test-token/nonce/37", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		var result data.GenericAPIResponse
		loadResponse(resp.Body, &result)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, expectedError.Error(), result.Error)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{
			IsOldStorageForTokenCalled: func(_ string, _ uint64) (bool, error) {
				return true, nil
			},
		}
		nodeGroup, _ := groups.NewNodeGroup(facade)

		ws := startProxyServer(nodeGroup, nodePath)

		req, _ := http.NewRequest("GET", "/node/old-storage-token/test-token/nonce/37", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		var result data.GenericAPIResponse
		loadResponse(resp.Body, &result)

		assert.Equal(t, http.StatusOK, resp.Code)
		fmt.Printf("%v\n", result.Data)
		assert.Equal(t, "map[isOldStorage:true]", fmt.Sprintf("%v", result.Data))
	})
}

func TestBaseNodeGroup_GetWaitingEpochsLeftForPublicKey(t *testing.T) {
	t.Parallel()

	t.Run("facade returns bad request", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{
			GetWaitingEpochsLeftForPublicKeyCalled: func(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
				return nil, errors.New("bad request")
			},
		}
		nodeGroup, err := groups.NewNodeGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(nodeGroup, nodePath)

		req, _ := http.NewRequest("GET", "/node/waiting-epochs-left/key", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("facade returns bad request", func(t *testing.T) {
		t.Parallel()

		providedData := data.WaitingEpochsLeftResponse{
			EpochsLeft: 10,
		}
		facade := &mock.FacadeStub{
			GetWaitingEpochsLeftForPublicKeyCalled: func(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
				return &data.WaitingEpochsLeftApiResponse{
					Data: providedData,
				}, nil
			},
		}
		nodeGroup, err := groups.NewNodeGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(nodeGroup, nodePath)

		req, _ := http.NewRequest("GET", "/node/waiting-epochs-left/key", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var result data.WaitingEpochsLeftApiResponse
		loadResponse(resp.Body, &result)
		assert.Equal(t, providedData, result.Data)
	})
}
