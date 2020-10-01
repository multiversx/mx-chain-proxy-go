package process

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// NetworkStatusPath represents the path where an observer exposes his network metrics
const NetworkStatusPath = "/network/status"

// NetworkConfigPath represents the path where an observer exposes his network metrics
const NetworkConfigPath = "/network/config"

// NetworkConfigPath represents the path where an observer exposes his node status metrics
const NodeStatusPath = "/node/status"

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

// GetNetworkStatusMetrics will simply forward the network status metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkStatusPath, &responseNetworkMetrics)
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
	observers, err := nsp.proc.GetAllObservers()
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkConfigPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getNodeStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NodeStatusPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("node status metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("node status metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetNetworkConfigMetrics will compute nonce of the latest block that can be returned
func (nsp *NodeStatusProcessor) GetLatestBlockNonce() (uint64, error) {
	observers, err := nsp.proc.GetAllObservers()
	if err != nil {
		return 0, err
	}

	shardsIDs := getShardsIDs(observers)
	nonces := make([]uint64, 0)
	for shardID := range shardsIDs {
		nodeStatusResponse, err := nsp.getNodeStatusMetrics(shardID)
		if err != nil {
			return 0, err
		}

		if nodeStatusResponse.Error != "" {
			return 0, errors.New(nodeStatusResponse.Error)
		}

		var nonce uint64
		var ok bool
		if shardID != core.MetachainShardId {
			nonce, ok = getNonceFromShard(nodeStatusResponse.Data)
		} else {
			nonce, ok = getNonceFromMeta(nodeStatusResponse.Data)
		}
		if !ok {
			return 0, ErrCannotParseNodeStatusMetrics
		}

		nonces = append(nonces, nonce)
	}

	return getMinValue(nonces), nil
}

func getMinValue(noncesSlice []uint64) uint64 {
	var min uint64
	for idx, value := range noncesSlice {
		if idx == 0 || value < min {
			min = value
		}
	}

	return min
}

func getShardsIDs(observers []*data.NodeData) map[uint32]struct{} {
	shardsIDs := make(map[uint32]struct{})
	for _, observer := range observers {
		shardsIDs[observer.ShardId] = struct{}{}
	}

	return shardsIDs
}

func getNonceFromShard(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, core.MetricCrossCheckBlockHeight)
	if !ok {
		return 0, false
	}

	return getNonceValue(metric)
}

func getNonceFromMeta(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, core.MetricNonce)
	if !ok {
		return 0, false
	}

	return getUint(metric), true
}

func getMetric(nodeStatusData interface{}, metric string) (interface{}, bool) {
	metrics, ok := nodeStatusData.(map[string]interface{})["metrics"].(map[string]interface{})
	if !ok {
		return nil, false
	}

	value, ok := metrics[metric]
	if !ok {
		return nil, false
	}

	return value, true
}

func getNonceValue(value interface{}) (uint64, bool) {
	valueStr, ok := value.(string)
	if !ok {
		return 0, false
	}

	// metric looks like that
	// "meta 886717"
	values := strings.Split(valueStr, " ")
	if len(values) < 2 {
		return 0, false
	}

	nonce, err := strconv.ParseUint(values[1], 10, 64)
	if err != nil {
		return 0, false
	}

	return nonce, true
}

func getUint(value interface{}) uint64 {
	valueFloat, ok := value.(float64)
	if !ok {
		return 0
	}

	return uint64(valueFloat)
}
