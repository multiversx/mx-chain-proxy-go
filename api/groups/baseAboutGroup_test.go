package groups_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type aboutResponse struct {
	Data  data.AboutInfo `json:"data"`
	Error string         `json:"error"`
	Code  string         `json:"code"`
}

type nodesVersionsResponse struct {
	Data  data.NodesVersionProxyResponseData `json:"data"`
	Error string                             `json:"error"`
	Code  string                             `json:"code"`
}

func TestNewAboutGroup(t *testing.T) {
	t.Parallel()

	t.Run("wrong facade, should fail", func(t *testing.T) {
		t.Parallel()

		wrongFacade := &mock.WrongFacade{}
		group, err := groups.NewAboutGroup(wrongFacade)
		require.Nil(t, group)
		require.Equal(t, groups.ErrWrongTypeAssertion, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		group, err := groups.NewAboutGroup(&mock.FacadeStub{})
		require.Nil(t, err)
		require.NotNil(t, group)
	})
}

func TestAboutGroup_GetAboutInfo(t *testing.T) {
	t.Parallel()

	commitID := "commitID"
	version := "appVersion"

	facade := &mock.FacadeStub{
		GetAboutInfoCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: data.AboutInfo{AppVersion: version, CommitID: commitID},
			}, nil
		},
	}
	aboutGroup, err := groups.NewAboutGroup(facade)
	require.NoError(t, err)

	ws := startProxyServer(aboutGroup, "/about")

	req, _ := http.NewRequest("GET", "/about", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	apiResp := aboutResponse{}
	loadResponse(resp.Body, &apiResp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, apiResp.Data.AppVersion, version)
	assert.Equal(t, apiResp.Data.CommitID, commitID)
	assert.Empty(t, apiResp.Error)
}

func TestAboutGroup_GetNodesVersions(t *testing.T) {
	t.Parallel()

	t.Run("facade error, should err", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("error")
		facade := &mock.FacadeStub{
			GetNodesVersionsCalled: func() (*data.GenericAPIResponse, error) {
				return nil, expectedErr
			},
		}
		aboutGroup, err := groups.NewAboutGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(aboutGroup, "/about")

		req, _ := http.NewRequest("GET", "/about/nodes-versions", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := nodesVersionsResponse{}
		loadResponse(resp.Body, &apiResp)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, apiResp.Error, expectedErr.Error())
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		expectedVersions := map[uint32][]string{
			0:                     {"v1", "v2"},
			1:                     {"v2"},
			2:                     {"v3"},
			core.MetachainShardId: {"v4"},
		}
		facade := &mock.FacadeStub{
			GetNodesVersionsCalled: func() (*data.GenericAPIResponse, error) {
				return &data.GenericAPIResponse{
					Data: data.NodesVersionProxyResponseData{
						Versions: expectedVersions,
					},
				}, nil
			},
		}
		aboutGroup, err := groups.NewAboutGroup(facade)
		require.NoError(t, err)

		ws := startProxyServer(aboutGroup, "/about")

		req, _ := http.NewRequest("GET", "/about/nodes-versions", nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		apiResp := nodesVersionsResponse{}
		loadResponse(resp.Body, &apiResp)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expectedVersions, apiResp.Data.Versions)
	})
}
