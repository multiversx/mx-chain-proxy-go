package groups

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestPagingRegistryMethodIdentityAndRemoval(t *testing.T) {
	t.Parallel()
	group := &baseGroup{}
	get := func(c *gin.Context) { c.String(200, "get") }
	post := func(c *gin.Context) { c.String(200, "post") }
	require.NoError(t, group.AddEndpoint("/page", data.EndpointHandlerData{Method: "GET", Handler: get}))
	require.NoError(t, group.AddEndpoint("/page", data.EndpointHandlerData{Method: "POST", Handler: post}))
	require.ErrorIs(t, group.AddEndpoint("/page", data.EndpointHandlerData{Method: "POST", Handler: post}), ErrEndpointAlreadyRegistered)
	require.NoError(t, group.UpdateEndpoint("/page", data.EndpointHandlerData{Method: "POST", Handler: func(c *gin.Context) { c.String(200, "updated") }}))
	router := gin.New()
	group.RegisterRoutes(router.Group("/group"), data.ApiRoutesConfig{}, nil, nil, nil)
	for method, want := range map[string]string{"GET": "get", "POST": "updated"} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(method, "/group/page", nil))
		require.Equal(t, want, w.Body.String())
	}
	require.NoError(t, group.RemoveEndpoint("/page"))
	require.Empty(t, group.endpoints)
	require.ErrorIs(t, group.RemoveEndpoint("/page"), ErrHandlerDoesNotExist)
}

func TestPagingRegistrySharesRouteConfiguration(t *testing.T) {
	t.Parallel()
	for _, open := range []bool{false, true} {
		group := &baseGroup{}
		for _, method := range []string{"GET", "POST"} {
			require.NoError(t, group.AddEndpoint("/page", data.EndpointHandlerData{Method: method, Handler: func(c *gin.Context) { c.Status(http.StatusOK) }}))
		}
		auth, limits := 0, 0
		config := data.ApiRoutesConfig{APIPackages: map[string]data.APIPackageConfig{"group": {Routes: []data.RouteConfig{{Name: "/page", Open: open, Secured: true, RateLimit: 1}}}}}
		router := gin.New()
		group.RegisterRoutes(router.Group("/group"), config, func(c *gin.Context) { auth++ }, func(c *gin.Context) { limits++ }, func(c *gin.Context) {})
		for _, method := range []string{"GET", "POST"} {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(method, "/group/page", nil))
			if open {
				require.Equal(t, 200, w.Code)
			} else {
				require.Equal(t, 404, w.Code)
			}
		}
		if open {
			require.Equal(t, 2, auth)
			require.Equal(t, 2, limits)
		} else {
			require.Zero(t, auth)
			require.Zero(t, limits)
		}
	}
}
