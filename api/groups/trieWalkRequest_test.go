package groups

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func pagingContext(method, query, body string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(method, "/address/alice/keys?"+query, strings.NewReader(body))
	return c
}

func TestParseTrieWalkRequest(t *testing.T) {
	t.Parallel()
	cp := `{"blockNonce":9007199254740993,"epoch":7,"iteratorState":["YQ=="],"extra":18446744073709551615}`
	encoded := base64.RawURLEncoding.EncodeToString([]byte(cp))
	for _, method := range []string{http.MethodGet, http.MethodPost} {
		options, page, err := parseTrieWalkRequest(pagingContext(method, "numKeys=0&iteratorState="+encoded+"&blockRootHash=aa&hintEpoch=2", ""), "alice", false)
		require.NoError(t, err)
		require.True(t, page.Paginated)
		require.True(t, page.IsContinuation())
		require.True(t, page.NumKeys.HasValue)
		require.Zero(t, page.NumKeys.Value)
		require.Equal(t, cp, string(page.IteratorState))
		require.Equal(t, []byte{0xaa}, options.BlockRootHash)
		require.EqualValues(t, 2, options.HintEpoch.Value)
	}
	options, page, err := parseTrieWalkRequest(pagingContext("POST", "numKeys=7&onFinalBlock=true&blockHash=aa&iteratorState="+encoded, `{"numKeys":0,"onFinalBlock":false,"blockHash":"","iteratorState":null,"onStartOfEpoch":0,"blockNonce":18446744073709551615,"hintEpoch":0,"blockRootHash":"bb"}`), "alice", false)
	require.NoError(t, err)
	require.Zero(t, page.NumKeys.Value)
	require.False(t, options.OnFinalBlock)
	require.Equal(t, []byte{0xaa}, options.BlockHash)
	require.Equal(t, []byte{0xbb}, options.BlockRootHash)
	require.True(t, options.OnStartOfEpoch.HasValue)
	require.True(t, options.HintEpoch.HasValue)
	require.Equal(t, ^uint64(0), options.BlockNonce.Value)
	require.Equal(t, cp, string(page.IteratorState))
	_, page, err = parseTrieWalkRequest(pagingContext("POST", "numKeys=7&iteratorState="+encoded, `{"numKeys":null,"iteratorState":{"blockNonce":18446744073709551615},"ignored":1}`), "alice", false)
	require.NoError(t, err)
	require.EqualValues(t, 7, page.NumKeys.Value)
	require.JSONEq(t, `{"blockNonce":18446744073709551615}`, string(page.IteratorState))
	for _, body := range []string{"", "null", "{}"} {
		_, page, err = parseTrieWalkRequest(pagingContext("POST", "", body), "alice", false)
		require.NoError(t, err)
		require.True(t, page.Paginated)
		require.False(t, page.IsContinuation())
	}
	_, page, err = parseTrieWalkRequest(pagingContext("GET", "numKeys=&iteratorState=", "not json"), "alice", false)
	require.NoError(t, err)
	require.False(t, page.Paginated)
	_, _, err = parseTrieWalkRequest(pagingContext("POST", "numKeys=-1", `{"numKeys":2}`), "alice", false)
	require.Error(t, err, "invalid query must not be hidden by body override")
}

func TestParseTrieWalkRequestRejectsInvalidInput(t *testing.T) {
	t.Parallel()
	for _, query := range []string{"numKeys=-1", "numKeys=x", "numKeys=18446744073709551616", "iteratorState=***", "iteratorState=" + base64.RawURLEncoding.EncodeToString([]byte("[]")), "iteratorState=" + base64.RawURLEncoding.EncodeToString([]byte("null")), "iteratorState=e30=", "blockHash=zz", "blockNonce=x", "hintEpoch=-1", "onFinalBlock=x"} {
		_, _, err := parseTrieWalkRequest(pagingContext("GET", query, ""), "alice", false)
		require.Error(t, err, query)
	}
	for _, body := range []string{"not json", `{"numKeys":-1}`, `{"numKeys":18446744073709551616}`, `{"numKeys":"2"}`, `{"iteratorState":[]}`, `{"iteratorState":"abc"}`, `{"blockHash":"zz"}`, `{"blockRootHash":"zz"}`, `{"onStartOfEpoch":-1}`, `{"blockNonce":-1}`} {
		_, _, err := parseTrieWalkRequest(pagingContext("POST", "", body), "alice", false)
		require.Error(t, err, body)
	}
	_, _, err := parseTrieWalkRequest(pagingContext("GET", "forced-shard-id=0", ""), "alice", false)
	require.ErrorIs(t, err, ErrForcedShardIDCannotBeProvided)
	options, _, err := parseTrieWalkRequest(pagingContext("GET", "forced-shard-id=0", ""), SystemAccountAddressBech, false)
	require.NoError(t, err)
	require.True(t, options.ForcedShardID.HasValue)
	options, _, err = parseTrieWalkRequest(pagingContext("GET", "forced-shard-id=invalid", ""), "", true)
	require.NoError(t, err)
	require.False(t, options.ForcedShardID.HasValue)
}

func pagingRouter(t *testing.T, facade *mock.FacadeStub) *gin.Engine {
	t.Helper()
	router := gin.New()
	ag, err := NewAccountsGroup(facade)
	require.NoError(t, err)
	ng, err := NewNetworkGroup(facade)
	require.NoError(t, err)
	ag.RegisterRoutes(router.Group("/address"), data.ApiRoutesConfig{}, nil, nil, nil)
	ng.RegisterRoutes(router.Group("/network"), data.ApiRoutesConfig{}, nil, nil, nil)
	return router
}

type pagingContextKey struct{}

func TestPagingRoutePairsPreservePageAndContext(t *testing.T) {
	t.Parallel()
	cp := json.RawMessage(`{"blockNonce":9007199254740993,"blockRootHash":"aa"}`)
	for _, method := range []string{"GET", "POST"} {
		for _, path := range []string{
			"/address/alice?withKeys=true", "/address/alice/keys", "/address/alice/esdt", "/address/alice/esdts/roles", "/address/alice/esdts-with-role/role", "/address/alice/registered-nfts",
			"/network/esdts", "/network/esdt/fungible-tokens", "/network/esdt/semi-fungible-tokens", "/network/esdt/non-fungible-tokens",
		} {
			t.Run(method+path, func(t *testing.T) {
				calls := 0
				callback := func(ctx context.Context, address string, options common.AccountQueryOptions, page common.TrieWalkPage) (*data.GenericAPIResponse, error) {
					calls++
					require.Equal(t, "value", ctx.Value(pagingContextKey{}))
					require.True(t, page.Paginated)
					require.True(t, page.NumKeys.HasValue)
					require.EqualValues(t, 2, page.NumKeys.Value)
					return &data.GenericAPIResponse{Data: json.RawMessage(`{"pairs":{},"iteratorState":` + string(cp) + `}`), Code: data.ReturnCodeSuccess}, nil
				}
				facade := &mock.FacadeStub{GetKeyValuePairsHandler: callback, GetAllESDTTokensCalled: callback, GetESDTsRolesCalled: callback, GetNFTTokenIDsRegisteredByAddressCalled: callback,
					GetESDTsWithRoleCalled: func(ctx context.Context, address, role string, options common.AccountQueryOptions, page common.TrieWalkPage) (*data.GenericAPIResponse, error) {
						require.Equal(t, "role", role)
						return callback(ctx, address, options, page)
					},
					GetAllIssuedESDTsHandler: func(ctx context.Context, kind string, options common.AccountQueryOptions, page common.TrieWalkPage) (*data.GenericAPIResponse, error) {
						return callback(ctx, "", options, page)
					},
					GetAccountWithKeysCalled: func(ctx context.Context, address string, options common.AccountQueryOptions, page common.TrieWalkPage) (*data.AccountModel, error) {
						_, err := callback(ctx, address, options, page)
						return &data.AccountModel{IteratorState: cp, BlockInfo: data.BlockInfo{RootHash: "aa"}}, err
					},
				}
				query := "?numKeys=2"
				if strings.Contains(path, "?") {
					query = "&numKeys=2"
				}
				req := httptest.NewRequest(method, path+query, nil).WithContext(context.WithValue(context.Background(), pagingContextKey{}, "value"))
				w := httptest.NewRecorder()
				pagingRouter(t, facade).ServeHTTP(w, req)
				require.Equal(t, 200, w.Code)
				require.Equal(t, 1, calls)
				require.Contains(t, w.Body.String(), "9007199254740993")
				require.Contains(t, w.Body.String(), `"code":"successful"`)
				if path == "/address/alice?withKeys=true" {
					require.Contains(t, w.Body.String(), `"nonce":0`)
					require.Contains(t, w.Body.String(), `"hash":""`)
				}
			})
		}
	}
}

func TestAccountPagingBranchAndErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		method, path, body string
		paged              bool
		status             int
	}{
		{"GET", "/address/alice?numKeys=-1", "", false, 200},
		{"POST", "/address/alice", "", true, 200},
		{"POST", "/address/alice?withKeys=", "", true, 200},
		{"POST", "/address/alice?withKeys=false&numKeys=-1", "not json", false, 200},
		{"GET", "/address/alice?withKeys=abc", "", false, 400},
		{"GET", "/address/alice/balance?withKeys=abc", "", false, 200},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			plain, paged := 0, 0
			facade := &mock.FacadeStub{GetAccountHandler: func(string, common.AccountQueryOptions) (*data.AccountModel, error) {
				plain++
				return &data.AccountModel{}, nil
			}, GetAccountWithKeysCalled: func(context.Context, string, common.AccountQueryOptions, common.TrieWalkPage) (*data.AccountModel, error) {
				paged++
				return &data.AccountModel{}, nil
			}}
			w := httptest.NewRecorder()
			pagingRouter(t, facade).ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
			require.Equal(t, tc.status, w.Code)
			if tc.status == 200 {
				if tc.paged {
					require.Equal(t, 1, paged)
					require.Zero(t, plain)
				} else {
					require.Equal(t, 1, plain)
					require.Zero(t, paged)
				}
			}
		})
	}
}

func TestPagingHandlerErrorStatusAndCancellation(t *testing.T) {
	t.Parallel()
	for _, status := range []int{400, 413, 429, 500} {
		for _, path := range []string{"/address/alice/keys", "/network/esdts"} {
			failure := &data.ObserverAPIError{StatusCode: status, Message: "node failure"}
			facade := &mock.FacadeStub{GetKeyValuePairsHandler: func(context.Context, string, common.AccountQueryOptions, common.TrieWalkPage) (*data.GenericAPIResponse, error) {
				return nil, failure
			}, GetAllIssuedESDTsHandler: func(context.Context, string, common.AccountQueryOptions, common.TrieWalkPage) (*data.GenericAPIResponse, error) {
				return nil, failure
			}}
			w := httptest.NewRecorder()
			pagingRouter(t, facade).ServeHTTP(w, httptest.NewRequest("POST", path, nil))
			require.Equal(t, status, w.Code)
			code := "bad_request"
			if status == 429 {
				code = "system_busy"
			}
			if status == 500 {
				code = "internal_issue"
			}
			require.Contains(t, w.Body.String(), `"code":"`+code+`"`)
			if path == "/network/esdts" {
				require.Contains(t, w.Body.String(), `"error":"node failure"`)
			}
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	facade := &mock.FacadeStub{GetKeyValuePairsHandler: func(ctx context.Context, _ string, _ common.AccountQueryOptions, _ common.TrieWalkPage) (*data.GenericAPIResponse, error) {
		return nil, ctx.Err()
	}}
	w := httptest.NewRecorder()
	pagingRouter(t, facade).ServeHTTP(w, httptest.NewRequest("POST", "/address/alice/keys", nil).WithContext(ctx))
	require.Empty(t, w.Body.String())
	// A normal facade error retains the existing 500 mapping.
	facade.GetKeyValuePairsHandler = func(context.Context, string, common.AccountQueryOptions, common.TrieWalkPage) (*data.GenericAPIResponse, error) {
		return nil, errors.New("plain failure")
	}
	w = httptest.NewRecorder()
	pagingRouter(t, facade).ServeHTTP(w, httptest.NewRequest("POST", "/address/alice/keys", nil))
	require.Equal(t, 500, w.Code)
}

func TestPagingDoesNotShadowStaticPostRoutes(t *testing.T) {
	t.Parallel()
	facade := &mock.FacadeStub{
		GetAccountsHandler: func([]string, common.AccountQueryOptions) (*data.AccountsModel, error) {
			return &data.AccountsModel{Accounts: map[string]*data.Account{"bulk": {}}}, nil
		},
		IterateKeysCalled: func(string, uint, [][]byte, common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: map[string]string{"legacy": "iterate-keys"}, Code: data.ReturnCodeSuccess}, nil
		},
	}
	router := pagingRouter(t, facade)
	for _, tc := range []struct{ path, body, expected string }{
		{"/address/bulk", `[]`, "bulk"},
		{"/address/iterate-keys", `{"address":"alice","numKeys":2}`, "iterate-keys"},
	} {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest("POST", tc.path, strings.NewReader(tc.body)))
		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), tc.expected)
	}
}
