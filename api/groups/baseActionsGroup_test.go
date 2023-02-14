package groups_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const actionsPath = "/actions"

func TestNewActionsGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewActionsGroup(wrongFacade)

	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestActions_ReloadObserversFailWithBadRequest(t *testing.T) {
	t.Parallel()

	expectedErrMsg := "request err"
	descriptionErr := "description for issue"
	facade := &mock.FacadeStub{
		ReloadObserversCalled: func() data.NodesReloadResponse {
			return data.NodesReloadResponse{
				OkRequest:   false,
				Error:       expectedErrMsg,
				Description: descriptionErr,
			}
		},
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(actionsGroup, actionsPath)

	req, _ := http.NewRequest("POST", "/actions/reload-observers", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, response)
	assert.Equal(t, descriptionErr, response.Data.(string))
	assert.Equal(t, expectedErrMsg, response.Error)
}

func TestActions_ReloadObserversFailWithInternalError(t *testing.T) {
	t.Parallel()

	expectedErrMsg := "internal err"
	descriptionErr := "description for issue"
	facade := &mock.FacadeStub{
		ReloadObserversCalled: func() data.NodesReloadResponse {
			return data.NodesReloadResponse{
				OkRequest:   true,
				Error:       expectedErrMsg,
				Description: descriptionErr,
			}
		},
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(actionsGroup, actionsPath)

	req, _ := http.NewRequest("POST", "/actions/reload-observers", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, response)
	assert.Equal(t, descriptionErr, response.Data.(string))
	assert.Equal(t, expectedErrMsg, response.Error)
}

func TestActions_ReloadObserversShouldWork(t *testing.T) {
	t.Parallel()

	description := "description"
	facade := &mock.FacadeStub{
		ReloadObserversCalled: func() data.NodesReloadResponse {
			return data.NodesReloadResponse{
				OkRequest:   true,
				Error:       "",
				Description: description,
			}
		},
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(actionsGroup, actionsPath)

	req, _ := http.NewRequest("POST", "/actions/reload-observers", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, response)
	assert.Equal(t, description, response.Data.(string))
	assert.Equal(t, "", response.Error)
}

func TestActions_ReloadFullHistoryObserversFailWithBadRequest(t *testing.T) {
	t.Parallel()

	expectedErrMsg := "request err"
	descriptionErr := "description for issue"
	facade := &mock.FacadeStub{
		ReloadFullHistoryObserversCalled: func() data.NodesReloadResponse {
			return data.NodesReloadResponse{
				OkRequest:   false,
				Error:       expectedErrMsg,
				Description: descriptionErr,
			}
		},
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(actionsGroup, actionsPath)

	req, _ := http.NewRequest("POST", "/actions/reload-full-history-observers", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, response)
	assert.Equal(t, descriptionErr, response.Data.(string))
	assert.Equal(t, expectedErrMsg, response.Error)
}

func TestActions_ReloadFullHistoryObserversFailWithInternalError(t *testing.T) {
	t.Parallel()

	expectedErrMsg := "internal err"
	descriptionErr := "description for issue"
	facade := &mock.FacadeStub{
		ReloadFullHistoryObserversCalled: func() data.NodesReloadResponse {
			return data.NodesReloadResponse{
				OkRequest:   true,
				Error:       expectedErrMsg,
				Description: descriptionErr,
			}
		},
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(actionsGroup, actionsPath)

	req, _ := http.NewRequest("POST", "/actions/reload-full-history-observers", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, response)
	assert.Equal(t, descriptionErr, response.Data.(string))
	assert.Equal(t, expectedErrMsg, response.Error)
}

func TestActions_ReloadFullHistoryObserversShouldWork(t *testing.T) {
	t.Parallel()

	description := "description"
	facade := &mock.FacadeStub{
		ReloadFullHistoryObserversCalled: func() data.NodesReloadResponse {
			return data.NodesReloadResponse{
				OkRequest:   true,
				Error:       "",
				Description: description,
			}
		},
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(actionsGroup, actionsPath)

	req, _ := http.NewRequest("POST", "/actions/reload-full-history-observers", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, response)
	assert.Equal(t, description, response.Data.(string))
	assert.Equal(t, "", response.Error)
}
