package process_test

import (
	"errors"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(nil, &mock.KeygenStub{}, &mock.SignerStub{})

	assert.Nil(t, tp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewTransactionProcessor_NilKeygenShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, nil, &mock.SignerStub{})

	assert.Nil(t, tp)
	assert.Equal(t, process.ErrNilKeyGen, err)
}

func TestNewTransactionProcessor_NilSingleSignerShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.KeygenStub{}, nil)

	assert.Nil(t, tp)
	assert.Equal(t, process.ErrNilSingleSigner, err)
}

func TestNewTransactionProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.KeygenStub{}, &mock.SignerStub{})

	assert.NotNil(t, tp)
	assert.Nil(t, err)
}

//------- SendTransaction

func TestTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.KeygenStub{}, &mock.SignerStub{})
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
		&mock.KeygenStub{},
		&mock.SignerStub{},
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
		&mock.KeygenStub{},
		&mock.SignerStub{},
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
		&mock.KeygenStub{},
		&mock.SignerStub{},
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
		&mock.KeygenStub{},
		&mock.SignerStub{},
	)
	address := "DEADBEEF"
	resultedTxHash, err := tp.SendTransaction(&data.Transaction{
		Sender: address,
	})

	assert.Equal(t, resultedTxHash, txHash)
	assert.Nil(t, err)
}

//------- SendUserFunds

func TestTransactionProcessor_SendUserFundsInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.KeygenStub{}, &mock.SignerStub{})
	err := tp.SendUserFunds("invalid hex number", big.NewInt(10))

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestTransactionProcessor_SendUserFundsGetObserversFailsShouldErr(t *testing.T) {
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
		&mock.KeygenStub{},
		&mock.SignerStub{},
	)
	address := "DEADBEEF"
	err := tp.SendUserFunds(address, big.NewInt(10))

	assert.Equal(t, errExpected, err)
}

func TestTransactionProcessor_SendUserFundsComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	},
		&mock.KeygenStub{},
		&mock.SignerStub{},
	)
	address := "DEADBEEF"
	err := tp.SendUserFunds(address, big.NewInt(10))

	assert.Equal(t, errExpected, err)
}

func TestTransactionProcessor_SendUserFundsSendingFailsOnAllObserversShouldErr(t *testing.T) {
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
	},
		&mock.KeygenStub{},
		&mock.SignerStub{},
	)
	address := "DEADBEEF"
	err := tp.SendUserFunds(address, big.NewInt(10))

	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestTransactionProcessor_SendUserFundsSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
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
	},
		&mock.KeygenStub{},
		&mock.SignerStub{},
	)
	address := "DEADBEEF"
	err := tp.SendUserFunds(address, big.NewInt(10))

	assert.Nil(t, err)
}

//------- SignAndSendTransaction

func TestTransactionProcessor_SignAndSendTransactionInvalidPrivKeyShouldErr(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("error")

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{},
		&mock.KeygenStub{
			PrivateKeyFromByteArrayCalled: func(b []byte) (key crypto.PrivateKey, e error) {
				return nil, expectedErr
			},
		},
		&mock.SignerStub{},
	)

	_, err := tp.SignAndSendTransaction(&data.Transaction{}, []byte("sk"))
	assert.Equal(t, expectedErr, err)
}

func TestTransactionProcessor_SignAndSendTransaction(t *testing.T) {
	t.Parallel()

	signWasCalled := false
	callEndpointWasCalled := false

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) ([]*data.Observer, error) {
			return []*data.Observer{
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) error {
			callEndpointWasCalled = true
			return nil
		},
	},
		&mock.KeygenStub{
			PrivateKeyFromByteArrayCalled: func(b []byte) (crypto.PrivateKey, error) {
				return nil, nil
			},
		},
		&mock.SignerStub{
			SignCalled: func(private crypto.PrivateKey, msg []byte) ([]byte, error) {
				signWasCalled = true
				return nil, nil
			},
		},
	)

	resp, err := tp.SignAndSendTransaction(&data.Transaction{}, []byte("sk"))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.True(t, signWasCalled)
	assert.True(t, callEndpointWasCalled)
}
