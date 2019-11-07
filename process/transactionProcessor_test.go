package process_test

import (
	"encoding/hex"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(nil)

	assert.Nil(t, tp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewTransactionProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{})

	assert.NotNil(t, tp)
	assert.Nil(t, err)
}

//------- SendTransaction

func TestTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{})
	txHash, err := tp.SendTransaction(&data.Transaction{
		Sender: "invalid hex number",
	})

	assert.Empty(t, txHash)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestTransactionProcessor_SendTransactionComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	},
	)
	txHash, err := tp.SendTransaction(&data.Transaction{})

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
}

func TestTransactionProcessor_SendTransactionGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	},
	)
	address := "DEADBEEF"
	txHash, err := tp.SendTransaction(&data.Transaction{
		Sender: address,
	})

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
}

func TestTransactionProcessor_SendTransactionSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errExpected
		},
	},
	)
	address := "DEADBEEF"
	txHash, err := tp.SendTransaction(&data.Transaction{
		Sender: address,
	})

	assert.Empty(t, txHash)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestTransactionProcessor_SendTransactionSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	txHash := "DEADBEEF01234567890"
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: addressFail, ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) error {
			txResponse := response.(*data.ResponseTransaction)
			txResponse.TxHash = txHash
			return nil
		},
	},
	)
	address := "DEADBEEF"
	resultedTxHash, err := tp.SendTransaction(&data.Transaction{
		Sender: address,
	})

	assert.Equal(t, resultedTxHash, txHash)
	assert.Nil(t, err)
}

////------- SendMultipleTransactions

func TestTransactionProcessor_SendMultipleTransactionsShouldWork(t *testing.T) {
	t.Parallel()

	var txsToSend []*data.Transaction
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "rcvr1", Sender: hex.EncodeToString([]byte("sndr1"))})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "rcvr2", Sender: hex.EncodeToString([]byte("sndr2"))})

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
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) error {
				receivedTxs, ok := value.([]*data.Transaction)
				assert.True(t, ok)
				assert.Equal(t, txsToSend, receivedTxs)
				resp := response.(*data.ResponseMultiTransactions)
				resp.NumOfTxs = uint64(len(receivedTxs))
				response = resp
				return nil
			},
		},
	)

	numOfSentTxs, err := tp.SendMultipleTransactions(txsToSend)
	assert.Equal(t, uint64(len(txsToSend)), numOfSentTxs)
	assert.Nil(t, err)
}

func TestTransactionProcessor_SendMultipleTransactionsShouldWorkAndSendTxsByShard(t *testing.T) {
	t.Parallel()

	var txsToSend []*data.Transaction
	sndrShard0 := hex.EncodeToString([]byte("sender shard 0"))
	sndrShard1 := hex.EncodeToString([]byte("sender shard 1"))
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "rcvr1", Sender: sndrShard0})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "rcvr2", Sender: sndrShard0})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "rcvr3", Sender: sndrShard1})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "rcvr4", Sender: sndrShard1})
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
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) error {
				atomic.AddUint32(&numOfTimesPostEndpointWasCalled, 1)
				resp := response.(*data.ResponseMultiTransactions)
				resp.NumOfTxs = uint64(2)
				response = resp
				return nil
			},
		},
	)

	numOfSentTxs, err := tp.SendMultipleTransactions(txsToSend)
	assert.Equal(t, uint64(len(txsToSend)), numOfSentTxs)
	assert.Nil(t, err)
	assert.Equal(t, uint32(2), atomic.LoadUint32(&numOfTimesPostEndpointWasCalled))
}
