package wallet_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/wallet"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	walletRoute := ws.Group("/wallet")
	wallet.Routes(walletRoute)
	return ws
}

func startNodeServer(handler wallet.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	walletRoute := ws.Group("/wallet")
	if handler != nil {
		walletRoute.Use(api.WithElrondProxyFacade(handler))
	}
	wallet.Routes(walletRoute)
	return ws
}

func loadResponse(rsp io.Reader, destination interface{}) {
	jsonParser := json.NewDecoder(rsp)
	err := jsonParser.Decode(destination)
	if err != nil {
		logError(err)
	}
}

func logError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type publicKeyResp struct {
	PublicKey string `json:"publicKey"`
}

func TestPublicKey_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("POST", "/wallet/publickey", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestPublicKey_InvalidRequestShouldError(t *testing.T) {
	t.Parallel()

	handler := &mock.Facade{}
	ws := startNodeServer(handler)

	req, _ := http.NewRequest("POST", "/wallet/publickey", bytes.NewBuffer([]byte("!!invalid!!")))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestPublicKey_FacadeFailsShouldErr(t *testing.T) {
	t.Parallel()

	handler := &mock.Facade{
		PublicKeyFromPrivateKeyCalled: func(privateKeyHex string) (string, error) {
			return "", errors.New("error in facade")
		},
	}
	ws := startNodeServer(handler)

	req, _ := http.NewRequest("POST", "/wallet/publickey", bytes.NewBuffer([]byte(`{"privateKey":"privkey"}`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	var statusRsp publicKeyResp
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestPublicKey_ShouldWork(t *testing.T) {
	t.Parallel()

	handler := &mock.Facade{
		PublicKeyFromPrivateKeyCalled: func(privateKeyHex string) (string, error) {
			return "pubKey", nil
		},
	}
	ws := startNodeServer(handler)

	req, _ := http.NewRequest("POST", "/wallet/publickey", bytes.NewBuffer([]byte(`{"privateKey":"privkey"}`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	var statusRsp publicKeyResp
	loadResponse(resp.Body, &statusRsp)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "pubKey", statusRsp.PublicKey)
}

func TestSendAndSignTx_FailsWithWrongFacadeTypeConversion(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("POST", "/wallet/send", strings.NewReader("invalid"))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestSendAndSignTx_InvalidRequestShouldError(t *testing.T) {
	t.Parallel()

	handler := &mock.Facade{}
	ws := startNodeServer(handler)

	req, _ := http.NewRequest("POST", "/wallet/send", bytes.NewBuffer([]byte("!!invalid!!")))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestSendAndSignTx_FacadeFailsShouldErr(t *testing.T) {
	t.Parallel()

	handler := &mock.Facade{
		SignAndSendTransactionCalled: func(tx *data.Transaction, sk []byte) (string, error) {
			return "", errors.New("error")
		},
	}
	ws := startNodeServer(handler)

	req, _ := http.NewRequest("POST", "/wallet/send", bytes.NewBuffer([]byte(`{"sender":"sndr"}`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

type txSendResp struct {
	TxHash string `json:"txHash"`
}

func TestSendAndSignTx_ShouldWork(t *testing.T) {
	t.Parallel()

	expectedTxHash := "txHash"
	handler := &mock.Facade{
		SignAndSendTransactionCalled: func(tx *data.Transaction, sk []byte) (string, error) {
			return expectedTxHash, nil
		},
	}
	ws := startNodeServer(handler)

	req, _ := http.NewRequest("POST", "/wallet/send", bytes.NewBuffer([]byte(`{"sender":"sndr"}`)))
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	var respTx txSendResp
	loadResponse(resp.Body, &respTx)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, expectedTxHash, respTx.TxHash)
}
