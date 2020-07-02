package address_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/api/address"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
}

//addressResponse structure
type addressResponse struct {
	GeneralResponse
	Balance *big.Int `json:"balance"`
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

type nonceResponseData struct {
	Nonce uint64 `json:"nonce"`
}

// nonceResponse contains the nonce and GeneralResponse fields
type nonceResponse struct {
	GeneralResponse
	Data nonceResponseData
}

func init() {
	gin.SetMode(gin.TestMode)
}

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	addressRoute := ws.Group("/address")
	address.Routes(addressRoute)
	return ws
}

func startNodeServer(handler address.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	addressRoutes := ws.Group("/address")
	if handler != nil {
		addressRoutes.Use(api.WithElrondProxyFacade(handler))
	}
	address.Routes(addressRoutes)
	return ws
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	logError(err)
}

func TestAddressRoute_EmptyTrailReturns404(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	req, _ := http.NewRequest("GET", "/address", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

//------- GetAccount

func TestGetAccount_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/address/empty", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := addressResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetAccount_FailWhenFacadeGetAccountFails(t *testing.T) {
	t.Parallel()

	returnedError := "i am an error"
	facade := mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return nil, errors.New(returnedError)
		},
	}
	ws := startNodeServer(&facade)

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

	facade := mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: address,
				Nonce:   1,
				Balance: "100",
			}, nil
		},
	}
	ws := startNodeServer(&facade)

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

func TestGetBalance_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/address/empty/balance", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := balanceResponse{}
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetBalance_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: address,
				Nonce:   1,
				Balance: "100",
			}, nil
		},
	}
	ws := startNodeServer(&facade)

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

//------- GetNonce

func TestGetNonce_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("GET", "/address/empty/nonce", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	statusRsp := nonceResponse{}
	loadResponse(resp.Body, &statusRsp)
	assert.Equal(t, resp.Code, http.StatusInternalServerError)
	assert.Equal(t, statusRsp.Error, apiErrors.ErrInvalidAppContext.Error())
}

func TestGetNonce_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	facade := mock.Facade{
		GetAccountHandler: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: address,
				Nonce:   1,
				Balance: "100",
			}, nil
		},
	}
	ws := startNodeServer(&facade)

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
