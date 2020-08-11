package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type ObserversProviderStub struct {
	GetNodesByShardIdCalled func(shardId uint32) ([]*data.NodeData, error)
	GetAllNodesCalled       func() ([]*data.NodeData, error)
}

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

func (ops *ObserversProviderStub) IsInterfaceNil() bool {
	return ops == nil
}
