package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetConfigMetricsCalled        func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsCalled       func(shardID uint32) (*data.GenericAPIResponse, error)
	GetLatestBlockNonceCalled     func() (uint64, error)
	GetEconomicsDataMetricsCalled func() (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsCalled       func(tokenType string) (*data.GenericAPIResponse, error)
	GetDirectStakedInfoCalled     func() (*data.GenericAPIResponse, error)
	GetDelegatedInfoCalled        func() (*data.GenericAPIResponse, error)
	GetEnableEpochsMetricsCalled  func() (*data.GenericAPIResponse, error)
	GetRatingsConfigCalled        func() (*data.GenericAPIResponse, error)
	GetGenesisNodesPubKeysCalled  func() (*data.GenericAPIResponse, error)
	GetGasConfigsCalled           func() (*data.GenericAPIResponse, error)
	GetTriesStatisticsCalled      func(shardID uint32) (*data.TrieStatisticsAPIResponse, error)
}

// GetNetworkConfigMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if nsps.GetConfigMetricsCalled != nil {
		return nsps.GetConfigMetricsCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetNetworkStatusMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	if nsps.GetNetworkMetricsCalled != nil {
		return nsps.GetNetworkMetricsCalled(shardID)
	}
	return &data.GenericAPIResponse{}, nil
}

// GetEconomicsDataMetrics --
func (nsps *NodeStatusProcessorStub) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	if nsps.GetEconomicsDataMetricsCalled != nil {
		return nsps.GetEconomicsDataMetricsCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetLatestFullySynchronizedHyperblockNonce -
func (nsps *NodeStatusProcessorStub) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	if nsps.GetLatestBlockNonceCalled != nil {
		return nsps.GetLatestBlockNonceCalled()
	}
	return 0, nil
}

// GetAllIssuedESDTs -
func (nsps *NodeStatusProcessorStub) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	if nsps.GetAllIssuedESDTsCalled != nil {
		return nsps.GetAllIssuedESDTsCalled(tokenType)
	}
	return &data.GenericAPIResponse{}, nil
}

// GetDirectStakedInfo -
func (nsps *NodeStatusProcessorStub) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	if nsps.GetDirectStakedInfoCalled != nil {
		return nsps.GetDirectStakedInfoCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetDelegatedInfo -
func (nsps *NodeStatusProcessorStub) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	if nsps.GetDelegatedInfoCalled != nil {
		return nsps.GetDelegatedInfoCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetEnableEpochsMetrics -
func (nsps *NodeStatusProcessorStub) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	if nsps.GetEnableEpochsMetricsCalled != nil {
		return nsps.GetEnableEpochsMetricsCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetRatingsConfig -
func (nsps *NodeStatusProcessorStub) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	if nsps.GetRatingsConfigCalled != nil {
		return nsps.GetRatingsConfigCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetGenesisNodesPubKeys -
func (nsps *NodeStatusProcessorStub) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	if nsps.GetGenesisNodesPubKeysCalled != nil {
		return nsps.GetGenesisNodesPubKeysCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetGasConfigs -
func (nsps *NodeStatusProcessorStub) GetGasConfigs() (*data.GenericAPIResponse, error) {
	if nsps.GetGasConfigsCalled != nil {
		return nsps.GetGasConfigsCalled()
	}
	return &data.GenericAPIResponse{}, nil
}

// GetTriesStatistics -
func (nsps *NodeStatusProcessorStub) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	if nsps.GetTriesStatisticsCalled != nil {
		return nsps.GetTriesStatisticsCalled(shardID)
	}
	return &data.TrieStatisticsAPIResponse{}, nil
}
