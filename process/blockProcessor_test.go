package process_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go/api/block"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
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
		GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getFullHistoryNodesCalled = true
			return nil, nil
		},
		GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
			getObserversCalled = true
			return nil, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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
			valResp := value.(*data.GenericAPIResponse)
			valResp.Data = block.APIBlock{Nonce: nonce}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", false)
	require.NoError(t, err)
	require.NotNil(t, res)

	blck, ok := res.Data.(block.APIBlock)
	require.True(t, ok)
	require.Equal(t, nonce, blck.Nonce)
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
			valResp := value.(*data.GenericAPIResponse)
			valResp.Data = block.APIBlock{Nonce: nonce}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByHash(0, "hash", true)
	require.NoError(t, err)
	require.NotNil(t, res)

	blck, ok := res.Data.(block.APIBlock)
	require.True(t, ok)
	require.Equal(t, nonce, blck.Nonce)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
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
			valResp := value.(*data.GenericAPIResponse)
			valResp.Data = block.APIBlock{Nonce: nonce}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, nonce, false)
	require.NoError(t, err)
	require.NotNil(t, res)

	blck, ok := res.Data.(block.APIBlock)
	require.True(t, ok)
	require.Equal(t, nonce, blck.Nonce)
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
			valResp := value.(*data.GenericAPIResponse)
			valResp.Data = block.APIBlock{Nonce: nonce}
			return 200, nil
		},
	}

	bp, _ := process.NewBlockProcessor(&mock.ExternalStorageConnectorStub{}, proc)
	require.NotNil(t, bp)

	res, err := bp.GetBlockByNonce(0, 3, true)
	require.NoError(t, err)
	require.NotNil(t, res)

	blck, ok := res.Data.(block.APIBlock)
	require.True(t, ok)
	require.Equal(t, nonce, blck.Nonce)
	require.True(t, isAddressCorrect)
}
