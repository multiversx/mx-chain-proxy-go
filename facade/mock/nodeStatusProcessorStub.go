package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetConfigMetricsCalled        func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsCalled       func(shardID uint32) (*data.GenericAPIResponse, error)
	GetLatestBlockNonceCalled     func() (uint64, error)
	GetEconomicsDataMetricsCalled func() (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsCalled       func() (*data.GenericAPIResponse, error)
	GetDirectStakedInfoCalled     func() (*data.GenericAPIResponse, error)
	GetDelegatedInfoCalled        func() (*data.GenericAPIResponse, error)
	CreateSnapshotCalled          func(timestamp string) (*data.GenericAPIResponse, error)
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
func (nsps *NodeStatusProcessorStub) GetAllIssuedESDTs() (*data.GenericAPIResponse, error) {
	return nsps.GetAllIssuedESDTsCalled()
}

// GetDirectStakedInfo -
func (nsps *NodeStatusProcessorStub) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	return nsps.GetDirectStakedInfoCalled()
}

// GetDelegatedInfo-
func (nsps *NodeStatusProcessorStub) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	return nsps.GetDelegatedInfoCalled()
}

// GetDelegatedInfo-
func (nsps *NodeStatusProcessorStub) CreateSnapshot(timestamp string) (*data.GenericAPIResponse, error) {
	return nsps.CreateSnapshotCalled(timestamp)
}
