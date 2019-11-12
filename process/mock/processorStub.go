package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/pkg/errors"
)

var errNotImplemented = errors.New("not implemented")

type ProcessorStub struct {
	ApplyConfigCalled               func(cfg *config.Config) error
	GetObserversCalled              func(shardId uint32) ([]*data.Observer, error)
	ComputeShardIdCalled            func(addressBuff []byte) (uint32, error)
	CallGetRestEndPointCalled       func(address string, path string, value interface{}) error
	CallPostRestEndPointCalled      func(address string, path string, data interface{}, response interface{}) (int, error)
	GetFirstAvailableObserverCalled func() (*data.Observer, error)
	GetAllObserversCalled           func() ([]*data.Observer, error)
}

// ApplyConfig will call the ApplyConfigCalled handler if not nil
func (ps *ProcessorStub) ApplyConfig(cfg *config.Config) error {
	if ps.ApplyConfigCalled != nil {
		return ps.ApplyConfigCalled(cfg)
	}

	return errNotImplemented
}

// GetObservers will call the GetObserversCalled handler if not nil
func (ps *ProcessorStub) GetObservers(shardId uint32) ([]*data.Observer, error) {
	if ps.GetObserversCalled != nil {
		return ps.GetObserversCalled(shardId)
	}

	return nil, errNotImplemented
}

// ComputeShardId will call the ComputeShardIdCalled if not nil
func (ps *ProcessorStub) ComputeShardId(addressBuff []byte) (uint32, error) {
	if ps.ComputeShardIdCalled != nil {
		return ps.ComputeShardIdCalled(addressBuff)
	}

	return 0, errNotImplemented
}

// CallGetRestEndPoint will call the CallGetRestEndPointCalled if not nil
func (ps *ProcessorStub) CallGetRestEndPoint(address string, path string, value interface{}) error {
	if ps.CallGetRestEndPointCalled != nil {
		return ps.CallGetRestEndPointCalled(address, path, value)
	}

	return errNotImplemented
}

// CallPostRestEndPoint will call the CallPostRestEndPoint if not nil
func (ps *ProcessorStub) CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error) {
	if ps.CallPostRestEndPointCalled != nil {
		return ps.CallPostRestEndPointCalled(address, path, data, response)
	}

	return 0, errNotImplemented
}

// GetAllObservers will call the GetAllObservers if not nil
func (ps *ProcessorStub) GetAllObservers() ([]*data.Observer, error) {
	if ps.GetAllObserversCalled != nil {
		return ps.GetAllObserversCalled()
	}

	return nil, errNotImplemented
}
