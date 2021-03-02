package groups_test

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/mock"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const transactionsPath = "/transaction"

type txHashResponseData struct {
	Message string `json:"message"`
}

// TxHashResponse structure
type TxHashResponse struct {
	GeneralResponse
	Data txHashResponseData
}

type numOfSentTxsResponseData struct {
	Num uint64 `json:"numOfSentTxs"`
}

// MultiTxsResponse structure
type MultiTxsResponse struct {
	GeneralResponse
	Data numOfSentTxsResponseData `json:"data"`
}

func TestNewTransactionGroup_WrongFacadeShouldErr(t *testing.T) {
	wrongFacade := &mock.WrongFacade{}
	group, err := groups.NewTransactionGroup(wrongFacade)
	require.Nil(t, group)
	require.Equal(t, groups.ErrWrongTypeAssertion, err)
}

func TestSendTransaction_WrongParametersShouldErrorOnValidation(t *testing.T) {
	t.Parallel()

	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "ishouldbeint"
	dataField := "data"

	facade := &mock.Facade{}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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

func TestSendTransaction_ErrorWhenFacadeSendTransactionError(t *testing.T) {
	t.Parallel()
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"
	errorString := "send transaction error"

	facade := &mock.Facade{
		SendTransactionHandler: func(tx *data.Transaction) (int, string, error) {
			return http.StatusInternalServerError, "", errors.New(errorString)
		},
	}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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

func TestSendTransaction_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"
	txHash := "tx hash"

	facade := &mock.Facade{
		SendTransactionHandler: func(tx *data.Transaction) (int, string, error) {
			return 0, txHash, nil
		},
	}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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
	assert.Equal(t, string(data.ReturnCodeSuccess), response.GeneralResponse.Code)
}

func TestSimulateTransaction_WrongParametersShouldErrorOnValidation(t *testing.T) {
	t.Parallel()

	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "ishouldbeint"
	dataField := "data"

	facade := &mock.Facade{}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":%s, "data":"%s"}`,
		sender,
		receiver,
		value,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/simulate", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, response.Error, apiErrors.ErrValidation.Error())
}

func TestSimulateTransaction_ErrorWhenFacadeSimulateTransactionError(t *testing.T) {
	t.Parallel()

	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"
	errorString := "simulate transaction error"

	facade := &mock.Facade{
		SimulateTransactionHandler: func(tx *data.Transaction, _ bool) (*data.GenericAPIResponse, error) {
			return nil, errors.New(errorString)
		},
	}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	jsonStr := fmt.Sprintf(
		`{"sender":"%s", "receiver":"%s", "value":"%s", "signature":"%s",  "data":"%s"}`,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/simulate", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, response.Error, errorString)
}

func TestSimulateTransaction_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := "10"
	dataField := "data"
	signature := "aabbccdd"

	expectedResult := data.GenericAPIResponse{
		Data: data.TransactionSimulationResponseData{
			Result: data.TransactionSimulationResults{FailReason: "reason"},
		},
		Code: data.ReturnCodeSuccess,
	}
	facade := &mock.Facade{
		SimulateTransactionHandler: func(tx *data.Transaction, _ bool) (*data.GenericAPIResponse, error) {
			return &expectedResult, nil
		},
	}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	jsonStr := fmt.Sprintf(
		`{"nonce": %d, "sender": "%s", "receiver": "%s", "value": "%s", "signature": "%s", "data": "%s"	}`,
		nonce,
		sender,
		receiver,
		value,
		signature,
		dataField,
	)
	req, _ := http.NewRequest("POST", "/transaction/simulate", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := data.ResponseTransactionSimulation{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Empty(t, response.Error)
	assert.Equal(t, data.ReturnCodeSuccess, response.Code)
	assert.Equal(t, expectedResult.Data, response.Data)
}

func TestSendMultipleTransactions_WrongParametersShouldErrorOnValidation(t *testing.T) {
	t.Parallel()

	sender := "addr1"
	receiver := "addr2"
	value := "ishouldbeint"
	dataField := "data"

	facade := &mock.Facade{}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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

func TestSendMultipleTransactions_ReturnsSuccessfully(t *testing.T) {
	t.Parallel()

	nonce := uint64(1)
	sender := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	value := big.NewInt(10)
	dataField := "data"
	signature := "aabbccdd"
	txHash := "tx hash"

	facade := &mock.Facade{
		SendTransactionHandler: func(tx *data.Transaction) (int, string, error) {
			return 0, txHash, nil
		},
		SendMultipleTransactionsHandler: func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
			return data.MultipleTransactionsResponseData{
				NumOfTxs:  10,
				TxsHashes: nil,
			}, nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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
	assert.Equal(t, uint64(10), response.Data.Num)
}

func TestSendUserFunds_ErrorWhenFacadeSendUserFundsError(t *testing.T) {
	t.Parallel()

	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"
	errorString := "send user funds error"

	facade := &mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			return errors.New(errorString)
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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

	facade := &mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			return nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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
	facade := &mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			callValue = value
			return nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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
	facade := &mock.Facade{
		SendUserFundsCalled: func(receiver string, value *big.Int) error {
			callValue = value
			return nil
		},
	}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

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

func TestSendUserFunds_FaucetNotEnabled(t *testing.T) {
	t.Parallel()

	receiver := "05702a5fd947a9ddb861ce7ffebfea86c2ca8906df3065ae295f283477ae4e43"

	facade := &mock.Facade{
		IsFaucetEnabledHandler: func() bool {
			return false
		},
	}
	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	value := "100000000000000"
	jsonStr := fmt.Sprintf(
		`{"receiver":"%s", "value": %s}`, receiver, value)

	req, _ := http.NewRequest("POST", "/transaction/send-user-funds", bytes.NewBuffer([]byte(jsonStr)))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := GeneralResponse{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, apiErrors.ErrFaucetNotEnabled.Error(), response.Error)
}
