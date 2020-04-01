package mock

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetShardStatusCalled  func(shardID uint32) (map[string]interface{}, error)
	GetEpochMetricsCalled func(shardID uint32) (map[string]interface{}, error)
}

// GetEpochMetrics --
func (nsps *NodeStatusProcessorStub) GetEpochMetrics(shardID uint32) (map[string]interface{}, error) {
	return nsps.GetEpochMetricsCalled(shardID)
}

// GetShardStatus --
func (nsps *NodeStatusProcessorStub) GetShardStatus(shardID uint32) (map[string]interface{}, error) {
	return nsps.GetShardStatusCalled(shardID)
}
