package process

import (
	"errors"
	"net/http"
	"testing"

	coreProcess "github.com/ElrondNetwork/elrond-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

func TestNewSCQueryServiceProxy_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	proxy, err := NewSCQueryServiceProxy(nil)
	require.Nil(t, proxy)
	require.Equal(t, ErrNilCoreProcessor, err)
}

func TestNewSCQueryServiceProxy_WithCoreProcessor(t *testing.T) {
	t.Parallel()

	proxy, err := NewSCQueryServiceProxy(&mock.ProcessorStub{})
	require.NotNil(t, proxy)
	require.Nil(t, err)
}

func TestSCQueryServiceProxy_ExecuteQuery_ComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	proxy, _ := NewSCQueryServiceProxy(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	})

	value, err := proxy.ExecuteQuery(&coreProcess.SCQuery{})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}

func TestSCQueryServiceProxy_ExecuteQuery_GetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	proxy, _ := NewSCQueryServiceProxy(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	})

	value, err := proxy.ExecuteQuery(&coreProcess.SCQuery{})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}

func TestSCQueryServiceProxy_ExecuteQuery_SendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	proxy, _ := NewSCQueryServiceProxy(&mock.ProcessorStub{
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
			return http.StatusNotFound, errExpected
		},
	})

	value, err := proxy.ExecuteQuery(&coreProcess.SCQuery{})
	require.Empty(t, value)
	require.Equal(t, ErrSendingRequest, err)
}

func TestSCQueryServiceProxy_ExecuteQuery(t *testing.T) {
	t.Parallel()

	proxy, _ := NewSCQueryServiceProxy(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) (int, error) {
			response.(*data.ResponseVmValue).Data = &vmcommon.VMOutput{
				ReturnData: [][]byte{[]byte{42}},
			}

			return http.StatusOK, nil
		},
	})

	value, err := proxy.ExecuteQuery(&coreProcess.SCQuery{
		ScAddress: []byte("address"),
		FuncName:  "function",
		Arguments: [][]byte{[]byte("aa")},
	})

	require.Nil(t, err)
	require.Equal(t, byte(42), value.ReturnData[0][0])
}

func TestSCQueryServiceProxy_ExecuteQuery_FailsOnRandomError(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	proxy, _ := NewSCQueryServiceProxy(&mock.ProcessorStub{
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
	})

	value, err := proxy.ExecuteQuery(&coreProcess.SCQuery{})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}

func TestSCQueryServiceProxy_ExecuteQuery_FailsOnBadRequestWithExplicitError(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("this error")
	proxy, _ := NewSCQueryServiceProxy(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) (int, error) {
			response.(*data.ResponseVmValue).Error = errExpected.Error()
			return http.StatusBadRequest, nil
		},
	})

	value, err := proxy.ExecuteQuery(&coreProcess.SCQuery{})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}
