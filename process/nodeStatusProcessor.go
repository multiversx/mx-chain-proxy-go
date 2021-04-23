package process

import (
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const maiarListUrl = "https://internal-tools.maiar.com/mex-distribution/eligible-addresses"

const (
	// NetworkStatusPath represents the path where an observer exposes his network metrics
	NetworkStatusPath = "/network/status"

	// NetworkConfigPath represents the path where an observer exposes his network metrics
	NetworkConfigPath = "/network/config"

	// NetworkConfigPath represents the path where an observer exposes his node status metrics
	NodeStatusPath = "/node/status"

	// NodeStatusPath represents the path where an observer exposes all the issued ESDTs
	AllIssuedESDTsPath = "/network/esdts"

	// DelegatedInfoPath represents the path where an observer exposes his network delegated info
	DelegatedInfoPath = "/network/delegated-info"

	// DirectStakedPath represents the path where an observer exposes his network direct staked info
	DirectStakedPath = "/network/direct-staked-info"

	// QueryPath represents the path for a general vm-query
	QueryPath = "/vm-values/query"
)

// NodeStatusProcessor handles the action needed for fetching data related to status metrics from nodes
type NodeStatusProcessor struct {
	proc                  Processor
	economicMetricsCacher GenericApiResponseCacheHandler
	cacheValidityDuration time.Duration
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

		log.Info("network metrics request", "shard ID", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetAllIssuedESDTs will simply forward all the issued ESDTs from an observer in the metachain
func (nsp *NodeStatusProcessor) GetAllIssuedESDTs() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseAllIssuedESDTs *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, AllIssuedESDTsPath, &responseAllIssuedESDTs)
		if err != nil {
			log.Error("all issued esdts request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("all issued esdts request", "shard ID", observer.ShardId, "observer", observer.Address)
		return responseAllIssuedESDTs, nil

	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getAccountList() ([]*data.AccountBalance, error) {
	shardIds := nsp.proc.GetShardIDs()
	accountList := make([]*data.AccountBalance, 0)

	for _, shardId := range shardIds {
		if shardId == core.MetachainShardId {
			continue
		}

		observers, err := nsp.proc.GetObservers(shardId)
		if err != nil {
			return nil, err
		}

		for _, observer := range observers {
			var accountListResponse *data.GenericAPIResponse

			_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DelegatedInfoPath, &accountListResponse)
			if err != nil {
				log.Error("network delegated info request", "observer", observer.Address, "error", err.Error())
				continue
			}

			log.Info("network delegated info request", "shard id", observer.ShardId, "observer", observer.Address)
			if len(accountListResponse.Error) > 0 {
				return nil, errors.New("network delegated info request on observer: " + observer.Address + " - " + accountListResponse.Error)
			}

			accounts, ok := accountListResponse.Data.(*data.AccountBalanceListResponse)
			if !ok {
				return nil, errors.New("network delegated info request on observer: " + observer.Address + " - could not decode response data")
			}

			accountList = append(accountList, accounts.List...)
			break
		}
	}

	return accountList, nil
}

// GetDelegatedInfo returns the delegated info from nodes
func (nsp *NodeStatusProcessor) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var delegatedInfoResponse *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DelegatedInfoPath, &delegatedInfoResponse)
		if err != nil {
			log.Error("network delegated info request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network delegated info request", "shard id", observer.ShardId, "observer", observer.Address)
		return delegatedInfoResponse, nil

	}

	return nil, ErrSendingRequest
}

// GetDelegatedInfo returns the delegated info from nodes
func (nsp *NodeStatusProcessor) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var directStakedResponse *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DirectStakedPath, &directStakedResponse)
		if err != nil {
			log.Error("network direct staked request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network direct staked request", "shard id", observer.ShardId, "observer", observer.Address)
		return directStakedResponse, nil

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

// GetEligibleAddresses returns the list of eligible addresses for Maiar MEX multiplier
func (nsp *NodeStatusProcessor) getEligibleAddresses() (*data.MaiarReferalApiResponse, error) {
	maiarEligibleList := &data.MaiarReferalApiResponse{}
	status, err := nsp.proc.CallGetRestEndPoint(maiarListUrl, "", maiarEligibleList)
	if err != nil {
		log.Error("error fetching maiar eligible list", err.Error())
		return nil, ErrSendingRequest
	}
	if status != http.StatusOK {
		log.Error("error status received fetching maiar eligible list", status)
		return nil, ErrSendingRequest
	}

	return maiarEligibleList, nil
}

func (nsp *NodeStatusProcessor) CreateSnapshot() (*data.GenericAPIResponse, error) {
	// Create final file - do this first, since if it errors, there's no point in doing all the work
	file, err:= core.CreateFile(core.ArgCreateFileArgument{
		Directory: "/home/ubuntu/snapshots",
		Prefix: "snapshot-10",
		FileExtension: "json",
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		fileCloseErr := file.Close()
		log.Error("error closing snapshot file", fileCloseErr)
	}()

	// 1. Gather Data
	// 1.1 Fetch maiar list - done
	//maiarData, err := nsp.getEligibleAddresses()
	//if err != nil {
	//	return nil, err
	//}

	maiarData := &data.MaiarReferalApiResponse{
		Total: 0,
		Addresses: make([]string, 0),
	}

	// 1.2 Fetch delegation manager data - done
	delegatedInfo, err := nsp.getDecodedDelegatedList()
	if err != nil {
		return nil, err
	}

	// 1.3 Fetch delegation legacy data
	legacyDelegatedInfo, err := nsp.getLegacyDelegationData()
	if err != nil {
		return nil, err
	}

	// 1.4 Fetch staking data
	stakingData, err := nsp.getDecodedDirectStakedInfo()
	if err != nil {
		return nil, err
	}

	// 1.5 Fetch all accounts data
	accountBalances, err := nsp.getAccountList()
	if err != nil {
		return nil, err
	}


	// 2. Merge data
	snapshotList := make([]*data.SnapshotItem, len(accountBalances))
	for i, accountBalance := range accountBalances {
		snapshotList[i] = nsp.buildSnapshotItem(
			accountBalance,
			maiarData,
			delegatedInfo,
			legacyDelegatedInfo,
			stakingData,
		)
	}

	jsonEncoded, err := json.Marshal(snapshotList)
	if err != nil {
		return nil, err
	}
	_, err = file.Write(jsonEncoded)
	if err != nil {
		return nil, err
	}

	return &data.GenericAPIResponse{
		Data: "ok",
		Error: "",
		Code: data.ReturnCodeSuccess,
	}, nil
}

func (nsp *NodeStatusProcessor) buildSnapshotItem(
	accountBalance *data.AccountBalance,
	maiarData *data.MaiarReferalApiResponse,
	delegatedInfo *data.DelegationList,
	legacyDelegatedInfo *data.DelegationList,
	stakingData *data.DirectStakedValueList,
	) *data.SnapshotItem {
	si := &data.SnapshotItem{
		Address: accountBalance.Address,
		Balance: accountBalance.Balance,
		Staked: "0",
		Waiting: "0",
		Unstaked: "0",
		Unclaimed: "0",
		IsMaiarEligible: false,
	}

	for _, maiarEligible := range maiarData.Addresses {
		if accountBalance.Address == maiarEligible {
			si.IsMaiarEligible = true
			break
		}
	}

	for _, delegationInfo := range delegatedInfo.List {
		if accountBalance.Address == delegationInfo.DelegatorAddress {
			staked, _ := big.NewInt(0).SetString(si.Staked, 10)
			unstaked, _ := big.NewInt(0).SetString(si.Unstaked, 10)
			unclaimed, _ := big.NewInt(0).SetString(si.Unclaimed, 10)

			newStaked, _ := big.NewInt(0).SetString(delegationInfo.Total, 10)
			newUnstaked, _ := big.NewInt(0).SetString(delegationInfo.UndelegatedTotal, 10)
			newUnclaimed, _ := big.NewInt(0).SetString(delegationInfo.UnclaimedTotal, 10)

			si.Staked = big.NewInt(0).Add(staked, newStaked).String()
			si.Unstaked = big.NewInt(0).Add(unstaked, newUnstaked).String()
			si.Unclaimed = big.NewInt(0).Add(unclaimed, newUnclaimed).String()

			break
		}
	}

	for _, delegationInfo := range legacyDelegatedInfo.List {
		if accountBalance.Address == delegationInfo.DelegatorAddress {
			staked, _ := big.NewInt(0).SetString(si.Staked, 10)
			unstaked, _ := big.NewInt(0).SetString(si.Unstaked, 10)
			unclaimed, _ := big.NewInt(0).SetString(si.Unclaimed, 10)
			waiting, _ := big.NewInt(0).SetString(si.Waiting, 10)

			newStaked, _ := big.NewInt(0).SetString(delegationInfo.Total, 10)
			newUnstaked, _ := big.NewInt(0).SetString(delegationInfo.UndelegatedTotal, 10)
			newUnclaimed, _ := big.NewInt(0).SetString(delegationInfo.UnclaimedTotal, 10)
			newWaiting, _ := big.NewInt(0).SetString(delegationInfo.WaitingTotal, 10)

			si.Staked = big.NewInt(0).Add(staked, newStaked).String()
			si.Unstaked = big.NewInt(0).Add(unstaked, newUnstaked).String()
			si.Unclaimed = big.NewInt(0).Add(unclaimed, newUnclaimed).String()
			si.Waiting = big.NewInt(0).Add(waiting, newWaiting).String()

			break
		}
	}

	for _, stakeData := range stakingData.List {
		if accountBalance.Address == stakeData.Address {
			staked, _ := big.NewInt(0).SetString(si.Staked, 10)
			unstaked, _ := big.NewInt(0).SetString(si.Unstaked, 10)

			newStaked, _ := big.NewInt(0).SetString(stakeData.Total, 10)
			newUnstaked, _ := big.NewInt(0).SetString(stakeData.Unstaked, 10)

			si.Staked = big.NewInt(0).Add(staked, newStaked).String()
			si.Unstaked = big.NewInt(0).Add(unstaked, newUnstaked).String()
			break
		}
	}

	return si
}

func (nsp *NodeStatusProcessor) getDecodedDelegatedList() (*data.DelegationList, error) {
	delegatedInfo, err := nsp.GetDelegatedInfo()
	if err != nil {
		return nil, err
	}

	if delegatedInfo.Error != "" {
		return nil, errors.New(delegatedInfo.Error)
	}

	decodedList, ok := delegatedInfo.Data.(*data.DelegationList)
	if !ok {
		return nil, ErrInvalidDelegationListReceived
	}

	return decodedList, nil
}

func (nsp *NodeStatusProcessor) getDecodedDirectStakedInfo() (*data.DirectStakedValueList, error) {
	stakedInfo, err := nsp.GetDirectStakedInfo()
	if err != nil {
		return nil, err
	}

	if stakedInfo.Error != "" {
		return nil, errors.New(stakedInfo.Error)
	}

	decodedList, ok := stakedInfo.Data.(*data.DirectStakedValueList)
	if !ok {
		return nil, ErrInvalidDirectStakeListReceived
	}

	return decodedList, nil
}

func (nsp *NodeStatusProcessor) getLegacyDelegationData() (*data.DelegationList, error) {


	return nil, nil
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
	observers, err := nsp.proc.GetAllObservers()
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
	metric, ok := getMetric(nodeStatusData, core.MetricCrossCheckBlockHeight)
	if !ok {
		return 0, false
	}

	return parseMetricCrossCheckBlockHeight(metric)
}

func getNonceFromMetachainStatus(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, core.MetricNonce)
	if !ok {
		return 0, false
	}

	return getUint(metric), true
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
