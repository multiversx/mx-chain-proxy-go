package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetConfigMetricsCalled        func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsCalled       func(shardID uint32) (*data.GenericAPIResponse, error)
	GetLatestBlockNonceCalled     func() (uint64, error)
	GetEconomicsDataMetricsCalled func() (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsCalled       func(tokenType string) (*data.GenericAPIResponse, error)
	GetEnableEpochsMetricsCalled  func() (*data.GenericAPIResponse, error)
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

// GetAllIssuedESDTs -
func (nsps *NodeStatusProcessorStub) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	return nsps.GetAllIssuedESDTsCalled(tokenType)
}

// GetEnableEpochsMetrics -
func (nsps *NodeStatusProcessorStub) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	return nsps.GetEnableEpochsMetricsCalled()
}
