package process_test

import (
	"encoding/hex"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(nil, &mock.PubKeyConverterMock{})

	assert.Nil(t, tp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewTransactionProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, nil)

	assert.Nil(t, tp)
	assert.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewTransactionProcessor_OkValuesShouldWork(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})

	assert.NotNil(t, tp)
	assert.Nil(t, err)
}

//------- SendTransaction

func TestTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	rc, txHash, err := tp.SendTransaction(&data.ApiTransaction{
		Sender: "invalid hex number",
	})

	assert.Empty(t, txHash)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
	assert.Equal(t, http.StatusBadRequest, rc)
}

func TestTransactionProcessor_SendTransactionComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	rc, txHash, err := tp.SendTransaction(&data.ApiTransaction{})

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
	assert.Equal(t, http.StatusInternalServerError, rc)
}

func TestTransactionProcessor_SendTransactionGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return nil, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	rc, txHash, err := tp.SendTransaction(&data.ApiTransaction{
		Sender: address,
	})

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
	assert.Equal(t, http.StatusInternalServerError, rc)
}

func TestTransactionProcessor_SendTransactionSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return []*data.Observer{
					{Address: "address1", ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) (int, error) {
				return http.StatusInternalServerError, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	rc, txHash, err := tp.SendTransaction(&data.ApiTransaction{
		Sender: address,
	})

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
	assert.Equal(t, http.StatusInternalServerError, rc)
}

func TestTransactionProcessor_SendTransactionSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	txHash := "DEADBEEF01234567890"
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return []*data.Observer{
					{Address: addressFail, ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				txResponse := response.(*data.ResponseTransaction)
				txResponse.TxHash = txHash
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	rc, resultedTxHash, err := tp.SendTransaction(&data.ApiTransaction{
		Sender: address,
	})

	assert.Equal(t, resultedTxHash, txHash)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rc)
}

////------- SendMultipleTransactions

func TestTransactionProcessor_SendMultipleTransactionsShouldWork(t *testing.T) {
	t.Parallel()

	var txsToSend []*data.ApiTransaction
	txsToSend = append(txsToSend, &data.ApiTransaction{Receiver: "aaaaaa", Sender: hex.EncodeToString([]byte("cccccc"))})
	txsToSend = append(txsToSend, &data.ApiTransaction{Receiver: "bbbbbb", Sender: hex.EncodeToString([]byte("dddddd"))})

	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return []*data.Observer{
					{Address: "observer1", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				receivedTxs, ok := value.([]*data.Transaction)
				assert.True(t, ok)
				resp := response.(*data.ResponseMultiTransactions)
				resp.NumOfTxs = uint64(len(receivedTxs))
				response = resp
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	numOfSentTxs, err := tp.SendMultipleTransactions(txsToSend)
	assert.Equal(t, uint64(len(txsToSend)), numOfSentTxs)
	assert.Nil(t, err)
}

func TestTransactionProcessor_SendMultipleTransactionsShouldWorkAndSendTxsByShard(t *testing.T) {
	t.Parallel()

	var txsToSend []*data.ApiTransaction
	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))
	txsToSend = append(txsToSend, &data.ApiTransaction{Receiver: "aaaaaa", Sender: sndrShard0})
	txsToSend = append(txsToSend, &data.ApiTransaction{Receiver: "aaaaaa", Sender: sndrShard0})
	txsToSend = append(txsToSend, &data.ApiTransaction{Receiver: "aaaaaa", Sender: sndrShard1})
	txsToSend = append(txsToSend, &data.ApiTransaction{Receiver: "aaaaaa", Sender: sndrShard1})
	numOfTimesPostEndpointWasCalled := uint32(0)

	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				sndrHex := hex.EncodeToString(addressBuff)
				if sndrHex == sndrShard0 {
					return uint32(0), nil
				}
				if sndrHex == sndrShard1 {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return []*data.Observer{
					{Address: "observer1", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				atomic.AddUint32(&numOfTimesPostEndpointWasCalled, 1)
				resp := response.(*data.ResponseMultiTransactions)
				resp.NumOfTxs = uint64(2)
				response = resp
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	numOfSentTxs, err := tp.SendMultipleTransactions(txsToSend)
	assert.Equal(t, uint64(len(txsToSend)), numOfSentTxs)
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), atomic.LoadUint32(&numOfTimesPostEndpointWasCalled))
}

func TestParseTxStatusResponses(t *testing.T) {
	t.Parallel()

	responses1 := map[uint32][]string{
		0: {"Ok", "Ok", "Ok"},
		1: {process.UnknownStatusTx, process.UnknownStatusTx},
		2: {"Ok"},
	}

	_, err := process.ParseTxStatusResponses(responses1)
	assert.Equal(t, process.ErrCannotGetTransactionStatus, err)

	responses2 := map[uint32][]string{
		0: {process.UnknownStatusTx, process.UnknownStatusTx, process.UnknownStatusTx},
		1: {process.UnknownStatusTx, process.UnknownStatusTx},
		2: {process.UnknownStatusTx},
	}

	status, err := process.ParseTxStatusResponses(responses2)
	assert.NoError(t, err)
	assert.Equal(t, process.UnknownStatusTx, status)

	responses3 := map[uint32][]string{
		0: {"Ok"},
		1: {process.UnknownStatusTx, process.UnknownStatusTx},
		2: {process.UnknownStatusTx},
	}

	status, err = process.ParseTxStatusResponses(responses3)
	assert.NoError(t, err)
	assert.Equal(t, "Ok", status)

	responses4 := map[uint32][]string{
		0: {"Ok", "NotOk"},
		1: {process.UnknownStatusTx, process.UnknownStatusTx},
		2: {process.UnknownStatusTx},
	}

	_, err = process.ParseTxStatusResponses(responses4)
	assert.Equal(t, process.ErrCannotGetTransactionStatus, err)
}
