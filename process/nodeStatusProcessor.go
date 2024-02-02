package process

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const (
	// NetworkStatusPath represents the path where an observer exposes his network metrics
	NetworkStatusPath = "/network/status"

	// NetworkConfigPath represents the path where an observer exposes his network metrics
	NetworkConfigPath = "/network/config"

	// NodeStatusPath represents the path where an observer exposes his node status metrics
	NodeStatusPath = "/node/status"

	// AllIssuedESDTsPath represents the path where an observer exposes all the issued ESDTs
	AllIssuedESDTsPath = "/network/esdts"

	// NetworkEsdtTokensPrefix represents the prefix for the path where an observer exposes ESDT tokens of a kind
	NetworkEsdtTokensPrefix = "/network/esdt"

	// DelegatedInfoPath represents the path where an observer exposes his network delegated info
	DelegatedInfoPath = "/network/delegated-info"

	// DirectStakedPath represents the path where an observer exposes his network direct staked info
	DirectStakedPath = "/network/direct-staked-info"

	// RatingsConfigPath represents the path where an observer exposes his ratings metrics
	RatingsConfigPath = "/network/ratings"

	// GenesisNodesConfigPath represents the path where an observer exposes genesis nodes config
	GenesisNodesConfigPath = "/network/genesis-nodes"

	// GasConfigsPath represents the path where an observer exposes gas configs
	GasConfigsPath = "/network/gas-configs"

	// EnableEpochsPath represents the path where an observer exposes all the activation epochs
	EnableEpochsPath = "/network/enable-epochs"

	// MetricCrossCheckBlockHeight is the metric that stores cross block height
	MetricCrossCheckBlockHeight = "erd_cross_check_block_height"

	// MetricAccountsSnapshotNumNodes is the metric that outputs the number of trie nodes written for accounts after snapshot
	MetricAccountsSnapshotNumNodes = "erd_accounts_snapshot_num_nodes"

	// MetricNonce is the metric for monitoring the nonce of a node
	MetricNonce = "erd_nonce"
)

// NodeStatusProcessor handles the action needed for fetching data related to status metrics from nodes
type NodeStatusProcessor struct {
	proc                  Processor
	economicMetricsCacher GenericApiResponseCacheHandler
	cacheValidityDuration time.Duration
	cancelFunc            func()
}

// NewNodeStatusProcessor creates a new instance of NodeStatusProcessor
func NewNodeStatusProcessor(
	processor Processor,
	economicMetricsCacher GenericApiResponseCacheHandler,
	cacheValidityDuration time.Duration,
) (*NodeStatusProcessor, error) {
	if check.IfNil(processor) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(economicMetricsCacher) {
		return nil, ErrNilEconomicMetricsCacher
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}

	return &NodeStatusProcessor{
		proc:                  processor,
		economicMetricsCacher: economicMetricsCacher,
		cacheValidityDuration: cacheValidityDuration,
	}, nil
}

// GetNetworkStatusMetrics will simply forward the network status metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseNetworkMetrics := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkStatusPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseNetworkMetrics, nil

	}

	return nil, WrapObserversError(responseNetworkMetrics.Error)
}

// GetNetworkConfigMetrics will simply forward the network config metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetAllObservers(data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseNetworkMetrics := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err = nsp.proc.CallGetRestEndPoint(observer.Address, NetworkConfigPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseNetworkMetrics, nil

	}

	return nil, WrapObserversError(responseNetworkMetrics.Error)
}

// GetEnableEpochsMetrics will simply forward the activation epochs config metrics from an observer
func (nsp *NodeStatusProcessor) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetAllObservers(data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseEnableEpochsMetrics := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, EnableEpochsPath, &responseEnableEpochsMetrics)
		if err != nil {
			log.Error("enable epochs metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("enable epochs metrics request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseEnableEpochsMetrics, nil
	}

	return nil, WrapObserversError(responseEnableEpochsMetrics.Error)
}

// GetAllIssuedESDTs will forward the issued ESDTs based on the provided type
func (nsp *NodeStatusProcessor) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	if !data.IsValidEsdtPath(tokenType) && tokenType != "" {
		return nil, ErrInvalidTokenType
	}

	observers, err := nsp.proc.GetObservers(core.MetachainShardId, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseAllIssuedESDTs := data.GenericAPIResponse{}
	for _, observer := range observers {

		path := AllIssuedESDTsPath
		if tokenType != "" {
			path = fmt.Sprintf("%s/%s", NetworkEsdtTokensPrefix, tokenType)
		}
		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, path, &responseAllIssuedESDTs)
		if err != nil {
			log.Error("all issued esdts request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("all issued esdts request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseAllIssuedESDTs, nil

	}

	return nil, WrapObserversError(responseAllIssuedESDTs.Error)
}

// GetDelegatedInfo returns the delegated info from nodes
func (nsp *NodeStatusProcessor) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	delegatedInfoResponse := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DelegatedInfoPath, &delegatedInfoResponse)
		if err != nil {
			log.Error("network delegated info request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network delegated info request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &delegatedInfoResponse, nil

	}

	return nil, WrapObserversError(delegatedInfoResponse.Error)
}

// GetDirectStakedInfo returns the delegated info from nodes
func (nsp *NodeStatusProcessor) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	directStakedResponse := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DirectStakedPath, &directStakedResponse)
		if err != nil {
			log.Error("network direct staked request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network direct staked request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &directStakedResponse, nil

	}

	return nil, WrapObserversError(directStakedResponse.Error)
}

// GetRatingsConfig will simply forward the ratings configuration from an observer
func (nsp *NodeStatusProcessor) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetAllObservers(data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseRatingsConfig := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err = nsp.proc.CallGetRestEndPoint(observer.Address, RatingsConfigPath, &responseRatingsConfig)
		if err != nil {
			log.Error("ratings metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("ratings metrics request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseRatingsConfig, nil

	}

	return nil, WrapObserversError(responseRatingsConfig.Error)
}

func (nsp *NodeStatusProcessor) getNodeStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseNetworkMetrics := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err = nsp.proc.CallGetRestEndPoint(observer.Address, NodeStatusPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("node status metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("node status metrics request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseNetworkMetrics, nil

	}

	return nil, WrapObserversError(responseNetworkMetrics.Error)
}

// GetLatestFullySynchronizedHyperblockNonce will compute nonce of the latest hyperblock that can be returned
func (nsp *NodeStatusProcessor) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	shardsIDs, err := nsp.getShardsIDs()
	if err != nil {
		return 0, err
	}

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
		if shardID == core.MetachainShardId {
			nonce, ok = getNonceFromMetachainStatus(nodeStatusResponse.Data)
		} else {
			nonce, ok = getNonceFromShardStatus(nodeStatusResponse.Data)
		}
		if !ok {
			return 0, ErrCannotParseNodeStatusMetrics
		}

		nonces = append(nonces, nonce)
	}

	return getMinNonce(nonces), nil
}

// GetTriesStatistics will return trie statistics
func (nsp *NodeStatusProcessor) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	nodeStatusResponse, err := nsp.getNodeStatusMetrics(shardID)
	if err != nil {
		return nil, err
	}

	return getTrieStatistics(nodeStatusResponse.Data)
}

func getMinNonce(noncesSlice []uint64) uint64 {
	// initialize min with max uint64 value
	min := uint64(math.MaxUint64)
	for _, value := range noncesSlice {
		if value < min {
			min = value
		}
	}

	return min
}

func (nsp *NodeStatusProcessor) getShardsIDs() (map[uint32]struct{}, error) {
	observers, err := nsp.proc.GetAllObservers(data.AvailabilityAll)
	if err != nil {
		return nil, err
	}

	shardsIDs := make(map[uint32]struct{})
	for _, observer := range observers {
		shardsIDs[observer.ShardId] = struct{}{}
	}

	if len(shardsIDs) == 0 {
		return nil, ErrMissingObserver
	}

	return shardsIDs, nil
}

func getNonceFromShardStatus(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, MetricCrossCheckBlockHeight)
	if !ok {
		return 0, false
	}

	return parseMetricCrossCheckBlockHeight(metric)
}

func getNonceFromMetachainStatus(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, MetricNonce)
	if !ok {
		return 0, false
	}

	return getUint(metric), true
}

func getTrieStatistics(nodeStatusData interface{}) (*data.TrieStatisticsAPIResponse, error) {
	trieStatistics := &data.TrieStatisticsAPIResponse{}
	numNodesMetric, ok := getMetric(nodeStatusData, MetricAccountsSnapshotNumNodes)
	if !ok {
		return nil, ErrCannotParseNodeStatusMetrics
	}

	trieStatistics.Data.AccountsSnapshotNumNodes = getUint(numNodesMetric)
	return trieStatistics, nil
}

func getMetric(nodeStatusData interface{}, metric string) (interface{}, bool) {
	metricsMapI, ok := nodeStatusData.(map[string]interface{})
	if !ok {
		return nil, false
	}

	metricsMap, ok := metricsMapI["metrics"]
	if !ok {
		return nil, false
	}

	metrics, ok := metricsMap.(map[string]interface{})
	if !ok {
		return nil, false
	}

	value, ok := metrics[metric]
	if !ok {
		return nil, false
	}

	return value, true
}

func parseMetricCrossCheckBlockHeight(value interface{}) (uint64, bool) {
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

// GetGenesisNodesPubKeys will return genesis nodes public keys
func (nsp *NodeStatusProcessor) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetAllObservers(data.AvailabilityAll)
	if err != nil {
		return nil, err
	}

	response := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err = nsp.proc.CallGetRestEndPoint(observer.Address, GenesisNodesConfigPath, &response)
		if err != nil {
			log.Error("genesis nodes request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("genesis nodes request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &response, nil

	}

	return nil, WrapObserversError(response.Error)
}

// GetGasConfigs will return gas configs
func (nsp *NodeStatusProcessor) GetGasConfigs() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetAllObservers(data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	responseGenesisNodesConfig := data.GenericAPIResponse{}
	for _, observer := range observers {

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, GasConfigsPath, &responseGenesisNodesConfig)
		if err != nil {
			log.Error("gas configs request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("gas configs request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseGenesisNodesConfig, nil

	}

	return nil, WrapObserversError(responseGenesisNodesConfig.Error)
}

// GetEpochStartData will return the epoch-start data for the given epoch and shard
func (nsp *NodeStatusProcessor) GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID, data.AvailabilityAll)
	if err != nil {
		return nil, err
	}

	responseEpochStartData := data.GenericAPIResponse{}
	path := fmt.Sprintf("/node/epoch-start/%d", epoch)
	for _, observer := range observers {

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, path, &responseEpochStartData)
		if err != nil {
			log.Error("epoch start data request", "observer", observer.Address, "shard ID", observer.ShardId, "error", err)
			continue
		}

		log.Info("epoch start data request", "shard ID", observer.ShardId, "observer", observer.Address)
		return &responseEpochStartData, nil
	}

	return nil, WrapObserversError(responseEpochStartData.Error)
}
