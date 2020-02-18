package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type ObserversProviderStub struct {
	GetObserversByShardIdCalled func(shardId uint32) ([]*data.Observer, error)
	GetAllObserversCalled       func() []*data.Observer
}

func (ops *ObserversProviderStub) GetObserversByShardId(shardId uint32) ([]*data.Observer, error) {
	if ops.GetObserversByShardIdCalled != nil {
		return ops.GetObserversByShardIdCalled(shardId)
	}

	return []*data.Observer{
		{
			Address: "address",
			ShardId: 0,
		},
	}, nil
}

func (ops *ObserversProviderStub) GetAllObservers() []*data.Observer {
	if ops.GetAllObserversCalled != nil {
		return ops.GetAllObserversCalled()
	}

	return []*data.Observer{
		{
			Address: "address",
			ShardId: 0,
		},
	}
}

func (ops *ObserversProviderStub) IsInterfaceNil() bool {
	return ops == nil
}
