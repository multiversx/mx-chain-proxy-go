package mock

import "github.com/multiversx/mx-chain-proxy-go/data"

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetConfigMetricsCalled                          func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsCalled                         func(shardID uint32) (*data.GenericAPIResponse, error)
	GetLatestFullySynchronizedHyperblockNonceCalled func() (uint64, error)
	GetEconomicsDataMetricsCalled                   func() (*data.GenericAPIResponse, error)
	GetAllIssuedESDTsCalled                         func(tokenType string) (*data.GenericAPIResponse, error)
	GetDirectStakedInfoCalled                       func() (*data.GenericAPIResponse, error)
	GetDelegatedInfoCalled                          func() (*data.GenericAPIResponse, error)
	GetEnableEpochsMetricsCalled                    func() (*data.GenericAPIResponse, error)
	GetRatingsConfigCalled                          func() (*data.GenericAPIResponse, error)
	GetGenesisNodesPubKeysCalled                    func() (*data.GenericAPIResponse, error)
	GetGasConfigsCalled                             func() (*data.GenericAPIResponse, error)
	GetTriesStatisticsCalled                        func(shardID uint32) (*data.TrieStatisticsAPIResponse, error)
	GetEpochStartDataCalled                         func(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error)
}

// GetNetworkConfigMetrics --
func (stub *NodeStatusProcessorStub) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if stub.GetConfigMetricsCalled != nil {
		return stub.GetConfigMetricsCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetNetworkStatusMetrics --
func (stub *NodeStatusProcessorStub) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	if stub.GetNetworkMetricsCalled != nil {
		return stub.GetNetworkMetricsCalled(shardID)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetEconomicsDataMetrics --
func (stub *NodeStatusProcessorStub) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	if stub.GetEconomicsDataMetricsCalled != nil {
		return stub.GetEconomicsDataMetricsCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetLatestFullySynchronizedHyperblockNonce -
func (stub *NodeStatusProcessorStub) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	if stub.GetLatestFullySynchronizedHyperblockNonceCalled != nil {
		return stub.GetLatestFullySynchronizedHyperblockNonceCalled()
	}

	return 0, nil
}

// GetAllIssuedESDTs -
func (stub *NodeStatusProcessorStub) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	if stub.GetAllIssuedESDTsCalled != nil {
		return stub.GetAllIssuedESDTsCalled(tokenType)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetDirectStakedInfo -
func (stub *NodeStatusProcessorStub) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	if stub.GetDirectStakedInfoCalled != nil {
		return stub.GetDirectStakedInfoCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetDelegatedInfo -
func (stub *NodeStatusProcessorStub) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	if stub.GetDelegatedInfoCalled != nil {
		return stub.GetDelegatedInfoCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetEnableEpochsMetrics -
func (stub *NodeStatusProcessorStub) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	if stub.GetEnableEpochsMetricsCalled != nil {
		return stub.GetEnableEpochsMetricsCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetRatingsConfig -
func (stub *NodeStatusProcessorStub) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	if stub.GetRatingsConfigCalled != nil {
		return stub.GetRatingsConfigCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetGenesisNodesPubKeys -
func (stub *NodeStatusProcessorStub) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	if stub.GetGenesisNodesPubKeysCalled != nil {
		return stub.GetGenesisNodesPubKeysCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetGasConfigs -
func (stub *NodeStatusProcessorStub) GetGasConfigs() (*data.GenericAPIResponse, error) {
	if stub.GetGasConfigsCalled != nil {
		return stub.GetGasConfigsCalled()
	}

	return &data.GenericAPIResponse{}, nil
}

// GetEpochStartData -
func (stub *NodeStatusProcessorStub) GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error) {
	if stub.GetEpochStartDataCalled != nil {
		return stub.GetEpochStartDataCalled(epoch, shardID)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetTriesStatistics -
func (stub *NodeStatusProcessorStub) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	if stub.GetTriesStatisticsCalled != nil {
		return stub.GetTriesStatisticsCalled(shardID)
	}
	return &data.TrieStatisticsAPIResponse{}, nil
}
