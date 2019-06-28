package process_test

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewGetValuesProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	vvp, err := process.NewVmValuesProcessor(nil)

	assert.Nil(t, vvp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewGetValuesProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	vvp, err := process.NewVmValuesProcessor(&mock.ProcessorStub{})

	assert.NotNil(t, vvp)
	assert.Nil(t, err)
}

//------- GetValues

func TestGetValuesProcessor_GetDataValueInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	vvp, _ := process.NewVmValuesProcessor(&mock.ProcessorStub{})
	funcName := "func"
	value, err := vvp.GetVmValue("invalid hex number", funcName)

	assert.Empty(t, value)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestGetValuesProcessor_GetDataValueComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	vvp, _ := process.NewVmValuesProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	})
	address := "DEADBEEF"
	funcName := "func"
	value, err := vvp.GetVmValue(address, funcName)

	assert.Empty(t, value)
	assert.Equal(t, errExpected, err)
}

func TestGetValuesProcessor_GetDataValueGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	vvp, _ := process.NewVmValuesProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	})
	address := "DEADBEEF"
	funcName := "func"
	value, err := vvp.GetVmValue(address, funcName)

	assert.Empty(t, value)
	assert.Equal(t, errExpected, err)
}

func TestGetValuesProcessor_GetDataValueSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	vvp, _ := process.NewVmValuesProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) error {
			return errExpected
		},
	})
	address := "DEADBEEF"
	funcName := "func"
	value, err := vvp.GetVmValue(address, funcName)

	assert.Empty(t, value)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestGetValuesProcessor_GetDataValueReceivedNonHexOneObserverShouldErr(t *testing.T) {
	t.Parallel()

	vvp, _ := process.NewVmValuesProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) error {
			response.(*data.ResponseVmValue).HexData = "not a hex data"

			return nil
		},
	})
	address := "DEADBEEF"
	funcName := "func"
	value, err := vvp.GetVmValue(address, funcName)

	assert.Empty(t, value)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestGetValuesProcessor_GetDataValueShouldWork(t *testing.T) {
	t.Parallel()

	expectedValueHex := "DEADBEEFDEADBEEFDEADBEEF"
	expectedValue, _ := hex.DecodeString(expectedValueHex)
	vvp, _ := process.NewVmValuesProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) error {
			response.(*data.ResponseVmValue).HexData = expectedValueHex

			return nil
		},
	})
	address := "DEADBEEF"
	funcName := "func"
	value, err := vvp.GetVmValue(address, funcName)

	assert.Nil(t, err)
	assert.Equal(t, expectedValue, value)
}
