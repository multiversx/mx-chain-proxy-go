package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetConfigMetricsCalled    func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsCalled   func(shardID uint32) (*data.GenericAPIResponse, error)
	GetLatestBlockNonceCalled func() (uint64, error)
}

// GetNetworkConfigMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	return nsps.GetConfigMetricsCalled()
}

// GetNetworkStatusMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	return nsps.GetNetworkMetricsCalled(shardID)
}

// GetLatestBlockNonce -
func (nsps *NodeStatusProcessorStub) GetLatestBlockNonce() (uint64, error) {
	return nsps.GetLatestBlockNonceCalled()
}
