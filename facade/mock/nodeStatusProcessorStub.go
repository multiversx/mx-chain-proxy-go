package mock

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetEpochMetricsCalled   func(shardID uint32) (map[string]interface{}, error)
	GetConfigMetricsCalled  func() (map[string]interface{}, error)
	GetNetworkMetricsCalled func(shardID uint32) (map[string]interface{}, error)
}

// GetNetworkConfigMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkConfigMetrics() (map[string]interface{}, error) {
	return nsps.GetConfigMetricsCalled()
}

// GetNetworkStatusMetrics --
func (nsps *NodeStatusProcessorStub) GetNetworkStatusMetrics(shardID uint32) (map[string]interface{}, error) {
	return nsps.GetNetworkMetricsCalled(shardID)
}

// GetEpochMetrics --
func (nsps *NodeStatusProcessorStub) GetEpochMetrics(shardID uint32) (map[string]interface{}, error) {
	return nsps.GetEpochMetricsCalled(shardID)
}
