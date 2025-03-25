package groups_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/mock"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const addressPath = "/address"

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

type accountResponseData struct {
	Account data.Account `json:"account"`
}

// accountResponse contains the account data and GeneralResponse fields
type accountResponse struct {
	GeneralResponse
	Data accountResponseData
}

type balanceResponseData struct {
	Balance string `json:"balance"`
}

// balanceResponse contains the balance and GeneralResponse fields
type balanceResponse struct {
	GeneralResponse
	Data balanceResponseData
}

type accountsResponseData struct {
	Accounts map[string]*data.Account `json:"accounts"`
}

type accountsResponse struct {
	GeneralResponse
	Data accountsResponseData `json:"data"`
}

type usernameResponseData struct {
	Username string `json:"username"`
}

// usernameResponse contains the username and GeneralResponse fields
type usernameResponse struct {
	GeneralResponse
	Data usernameResponseData
}

type getShardResponseData struct {
	ShardID uint32 `json:"shardID"`
}

type getShardResponse struct {
	GeneralResponse
	Data getShardResponseData
}

type guardianDataApiResponse struct {
	GeneralResponse
	Data guardianDataApiResponseData `json:"data"`
}

type guardianDataApiResponseData struct {
	GuardianData guardianData `json:"guardianData"`
}

type guardianData struct {
	ActiveGuardian  guardian `json:"activeGuardian,omitempty"`
	PendingGuardian guardian `json:"pendingGuardian,omitempty"`
	Guarded         bool     `json:"guarded,omitempty"`
}

type guardian struct {
	Address         string `json:"address"`
	ActivationEpoch uint32 `json:"activationEpoch"`
	ServiceUID      string `json:"serviceUID"`
}

type getEsdtTokensResponseData struct {
	Tokens []string `json:"tokens"`
}

type getEsdtTokensResponse struct {
	GeneralResponse
	Data getEsdtTokensResponseData
}

type esdtTokenData struct {
	TokenIdentifier string `json:"tokenIdentifier"`
	Balance         string `json:"balance"`
	Properties      string `json:"properties"`
}

type esdtNftData struct {
	TokenIdentifier string   `json:"tokenIdentifier"`
	Balance         string   `json:"balance"`
	Properties      string   `json:"properties"`
	Name            string   `json:"name"`
	Creator         string   `json:"creator"`
	Royalties       string   `json:"royalties"`
	Hash            []byte   `json:"hash"`
	URIs            [][]byte `json:"uris"`
	Attributes      []byte   `json:"attributes"`
}

type getEsdtTokenDataResponseData struct {
	TokenData esdtTokenData `json:"tokenData"`
}

type getEsdtTokenDataResponse struct {
	GeneralResponse
	Data getEsdtTokenDataResponseData
}

type getEsdtNftTokenDataResponseData struct {
	TokenData esdtNftData `json:"tokenData"`
}

type getEsdtNftTokenDataResponse struct {
	GeneralResponse
	Data getEsdtNftTokenDataResponseData
}

type getESDTsRolesResponseData struct {
	Roles map[string][]string `json:"roles"`
}

type getESDTsRolesResponse struct {
	GeneralResponse
	Data getESDTsRolesResponseData
}

type getEsdtsWithRoleResponseData struct {
	Tokens []string `json:"tokenData"`
}

type getEsdtsWithRoleResponse struct {
	GeneralResponse
	Data getEsdtsWithRoleResponseData
}

type nonceResponseData struct {
	Nonce uint64 `json:"nonce"`
}

// nonceResponse contains the nonce and GeneralResponse fields
type nonceResponse struct {
	GeneralResponse
	Data nonceResponseData
}

func TestNewAccountGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewAccountsGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestAddressRoute_EmptyTrailReturns404(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("GET", "/address", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

//------- GetAccount

func TestGetAccount_FailWhenFacadeGetAccountFails(t *testing.T) {
	t.Parallel()

	returnedError := "i am an error"
	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return nil, errors.New(returnedError)
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("GET", "/address/test", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountResponse := accountResponse{}
	loadResponse(resp.Body, &accountResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, accountResponse.Data)
	assert.Contains(t, accountResponse.Error, returnedError)
}

func TestGetAccount_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountResponse := accountResponse{}
	loadResponse(resp.Body, &accountResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, accountResponse.Data.Account.Address, reqAddress)
	assert.Equal(t, accountResponse.Data.Account.Nonce, uint64(1))
	assert.Equal(t, accountResponse.Data.Account.Balance, "100")
	assert.Empty(t, accountResponse.Error)
}

//------- GetAccounts

func TestGetAccount_FailsWhenInvalidRequest(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("POST", "/address/bulk", bytes.NewBuffer([]byte(`invalid request`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountsResponse := accountsResponse{}
	loadResponse(resp.Body, &accountsResponse)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Empty(t, accountsResponse.Data)
	assert.Equal(t, accountsResponse.Error, apiErrors.ErrInvalidAddressesArray.Error())
}

func TestGetAccount_FailWhenFacadeGetAccountsFails(t *testing.T) {
	t.Parallel()

	returnedError := "i am an error"
	facade := &mock.FacadeStub{
		GetAccountsHandler: func(addresses []string, _ common.AccountQueryOptions) (*data.AccountsModel, error) {
			return nil, errors.New(returnedError)
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	req, _ := http.NewRequest("POST", "/address/bulk", bytes.NewBuffer([]byte(`["test", "test"]`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountsResponse := accountsResponse{}
	loadResponse(resp.Body, &accountsResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Empty(t, accountsResponse.Data)
	assert.Contains(t, accountsResponse.Error, returnedError)
}

func TestGetAccounts_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	accounts := map[string]*data.Account{
		"erd1alice": {
			Address: "erd1alice",
			Nonce:   1,
			Balance: "100",
		},
		"erd1bob": {
			Address: "erd1bob",
			Nonce:   1,
			Balance: "101",
		},
	}
	facade := &mock.FacadeStub{
		GetAccountsHandler: func(addresses []string, _ common.AccountQueryOptions) (*data.AccountsModel, error) {
			return &data.AccountsModel{
				Accounts: accounts,
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddresses := []string{"erd1alice", "erd1bob"}
	addressBytes, _ := json.Marshal(reqAddresses)
	fmt.Println(string(addressBytes))
	req, _ := http.NewRequest("POST", "/address/bulk", bytes.NewBuffer(addressBytes))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	accountsResponse := accountsResponse{}
	loadResponse(resp.Body, &accountsResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, accountsResponse.Data.Accounts, accounts)
	assert.Empty(t, accountsResponse.Error)
}

//------- GetBalance

func TestGetBalance_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/balance", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	balanceResponse := balanceResponse{}
	loadResponse(resp.Body, &balanceResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, balanceResponse.Data.Balance, "100")
	assert.Empty(t, balanceResponse.Error)
}

//------- GetUsername

func TestGetUsername_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedUsername := "testUser"
	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address:  address,
					Nonce:    1,
					Balance:  "100",
					Username: expectedUsername,
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/username", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	usernameResponse := usernameResponse{}
	loadResponse(resp.Body, &usernameResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedUsername, usernameResponse.Data.Username)
	assert.Empty(t, usernameResponse.Error)
}

//------- GetNonce

func TestGetNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{
		GetAccountHandler: func(address string, _ common.AccountQueryOptions) (*data.AccountModel, error) {
			return &data.AccountModel{
				Account: data.Account{
					Address: address,
					Nonce:   1,
					Balance: "100",
				},
			}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nonce", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	nonceResponse := nonceResponse{}
	loadResponse(resp.Body, &nonceResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, uint64(1), nonceResponse.Data.Nonce)
	assert.Empty(t, nonceResponse.Error)
}

// ---- GetShard

func TestGetShard_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("cannot compute shard ID")
	facade := &mock.FacadeStub{
		GetShardIDForAddressHandler: func(_ string) (uint32, error) {
			return 0, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/shard", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getShardResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetShard_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedShardID := uint32(37)
	facade := &mock.FacadeStub{
		GetShardIDForAddressHandler: func(_ string) (uint32, error) {
			return expectedShardID, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/shard", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getShardResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.ShardID, expectedShardID)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetESDTTokens

func TestGetESDTTokens_FailsWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetAllESDTTokensCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}

	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdt", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getEsdtTokensResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetESDTTokens_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokens := []string{"abc", "def"}
	facade := &mock.FacadeStub{
		GetAllESDTTokensCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtTokensResponseData{Tokens: expectedTokens}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdt", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getEsdtTokensResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedTokens, shardResponse.Data.Tokens)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetGuardianData

func TestGetGuardianData(t *testing.T) {
	t.Parallel()

	expectedGuardianData := guardianDataApiResponseData{
		GuardianData: guardianData{
			ActiveGuardian:  guardian{Address: "address1", ActivationEpoch: 0, ServiceUID: "serviceUID"},
			PendingGuardian: guardian{Address: "address2", ActivationEpoch: 1, ServiceUID: "serviceUID2"},
			Guarded:         false,
		}}

	expectedErr := errors.New("expected error")

	t.Run("internal error", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{
			GetGuardianDataCalled: func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return nil, expectedErr
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)
		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/guardian-data", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)
		shardResponse := data.GenericAPIResponse{}
		loadResponse(resp.Body, &shardResponse)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
	})
	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		facade := &mock.FacadeStub{
			GetGuardianDataCalled: func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return &data.GenericAPIResponse{
					Data: expectedGuardianData,
				}, nil
			},
		}

		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)
		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/guardian-data", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)
		shardResponse := guardianDataApiResponse{}
		loadResponse(resp.Body, &shardResponse)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expectedGuardianData, shardResponse.Data)
		assert.Empty(t, shardResponse.Error)
	})
}

// ---- GetESDTsRoles

func TestGetESDTsRoles_FailsWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetESDTsRolesCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}

	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdts/roles", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getESDTsRolesResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetESDTsRoles_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedRoles := map[string][]string{
		"tkn0": {"role0", "role1"},
		"tkn1": {"role1"},
	}
	facade := &mock.FacadeStub{
		GetESDTsRolesCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getESDTsRolesResponseData{Roles: expectedRoles}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdts/roles", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getESDTsRolesResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.Roles, expectedRoles)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetESDTTokenData

func TestGetESDTTokenData_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetESDTTokenDataCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdt/tkn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getEsdtTokenDataResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetESDTTokenData_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokenData := esdtTokenData{
		TokenIdentifier: "name",
		Balance:         "123",
		Properties:      "1",
	}
	facade := &mock.FacadeStub{
		GetESDTTokenDataCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtTokenDataResponseData{TokenData: expectedTokenData}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdt/tkn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getEsdtTokenDataResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.TokenData, expectedTokenData)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetESDTNftTokenData

func TestGetESDTNftTokenData_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetESDTNftTokenDataCalled: func(_ string, _ string, _ uint64, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nft/tkn/nonce/0", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getEsdtNftTokenDataResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(shardResponse.Error, expectedErr.Error()))
}

func TestGetESDTNftTokenData_FailWhenNonceParamIsInvalid(t *testing.T) {
	t.Parallel()

	facade := &mock.FacadeStub{}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nft/tkn/nonce/qq", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getEsdtNftTokenDataResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.True(t, strings.Contains(response.Error, apiErrors.ErrCannotParseNonce.Error()))
}

func TestGetESDTNftTokenData_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokenData := esdtNftData{
		TokenIdentifier: "name",
		Balance:         "123",
		Properties:      "1",
		Royalties:       "10000",
	}
	facade := &mock.FacadeStub{
		GetESDTNftTokenDataCalled: func(_ string, _ string, _ uint64, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtNftTokenDataResponseData{TokenData: expectedTokenData}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/nft/tkn/nonce/0", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getEsdtNftTokenDataResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.TokenData, expectedTokenData)
	assert.Empty(t, response.Error)
}

// ---- GetESDTsWithRole

func TestGetESDTsWithRole_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetESDTsWithRoleCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdts-with-role/ESDTRoleNFTBurn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	esdtsWithRoleResponse := getEsdtsWithRoleResponse{}
	loadResponse(resp.Body, &esdtsWithRoleResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(esdtsWithRoleResponse.Error, expectedErr.Error()))
}

func TestGetESDTsWithRole_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokens := []string{"FDF-00rr44", "CVC-2598v7"}
	facade := &mock.FacadeStub{
		GetESDTsWithRoleCalled: func(_ string, _ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtsWithRoleResponseData{Tokens: expectedTokens}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdts-with-role/ESDTRoleNFTBurn", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getEsdtsWithRoleResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.Tokens, expectedTokens)
	assert.Empty(t, response.Error)
}

// ---- GetNFTTokenIDsRegisteredByAddress

func TestGetNFTTokenIDsRegisteredByAddress_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetNFTTokenIDsRegisteredByAddressCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/registered-nfts", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	tokensResponse := getEsdtsWithRoleResponse{}
	loadResponse(resp.Body, &tokensResponse)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(tokensResponse.Error, expectedErr.Error()))
}

func TestGetNFTTokenIDsRegisteredByAddress_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedTokens := []string{"FDF-00rr44", "CVC-2598v7"}
	facade := &mock.FacadeStub{
		GetNFTTokenIDsRegisteredByAddressCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtsWithRoleResponseData{Tokens: expectedTokens}}, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/registered-nfts", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := getEsdtsWithRoleResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Data.Tokens, expectedTokens)
	assert.Empty(t, response.Error)
}

// ---- GetKeyValuePairs

func TestGetKeyValuePairs_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetKeyValuePairsHandler: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/keys", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
}

func TestGetKeyValuePairs_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedResponse := &data.GenericAPIResponse{
		Data: map[string]interface{}{"pairs": map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		}},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}
	facade := &mock.FacadeStub{
		GetKeyValuePairsHandler: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return expectedResponse, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/keys", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	actualResponse := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &actualResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Empty(t, actualResponse.Error)
}

// ---- get code hash

func TestGetCodeHash_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := &mock.FacadeStub{
		GetCodeHashCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/code-hash", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
}

func TestGetCodeHash_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	expectedResponse := &data.GenericAPIResponse{
		Data:  "code hash",
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}
	facade := &mock.FacadeStub{
		GetCodeHashCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
			return expectedResponse, nil
		},
	}
	addressGroup, err := groups.NewAccountsGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(addressGroup, addressPath)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/code-hash", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	actualResponse := &data.GenericAPIResponse{}
	loadResponse(resp.Body, &actualResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedResponse, actualResponse)
	assert.Empty(t, actualResponse.Error)
}

func TestAccountsGroup_IsDataTrieMigrated(t *testing.T) {
	t.Parallel()

	t.Run("should return error when facade returns error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("internal err")
		facade := &mock.FacadeStub{
			IsDataTrieMigratedCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return nil, expectedErr
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)

		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/is-data-trie-migrated", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		response := &data.GenericAPIResponse{}
		loadResponse(resp.Body, &response)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
	})

	t.Run("should return successfully", func(t *testing.T) {
		t.Parallel()

		expectedResponse := &data.GenericAPIResponse{
			Data: map[string]interface{}{
				"isMigrated": "true",
			},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		}
		facade := &mock.FacadeStub{
			IsDataTrieMigratedCalled: func(_ string, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return expectedResponse, nil
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)

		reqAddress := "test"
		req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/is-data-trie-migrated", reqAddress), nil)
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		actualResponse := &data.GenericAPIResponse{}
		loadResponse(resp.Body, &actualResponse)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, expectedResponse, actualResponse)
		assert.Empty(t, actualResponse.Error)
	})
}

func TestAccountsGroup_IterateKeys(t *testing.T) {
	t.Parallel()

	t.Run("should return error when facade returns error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("internal err")
		facade := &mock.FacadeStub{
			IterateKeysCalled: func(_ string, _ uint, _ [][]byte, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return nil, expectedErr
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)

		requestData := data.IterateKeysRequest{
			Address:       "asd",
			NumKeys:       10,
			IteratorState: [][]byte{{1, 2, 3}, {4, 5, 6}},
		}
		requestBytes, err := json.Marshal(requestData)
		assert.Nil(t, err)

		req, _ := http.NewRequest("POST", "/address/iterate-keys", bytes.NewBuffer(requestBytes))
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		response := &data.GenericAPIResponse{}
		loadResponse(resp.Body, response)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.True(t, strings.Contains(response.Error, expectedErr.Error()))
	})
	t.Run("should return successfully", func(t *testing.T) {
		t.Parallel()

		pairs := "pairs"
		pairsMap := map[string]string{"a": "b", "c": "d"}
		newIteratorState := "newIteratorState"
		iterState := [][]byte{{7}, {10}}
		expectedResponse := &data.GenericAPIResponse{
			Data: map[string]interface{}{
				pairs:            pairsMap,
				newIteratorState: iterState,
			},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		}
		facade := &mock.FacadeStub{
			IterateKeysCalled: func(_ string, _ uint, _ [][]byte, _ common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
				return expectedResponse, nil
			},
		}
		addressGroup, err := groups.NewAccountsGroup(facade)
		require.NoError(t, err)
		ws := startProxyServer(addressGroup, addressPath)

		requestData := data.IterateKeysRequest{
			Address:       "asd",
			NumKeys:       10,
			IteratorState: [][]byte{{1, 2, 3}, {4, 5, 6}},
		}
		requestBytes, err := json.Marshal(requestData)
		assert.Nil(t, err)

		req, _ := http.NewRequest("POST", "/address/iterate-keys", bytes.NewBuffer(requestBytes))
		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		response := &data.GenericAPIResponse{}
		loadResponse(resp.Body, response)

		assert.Equal(t, http.StatusOK, resp.Code)
		responseMap, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)

		respPairsMap, ok := responseMap[pairs].(map[string]interface{})
		assert.True(t, ok)
		for k, v := range pairsMap {
			assert.Equal(t, v, respPairsMap[k])
		}

		respIterState, ok := responseMap[newIteratorState].([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 2, len(respIterState))
	})
}
