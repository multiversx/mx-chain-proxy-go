package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ObserversProviderStub -
type ObserversProviderStub struct {
	GetNodesByShardIdCalled func(shardId uint32) ([]*data.NodeData, error)
	GetAllNodesCalled       func() ([]*data.NodeData, error)
	ReloadNodesCalled       func(nodesType data.NodeType) data.NodesReloadResponse
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
