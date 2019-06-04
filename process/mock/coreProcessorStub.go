package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/pkg/errors"
)

var errNotImplemented = errors.New("not implemented")

type CoreProcessorStub struct {
	ApplyConfigCalled         func(cfg *config.Config) error
	GetObserversCalled        func(shardId uint32) ([]*data.Observer, error)
	ComputeShardIdCalled      func(addressBuff []byte) (uint32, error)
	CallGetRestEndPointCalled func(address string, path string, value interface{}) error
}

func (cps *CoreProcessorStub) ApplyConfig(cfg *config.Config) error {
	if cps.ApplyConfigCalled != nil {
		return cps.ApplyConfigCalled(cfg)
	}

	return errNotImplemented
}

func (cps *CoreProcessorStub) GetObservers(shardId uint32) ([]*data.Observer, error) {
	if cps.GetObserversCalled != nil {
		return cps.GetObserversCalled(shardId)
	}

	return nil, errNotImplemented
}

func (cps *CoreProcessorStub) ComputeShardId(addressBuff []byte) (uint32, error) {
	if cps.ComputeShardIdCalled != nil {
		return cps.ComputeShardIdCalled(addressBuff)
	}

	return 0, errNotImplemented
}

func (cps *CoreProcessorStub) CallGetRestEndPoint(address string, path string, value interface{}) error {
	if cps.CallGetRestEndPointCalled != nil {
		return cps.CallGetRestEndPointCalled(address, path, value)
	}

	return errNotImplemented
}
