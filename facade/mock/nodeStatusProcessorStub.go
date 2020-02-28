package mock

// NodeStatusProcessorStub --
type NodeStatusProcessorStub struct {
	GetNodeStatusDataCalled func(shardId string) (map[string]interface{}, error)
}

// GetNodeStatusData --
func (nsps *NodeStatusProcessorStub) GetNodeStatusData(shardId string) (map[string]interface{}, error) {
	return nsps.GetNodeStatusDataCalled(shardId)
}
