package process

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

var testPubKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "erd")
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

	value, _, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
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
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
			return nil, errExpected
		},
	}, testPubKeyConverter)

	value, _, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
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
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) (int, error) {
			return http.StatusNotFound, errExpected
		},
	}, testPubKeyConverter)

	value, _, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.True(t, errors.Is(err, ErrSendingRequest))
}

func TestSCQueryProcessor_ExecuteQuery(t *testing.T) {
	t.Parallel()

	providedBlockInfo := data.BlockInfo{
		Nonce:    123,
		Hash:     "block hash",
		RootHash: "block rootHash",
	}
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "adress1", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) (int, error) {
			response.(*data.ResponseVmValue).Data.Data = &vm.VMOutputApi{
				ReturnData: [][]byte{{42}},
			}
			response.(*data.ResponseVmValue).Data.BlockInfo = providedBlockInfo

			return http.StatusOK, nil
		},
	}, testPubKeyConverter)

	value, blockInfo, err := processor.ExecuteQuery(&data.SCQuery{
		ScAddress: dummyScAddress,
		FuncName:  "function",
		Arguments: [][]byte{[]byte("aa")},
	})

	require.Nil(t, err)
	require.Equal(t, byte(42), value.ReturnData[0][0])
	require.Equal(t, providedBlockInfo, blockInfo)
}

func TestSCQueryProcessor_ExecuteQueryWithCoordinates(t *testing.T) {
	t.Parallel()

	providedNonce := uint64(123)
	providedHash := []byte("block hash")
	providedBlockInfo := data.BlockInfo{
		Nonce:    providedNonce,
		Hash:     string(providedHash),
		RootHash: "block rootHash",
	}
	providedAddr := "address1"
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: providedAddr, ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, dataValue interface{}, response interface{}) (int, error) {
			expectedPath := fmt.Sprintf("%s/vm-values/query?blockHash=%s&blockNonce=%d", providedAddr, hex.EncodeToString(providedHash), providedNonce)
			require.Equal(t, expectedPath, address+path)

			response.(*data.ResponseVmValue).Data.Data = &vm.VMOutputApi{
				ReturnData: [][]byte{{42}},
			}
			response.(*data.ResponseVmValue).Data.BlockInfo = providedBlockInfo

			return http.StatusOK, nil
		},
	}, testPubKeyConverter)

	value, blockInfo, err := processor.ExecuteQuery(&data.SCQuery{
		ScAddress: dummyScAddress,
		FuncName:  "function",
		Arguments: [][]byte{[]byte("aa")},
		BlockNonce: core.OptionalUint64{
			Value:    providedNonce,
			HasValue: true,
		},
		BlockHash: providedHash,
	})

	require.Nil(t, err)
	require.Equal(t, byte(42), value.ReturnData[0][0])
	require.Equal(t, providedBlockInfo, blockInfo)
}

func TestSCQueryProcessor_ExecuteQueryFailsOnRandomErrorShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	processor, _ := NewSCQueryProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
			return []*data.NodeData{
				{Address: "address1", ShardId: 0},
				{Address: "address2", ShardId: 0},
			}, nil
		},
		CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) (int, error) {
			return http.StatusInternalServerError, errExpected
		},
	}, testPubKeyConverter)

	value, _, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
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
		GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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

	value, _, err := processor.ExecuteQuery(&data.SCQuery{ScAddress: dummyScAddress})
	require.Empty(t, value)
	require.Equal(t, errExpected, err)
}
