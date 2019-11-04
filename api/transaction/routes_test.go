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
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// General response structure
type GeneralResponse struct {
	Error string `json:"error"`
}

// TxHashResponse structure
type TxHashResponse struct {
	Error  string `json:"error"`
	TxHash string `json:"txHash"`
}

// MultiTxsResponse structure
type MultiTxsResponse struct {
	Error    string `json:"error"`
	NumOfTxs uint64 `json:"numOfSentTxs"`
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

	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "ishouldbeint"
	dataField := "data"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":%s, "data":"%s"}`,
		sender,
		receiver,
		value,
		dataField,
	)
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

	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "not#only$hex%characters^"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s", "data":"%s"}`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
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
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"
	errorString := "send transaction error"

	facade := mock.Facade{
		SendTransactionHandler: func(tx *data.Transaction) (string, error) {
			return "", errors.New(errorString)
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s",  "data":"%s"}`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, response.Error, errorString)
}

func TestSendTransaction_ErrorWhenInvalidSender(t *testing.T) {
	t.Parallel()
	sender := "sender"
	receiver := "receiver"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s", "data":"%s"}`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrInvalidSenderAddress.Error())
}

func TestSendTransaction_ErrorWhenInvalidReceiver(t *testing.T) {
	t.Parallel()
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "receiver"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s", "data":"%s"}`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrInvalidReceiverAddress.Error())
}

func TestSendTransaction_ErrorWhenInvalidSignature(t *testing.T) {
	t.Parallel()
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aftgyi"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s", "data":"%s"}`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrInvalidSignatureHex.Error())
}

func TestSendTransaction_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"
	txHash := "tx hash"

	facade := mock.Facade{
		SendTransactionHandler: func(tx *data.Transaction) (string, error) {
			return txHash, nil
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"nonce": %d, "sender": "%s", "receiver": "%s", "value": "%s", "signature": "%s", "data": "%s"	}`,
		nonce,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := TxHashResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)
	assert.Equal(t, txHash, response.TxHash)
}

func TestSendMultipleTransactions_ErrorWithWrongFacade(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("POST", "/transaction/send-multiple", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}

func TestSendMultipleTransactions_WrongParametersShouldErrorOnValidation(t *testing.T) {
	t.Parallel()

	sender := "addr1"
	receiver := "addr2"
	value := "ishouldbeint"
	dataField := "data"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`[{"sender":"%s", "receiver":"%s", "value":%s, "data":"%s"}]`,
		sender,
		receiver,
		value,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send-multiple", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrValidation.Error())
}

func TestSendMultipleTransactions_ErrorWhenInvalidSignature(t *testing.T) {
	t.Parallel()
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := big.NewInt(10)
	dataField := "data"
	signature := "aftgyi"

	facade := mock.Facade{}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`[{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s", "data":"%s"}]`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send-multiple", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrInvalidSignatureHex.Error())
}

func TestSendMultipleTransactions_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := big.NewInt(10)
	dataField := "data"
	signature := "aabbccdd"
	txHash := "tx hash"

	facade := mock.Facade{
		SendTransactionHandler: func(tx *data.Transaction) (string, error) {
			return txHash, nil
		},
		SendMultipleTransactionsHandler: func(txs []*data.Transaction) (uint64, error) {
			return uint64(10), nil
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`[{"nonce": %d, "sender": "%s", "receiver": "%s", "value": "%s", "signature": "%s", "data": "%s"	}]`,
		nonce,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/send-multiple", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := MultiTxsResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)
	assert.Equal(t, uint64(10), response.NumOfTxs)
}

func TestSendUserFunds_ErrorWithWrongFacade(t *testing.T) {
	t.Parallel()

	ws := startNodeServerWrongFacade()
	req, _ := http.NewRequest("POST", "/transaction/send-user-funds", nil)
	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, resp.Code, http.StatusInternalServerError)
}

func TestSendUserFunds_ErrorWhenFacadeSendUserFundsError(t *testing.T) {
	t.Parallel()
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	errorString := "send user funds error"

	facade := mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			return errors.New(errorString)
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"receiver":"%s"}`, receiver)

	req, _ := http.NewRequest("POST", "/transaction/send-user-funds", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, response.Error, errorString)
}

func TestSendUserFunds_ReturnsSuccesfully(t *testing.T) {
	t.Parallel()
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"

	facade := mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			return nil
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"receiver":"%s"}`, receiver)

	req, _ := http.NewRequest("POST", "/transaction/send-user-funds", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Error, "")
}

func TestSendUserFunds_NilValue(t *testing.T) {
	t.Parallel()
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"

	var callValue *big.Int
	facade := mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			callValue = value
			return nil
		},
	}
	ws := startNodeServer(&facade)

	jsonStr := fmt.Sprintf(
		`{"receiver":"%s"}`, receiver)

	req, _ := http.NewRequest("POST", "/transaction/send-user-funds", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Nil(t, callValue)
}

func TestSendUserFunds_CorrectValue(t *testing.T) {
	t.Parallel()
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"

	var callValue *big.Int
	facade := mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			callValue = value
			return nil
		},
	}
	ws := startNodeServer(&facade)

	expectedValue, _ := big.NewInt(0).SetString("100000000000000", 10)
	jsonStr := fmt.Sprintf(
		`{"receiver":"%s", "value": %d}`, receiver, expectedValue)

	req, _ := http.NewRequest("POST", "/transaction/send-user-funds", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, expectedValue, callValue)
}
