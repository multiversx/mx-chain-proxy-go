package process

// NodeStatusPath represents the path where an observer exposes his nodeStatus
const NodeStatusPath = "/node/status"

type NodeStatusProcessor struct {
	proc Processor
}

// NewNodeStatusProcessor creates a new instance of NodeStatusProcessor
func NewNodeStatusProcessor(processor Processor) (*NodeStatusProcessor, error) {
	if processor == nil {
		return nil, ErrNilCoreProcessor
	}

	return &NodeStatusProcessor{
		proc: processor,
	}, nil
}

// GetShardStatus will simply forward the node status from an observer
func (nsp *NodeStatusProcessor) GetShardStatus(shardID uint32) (map[string]interface{}, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNodeStatus map[string]interface{}

		err = nsp.proc.CallGetRestEndPoint(observer.Address, NodeStatusPath, &responseNodeStatus)
		if err != nil {
			log.Error("nodeStatus status request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("nodeStatus status request", "shard id", shardID, "observer", observer.Address)
		return responseNodeStatus, nil

	}

	return nil, ErrSendingRequest
}
