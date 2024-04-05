package mock

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/config"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer"
	"github.com/pkg/errors"
)

var errNotImplemented = errors.New("not implemented")

type ProcessorStub struct {
	ApplyConfigCalled                    func(cfg *config.Config) error
	GetObserversCalled                   func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllObserversCalled                func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetObserversOnePerShardCalled        func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetFullHistoryNodesOnePerShardCalled func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetFullHistoryNodesCalled            func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetAllFullHistoryNodesCalled         func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error)
	GetShardIDsCalled                    func() []uint32
	ComputeShardIdCalled                 func(addressBuff []byte) (uint32, error)
	CallGetRestEndPointCalled            func(address string, path string, value interface{}) (int, error)
	CallPostRestEndPointCalled           func(address string, path string, data interface{}, response interface{}) (int, error)
	GetShardCoordinatorCalled            func() common.Coordinator
	GetPubKeyConverterCalled             func() core.PubkeyConverter
	GetObserverProviderCalled            func() observer.NodesProviderHandler
	GetFullHistoryNodesProviderCalled    func() observer.NodesProviderHandler
}

// GetShardCoordinator -
func (ps *ProcessorStub) GetShardCoordinator() common.Coordinator {
	if ps.GetShardCoordinatorCalled != nil {
		return ps.GetShardCoordinatorCalled()
	}

	return &ShardCoordinatorMock{}
}

// GetPubKeyConverter -
func (ps *ProcessorStub) GetPubKeyConverter() core.PubkeyConverter {
	if ps.GetPubKeyConverterCalled != nil {
		return ps.GetPubKeyConverterCalled()
	}

	return &PubKeyConverterMock{}
}

// GetObserverProvider -
func (ps *ProcessorStub) GetObserverProvider() observer.NodesProviderHandler {
	if ps.GetObserverProviderCalled != nil {
		return ps.GetObserverProviderCalled()
	}

	return &ObserversProviderStub{}
}

// GetFullHistoryNodesProvider -
func (ps *ProcessorStub) GetFullHistoryNodesProvider() observer.NodesProviderHandler {
	if ps.GetFullHistoryNodesProviderCalled != nil {
		return ps.GetFullHistoryNodesProviderCalled()
	}

	return &ObserversProviderStub{}
}

// ApplyConfig will call the ApplyConfigCalled handler if not nil
func (ps *ProcessorStub) ApplyConfig(cfg *config.Config) error {
	if ps.ApplyConfigCalled != nil {
		return ps.ApplyConfigCalled(cfg)
	}

	return errNotImplemented
}

// GetObservers will call the GetObserversCalled handler if not nil
func (ps *ProcessorStub) GetObservers(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	if ps.GetObserversCalled != nil {
		return ps.GetObserversCalled(shardID, dataAvailability)
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
func (ps *ProcessorStub) CallGetRestEndPoint(address string, path string, value interface{}) (int, error) {
	if ps.CallGetRestEndPointCalled != nil {
		return ps.CallGetRestEndPointCalled(address, path, value)
	}

	return 0, errNotImplemented
}

// CallPostRestEndPoint will call the CallPostRestEndPoint if not nil
func (ps *ProcessorStub) CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error) {
	if ps.CallPostRestEndPointCalled != nil {
		return ps.CallPostRestEndPointCalled(address, path, data, response)
	}

	return 0, errNotImplemented
}

// GetShardIDs will call the GetShardIDsCalled if not nil
func (ps *ProcessorStub) GetShardIDs() []uint32 {
	if ps.GetShardIDsCalled != nil {
		return ps.GetShardIDsCalled()
	}

	return nil
}

// GetAllObservers will call the GetAllNodesCalled if not nil
func (ps *ProcessorStub) GetAllObservers(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	if ps.GetAllObserversCalled != nil {
		return ps.GetAllObserversCalled(dataAvailability)
	}

	return nil, nil
}

// GetObserversOnePerShard will call the GetObserversOnePerShardCalled if not nil
func (ps *ProcessorStub) GetObserversOnePerShard(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	if ps.GetObserversOnePerShardCalled != nil {
		return ps.GetObserversOnePerShardCalled(dataAvailability)
	}

	return nil, nil
}

// GetFullHistoryNodesOnePerShard will call the GetFullHistoryNodesOnePerShardCalled if not nil
func (ps *ProcessorStub) GetFullHistoryNodesOnePerShard(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	if ps.GetFullHistoryNodesOnePerShardCalled != nil {
		return ps.GetFullHistoryNodesOnePerShardCalled(dataAvailability)
	}

	return nil, nil
}

// GetFullHistoryNodes will call the GetFullHistoryNodes handler if not nil
func (ps *ProcessorStub) GetFullHistoryNodes(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	if ps.GetFullHistoryNodesCalled != nil {
		return ps.GetFullHistoryNodesCalled(shardID, dataAvailability)
	}

	return nil, errNotImplemented
}

// GetAllFullHistoryNodes will call the GetAllFullHistoryNodes handler if not nil
func (ps *ProcessorStub) GetAllFullHistoryNodes(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
	if ps.GetAllFullHistoryNodesCalled != nil {
		return ps.GetAllFullHistoryNodesCalled(dataAvailability)
	}

	return nil, errNotImplemented
}

// IsInterfaceNil -
func (ps *ProcessorStub) IsInterfaceNil() bool {
	return ps == nil
}
