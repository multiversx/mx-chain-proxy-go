package mock

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetShardStatusCalled func(shardID uint32) (map[string]interface{}, error)
}

// GetShardStatus --
func (nsps *NodeStatusProcessorStub) GetShardStatus(shardID uint32) (map[string]interface{}, error) {
	return nsps.GetShardStatusCalled(shardID)
}
