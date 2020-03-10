package process

import (
	"strconv"
)

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

// GetNodeStatusData will simply forward the node status from an observer
func (nsp *NodeStatusProcessor) GetNodeStatusData(shardId string) (map[string]interface{}, error) {
	if len(shardId) == 0 {
		return nil, ErrInvalidShardId
	}

	shardIdUint, err := strconv.ParseUint(shardId, 10, 32)
	if err != nil {
		return nil, err
	}

	observers, err := nsp.proc.GetObservers(uint32(shardIdUint))
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

		log.Info("nodeStatus status request", "shard id", shardId, "observer", observer.Address)
		return responseNodeStatus, nil

	}

	return nil, ErrSendingRequest
}
