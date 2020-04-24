package process

import "github.com/ElrondNetwork/elrond-go/core/check"

// NodeStatusPath represents the path where an observer exposes his nodeStatus
const NodeStatusPath = "/node/status"

// NodeEpochPath represents the path where an observer exposes his epoch metrics
const NodeEpochPath = "/node/epoch"

// NodeConfigPath represents the path where an observer exposes his configuration metrics
const NodeConfigPath = "/node/config"

// NetworkPath represents the path where an observer exposes his network metrics
const NetworkPath = "/network"

// NodeStatusProcessor handles the action needed for fetching data related to status metrics from nodes
type NodeStatusProcessor struct {
	proc Processor
}

// NewNodeStatusProcessor creates a new instance of NodeStatusProcessor
func NewNodeStatusProcessor(processor Processor) (*NodeStatusProcessor, error) {
	if check.IfNil(processor) {
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

// GetEpochMetrics will simply forward the epoch metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetEpochMetrics(shardID uint32) (map[string]interface{}, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseEpochMetrics map[string]interface{}

		err = nsp.proc.CallGetRestEndPoint(observer.Address, NodeEpochPath, &responseEpochMetrics)
		if err != nil {
			log.Error("epoch metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("epoch metrics request", "shard id", shardID, "observer", observer.Address)
		return responseEpochMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetConfigMetrics will simply forward the network configuration metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetConfigMetrics() (map[string]interface{}, error) {
	observers := nsp.proc.GetAllObservers()

	for _, observer := range observers {
		var responseConfigMetrics map[string]interface{}

		err := nsp.proc.CallGetRestEndPoint(observer.Address, NodeConfigPath, &responseConfigMetrics)
		if err != nil {
			log.Error("configuration metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("configuration metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseConfigMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetNetworkMetrics will simply forward the network metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkMetrics(shardID uint32) (map[string]interface{}, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics map[string]interface{}

		err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}
