package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetConfigMetricsCalled        func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsCalled       func(shardID uint32) (*data.GenericAPIResponse, error)
	GetLatestBlockNonceCalled     func() (uint64, error)
	GetEconomicsDataMetricsCalled func() (*data.GenericAPIResponse, error)
	GetTotalStakedCalled          func() (*data.GenericAPIResponse, error)
}

// GetNetworkConfigMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	return nsps.GetConfigMetricsCalled()
}

// GetNetworkStatusMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	return nsps.GetNetworkMetricsCalled(shardID)
}

// GetEconomicsDataMetrics --
func (nsps *NodeStatusProcessorStub) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	return nsps.GetEconomicsDataMetricsCalled()
}

// GetLatestBlockNonce -
func (nsps *NodeStatusProcessorStub) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	return nsps.GetLatestBlockNonceCalled()
}

func (nsps *NodeStatusProcessorStub) GetTotalStaked() (*data.GenericAPIResponse, error) {
	return nsps.GetTotalStakedCalled()
}
