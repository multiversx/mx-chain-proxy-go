package process_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction_NilCoreProcessorShouldErr(t *testing.T) {
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

func TestNewTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{})
	sig := make([]byte, 0)
	txHash, err := tp.SendTransaction(0, "invalid hex number", "FF", big.NewInt(0), "", sig)

	assert.Empty(t, txHash)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestNewTransactionProcessor_SendTransactionComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	})
	address := "DEADBEEF"
	sig := make([]byte, 0)
	txHash, err := tp.SendTransaction(0, address, address, big.NewInt(0), "", sig)

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
}

func TestNewTransactionProcessor_SendTransactionGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	})
	address := "DEADBEEF"
	sig := make([]byte, 0)
	txHash, err := tp.SendTransaction(0, address, address, big.NewInt(0), "", sig)

	assert.Empty(t, txHash)
	assert.Equal(t, errExpected, err)
}

func TestNewTransactionProcessor_SendTransactionSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errExpected
		},
	})
	address := "DEADBEEF"
	sig := make([]byte, 0)
	txHash, err := tp.SendTransaction(0, address, address, big.NewInt(0), "", sig)

	assert.Empty(t, txHash)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestNewTransactionProcessor_SendTransactionSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
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
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) error {
			txResponse := response.(*data.ResponseTransaction)
			txResponse.TxHash = txHash
			return nil
		},
	})
	address := "DEADBEEF"
	sig := make([]byte, 0)
	resultedTxHash, err := tp.SendTransaction(0, address, address, big.NewInt(0), "", sig)

	assert.Equal(t, resultedTxHash, txHash)
	assert.Nil(t, err)
}

//------- SendUserFunds

func TestNewTransactionProcessor_SendUserFundsInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{})
	err := tp.SendUserFunds("invalid hex number")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestNewTransactionProcessor_SendUserFundsGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	})
	address := "DEADBEEF"
	err := tp.SendUserFunds(address)

	assert.Equal(t, errExpected, err)
}

func TestNewTransactionProcessor_SendUserFundsComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	})
	address := "DEADBEEF"
	err := tp.SendUserFunds(address)

	assert.Equal(t, errExpected, err)
}

func TestNewTransactionProcessor_SendUserFundsSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errExpected
		},
	})
	address := "DEADBEEF"
	err := tp.SendUserFunds(address)

	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestNewTransactionProcessor_SendUserFundsSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: addressFail, ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) error {
			return nil
		},
	})
	address := "DEADBEEF"
	err := tp.SendUserFunds(address)

	assert.Nil(t, err)
}