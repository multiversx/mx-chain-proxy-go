package process

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ALTree/bigfloat"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const contractPrefix = "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqq"
const legacyDelegationContract = "erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt"
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

	// AccountsListPath represents the path where an observer exposes the path to return a full list of accounts
	AccountsListPath = "/network/accounts-info"
)

// NodeStatusProcessor handles the action needed for fetching data related to status metrics from nodes
type NodeStatusProcessor struct {
	proc                  Processor
	economicMetricsCacher GenericApiResponseCacheHandler
	cacheValidityDuration time.Duration
	pubKeyConverter core.PubkeyConverter
	undelagatedSnapshots []string
	snapshots []string
}

// NewNodeStatusProcessor creates a new instance of NodeStatusProcessor
func NewNodeStatusProcessor(
	processor Processor,
	economicMetricsCacher GenericApiResponseCacheHandler,
	cacheValidityDuration time.Duration,
	pubKeyConverter core.PubkeyConverter,
) (*NodeStatusProcessor, error) {
	if check.IfNil(processor) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(economicMetricsCacher) {
		return nil, ErrNilEconomicMetricsCacher
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}
	if cacheValidityDuration <= 0 {
		return nil, ErrInvalidCacheValidityDuration
	}

	return &NodeStatusProcessor{
		proc:                  processor,
		economicMetricsCacher: economicMetricsCacher,
		cacheValidityDuration: cacheValidityDuration,
		pubKeyConverter: pubKeyConverter,
		undelagatedSnapshots: []string{
			//"undelegated-10-2021-05-04-15-53-46.json",
			//"undelegated-10-2021-05-04-17-04-51.json",
			//"undelegated-10-2021-05-04-18-29-38.json",
			//"undelegated-10-2021-05-04-19-12-05.json",
			//"undelegated-10-2021-05-04-19-47-39.json",
			//"undelegated-10-2021-05-04-20-59-34.json",
			//"undelegated-10-2021-05-04-21-55-41.json",
		},
		snapshots: []string {
			"snapshot-10-2021-05-06-09-15-54.json",
			"snapshot-10-2021-05-06-10-58-56.json",
			"snapshot-10-2021-05-06-12-19-53.json",
			"snapshot-10-2021-05-06-13-30-59.json",
			"snapshot-10-2021-05-06-15-19-06.json",
			"snapshot-10-2021-05-06-17-49-57.json",
			"snapshot-10-2021-05-06-20-01-02.json",
			"snapshot-10-day4backup-2021-05-06-22-49-26.json",
		},
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

		if shardId != 0 {
			continue
		}

		observers, err := nsp.proc.GetObservers(shardId)
		if err != nil {
			return nil, err
		}

		for _, observer := range observers {
			var accountListResponse data.AccountBalanceListResponse

			_, err := nsp.proc.CallGetRestEndPoint(observer.Address, AccountsListPath, &accountListResponse)
			if err != nil {
				log.Error("get account list request", "observer", observer.Address, "error", err.Error())
				continue
			}

			log.Info("get account list request", "shard id", observer.ShardId, "observer", observer.Address)
			if len(accountListResponse.Error) > 0 {
				return nil, errors.New("get account list request: " + observer.Address + " - " + accountListResponse.Error)
			}

			log.Info("get account list request", "shard id", observer.ShardId, "fetched", len(accountListResponse.Data.List))

			accountList = append(accountList, accountListResponse.Data.List...)
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

func (nsp *NodeStatusProcessor) computeMexValues(snapshotItems []*data.SnapshotItem) ([]*data.MexItem, error) {
	// Step 1 - find out multiplier

	exponent := big.NewFloat(0.95)
	weekOneMexApproxValue, _ := big.NewFloat(0).SetString("151200000000000000000000000")
	// To find multiplier we need the sum of all eased values
	fullEasedSum := big.NewFloat(0)
	for _, snapshotItem := range snapshotItems {
		multiplierOneBalance := big.NewFloat(0)
		multiplierOneQuarterBalance := big.NewFloat(0)
		multiplierOneHalfBalance := big.NewFloat(0)


		balance, _ := big.NewFloat(0).SetString(snapshotItem.Balance)
		staked, _ := big.NewFloat(0).SetString(snapshotItem.Staked)
		waiting, _ := big.NewFloat(0).SetString(snapshotItem.Waiting)
		unstaked, _ := big.NewFloat(0).SetString(snapshotItem.Unstaked)
		unclaimed, _ := big.NewFloat(0).SetString(snapshotItem.Unclaimed)

		multiplierOneBalance = multiplierOneBalance.Add(unclaimed, unstaked)
		multiplierOneQuarterBalance.Set(waiting)
		if snapshotItem.IsMaiarEligible {
			multiplierOneQuarterBalance = multiplierOneQuarterBalance.Add(multiplierOneQuarterBalance, balance)
		} else {
			multiplierOneBalance = multiplierOneBalance.Add(multiplierOneBalance, balance)
		}

		multiplierOneHalfBalance.Set(staked)

		multiplierOneBalanceEased := bigfloat.Pow(multiplierOneBalance, exponent)
		multiplierOneQuarterBalanceEased := bigfloat.Pow(multiplierOneQuarterBalance, exponent)
		multiplierOneHalfBalanceEased := bigfloat.Pow(multiplierOneHalfBalance, exponent)

		fullEasedSum = fullEasedSum.Add(fullEasedSum, multiplierOneBalanceEased)
		fullEasedSum = fullEasedSum.Add(fullEasedSum, multiplierOneQuarterBalanceEased)
		fullEasedSum = fullEasedSum.Add(fullEasedSum, multiplierOneHalfBalanceEased)
	}

	mexMultiplier := big.NewFloat(0).Quo(weekOneMexApproxValue, fullEasedSum)
	log.Info("======= mex multiplier ========", "having", mexMultiplier.String())

	// Step 2 - compute mex values
	mexItems := make([]*data.MexItem, len(snapshotItems))
	for i := 0; i < len(snapshotItems); i++ {
		balance, _ := big.NewFloat(0).SetString(snapshotItems[i].Balance)
		staked, _ := big.NewFloat(0).SetString(snapshotItems[i].Staked)
		waiting, _ := big.NewFloat(0).SetString(snapshotItems[i].Waiting)
		unstaked, _ := big.NewFloat(0).SetString(snapshotItems[i].Unstaked)
		unclaimed, _ := big.NewFloat(0).SetString(snapshotItems[i].Unclaimed)

		multiplierOneBalance := big.NewFloat(0)
		multiplierOneQuarterBalance := big.NewFloat(0)
		multiplierOneHalfBalance := big.NewFloat(0)

		multiplierOneBalance = multiplierOneBalance.Add(unclaimed, unstaked)
		multiplierOneQuarterBalance.Set(waiting)
		if snapshotItems[i].IsMaiarEligible {
			multiplierOneQuarterBalance = multiplierOneQuarterBalance.Add(multiplierOneQuarterBalance, balance)
		} else {
			multiplierOneBalance = multiplierOneBalance.Add(multiplierOneBalance, balance)
		}

		multiplierOneHalfBalance.Set(staked)

		multiplierOneBalanceEased := bigfloat.Pow(multiplierOneBalance, exponent)
		multiplierOneQuarterBalanceEased := bigfloat.Pow(multiplierOneQuarterBalance, exponent)
		multiplierOneHalfBalanceEased := bigfloat.Pow(multiplierOneHalfBalance, exponent)

		oneMex := big.NewFloat(0).Mul(multiplierOneBalanceEased, mexMultiplier)
		oneQuarterMex := big.NewFloat(0).Mul(multiplierOneQuarterBalanceEased, mexMultiplier)
		oneHalfMex := big.NewFloat(0).Mul(multiplierOneHalfBalanceEased, mexMultiplier)

		mexFullVal := big.NewFloat(0).Add(oneMex, oneQuarterMex)
		mexFullVal = mexFullVal.Add(mexFullVal, oneHalfMex)
		mexFullValInt, _ := mexFullVal.Int(nil)

		mexItems[i] = &data.MexItem{
			Address: snapshotItems[i].Address,
			Value:   mexFullValInt.String(),
		}
	}

	return mexItems, nil
}

func (nsp *NodeStatusProcessor) loadUndelegatedSnapshots() ([][]*data.Delegator, error) {
	delegators := make([][]*data.Delegator, len(nsp.undelagatedSnapshots))
	for i := 0; i < len(nsp.undelagatedSnapshots); i++ {
		var delegationList data.DelegationListResponse
		err := core.LoadJsonFile(&delegationList, "/home/ubuntu/snapshots/undelegate/" + nsp.undelagatedSnapshots[i])
		if err != nil {
			log.Error("unable to load delegation file", "file", nsp.undelagatedSnapshots[i])
			return nil, err
		}

		if delegationList.Data.List == nil {
			return nil, errors.New("unable to load delegation file")
		}

		delegators[i] = delegationList.Data.List
	}

	return delegators, nil
}

func (nsp *NodeStatusProcessor) loadLocalSnapshots() ([][]*data.SnapshotItem, error) {
	snapshotList := make([][]*data.SnapshotItem, len(nsp.snapshots))
	for i := 0; i < len(nsp.snapshots); i++ {
		var snapshot []*data.SnapshotItem
		err := core.LoadJsonFile(&snapshot, "/home/ubuntu/snapshots/" + nsp.snapshots[i])
		if err != nil {
			log.Error("unable to load snapshots file", "file", nsp.snapshots[i])
			return nil, err
		}

		snapshotList[i] = snapshot
	}

	return snapshotList, nil
}

func (nsp *NodeStatusProcessor) mergeSnapshotsWithUndelegate(undelegated [][]*data.Delegator, snapshots [][]*data.SnapshotItem) ([][]*data.SnapshotItem, error) {
	for i := 0; i < len(snapshots); i++ {
		// For each item in the current snapshot
		for sInternalIndex := 0; sInternalIndex < len(snapshots[i]); sInternalIndex ++ {
			// Find an undelegated item that matches account
			for uIndex := 0; uIndex < len(undelegated[i]); uIndex++ {
				if snapshots[i][sInternalIndex].Address == undelegated[i][uIndex].DelegatorAddress {
					currentUndelegateTotal, ok := big.NewInt(0).SetString(snapshots[i][sInternalIndex].Unstaked, 10)
					legacyUndelegate, legOk := big.NewInt(0).SetString(undelegated[i][uIndex].UndelegatedTotal, 10)
					if !ok || !legOk {
						log.Error("could not decode unstaked value from snapshot",
							"current value", snapshots[i][sInternalIndex].Unstaked,
							"new value", undelegated[i][uIndex].UndelegatedTotal,

						)
						return nil, ErrSendingRequest
					}

					snapshots[i][sInternalIndex].Unstaked = big.NewInt(0).Add(currentUndelegateTotal, legacyUndelegate).String()
					break
				}
			}
		}
	}

	return snapshots, nil
}

func (nsp *NodeStatusProcessor) mergeSnapshotsTogether(snapshots [][]*data.SnapshotItem) ([]*data.SnapshotItem, error) {
	snapshotsMap := make(map[string]*data.SnapshotItem)
	for i := 0; i < len(snapshots); i++ {
		for j := 0; j < len(snapshots[i]); j++ {
			currentSnapshot, exists := snapshotsMap[snapshots[i][j].Address]
			if !exists {
				snapshotsMap[snapshots[i][j].Address] = snapshots[i][j]
				continue
			}

			currentBalance, _ := big.NewInt(0).SetString(currentSnapshot.Balance, 10)
			currentStaked, _ := big.NewInt(0).SetString(currentSnapshot.Staked, 10)
			currentWaiting, _ := big.NewInt(0).SetString(currentSnapshot.Waiting, 10)
			currentUnstaked, _ := big.NewInt(0).SetString(currentSnapshot.Unstaked, 10)
			currentUnclaimed, _ := big.NewInt(0).SetString(currentSnapshot.Unclaimed, 10)

			newBalance, _ := big.NewInt(0).SetString(snapshots[i][j].Balance, 10)
			newStaked, _ := big.NewInt(0).SetString(snapshots[i][j].Staked, 10)
			newWaiting, _ := big.NewInt(0).SetString(snapshots[i][j].Waiting, 10)
			newUnstaked, _ := big.NewInt(0).SetString(snapshots[i][j].Unstaked, 10)
			newUnclaimed, _ := big.NewInt(0).SetString(snapshots[i][j].Unclaimed, 10)

			snapshotsMap[snapshots[i][j].Address].Balance = big.NewInt(0).Add(currentBalance, newBalance).String()
			snapshotsMap[snapshots[i][j].Address].Staked = big.NewInt(0).Add(currentStaked, newStaked).String()
			snapshotsMap[snapshots[i][j].Address].Waiting = big.NewInt(0).Add(currentWaiting, newWaiting).String()
			snapshotsMap[snapshots[i][j].Address].Unstaked = big.NewInt(0).Add(currentUnstaked, newUnstaked).String()
			snapshotsMap[snapshots[i][j].Address].Unclaimed = big.NewInt(0).Add(currentUnclaimed, newUnclaimed).String()
		}
	}

	snapshotItems := make([]*data.SnapshotItem, 0)
	for _, item := range snapshotsMap {
		snapshotItems = append(snapshotItems, item)
	}

	return snapshotItems, nil
}


// Only undelegate values
//func (nsp *NodeStatusProcessor) CreateSnapshot(timestamp string) (*data.GenericAPIResponse, error) {
//	// Create final file - do this first, since if it errors, there's no point in doing all the work
//	file, err:= core.CreateFile(core.ArgCreateFileArgument{
//		Directory: "/home/ubuntu/snapshots/undelegate",
//		Prefix: "undelegated-10",
//		FileExtension: "json",
//	})
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		fileCloseErr := file.Close()
//		if fileCloseErr != nil {
//			log.Error("error closing snapshot file", fileCloseErr.Error())
//		}
//
//		log.Info("closed file...")
//	}()
//
//	unstakeList, err := nsp.getLegacyDelegation()
//	if err != nil {
//		return nil, err
//	}
//	jsonEncoded, err := json.Marshal(unstakeList)
//	if err != nil {
//		return nil, err
//	}
//	_, err = file.Write(jsonEncoded)
//	if err != nil {
//		return nil, err
//	}
//
//	return &data.GenericAPIResponse{
//		Data: "ok",
//		Error: "",
//		Code: data.ReturnCodeSuccess,
//	}, nil
//}

// CreateSnapshot - mex indexing - should remove undelegate snapshot
//func (nsp *NodeStatusProcessor) CreateSnapshot(timestamp string) (*data.GenericAPIResponse, error) {
//	indexer, err := NewSnapshotIndexer()
//	if err != nil {
//		return nil, err
//	}
//	// 1. get undelegated items and index them - only relevant for 1st week
//	//log.Info("started fetching undelegated data...")
//	//undelegatedData, err := nsp.loadUndelegatedSnapshots()
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//log.Info("indexing undelegate values...")
//	//for i := 0; i < len(undelegatedData); i++ {
//	//	err = indexer.IndexUndelegatedValues(undelegatedData[i], i)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//}
//
//	log.Info("started fetching local snapshots...")
//	localSnapshots, err := nsp.loadLocalSnapshots()
//	if err != nil {
//		return nil, err
//	}
//
//	//log.Info("merging snapshots with undelegate...")
//	//correctedSnapshots, err := nsp.mergeSnapshotsWithUndelegate(undelegatedData, localSnapshots)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	log.Info("merging all snapshots together...")
//	mexComputeList, err := nsp.mergeSnapshotsTogether(localSnapshots)
//	if err != nil {
//		return nil, err
//	}
//
//	balance := big.NewInt(0)
//	staked := big.NewInt(0)
//	waiting := big.NewInt(0)
//	unstaked := big.NewInt(0)
//	unclaimed := big.NewInt(0)
//	total := big.NewInt(0)
//	for _, snapshotItem := range mexComputeList {
//		balanceBig, _ := big.NewInt(0).SetString(snapshotItem.Balance, 10)
//		stakedBig, _ := big.NewInt(0).SetString(snapshotItem.Staked, 10)
//		waitingBig, _ := big.NewInt(0).SetString(snapshotItem.Waiting, 10)
//		unstakedBig, _ := big.NewInt(0).SetString(snapshotItem.Unstaked, 10)
//		unclaimedBig, _ := big.NewInt(0).SetString(snapshotItem.Unclaimed, 10)
//
//
//		balance = balance.Add(balance, balanceBig)
//		staked = staked.Add(staked, stakedBig)
//		waiting = waiting.Add(waiting, waitingBig)
//		unstaked = unstaked.Add(unstaked, unstakedBig)
//		unclaimed = unclaimed.Add(unclaimed, unclaimedBig)
//	}
//
//	total = total.Add(total, balance)
//	total = total.Add(total, staked)
//	total = total.Add(total, waiting)
//	total = total.Add(total, unstaked)
//	total = total.Add(total, unclaimed)
//
//	log.Info("egld value", "balance", balance.String())
//	log.Info("egld value", "staked", staked.String())
//	log.Info("egld value", "waiting", waiting.String())
//	log.Info("egld value", "unstaked", unstaked.String())
//	log.Info("egld value", "unclaimed", unclaimed.String())
//	log.Info("egld value", "total", total.String())
//
//	log.Info("computing actual mex values")
//	mexValues, err := nsp.computeMexValues(mexComputeList)
//	if err != nil {
//		return nil, err
//	}
//
//	fullVal := big.NewInt(0)
//	for _, item := range mexValues {
//		itemMex, _ := big.NewInt(0).SetString(item.Value, 10)
//		fullVal = fullVal.Add(fullVal, itemMex)
//	}
//
//	log.Info("gathered mex value", "val", fullVal.String())
//
//	log.Info("indexing mex values...", "having", len(mexValues))
//	err = indexer.IndexMexValues(mexValues)
//	if err != nil {
//		return nil, err
//	}
//
//	return &data.GenericAPIResponse{
//		Data: "ok",
//		Error: "",
//		Code: data.ReturnCodeSuccess,
//	}, nil
//}

func (nsp *NodeStatusProcessor) CreateSnapshot(timestamp string) (*data.GenericAPIResponse, error) {

	var snapshot []*data.SnapshotItem
	// LOAD FIRST DAY
	err := core.LoadJsonFile(&snapshot, "/home/ubuntu/snapshots/" + nsp.snapshots[3])
	if err != nil {
		log.Error("unable to load snapshots file", "file", nsp.snapshots[3])
		return nil, err
	}

	var snapshot2 []*data.SnapshotItem
	// LOAD FIRST DAY
	err = core.LoadJsonFile(&snapshot2, "/home/ubuntu/snapshots/" + nsp.snapshots[7])
	if err != nil {
		log.Error("unable to load snapshots file", "file", nsp.snapshots[7])
		return nil, err
	}

	snapshot = append(snapshot, snapshot2...)

	activeList, _ := nsp.getLegacyDelegationStakingList()
	waitingList, _ := nsp.getLegacyDelegationStakingList()

	for index, snapshotItem := range snapshot {
		if snapshotItem.Unstaked != "0" {
			continue
		}

		for _, activeItem := range activeList {
			if activeItem.DelegatorAddress != snapshotItem.Address {
				continue
			}

			currentActive, _ := big.NewInt(0).SetString(snapshot[index].Staked, 10)
			newActive, _ := big.NewInt(0).SetString(activeItem.Total, 10)
			snapshot[index].Staked = big.NewInt(0).Add(currentActive, newActive).String()
		}

		for _, waitingItem := range waitingList {
			if waitingItem.DelegatorAddress != snapshotItem.Address {
				continue
			}

			currentWaiting, _ := big.NewInt(0).SetString(snapshot[index].Waiting, 10)
			newWaiting, _ := big.NewInt(0).SetString(waitingItem.WaitingTotal, 10)
			snapshot[index].Waiting = big.NewInt(0).Add(currentWaiting, newWaiting).String()
		}
	}


	// Now save the mf thing back
	file, err:= core.CreateFile(core.ArgCreateFileArgument{
		Directory: "/home/ubuntu/snapshots/week2/fixed",
		Prefix: "snapshot-10",
		FileExtension: "json",
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		fileCloseErr := file.Close()
		if fileCloseErr != nil {
			log.Error("error closing snapshot file", fileCloseErr.Error())
		}

		log.Info("closed file...")
	}()

	jsonEncoded, err := json.Marshal(snapshot)
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

// Flull snapshot generator
//func (nsp *NodeStatusProcessor) CreateSnapshot(timestamp string) (*data.GenericAPIResponse, error) {
//	// Create final file - do this first, since if it errors, there's no point in doing all the work
//	file, err:= core.CreateFile(core.ArgCreateFileArgument{
//		Directory: "/home/ubuntu/snapshots/week2",
//		Prefix: "snapshot-10-day4backup",
//		FileExtension: "json",
//	})
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		fileCloseErr := file.Close()
//		if fileCloseErr != nil {
//			log.Error("error closing snapshot file", fileCloseErr.Error())
//		}
//
//		log.Info("closed file...")
//	}()
//
//	// 1. Gather Data
//	// 1.1 Fetch maiar list - done
//	maiarData, err := nsp.getEligibleAddresses()
//	if err != nil {
//		return nil, err
//	}
//
//	// 1.2 Fetch delegation manager data - done
//	delegatedInfo, err := nsp.getDecodedDelegatedList()
//	if err != nil {
//		return nil, err
//	}
//
//	// 1.3 Fetch delegation legacy data
//	legacyDelegatedInfo, err := nsp.getLegacyDelegation()
//	if err != nil {
//		return nil, err
//	}
//
//	// 1.4 Fetch staking data
//	stakingData, err := nsp.getDecodedDirectStakedInfo()
//	if err != nil {
//		return nil, err
//	}
//
//	// 1.5 Fetch all accounts data
//	accountBalances, err := nsp.getAccountList()
//	if err != nil {
//		return nil, err
//	}
//
//
//	log.Info("merging lists....", "having a list of", len(accountBalances))
//	// 2. Merge data
//	snapshotList := make([]*data.SnapshotItem, 0)
//	exceptions := getExceptions()
//	for _, accountBalance := range accountBalances {
//		if exceptions[accountBalance.Address] {
//			continue
//		}
//		if strings.HasPrefix(accountBalance.Address, contractPrefix) {
//			continue
//		}
//
//		sl := nsp.buildSnapshotItem(
//			accountBalance,
//			maiarData,
//			delegatedInfo,
//			legacyDelegatedInfo,
//			stakingData,
//		)
//
//		if sl.Balance == "0" &&
//			sl.Unstaked == "0" &&
//			sl.Staked == "0" &&
//			sl.Unclaimed == "0" &&
//			sl.Waiting == "0" {
//			continue
//		}
//
//		snapshotList = append(snapshotList, sl)
//	}
//
//	jsonEncoded, err := json.Marshal(snapshotList)
//	if err != nil {
//		return nil, err
//	}
//	_, err = file.Write(jsonEncoded)
//	if err != nil {
//		return nil, err
//	}
//
//	// Now that we have the snapshot saved, index it
//	es, err := NewSnapshotIndexer()
//	if err != nil {
//		return nil, err
//	}
//
//
//	log.Info("started indexing snapshot...", "having remaining length", len(snapshotList))
//	err = es.IndexSnapshot(snapshotList, timestamp)
//	if err != nil {
//		return nil, err
//	}
//
//	return &data.GenericAPIResponse{
//		Data: "ok",
//		Error: "",
//		Code: data.ReturnCodeSuccess,
//	}, nil
//}

func (nsp *NodeStatusProcessor) buildSnapshotItem(
	accountBalance *data.AccountBalance,
	maiarData *data.MaiarReferalApiResponse,
	delegatedInfo *data.DelegationListResponse,
	legacyDelegatedInfo *data.DelegationListResponse,
	stakingData *data.DirectStakedValueListResponse,
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

	for _, delegationInfo := range delegatedInfo.Data.List {
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

	for _, delegationInfo := range legacyDelegatedInfo.Data.List {
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

			// For legacy delegation we do not break - waiting list can contain multiple entries
		}
	}

	for _, stakeData := range stakingData.Data.List {
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

func (nsp *NodeStatusProcessor) indexUndelegated() error {


	return nil
}

func (nsp *NodeStatusProcessor) getDecodedDelegatedList() (*data.DelegationListResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var delegatedInfoResponse data.DelegationListResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DelegatedInfoPath, &delegatedInfoResponse)
		if err != nil {
			log.Error("network delegated info request", "observer", observer.Address, "error", err)
			continue
		}

		log.Info("network delegated info request", "shard id", observer.ShardId, "observer", observer.Address)

		if delegatedInfoResponse.Error != "" {
			log.Error("received err", delegatedInfoResponse.Error)
			return nil, errors.New(delegatedInfoResponse.Error)
		}

		log.Info("delegation info debugger", "delegation list first item total", delegatedInfoResponse.Data.List[0].Total)

		return &delegatedInfoResponse, nil
	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getDecodedDirectStakedInfo() (*data.DirectStakedValueListResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var directStakedResponse data.DirectStakedValueListResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, DirectStakedPath, &directStakedResponse)
		if err != nil {
			log.Error("network direct staked request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network direct staked request", "shard id", observer.ShardId, "observer", observer.Address)
		if directStakedResponse.Error != "" {
			return nil, errors.New(directStakedResponse.Error)
		}


		log.Info("direct stake info debugger", "direct stake list first item total", directStakedResponse.Data.List[0].Total)
		return &directStakedResponse, nil

	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getLegacyDelegation() (*data.DelegationListResponse, error) {
	delegationList := &data.DelegationListResponse{}
	delegationList.Data = struct {
		List []*data.Delegator `json:"list"`
	}(struct{ List []*data.Delegator }{List: make([]*data.Delegator, 0)})

	numUsers, err := nsp.getLegacyNumUsers()
	if err != nil {
		return nil, err
	}

	for i := int64(1); i <= numUsers; i++ {
		userAddress, err := nsp.getLegacyUserAddressByIndex(big.NewInt(i))
		if err != nil {
			return nil, err
		}

		userStakeValues, err := nsp.getLegacyUserStakeValues(userAddress)
		if err != nil {
			return nil, err
		}

		withdrawOnly := userStakeValues[0]
		waiting := userStakeValues[1]
		active := userStakeValues[2]
		unstaked := userStakeValues[3]
		deferred := userStakeValues[4]

		if len(withdrawOnly) == 0 &&
			len(waiting) == 0 &&
			len(active) == 0 &&
			len(unstaked) == 0 &&
			len(deferred) == 0 {
			continue
		}

		hexUserAddress, _ := hex.DecodeString(userAddress)
		undelegated := big.NewInt(0).Add(
			big.NewInt(0).SetBytes(unstaked),
			big.NewInt(0).SetBytes(deferred),
		)
		undelegated = undelegated.Add(undelegated, big.NewInt(0).SetBytes(withdrawOnly))

		if undelegated.String() == "0" {
			continue
		}

		delegationItem := &data.Delegator{
			DelegatorAddress: nsp.pubKeyConverter.Encode(hexUserAddress),
			DelegatedTo: []*data.DelegationItem{{
				DelegationScAddress: legacyDelegationContract,
				UnclaimedRewards: "0",
				UndelegatedValue: undelegated.String(),
				Value: big.NewInt(0).SetBytes(active).String(),
			}},
			Total: big.NewInt(0).SetBytes(active).String(),
			UnclaimedTotal: "0",
			UndelegatedTotal: undelegated.String(),
			WaitingTotal: big.NewInt(0).SetBytes(waiting).String(),
		}

		delegationList.Data.List = append(delegationList.Data.List, delegationItem)
	}

	return delegationList, nil
}

func (nsp *NodeStatusProcessor) getLegacyUserStakeValues(userAddress string) ([][]byte, error) {
	query := &data.VmValueRequest{
		Address: legacyDelegationContract,
		FuncName: "getUserStakeByType",
		CallerAddr: legacyDelegationContract,
		CallValue: "0",
		Args: []string{userAddress},
	}

	observers, err := nsp.proc.GetObservers(2)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		response := &data.ResponseVmValue{}

		httpStatus, err := nsp.proc.CallPostRestEndPoint(observer.Address, SCQueryServicePath, query, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Info("SC query sent successfully, received response", "observer", observer.Address, "shard", 2)
			if len(response.Data.Data.ReturnData) != 5 {
				log.Error("legacy delegation waiting list", "invalid response data length", len(response.Data.Data.ReturnData))
				return nil, errors.New("invalid response data length received from legacy delegation get user stake")
			}

			return response.Data.Data.ReturnData, nil
		}

		if responseHasExplicitError {
			return nil, fmt.Errorf(response.Error)
		}

		return nil, err
	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getLegacyUserAddressByIndex(index *big.Int) (string, error) {

	indexString := index.Text(16)
	if len(indexString) % 2 != 0 {
		indexString = "0" + indexString
	}

	query := &data.VmValueRequest{
		Address: legacyDelegationContract,
		FuncName: "getUserAddress",
		CallerAddr: legacyDelegationContract,
		CallValue: "0",
		Args: []string{indexString},
	}

	observers, err := nsp.proc.GetObservers(2)
	if err != nil {
		return "", err
	}

	for _, observer := range observers {
		response := &data.ResponseVmValue{}

		httpStatus, err := nsp.proc.CallPostRestEndPoint(observer.Address, SCQueryServicePath, query, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Info("SC query sent successfully, received response", "observer", observer.Address, "shard", 2)
			if len(response.Data.Data.ReturnData) != 1 {
				log.Error("legacy delegation waiting list", "invalid response data length", len(response.Data.Data.ReturnData))
				return "", errors.New("invalid response data length received from legacy delegation get user address")
			}

			// Decode response
			return hex.EncodeToString(response.Data.Data.ReturnData[0]), nil
		}

		if responseHasExplicitError {
			return "", fmt.Errorf(response.Error)
		}

		return "", err
	}

	return "", ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getLegacyNumUsers() (int64, error) {
	query := &data.VmValueRequest{
		Address: legacyDelegationContract,
		FuncName: "getNumUsers",
		CallerAddr: legacyDelegationContract,
		CallValue: "0",
		Args: make([]string, 0),
	}

	observers, err := nsp.proc.GetObservers(2)
	if err != nil {
		return 0, err
	}

	for _, observer := range observers {
		response := &data.ResponseVmValue{}

		httpStatus, err := nsp.proc.CallPostRestEndPoint(observer.Address, SCQueryServicePath, query, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Info("SC query sent successfully, received response", "observer", observer.Address, "shard", 2)
			if len(response.Data.Data.ReturnData) != 1 {
				log.Error("legacy delegation waiting list", "invalid response data length", len(response.Data.Data.ReturnData))
				return 0, errors.New("invalid response data length received from legacy delegation num users")
			}

			// Decode response
			return big.NewInt(0).SetBytes(response.Data.Data.ReturnData[0]).Int64(), nil
		}

		if responseHasExplicitError {
			return 0, fmt.Errorf(response.Error)
		}

		return 0, err
	}

	return 0, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getLegacyDelegationWaitingList() ([]*data.Delegator, error) {
	query := &data.VmValueRequest{
		Address: legacyDelegationContract,
		FuncName: "getFullWaitingList",
		CallerAddr: legacyDelegationContract,
		CallValue: "0",
		Args: make([]string, 0),
	}

	observers, err := nsp.proc.GetObservers(2)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		response := &data.ResponseVmValue{}

		httpStatus, err := nsp.proc.CallPostRestEndPoint(observer.Address, SCQueryServicePath, query, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Info("SC query sent successfully, received response", "observer", observer.Address, "shard", 2)
			if len(response.Data.Data.ReturnData) % 3 != 0 {
				log.Error("legacy delegation waiting list", "invalid response data length", len(response.Data.Data.ReturnData))
				return nil, errors.New("invalid response data length received from legacy delegation waiting list")
			}

			delegationList := make([]*data.Delegator, 0)
			for i := 0; i < len(response.Data.Data.ReturnData); i+=3 {
				addressBytes := response.Data.Data.ReturnData[i]
				amountBytes := response.Data.Data.ReturnData[i+1]
				// We don't care about the nonce from i+2

				bechAddress := nsp.pubKeyConverter.Encode(addressBytes)
				if len(bechAddress) == 0 {
					log.Error("legacy delegation waiting list", "could not decode delegator address", string(addressBytes))
					return nil, errors.New("something went wrong decoding delegator's address")
				}

				amountString := big.NewInt(0).SetBytes(amountBytes).String()
				delegationList = append(delegationList, &data.Delegator{
					DelegatorAddress: bechAddress,
					DelegatedTo: []*data.DelegationItem{{
						DelegationScAddress: legacyDelegationContract,
						UnclaimedRewards: "0",
						UndelegatedValue: "0",
						Value: amountString,
					}},
					Total: "0",
					UnclaimedTotal: "0",
					UndelegatedTotal: "0",
					WaitingTotal: amountString,
				})
			}
			// Decode response
			return delegationList, nil
		}

		if responseHasExplicitError {
			return nil, fmt.Errorf(response.Error)
		}

		return nil, err
	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getLegacyDelegationStakingList() ([]*data.Delegator, error) {
	query := &data.VmValueRequest{
		Address: legacyDelegationContract,
		FuncName: "getFullActiveList",
		CallerAddr: legacyDelegationContract,
		CallValue: "0",
		Args: make([]string, 0),
	}

	observers, err := nsp.proc.GetObservers(2)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		response := &data.ResponseVmValue{}

		httpStatus, err := nsp.proc.CallPostRestEndPoint(observer.Address, SCQueryServicePath, query, response)
		isObserverDown := httpStatus == http.StatusNotFound || httpStatus == http.StatusRequestTimeout
		isOk := httpStatus == http.StatusOK
		responseHasExplicitError := len(response.Error) > 0

		if isObserverDown {
			log.LogIfError(err)
			continue
		}

		if isOk {
			log.Info("SC query sent successfully, received response", "observer", observer.Address, "shard", 2)
			if len(response.Data.Data.ReturnData) % 2 != 0 {
				log.Error("legacy delegation active list", "invalid response data length", len(response.Data.Data.ReturnData))
				return nil, errors.New("invalid response data length received from legacy delegation active list")
			}

			delegationList := make([]*data.Delegator, 0)
			for i := 0; i < len(response.Data.Data.ReturnData); i+=2 {
				addressBytes := response.Data.Data.ReturnData[i]
				amountBytes := response.Data.Data.ReturnData[i+1]

				bechAddress := nsp.pubKeyConverter.Encode(addressBytes)
				if len(bechAddress) == 0 {
					log.Error("legacy delegation active list", "could not decode delegator address", string(addressBytes))
					return nil, errors.New("something went wrong decoding delegator's address")
				}

				amountString := big.NewInt(0).SetBytes(amountBytes).String()
				delegationList = append(delegationList, &data.Delegator{
					DelegatorAddress: bechAddress,
					DelegatedTo: []*data.DelegationItem{{
						DelegationScAddress: legacyDelegationContract,
						UnclaimedRewards: "0",
						UndelegatedValue: "0",
						Value: amountString,
					}},
					Total: amountString,
					UnclaimedTotal: "0",
					UndelegatedTotal: "0",
					WaitingTotal: "0",
				})
			}
			// Decode response
			return delegationList, nil
		}

		if responseHasExplicitError {
			return nil, fmt.Errorf(response.Error)
		}

		return nil, err
	}

	return nil, ErrSendingRequest
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

func getExceptions() map[string]bool {
	return map[string]bool {
		"erd18s4cfunrctf27ejp3jmvylff7psfdgdssgc7e5aal6yusac62xzqly0yh5": true,
		"erd1a56dkgcpwwx6grmcvw9w5vpf9zeq53w3w7n6dmxcpxjry3l7uh2s3h9dtr": true,
		"erd1qt827an62lztf74rx7cg2s6utx3dp6l8k9snlttd77zny4dlzr9qccdqgx": true,
		"erd1rm8pg3yrngzyhrjejkz3xq2lfp64mvnt64llj3fyft53d3t4ckjq0q8v4k": true,
		"erd1043dp0s3yw8vd44s5xvxklnp30ypp7y56mylm9t87vdhhgwcx24s2e2g5y": true,
		"erd1z27mr0ertnan43avl4uhrud67awtkqklfsxzpetkp5u26cscsrzqdl56j8": true,
		"erd1jfempey50xue4wa5hzwmle4p4y6g55dn4327m9pvrynttdscn2eqvaxcgy": true,
		"erd16x7le8dpkjsafgwjx0e5kw94evsqw039rwp42m2j9eesd88x8zzs75tzry": true,
		"erd1rf4hv70arudgzus0ymnnsnc4pml0jkywg2xjvzslg0mz4nn2tg7q7k0t6p": true,
		"erd18umqd6v045nww2g9kgneupj4dwme9lycrpjn293sfkrhpntx9z2ss4kvhg": true,
		"erd1tqun7ku6yrygd0gjezmmz42jffqzlhgtvl2tsch3cel7rfylwzxs2dhrcg": true,
		"erd1v4ms58e22zjcp08suzqgm9ajmumwxcy4hfkdc23gvynnegjdflmsj6gmaq": true,
		"erd15qltd5ccalm5smmgdc5wnx46ssda3p32xhsz4wpp6usldq7hq7xqq5fmn6": true,
		"erd1qr9av6ar4ymr05xj93jzdxyezdrp6r4hz6u0scz4dtzvv7kmlldse7zktc": true,
		"erd1vup7q384decm8l8mu4ehz75c5mfs089nd32fteru95tm8d0a8dqs8g0yst": true,
		"erd195fe57d7fm5h33585sc7wl8trqhrmy85z3dg6f6mqd0724ymljxq3zjemc": true,
		"erd16rp9ur5crj6lcjyttr0ftft8vspmcgq3kk00wmzkv6p7lnqg3v8quhqhh5": true,
		"erd1d64xqa84x52qgl0v476zdc4fzmdqdcv6qvhsvj8hzmh5p883qfasdkus4p": true,
		"erd1hcaps2cq6v3j2ke8ldnjxmnacuk6hhgwspuxpv3whpnnq4ldlxdq5cukqz": true,
		"erd1mxu8utup9kgwq6adrad7c70he2ujn7df22n4f8qen06rvws8dacsrcxfhw": true,
		"erd1eaua3avcw0ax99ncmtt0hacfn20jy2azp8j6ay2447nza83msmase863gk": true,
		"erd1039c78xauyf74xvmcvwumpu5lctcnwwa5a2wv287nx2834vuqj9qvqnjck": true,
		"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l": true,
		"erd1qqqqqqqqqqqqqpgqxwakt2g7u9atsnr03gqcgmhcv38pt7mkd94q6shuwt": true,
	}
}
