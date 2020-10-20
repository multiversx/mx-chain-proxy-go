package process_test

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func getLatestFullySynchronizedHyperblockNonceMock() (uint64, error) {
	return math.MaxUint64, nil
}

func TestNewBlockProcessor_NilExternalStorageConnectorShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(nil, &mock.ProcessorStub{}, getLatestFullySynchronizedHyperblockNonceMock)
	require.Nil(t, bp)
	require.Equal(t, process.ErrNilDatabaseConnector, err)
}

func TestNewBlockProcessor_NilProcessorShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, nil, getLatestFullySynchronizedHyperblockNonceMock)
	require.Nil(t, bp)
	require.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewBlockProcessor_ErrNilFuncToGetLatestFullySynchronizedHyberblockNonce(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, &mock.ProcessorStub{}, nil)
	require.Nil(t, bp)
	require.Equal(t, process.ErrNilGetLatestFullySynchronizedHyperblockNonceFunction, err)
}

func TestNewBlockProcessor_ShouldWork(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, &mock.ProcessorStub{}, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)
	require.NoError(t, err)
}

func TestBlockProcessor_GetAtlasBlockByShardIDAndNonce(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, &mock.ProcessorStub{}, getLatestFullySynchronizedHyperblockNonceMock)
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
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByHash(0, "hash", false)

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByHashShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByHash(0, "hash", false)

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByHashNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", false)
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetBlockByHashCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", false)
	require.Equal(t, process.ErrSendingRequest, err)
	require.Nil(t, res)
}

func TestBlockProcessor_GetBlockByHashShouldWork(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.BlockApiResponse)
			valResp.Data.Block = data.Block{Nonce: nonce}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", false)
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
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			isAddressCorrect = strings.Contains(path, "withTxs=true")
			valResp := value.(*data.BlockApiResponse)
			valResp.Data = data.BlockApiResponsePayload{Block: data.Block{Nonce: nonce}}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", true)
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
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByNonce(0, 0, false)

	require.True(t, getFullHistoryNodesCalled)
	require.False(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByNonceShouldGetObservers(t *testing.T) {
	t.Parallel()

	getFullHistoryNodesCalled := false
	getObserversCalled := false

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, errors.New("local err")
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	_, _ = bp.GetBlockByNonce(0, 1, false)

	require.True(t, getFullHistoryNodesCalled)
	require.True(t, getObserversCalled)
}

func TestBlockProcessor_GetBlockByNonceNoFullNodesOrObserversShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("local err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return nil, localErr
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return nil, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 1, false)
	require.Nil(t, res)
	require.Equal(t, localErr, err)
}

func TestBlockProcessor_GetBlockByNonceCallGetFailsShouldErr(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			return 500, localErr
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 0, false)
	require.Equal(t, process.ErrSendingRequest, err)
	require.Nil(t, res)
}

func TestBlockProcessor_GetBlockByNonceShouldWork(t *testing.T) {
	t.Parallel()

	nonce := uint64(37)
	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valResp := value.(*data.BlockApiResponse)
			valResp.Data = data.BlockApiResponsePayload{Block: data.Block{Nonce: nonce}}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, nonce, false)
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
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: "addr"}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			isAddressCorrect = strings.Contains(path, "withTxs=true")
			valResp := value.(*data.BlockApiResponse)
			valResp.Data = data.BlockApiResponsePayload{Block: data.Block{Nonce: nonce}}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 3, true)
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
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: fmt.Sprintf("http://observer-%d", shardId)}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			numGetBlockCalled++

			response := value.(*data.BlockApiResponse)
			response.Data = data.BlockApiResponsePayload{Block: data.Block{Nonce: 42}}

			if strings.Contains(address, "4294967295") {
				response.Data.Block.Hash = "abcd"
				response.Data.Block.NotarizedBlocks = []*data.NotarizedBlock{
					{Shard: 0, Nonce: 39, Hash: "zero"},
					{Shard: 1, Nonce: 40, Hash: "one"},
					{Shard: 2, Nonce: 41, Hash: "two"},
				}
			}

			return 200, nil
		},
	}

	processor, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonceMock)
	require.Nil(t, err)
	require.NotNil(t, processor)

	numGetBlockCalled = 0
	response, err := processor.GetHyperBlockByHash("abcd")
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 4, numGetBlockCalled, "get block should be called for metablock and for all notarized shard blocks")
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))
	require.Equal(t, "abcd", response.Data.Hyperblock.Hash)

	numGetBlockCalled = 0
	response, err = processor.GetHyperBlockByNonce(42)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 4, numGetBlockCalled, "get block should be called for metablock and for all notarized shard blocks")
	require.Equal(t, 42, int(response.Data.Hyperblock.Nonce))
	require.Equal(t, "abcd", response.Data.Hyperblock.Hash)
}

func TestBlockProcessor_GetHyperBlockByHash_NonceIsTooHigh(t *testing.T) {
	t.Parallel()

	proc := &mock.ProcessorStub{
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			return []*data.NodeData{{ShardId: shardId, Address: fmt.Sprintf("http://observer-%d", shardId)}}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			response := value.(*data.BlockApiResponse)
			response.Data = data.BlockApiResponsePayload{Block: data.Block{Nonce: 42}}

			return 200, nil
		},
	}

	getLatestFullySynchronizedHyperblockNonce := func() (uint64, error) {
		return 1, nil
	}
	processor, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonce)
	require.Nil(t, err)
	require.NotNil(t, processor)

	response, err := processor.GetHyperBlockByHash("abcd")
	require.Nil(t, response)
	require.True(t, errors.Is(err, process.ErrCannotGetHyperblock))
}

func TestBlockProcessor_GetHyperBlockByNonce_NonceIsTooHigh(t *testing.T) {
	t.Parallel()

	proc := &mock.ProcessorStub{}

	getLatestFullySynchronizedHyperblockNonce := func() (uint64, error) {
		return 1, nil
	}
	processor, err := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc, getLatestFullySynchronizedHyperblockNonce)
	require.Nil(t, err)
	require.NotNil(t, processor)

	response, err := processor.GetHyperBlockByNonce(10)
	require.Nil(t, response)
	require.True(t, errors.Is(err, process.ErrCannotGetHyperblock))
}
