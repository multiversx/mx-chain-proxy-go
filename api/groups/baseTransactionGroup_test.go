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

type txPool struct {
	TxPool data.TransactionsPool `json:"txPool"`
}

type txPoolResp struct {
	GeneralResponse
	Data txPool
}

type txPoolForSender struct {
	TxPool data.TransactionsPoolForSender `json:"txPool"`
}

type txPoolForSenderResp struct {
	GeneralResponse
	Data txPoolForSender
}

type lastNonceResp struct {
	GeneralResponse
	Data data.TransactionsPoolLastNonceForSender
}

type nonceGaps struct {
	NonceGaps data.TransactionsPoolNonceGaps `json:"nonceGaps"`
}

type nonceGapsResp struct {
	GeneralResponse
	Data nonceGaps
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

func TestGetTransactionsPool_InvalidOptions(t *testing.T) {
	t.Parallel()

	t.Run("", testInvalidParameters("?last-nonce=true,fields=sender", apiErrors.ErrBadUrlParams))
	t.Run("", testInvalidParameters("?last-nonce=true&fields=sender", apiErrors.ErrFetchingLatestNonceCannotIncludeFields))
	t.Run("", testInvalidParameters("?nonce-gaps=true&fields=sender", apiErrors.ErrFetchingNonceGapsCannotIncludeFields))
	t.Run("", testInvalidParameters("?last-nonce=true", apiErrors.ErrEmptySenderToGetLatestNonce))
	t.Run("", testInvalidParameters("?nonce-gaps=true", apiErrors.ErrEmptySenderToGetNonceGaps))
	t.Run("", testInvalidParameters("?fields=123", apiErrors.ErrInvalidFields))
	t.Run("", testInvalidParameters("?fields=_/+", apiErrors.ErrInvalidFields))
}

func testInvalidParameters(path string, expectedErr error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		transactionsGroup, err := groups.NewTransactionGroup(&mock.Facade{})
		require.NoError(t, err)
		ws := startProxyServer(transactionsGroup, transactionsPath)

		req, _ := http.NewRequest("GET", "/transaction/pool"+path, bytes.NewBuffer([]byte("")))

		resp := httptest.NewRecorder()
		ws.ServeHTTP(resp, req)

		response := GeneralResponse{}
		loadResponse(resp.Body, &response)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, response.Error, expectedErr.Error())
	}
}

func TestGetTransactionsPool_ReturnsSuccesfully(t *testing.T) {
	t.Parallel()

	providedTx := data.WrappedTransaction{
		TxFields: map[string]interface{}{
			"sender": "sender",
			"hash":   "hash",
		},
	}
	providedTxPool := &data.TransactionsPool{
		RegularTransactions: []data.WrappedTransaction{providedTx},
	}
	facade := &mock.Facade{
		GetTransactionsPoolHandler: func(fields string) (*data.TransactionsPool, error) {
			return providedTxPool, nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	req, _ := http.NewRequest("GET", "/transaction/pool", bytes.NewBuffer([]byte("")))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := txPoolResp{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Error, "")
	assert.Equal(t, providedTxPool, &response.Data.TxPool)
}

func TestGetTransactionsPoolForShard_ReturnsSuccesfully(t *testing.T) {
	t.Parallel()

	providedTx := data.WrappedTransaction{
		TxFields: map[string]interface{}{
			"sender": "sender",
			"hash":   "hash",
		},
	}
	providedTxPool := &data.TransactionsPool{
		RegularTransactions: []data.WrappedTransaction{providedTx},
	}
	facade := &mock.Facade{
		GetTransactionsPoolForShardHandler: func(shardID uint32, fields string) (*data.TransactionsPool, error) {
			return providedTxPool, nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	req, _ := http.NewRequest("GET", "/transaction/pool?shard-id=0", bytes.NewBuffer([]byte("")))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := txPoolResp{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Error, "")
	assert.Equal(t, providedTxPool, &response.Data.TxPool)
}

func TestGetTransactionsPoolForSender_ReturnsSuccesfully(t *testing.T) {
	t.Parallel()

	providedTx := data.WrappedTransaction{
		TxFields: map[string]interface{}{
			"sender": "sender",
			"hash":   "hash",
		},
	}
	providedTxPool := &data.TransactionsPoolForSender{
		Transactions: []data.WrappedTransaction{providedTx},
	}
	facade := &mock.Facade{
		GetTransactionsPoolForSenderHandler: func(sender, fields string) (*data.TransactionsPoolForSender, error) {
			return providedTxPool, nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	req, _ := http.NewRequest("GET", "/transaction/pool?by-sender=dummy", bytes.NewBuffer([]byte("")))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := txPoolForSenderResp{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Error, "")
	assert.Equal(t, providedTxPool, &response.Data.TxPool)
}

func TestLastPoolNonceForSender_ReturnsSuccesfully(t *testing.T) {
	t.Parallel()

	providedNonce := uint64(33)
	facade := &mock.Facade{
		GetLastPoolNonceForSenderHandler: func(sender string) (uint64, error) {
			return providedNonce, nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	req, _ := http.NewRequest("GET", "/transaction/pool?by-sender=dummy&last-nonce=true", bytes.NewBuffer([]byte("")))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := lastNonceResp{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Error, "")
	assert.Equal(t, providedNonce, response.Data.Nonce)
}

func TestGetTransactionsPoolPoolNonceGapsForSender_ReturnsSuccesfully(t *testing.T) {
	t.Parallel()

	providedGap := data.NonceGap{
		From: 15,
		To:   55,
	}
	providedNonceGaps := &data.TransactionsPoolNonceGaps{
		Gaps: []data.NonceGap{providedGap},
	}
	facade := &mock.Facade{
		GetTransactionsPoolNonceGapsForSenderHandler: func(sender string) (*data.TransactionsPoolNonceGaps, error) {
			return providedNonceGaps, nil
		},
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	require.NoError(t, err)
	ws := startProxyServer(transactionsGroup, transactionsPath)

	req, _ := http.NewRequest("GET", "/transaction/pool?by-sender=dummy&nonce-gaps=true", bytes.NewBuffer([]byte("")))

	resp := httptest.NewRecorder()
	ws.ServeHTTP(resp, req)

	response := nonceGapsResp{}
	loadResponse(resp.Body, &response)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, response.Error, "")
	assert.Equal(t, providedNonceGaps, &response.Data.NonceGaps)
}
