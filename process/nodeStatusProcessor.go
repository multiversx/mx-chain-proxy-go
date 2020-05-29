package process

import (
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// NodeStatusPath represents the path where an observer exposes his nodeStatus
const NodeStatusPath = "/node/status"

// NetworkStatusPath represents the path where an observer exposes his network metrics
const NetworkStatusPath = "/network/status"

// NetworkConfigPath represents the path where an observer exposes his network metrics
const NetworkConfigPath = "/network/config"

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
func (nsp *NodeStatusProcessor) GetShardStatus(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNodeStatus *data.GenericAPIResponse

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

// GetNetworkStatusMetrics will simply forward the network status metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkStatusPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetNetworkConfigMetrics will simply forward the network config metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	observers := nsp.proc.GetAllObservers()

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkConfigPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}
