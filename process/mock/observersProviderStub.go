package mock

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// ObserversProviderStub -
type ObserversProviderStub struct {
	GetNodesByShardIdCalled           func(shardId uint32) ([]*data.NodeData, error)
	GetAllNodesCalled                 func() ([]*data.NodeData, error)
	ReloadNodesCalled                 func(nodesType data.NodeType) data.NodesReloadResponse
	UpdateNodesBasedOnSyncStateCalled func(nodesWithSyncStatus []*data.NodeData)
	GetAllNodesWithSyncStateCalled    func() []*data.NodeData
}

// GetNodesByShardId -
func (ops *ObserversProviderStub) GetNodesByShardId(shardId uint32) ([]*data.NodeData, error) {
	if ops.GetNodesByShardIdCalled != nil {
		return ops.GetNodesByShardIdCalled(shardId)
	}

	return []*data.NodeData{
		{
			Address: "address",
			ShardId: 0,
		},
	}, nil
}

// GetAllNodes -
func (ops *ObserversProviderStub) GetAllNodes() ([]*data.NodeData, error) {
	if ops.GetAllNodesCalled != nil {
		return ops.GetAllNodesCalled()
	}

	return []*data.NodeData{
		{
			Address: "address",
			ShardId: 0,
		},
	}, nil
}

// RemoveOutOfSyncNodesIfNeeded -
func (ops *ObserversProviderStub) UpdateNodesBasedOnSyncState(nodesWithSyncStatus []*data.NodeData) {
	if ops.UpdateNodesBasedOnSyncStateCalled != nil {
		ops.UpdateNodesBasedOnSyncStateCalled(nodesWithSyncStatus)
	}
}

// GetAllNodesWithSyncState -
func (ops *ObserversProviderStub) GetAllNodesWithSyncState() []*data.NodeData {
	if ops.GetAllNodesWithSyncStateCalled != nil {
		return ops.GetAllNodesWithSyncStateCalled()
	}

	return make([]*data.NodeData, 0)
}

// ReloadNodes -
func (ops *ObserversProviderStub) ReloadNodes(nodesType data.NodeType) data.NodesReloadResponse {
	if ops.ReloadNodesCalled != nil {
		return ops.ReloadNodesCalled(nodesType)
	}

	return data.NodesReloadResponse{}
}

// IsInterfaceNil -
func (ops *ObserversProviderStub) IsInterfaceNil() bool {
	return ops == nil
}
