package groups_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/groups/v_next"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

type pagingNextFacade struct{}

func (pagingNextFacade) GetAccount(string) (*data.AccountModel, error)              { return nil, nil }
func (pagingNextFacade) GetTransactions(string) ([]data.DatabaseTransaction, error) { return nil, nil }
func (pagingNextFacade) GetShardIDForAddressV_next(string, int) (uint32, error)     { return 3, nil }
func (pagingNextFacade) GetValueForKey(string, string) (string, error)              { return "", nil }
func (pagingNextFacade) NextEndpointHandler() string                                { return "new" }

func TestPagingVersionedGroupStillBuilds(t *testing.T) {
	t.Parallel()
	base, err := groups.NewAccountsGroup(&mock.FacadeStub{})
	require.NoError(t, err)
	_, err = v_next.NewAccountsGroupV_next(base, pagingNextFacade{})
	require.NoError(t, err)
	router := gin.New()
	base.RegisterRoutes(router.Group("/address"), data.ApiRoutesConfig{}, nil, nil, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/address/alice/shard", nil))
	require.Equal(t, 200, w.Code)
	require.Contains(t, w.Body.String(), "3")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest("GET", "/address/alice/nonce", nil))
	require.Equal(t, 404, w.Code)
}
