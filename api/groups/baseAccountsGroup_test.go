package groups_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
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

type getEsdtTokenDataResponseData struct {
	TokenData esdtTokenData `json:"tokenData"`
}

type getEsdtTokenDataResponse struct {
	GeneralResponse
	Data getEsdtTokenDataResponseData
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

	facade := &mock.Facade{}
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
	facade := &mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
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
	assert.Equal(t, returnedError, accountResponse.Error)
}

func TestGetAccount_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: address,
				Nonce:   1,
				Balance: "100",
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

//------- GetBalance

func TestGetBalance_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := &mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: address,
				Nonce:   1,
				Balance: "100",
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
	facade := &mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address:  address,
				Nonce:    1,
				Balance:  "100",
				Username: expectedUsername,
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

	facade := &mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: address,
				Nonce:   1,
				Balance: "100",
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
	facade := &mock.Facade{
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
	facade := &mock.Facade{
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

func TestGetESDTTokens_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/address/address/esdt", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := getEsdtTokensResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetESDTTokens_FailsWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := mock.Facade{
		GetAllESDTTokensCalled: func(_ string) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	ws := startNodeServer(&facade)

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
	facade := mock.Facade{
		GetAllESDTTokensCalled: func(_ string) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtTokensResponseData{Tokens: expectedTokens}}, nil
		},
	}
	ws := startNodeServer(&facade)

	reqAddress := "test"
	req, _ := http.NewRequest("GET", fmt.Sprintf("/address/%s/esdt", reqAddress), nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	shardResponse := getEsdtTokensResponse{}
	loadResponse(resp.Body, &shardResponse)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, shardResponse.Data.Tokens, expectedTokens)
	assert.Empty(t, shardResponse.Error)
}

// ---- GetESDTTokenData

func TestGetESDTTokenData_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/address/address/esdt/tkn", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := getEsdtTokenDataResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetESDTTokenData_FailWhenFacadeErrors(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("internal err")
	facade := mock.Facade{
		GetESDTTokenDataCalled: func(_ string, _ string) (*data.GenericAPIResponse, error) {
			return nil, expectedErr
		},
	}
	ws := startNodeServer(&facade)

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
	facade := mock.Facade{
		GetESDTTokenDataCalled: func(_ string, _ string) (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{Data: getEsdtTokenDataResponseData{TokenData: expectedTokenData}}, nil
		},
	}
	ws := startNodeServer(&facade)

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
