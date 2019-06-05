package transaction_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/api"
	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/transaction"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
}

func startNodeServerWrongFacade() *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	ws.Use(func(c *gin.Context) {
		c.Set("elrondProxyFacade", mock.WrongFacade{})
	})
	transactionRoute := ws.Group("/transaction")
	transaction.Routes(transactionRoute)
	return ws
}

func startNodeServer(handler transaction.FacadeHandler) *gin.Engine {
	ws := gin.New()
	ws.Use(cors.Default())
	transactionRoute := ws.Group("/transaction")
	if handler != nil {
		transactionRoute.Use(api.WithElrondProxyFacade(handler))
	}
	transaction.Routes(transactionRoute)
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

func TestSendTransaction_ErrorWithWrongFacade(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("POST", "/transaction/send", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}

func TestSendTransaction_WrongParametersShouldErrorOnValidation(t *testing.T) {
	t.Parallel()

	sender := "sender"
	receiver := "receiver"
	value := "ishouldbeint"
	dataField := "data"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s",`+
			`"receiver":"%s",`+
			`"value":%s,`+
			`"data":"%s"}`, sender, receiver, value, dataField)

	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrValidation.Error())
}

func TestSendTransaction_InvalidHexSignatureShouldError(t *testing.T) {
	t.Parallel()

	sender := "sender"
	receiver := "receiver"
	value := big.NewInt(10)
	dataField := "data"
	signature := "not#only$hex%characters^"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s",`+
			`"receiver":"%s",`+
			`"value":%s,`+
			`"signature":"%s",`+
			`"data":"%s"}`, sender, receiver, value, signature, dataField)

	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrInvalidSignatureHex.Error())
}

func TestSendTransaction_ErrorWhenFacadeSendTransactionError(t *testing.T) {
	t.Parallel()
	sender := "sender"
	receiver := "receiver"
	value := big.NewInt(10)
	dataField := "data"
	signature := "aabbccdd"
	errorString := "send transaction error"

	facade := mock.Facade{
		SendTransactionHandler: func(nonce uint64, sender string, receiver string, value *big.Int,
			code string, signature []byte) error {
			return errors.New(errorString)
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s",`+
			`"receiver":"%s",`+
			`"value":%s,`+
			`"signature":"%s",`+
			`"data":"%s"}`, sender, receiver, value, signature, dataField)

	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, response.Error, errorString)
}

func TestSendTransaction_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	sender := "sender"
	receiver := "receiver"
	value := big.NewInt(10)
	dataField := "data"
	signature := "aabbccdd"

	facade := mock.Facade{
		SendTransactionHandler: func(nonce uint64, sender string, receiver string, value *big.Int,
			code string, signature []byte) error {
			return nil
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(`{
		"nonce": %d,
		"sender": "%s",
		"receiver": "%s",
		"value": %s,
		"signature": "%s",
		"data": "%s"
	}`, nonce, sender, receiver, value, signature, dataField)

	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)
}
