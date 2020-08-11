package process

import (
	"errors"
	"net/http"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/require"
)

var testPubKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32)
var dummyScAddress = "erd1l453hd0gt5gzdp7czpuall8ggt2dcv5zwmfdf3sd3lguxseux2fsmsgldz"

func TestNewSCQueryProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	processor, err := NewSCQueryProcessor(nil, testPubKeyConverter)
	require.Nil(t, processor)
	require.Equal(t, ErrNilCoreProcessor, err)
}

func TestNewSCQueryProcessor_NilPubConverterShouldErr(t *testing.T) {
	t.Parallel()

	processor, err := NewSCQueryProcessor(&mock.ProcessorStub{}, nil)
	require.Nil(t, processor)
	require.Equal(t, ErrNilPubKeyConverter, err)
}

func TestNewSCQueryProcessor_WithCoreProcessor(t *testing.T) {
	t.Parallel()

	processor, err := NewSCQueryProcessor(&mock.ProcessorStub{}, testPubKeyConverter)
	require.NotNil(t, processor)
	require.Nil(t, err)
}

func TestSCQueryProcessor_ExecuteQueryComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	}, testPubKeyConverter)

	value, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}

func TestSCQueryProcessor_ExecuteQueryGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
			return nil, errExpected
		},
	}, testPubKeyConverter)

	value, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}

func TestSCQueryProcessor_ExecuteQuerySendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) (int, error) {
			return http.StatusNotFound, errExpected
		},
	}, testPubKeyConverter)

	value, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.Equal(t, ErrSendingRequest, err)
}

func TestSCQueryProcessor_ExecuteQuery(t *testing.T) {
	t.Parallel()

	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "adress1", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) (int, error) {
			response.(*data.ResponseVmValue).Data.Data = &vmcommon.VMOutput{
				ReturnData: [][]byte{{42}},
			}

			return http.StatusOK, nil
		},
	}, testPubKeyConverter)

	value, err := processor.ExecuteQuery(&data.SCQuery{
		ScAddress: dummyScAddress,
		FuncName:  "function",
		Arguments: [][]byte{[]byte("aa")},
	})

	require.Nil(t, err)
	require.Equal(t, byte(42), value.ReturnData[0][0])
}

func TestSCQueryProcessor_ExecuteQueryFailsOnRandomErrorShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) (int, error) {
			return http.StatusInternalServerError, errExpected
		},
	}, testPubKeyConverter)

	value, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}

func TestSCQueryProcessor_ExecuteQueryFailsOnBadRequestWithExplicitErrorShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("this error")
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) (int, error) {
			response.(*data.ResponseVmValue).Error = errExpected.Error()
			return http.StatusBadRequest, nil
		},
	}, testPubKeyConverter)

	value, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}
