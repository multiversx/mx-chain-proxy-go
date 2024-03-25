package process_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBlockProcessor_NilExternalStorageConnectorShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(nil, &mock.ProcessorStub{})
	require.Nil(t, bp)
	require.Equal(t, process.ErrNilDatabaseConnector, err)
}

func TestNewBlockProcessor_NilProcessorShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, nil)
	require.Nil(t, bp)
	require.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewBlockProcessor_ShouldWork(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, &mock.ProcessorStub{})
	require.NotNil(t, bp)
	require.NoError(t, err)
}

func TestBlockProcessor_GetAtlasBlockByShardIDAndNonce(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, &mock.ProcessorStub{})
	require.NotNil(t, bp)

	res, err := bp.GetAtlasBlockByShardIDAndNonce(0, 1)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestBlockProcessor_GetBlockByHashShouldGetFullHistoryNodes(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByHash(0, "hash", common.BlockQueryOptions{})

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByHashShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByHash(0, "hash", common.BlockQueryOptions{})

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByHashNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", common.BlockQueryOptions{})
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetBlockByHashCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", common.BlockQueryOptions{})
	require.True(t, errors.Is(err, process.ErrSendingRequest))
	require.Nil(t, res)
}

func TestBlockProcessor_GetBlockByHashShouldWork(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.BlockApiResponse)
			valResp.Data.Block = api.Block{Nonce: nonce}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", common.BlockQueryOptions{})
	require.NoError(t, err)
	require.NotNil(t, res)

	block := res.Data.Block
	require.Equal(t, nonce, block.Nonce)
}

func TestBlockProcessor_GetBlockByHashShouldWorkAndIncludeAlsoTxs(t *testing.T) {
	t.Parallel()

	isAddressCorrect := false
	nonce := uint64(37)
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			isAddressCorrect = strings.Contains(path, "withTxs=true")
			valResp := value.(*data.BlockApiResponse)
			valResp.Data = data.BlockApiResponsePayload{Block: api.Block{Nonce: nonce}}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", common.BlockQueryOptions{WithTransactions: true})
	require.NoError(t, err)
	require.NotNil(t, res)

	block := res.Data.Block
	require.Equal(t, nonce, block.Nonce)
	require.True(t, isAddressCorrect)
}

func TestBlockProcessor_GetBlockByNonceShouldGetFullHistoryNodes(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByNonce(0, 0, common.BlockQueryOptions{})

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByNonceShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByNonce(0, 1, common.BlockQueryOptions{})

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByNonceNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 1, common.BlockQueryOptions{})
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetBlockByNonceCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 0, common.BlockQueryOptions{})
	require.True(t, errors.Is(err, process.ErrSendingRequest))
	require.Nil(t, res)
}

func TestBlockProcessor_GetBlockByNonceShouldWork(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.BlockApiResponse)
			valResp.Data = data.BlockApiResponsePayload{Block: api.Block{Nonce: nonce}}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, nonce, common.BlockQueryOptions{})
	require.NoError(t, err)
	require.NotNil(t, res)

	block := res.Data.Block
	require.Equal(t, nonce, block.Nonce)
}

func TestBlockProcessor_GetBlockByNonceShouldWorkAndIncludeAlsoTxs(t *testing.T) {
	t.Parallel()

	isAddressCorrect := false
	nonce := uint64(37)
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			isAddressCorrect = strings.Contains(path, "withTxs=true")
			valResp := value.(*data.BlockApiResponse)
			valResp.Data = data.BlockApiResponsePayload{Block: api.Block{Nonce: nonce}}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 3, common.BlockQueryOptions{WithTransactions: true})
	require.NoError(t, err)
	require.NotNil(t, res)

	block := res.Data.Block
	require.Equal(t, nonce, block.Nonce)
	require.True(t, isAddressCorrect)
}

func TestBlockProcessor_GetHyperBlock(t *testing.T) {
	t.Parallel()

	numGetBlockCalled := 0
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: fmt.Sprintf("observer-%d", shardId)}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			numGetBlockCalled++

			response := value.(*data.BlockApiResponse)
			response.Data = data.BlockApiResponsePayload{Block: api.Block{Nonce: 42}}

			if strings.Contains(address, "4294967295") {
				response.Data.Block.Hash = "abcd"
				response.Data.Block.NotarizedBlocks = []*api.NotarizedBlock{
					{Shard: 0, Nonce: 39, Hash: "zero"},
					{Shard: 1, Nonce: 40, Hash: "one"},
					{Shard: 2, Nonce: 41, Hash: "two"},
				}
			}

			return 200, nil
		},
	}

	processor, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.Nil(t, err)
	require.NotNil(t, processor)

	numGetBlockCalled = 0
	response, err := processor.GetHyperBlockByHash("abcd", common.HyperblockQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 4, numGetBlockCalled, "get block should be called for metablock and for all notarized shard blocks")
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))
	require.Equal(t, "abcd", response.Data.Hyperblock.Hash)

	numGetBlockCalled = 0
	response, err = processor.GetHyperBlockByNonce(42, common.HyperblockQueryOptions{})
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 4, numGetBlockCalled, "get block should be called for metablock and for all notarized shard blocks")
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))
	require.Equal(t, "abcd", response.Data.Hyperblock.Hash)
}

// GetInternalBlockByNonce

func TestBlockProcessor_GetInternalBlockByNonceInvalidOutputFormat_ShouldFail(t *testing.T) {
	t.Parallel()

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	blk, err := bp.GetInternalBlockByNonce(0, 0, 2)
	require.Nil(t, blk)
	assert.Equal(t, process.ErrInvalidOutputFormat, err)
}

func TestBlockProcessor_GetInternalBlockByNonceShouldGetFullHistoryNodes(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalBlockByNonce(0, 0, common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalBlockByNonceShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalBlockByNonce(0, 1, common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalBlockByNonceNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalBlockByNonce(0, 1, common.Internal)
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetInternalBlockByNonceCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalBlockByNonce(0, 0, common.Internal)
	require.True(t, errors.Is(err, process.ErrSendingRequest))
	require.Nil(t, res)
}

func TestBlockProcessor_GetInternalBlockByNonceShouldWork(t *testing.T) {
	t.Parallel()

	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be sent",
	}

	nonce := uint64(37)
	expectedData := data.InternalBlockApiResponsePayload{Block: ts}
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.InternalBlockApiResponse)
			valResp.Data = expectedData
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalBlockByNonce(0, nonce, common.Internal)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)

	res, err = bp.GetInternalBlockByNonce(core.MetachainShardId, nonce, common.Proto)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)
}

// GetInternalBlockByHash

func TestBlockProcessor_GetInternalBlockByHashInvalidOutputFormat_ShouldFail(t *testing.T) {
	t.Parallel()

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	blk, err := bp.GetInternalBlockByHash(0, "aaaa", 2)
	require.Nil(t, blk)
	assert.Equal(t, process.ErrInvalidOutputFormat, err)
}

func TestBlockProcessor_GetInternalBlockByHashShouldGetFullHistoryNodes(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalBlockByHash(0, "aaaa", common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalBlockByHashShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalBlockByHash(0, "aaaa", common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalBlockByHashNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalBlockByHash(0, "aaaa", common.Internal)
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetInternalBlockByHashCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalBlockByHash(0, "aaaa", common.Internal)
	require.True(t, errors.Is(err, process.ErrSendingRequest))
	require.Nil(t, res)
}

func TestBlockProcessor_GetInternalBlockByHashShouldWork(t *testing.T) {
	t.Parallel()

	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be sent",
	}

	expectedData := data.InternalBlockApiResponsePayload{Block: ts}
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.InternalBlockApiResponse)
			valResp.Data = expectedData
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalBlockByHash(0, "aaaa", common.Internal)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)

	res, err = bp.GetInternalBlockByHash(core.MetachainShardId, "aaaa", common.Proto)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)
}

// GetInternalMiniBlockByHash

func TestBlockProcessor_GetInternalMiniBlockByHashInvalidOutputFormat_ShouldFail(t *testing.T) {
	t.Parallel()

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	blk, err := bp.GetInternalMiniBlockByHash(0, "aaaa", 1, 2)
	require.Nil(t, blk)
	assert.Equal(t, process.ErrInvalidOutputFormat, err)
}

func TestBlockProcessor_GetInternalMiniBlockByHashShouldGetFullHistoryNodes(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalMiniBlockByHash(0, "aaaa", 1, common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalMiniBlockByHashShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalMiniBlockByHash(0, "aaaa", 1, common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalMiniBlockByHashNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalMiniBlockByHash(0, "aaaa", 1, common.Internal)
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetInternalMiniBlockByHashCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalMiniBlockByHash(0, "aaaa", 1, common.Internal)
	require.True(t, errors.Is(err, process.ErrSendingRequest))
	require.Nil(t, res)
}

func TestBlockProcessor_GetInternalMiniBlockByHashShouldWork(t *testing.T) {
	t.Parallel()

	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be sent",
	}

	expectedData := data.InternalMiniBlockApiResponsePayload{MiniBlock: ts}
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.InternalMiniBlockApiResponse)
			valResp.Data = expectedData
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalMiniBlockByHash(0, "aaaa", 1, common.Internal)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)

	res, err = bp.GetInternalMiniBlockByHash(core.MetachainShardId, "aaaa", 1, common.Proto)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)
}

// GetInternalStartOfEpochMetaBlock

func TestBlockProcessor_GetInternalStartOfEpochMetaBlockInvalidOutputFormat_ShouldFail(t *testing.T) {
	t.Parallel()

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	blk, err := bp.GetInternalStartOfEpochMetaBlock(0, 2)
	require.Nil(t, blk)
	assert.Equal(t, process.ErrInvalidOutputFormat, err)
}

func TestBlockProcessor_GetInternalStartOfEpochMetaBlockShouldGetFullHistoryNodes(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalStartOfEpochMetaBlock(0, common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalStartOfEpochMetaBlockShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	_, _ = bp.GetInternalStartOfEpochMetaBlock(0, common.Internal)

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetInternalStartOfEpochMetaBlockNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalStartOfEpochMetaBlock(0, common.Internal)
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetInternalStartOfEpochMetaBlockCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			assert.Equal(t, shardId, core.MetachainShardId)
			return nil, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalStartOfEpochMetaBlock(0, common.Internal)
	require.True(t, errors.Is(err, process.ErrSendingRequest))
	require.Nil(t, res)
}

func TestBlockProcessor_GetInternalStartOfEpochMetaBlockShouldWork(t *testing.T) {
	t.Parallel()

	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be sent",
	}

	expectedData := data.InternalBlockApiResponsePayload{Block: ts}
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.InternalBlockApiResponse)
			valResp.Data = expectedData
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalStartOfEpochMetaBlock(1, common.Internal)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)

	res, err = bp.GetInternalStartOfEpochMetaBlock(1, common.Proto)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)
}

func TestBlockProcessor_GetAlteredAccountsByNonce(t *testing.T) {
	t.Parallel()

	requestedShardID := uint32(1)
	alteredAcc := &alteredAccount.AlteredAccount{Address: "erd1q"}

	t.Run("could not get observers, should return error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("local error")
		proc := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return nil, expectedErr
			},
		}

		bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
		res, err := bp.GetAlteredAccountsByNonce(requestedShardID, 4, common.GetAlteredAccountsForBlockOptions{})
		require.Equal(t, expectedErr, err)
		require.Nil(t, res)
	})

	t.Run("could not get response from any observer, should return error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("local error")
		callGetEndpointCt := 0
		node1 := &data.NodeData{ShardId: requestedShardID, Address: "addr1"}
		node2 := &data.NodeData{ShardId: requestedShardID, Address: "addr2"}

		proc := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, requestedShardID, shardId)
				return []*data.NodeData{node1, node2}, nil
			},

			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				callGetEndpointCt++
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.True(t, address == node1.Address || address == node2.Address)
				require.Equal(t, "/block/altered-accounts/by-nonce/4", path)
				return 0, expectedErr
			},
		}

		bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
		res, err := bp.GetAlteredAccountsByNonce(requestedShardID, 4, common.GetAlteredAccountsForBlockOptions{})
		require.Equal(t, 2, callGetEndpointCt)
		require.True(t, errors.Is(err, process.ErrSendingRequest))
		require.Nil(t, res)
	})

	t.Run("should work getting data from first observer", func(t *testing.T) {
		t.Parallel()

		proc := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, requestedShardID, shardId)
				return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
			},

			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.Equal(t, "addr", address)
				require.Equal(t, "/block/altered-accounts/by-nonce/4", path)

				ret := value.(*data.AlteredAccountsApiResponse)
				ret.Error = ""
				ret.Code = "success"
				ret.Data.Accounts = []*alteredAccount.AlteredAccount{alteredAcc}
				return 0, nil
			},
		}

		bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
		res, err := bp.GetAlteredAccountsByNonce(requestedShardID, 4, common.GetAlteredAccountsForBlockOptions{})
		require.Nil(t, err)
		require.Equal(t, &data.AlteredAccountsApiResponse{
			Data: data.AlteredAccountsPayload{
				Accounts: []*alteredAccount.AlteredAccount{alteredAcc},
			},
			Error: "",
			Code:  "success",
		}, res)
	})
}

func TestBlockProcessor_GetAlteredAccountsByHash(t *testing.T) {
	t.Parallel()

	requestedShardID := uint32(1)
	alteredAcc := &alteredAccount.AlteredAccount{Address: "erd1q"}

	t.Run("could not get observers, should return error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("local error")
		proc := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return nil, expectedErr
			},
		}

		bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
		res, err := bp.GetAlteredAccountsByHash(requestedShardID, "hash", common.GetAlteredAccountsForBlockOptions{})
		require.Equal(t, expectedErr, err)
		require.Nil(t, res)
	})

	t.Run("could not get response from any observer, should return error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("local error")
		callGetEndpointCt := 0
		node1 := &data.NodeData{ShardId: requestedShardID, Address: "addr1"}
		node2 := &data.NodeData{ShardId: requestedShardID, Address: "addr2"}

		proc := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, requestedShardID, shardId)
				return []*data.NodeData{node1, node2}, nil
			},

			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				callGetEndpointCt++
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.True(t, address == node1.Address || address == node2.Address)
				require.Equal(t, "/block/altered-accounts/by-hash/hash", path)
				return 0, expectedErr
			},
		}

		bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
		res, err := bp.GetAlteredAccountsByHash(requestedShardID, "hash", common.GetAlteredAccountsForBlockOptions{})
		require.Equal(t, 2, callGetEndpointCt)
		require.True(t, errors.Is(err, process.ErrSendingRequest))
		require.Nil(t, res)
	})

	t.Run("should work getting data from first observer", func(t *testing.T) {
		t.Parallel()

		proc := &mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, requestedShardID, shardId)
				return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
			},

			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.Equal(t, "addr", address)
				require.Equal(t, "/block/altered-accounts/by-hash/hash", path)

				ret := value.(*data.AlteredAccountsApiResponse)
				ret.Error = ""
				ret.Code = "success"
				ret.Data.Accounts = []*alteredAccount.AlteredAccount{alteredAcc}
				return 0, nil
			},
		}

		bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
		res, err := bp.GetAlteredAccountsByHash(requestedShardID, "hash", common.GetAlteredAccountsForBlockOptions{})
		require.Nil(t, err)
		require.Equal(t, &data.AlteredAccountsApiResponse{
			Data: data.AlteredAccountsPayload{
				Accounts: []*alteredAccount.AlteredAccount{alteredAcc},
			},
			Error: "",
			Code:  "success",
		}, res)
	})
}

func TestBlockProcessor_GetHyperBlockByNonceWithAlteredAccounts(t *testing.T) {
	t.Parallel()

	observerAddr := "observerAddress"
	alteredAcc1 := &alteredAccount.AlteredAccount{Address: "erd1q"}
	alteredAcc2 := &alteredAccount.AlteredAccount{Address: "erd1w"}

	callGetEndpointCt := 0
	getObserversCt := 0
	proc := &mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			switch getObserversCt {
			case 0:
				require.Equal(t, core.MetachainShardId, shardId)
			case 1, 2:
				require.Equal(t, uint32(1), shardId)
			case 3, 4:
				require.Equal(t, uint32(2), shardId)
			}

			getObserversCt++
			return []*data.NodeData{{ShardId: shardId, Address: observerAddr}}, nil
		},

		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			require.Equal(t, observerAddr, address)

			switch callGetEndpointCt {
			case 0:
				require.Equal(t, &data.BlockApiResponse{}, value)
				require.Equal(t, "/block/by-nonce/4?withTxs=true", path)

				ret := value.(*data.BlockApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Block = api.Block{
					StateRootHash: "stateRootHash",
					NotarizedBlocks: []*api.NotarizedBlock{
						{
							Shard: 1,
							Hash:  "hash1",
						},
						{
							Shard: 2,
							Hash:  "hash2",
						},
					},
				}
			case 1:
				require.Equal(t, &data.BlockApiResponse{}, value)
				require.Equal(t, "/block/by-hash/hash1?withTxs=true", path)

				ret := value.(*data.BlockApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Block = api.Block{Hash: "hash1", Shard: 1}
			case 2:
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.Equal(t, "/block/altered-accounts/by-hash/hash1", path)

				ret := value.(*data.AlteredAccountsApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Accounts = []*alteredAccount.AlteredAccount{alteredAcc1}
			case 3:
				require.Equal(t, &data.BlockApiResponse{}, value)
				require.Equal(t, "/block/by-hash/hash2?withTxs=true", path)

				ret := value.(*data.BlockApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Block = api.Block{Hash: "hash2", Shard: 2}
			case 4:
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.Equal(t, "/block/altered-accounts/by-hash/hash2", path)

				ret := value.(*data.AlteredAccountsApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Accounts = []*alteredAccount.AlteredAccount{alteredAcc2}
			}

			callGetEndpointCt++
			return 0, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)

	res, err := bp.GetHyperBlockByNonce(4, common.HyperblockQueryOptions{WithAlteredAccounts: true})
	require.Nil(t, err)

	expectedHyperBlock := api.Hyperblock{
		StateRootHash: "stateRootHash",
		ShardBlocks: []*api.NotarizedBlock{
			{
				Shard:           1,
				Hash:            "hash1",
				AlteredAccounts: []*alteredAccount.AlteredAccount{alteredAcc1},
				MiniBlockHashes: make([]string, 0),
			},
			{
				Shard:           2,
				Hash:            "hash2",
				AlteredAccounts: []*alteredAccount.AlteredAccount{alteredAcc2},
				MiniBlockHashes: make([]string, 0),
			},
		},
		Transactions: make([]*transaction.ApiTransactionResult, 0),
	}

	require.Equal(t, &data.HyperblockApiResponse{
		Code: data.ReturnCodeSuccess,
		Data: data.HyperblockApiResponsePayload{
			Hyperblock: expectedHyperBlock,
		},
	}, res)
	require.NotNil(t, res)
	require.Equal(t, 5, callGetEndpointCt)
	require.Equal(t, 5, getObserversCt)
}

func TestBlockProcessor_GetHyperBlockByHashWithAlteredAccounts(t *testing.T) {
	t.Parallel()

	observerAddr := "observerAddress"
	alteredAcc1 := &alteredAccount.AlteredAccount{Address: "erd1q"}
	alteredAcc2 := &alteredAccount.AlteredAccount{Address: "erd1w"}

	callGetEndpointCt := 0
	getObserversCt := 0
	proc := &mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			switch getObserversCt {
			case 0:
				require.Equal(t, core.MetachainShardId, shardId)
			case 1, 2:
				require.Equal(t, uint32(1), shardId)
			case 3, 4:
				require.Equal(t, uint32(2), shardId)
			}

			getObserversCt++
			return []*data.NodeData{{ShardId: shardId, Address: observerAddr}}, nil
		},

		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			require.Equal(t, observerAddr, address)

			switch callGetEndpointCt {
			case 0:
				require.Equal(t, &data.BlockApiResponse{}, value)
				require.Equal(t, "/block/by-hash/abcdef?withTxs=true", path)

				ret := value.(*data.BlockApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Block = api.Block{
					StateRootHash: "stateRootHash",
					NotarizedBlocks: []*api.NotarizedBlock{
						{
							Shard: 1,
							Hash:  "hash1",
						},
						{
							Shard: 2,
							Hash:  "hash2",
						},
					},
				}
			case 1:
				require.Equal(t, &data.BlockApiResponse{}, value)
				require.Equal(t, "/block/by-hash/hash1?withTxs=true", path)

				ret := value.(*data.BlockApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Block = api.Block{Hash: "hash1", Shard: 1}
			case 2:
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.Equal(t, "/block/altered-accounts/by-hash/hash1", path)

				ret := value.(*data.AlteredAccountsApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Accounts = []*alteredAccount.AlteredAccount{alteredAcc1}
			case 3:
				require.Equal(t, &data.BlockApiResponse{}, value)
				require.Equal(t, "/block/by-hash/hash2?withTxs=true", path)

				ret := value.(*data.BlockApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Block = api.Block{Hash: "hash2", Shard: 2}
			case 4:
				require.Equal(t, &data.AlteredAccountsApiResponse{}, value)
				require.Equal(t, "/block/altered-accounts/by-hash/hash2", path)

				ret := value.(*data.AlteredAccountsApiResponse)
				ret.Code = data.ReturnCodeSuccess
				ret.Data.Accounts = []*alteredAccount.AlteredAccount{alteredAcc2}
			}

			callGetEndpointCt++
			return 0, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)

	res, err := bp.GetHyperBlockByHash("abcdef", common.HyperblockQueryOptions{WithAlteredAccounts: true})
	require.Nil(t, err)

	expectedHyperBlock := api.Hyperblock{
		StateRootHash: "stateRootHash",
		ShardBlocks: []*api.NotarizedBlock{
			{
				Shard:           1,
				Hash:            "hash1",
				AlteredAccounts: []*alteredAccount.AlteredAccount{alteredAcc1},
				MiniBlockHashes: make([]string, 0),
			},
			{
				Shard:           2,
				Hash:            "hash2",
				AlteredAccounts: []*alteredAccount.AlteredAccount{alteredAcc2},
				MiniBlockHashes: make([]string, 0),
			},
		},
		Transactions: make([]*transaction.ApiTransactionResult, 0),
	}

	require.Equal(t, &data.HyperblockApiResponse{
		Code: data.ReturnCodeSuccess,
		Data: data.HyperblockApiResponsePayload{
			Hyperblock: expectedHyperBlock,
		},
	}, res)
	require.NotNil(t, res)
	require.Equal(t, 5, callGetEndpointCt)
	require.Equal(t, 5, getObserversCt)
}

func TestBlockProcessor_GetInternalStartOfEpochValidatorsInfo(t *testing.T) {
	t.Parallel()

	ts := &testStruct{
		Nonce: 10,
		Name:  "a test struct to be sent",
	}

	expectedData := data.InternalStartOfEpochValidators{
		ValidatorsInfo: ts,
	}
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.ValidatorsInfoApiResponse)
			valResp.Data = expectedData
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetInternalStartOfEpochValidatorsInfo(1)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedData, res.Data)
}
