package facade

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/api/address"
	"github.com/ElrondNetwork/elrond-proxy-go/api/block"
	"github.com/ElrondNetwork/elrond-proxy-go/api/blockatlas"
	"github.com/ElrondNetwork/elrond-proxy-go/api/hyperblock"
	"github.com/ElrondNetwork/elrond-proxy-go/api/network"
	"github.com/ElrondNetwork/elrond-proxy-go/api/node"
	"github.com/ElrondNetwork/elrond-proxy-go/api/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/api/validator"
	"github.com/ElrondNetwork/elrond-proxy-go/api/vmValues"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// interfaces assertions. verifies that all API endpoint have their corresponding methods in the facade
var _ address.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ block.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ blockatlas.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ hyperblock.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ network.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ node.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ transaction.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ validator.FacadeHandler = (*ElrondProxyFacade)(nil)
var _ vmValues.FacadeHandler = (*ElrondProxyFacade)(nil)

// ElrondProxyFacade implements the facade used in api calls
type ElrondProxyFacade struct {
	accountProc    AccountProcessor
	txProc         TransactionProcessor
	scQueryService SCQueryService
	heartbeatProc  HeartbeatProcessor
	valStatsProc   ValidatorStatisticsProcessor
	faucetProc     FaucetProcessor
	nodeStatusProc NodeStatusProcessor
	blockProc      BlockProcessor

	pubKeyConverter core.PubkeyConverter
}

// NewElrondProxyFacade creates a new ElrondProxyFacade instance
func NewElrondProxyFacade(
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	scQueryService SCQueryService,
	heartbeatProc HeartbeatProcessor,
	valStatsProc ValidatorStatisticsProcessor,
	faucetProc FaucetProcessor,
	nodeStatusProc NodeStatusProcessor,
	blockProc BlockProcessor,
	pubKeyConverter core.PubkeyConverter,
) (*ElrondProxyFacade, error) {

	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}
	if scQueryService == nil {
		return nil, ErrNilSCQueryService
	}
	if heartbeatProc == nil {
		return nil, ErrNilHeartbeatProcessor
	}
	if valStatsProc == nil {
		return nil, ErrNilValidatorStatisticsProcessor
	}
	if faucetProc == nil {
		return nil, ErrNilFaucetProcessor
	}
	if nodeStatusProc == nil {
		return nil, ErrNilNodeStatusProcessor
	}
	if blockProc == nil {
		return nil, ErrNilBlockProcessor
	}

	return &ElrondProxyFacade{
		accountProc:     accountProc,
		txProc:          txProc,
		scQueryService:  scQueryService,
		heartbeatProc:   heartbeatProc,
		valStatsProc:    valStatsProc,
		faucetProc:      faucetProc,
		nodeStatusProc:  nodeStatusProc,
		blockProc:       blockProc,
		pubKeyConverter: pubKeyConverter,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ElrondProxyFacade) GetAccount(address string) (*data.Account, error) {
	return epf.accountProc.GetAccount(address)
}

// GetValueForKey returns the value for the given address and key
func (epf *ElrondProxyFacade) GetValueForKey(address string, key string) (string, error) {
	return epf.accountProc.GetValueForKey(address, key)
}

// GetShardIDForAddress returns the computed shard ID for the given address based on the current proxy's configuration
func (epf *ElrondProxyFacade) GetShardIDForAddress(address string) (uint32, error) {
	return epf.accountProc.GetShardIDForAddress(address)
}

// GetTransactions returns transactions by address
func (epf *ElrondProxyFacade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return epf.accountProc.GetTransactions(address)
}

// SendTransaction should send the transaction to the correct observer
func (epf *ElrondProxyFacade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return epf.txProc.SendTransaction(tx)
}

// SendMultipleTransactions should send the transactions to the correct observers
func (epf *ElrondProxyFacade) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	return epf.txProc.SendMultipleTransactions(txs)
}

// SimulateTransaction should send the transaction to the correct observer for simulation
func (epf *ElrondProxyFacade) SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, error) {
	return epf.txProc.SimulateTransaction(tx)
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (epf *ElrondProxyFacade) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return epf.txProc.TransactionCostRequest(tx)
}

// GetTransactionStatus should return transaction status
func (epf *ElrondProxyFacade) GetTransactionStatus(txHash string, sender string) (string, error) {
	return epf.txProc.GetTransactionStatus(txHash, sender)
}

// GetTransaction should return a transaction by hash
func (epf *ElrondProxyFacade) GetTransaction(txHash string, withResults bool) (*data.FullTransaction, error) {
	return epf.txProc.GetTransaction(txHash, withResults)
}

// GetTransactionByHashAndSenderAddress should return a transaction by hash and sender address
func (epf *ElrondProxyFacade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error) {
	return epf.txProc.GetTransactionByHashAndSenderAddress(txHash, sndAddr, withEvents)
}

type networkConfig struct {
	chainID               string
	minTransactionVersion uint32
}

// IsFaucetEnabled returns true if the faucet mechanism is enabled or false otherwise
func (epf *ElrondProxyFacade) IsFaucetEnabled() bool {
	return epf.faucetProc.IsEnabled()
}

// SendUserFunds should send a transaction to load one user's account with extra funds from an account in the pem file
func (epf *ElrondProxyFacade) SendUserFunds(receiver string, value *big.Int) error {
	senderSk, senderPk, err := epf.faucetProc.SenderDetailsFromPem(receiver)
	if err != nil {
		return err
	}

	senderAccount, err := epf.accountProc.GetAccount(senderPk)
	if err != nil {
		return err
	}

	networkConfig, err := epf.getNetworkConfig()
	if err != nil {
		return err
	}

	tx, err := epf.faucetProc.GenerateTxForSendUserFunds(
		senderSk,
		senderPk,
		senderAccount.Nonce,
		receiver,
		value,
		networkConfig.chainID,
		networkConfig.minTransactionVersion,
	)
	if err != nil {
		return err
	}

	_, _, err = epf.txProc.SendTransaction(tx)
	return err
}

func (epf *ElrondProxyFacade) getNetworkConfig() (*networkConfig, error) {
	netConfig, err := epf.nodeStatusProc.GetNetworkConfigMetrics()
	if err != nil {
		return nil, err
	}

	netConf, ok := netConfig.Data.(map[string]interface{})["config"].(map[string]interface{})
	if !ok {
		return nil, errors.New("cannot get network config. something went wrong")
	}

	chainID, ok := netConf[core.MetricChainId].(string)
	if !ok {
		return nil, errors.New("cannot get chainID. something went wrong")
	}

	version, ok := netConf[core.MetricMinTransactionVersion].(float64)
	if !ok {
		return nil, errors.New("cannot get version. something went wrong")
	}

	return &networkConfig{
		chainID:               chainID,
		minTransactionVersion: uint32(version),
	}, nil
}

// ExecuteSCQuery retrieves data from existing SC trie through the use of a VM
func (epf *ElrondProxyFacade) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, error) {
	return epf.scQueryService.ExecuteQuery(query)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (epf *ElrondProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return epf.heartbeatProc.GetHeartbeatData()
}

// GetNetworkConfigMetrics retrieves the node's configuration's metrics
func (epf *ElrondProxyFacade) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetNetworkConfigMetrics()
}

// GetNetworkStatusMetrics retrieves the node's network metrics for a given shard
func (epf *ElrondProxyFacade) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetNetworkStatusMetrics(shardID)
}

// GetNetworkStatusMetrics retrieves the node's network metrics for a given shard
func (epf *ElrondProxyFacade) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetEconomicsDataMetrics()
}

// GetBlockByHash retrieves the block by hash for a given shard
func (epf *ElrondProxyFacade) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	return epf.blockProc.GetBlockByHash(shardID, hash, withTxs)
}

// GetBlockByNonce retrieves the block by nonce for a given shard
func (epf *ElrondProxyFacade) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	return epf.blockProc.GetBlockByNonce(shardID, nonce, withTxs)
}

// GetHyperBlockByHash retrieves the hyperblock by hash
func (epf *ElrondProxyFacade) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	return epf.blockProc.GetHyperBlockByHash(hash)
}

// GetHyperBlockByNonce retrieves the block by nonce
func (epf *ElrondProxyFacade) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	return epf.blockProc.GetHyperBlockByNonce(nonce)
}

// ValidatorStatistics will return the statistics from an observer
func (epf *ElrondProxyFacade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	valStats, err := epf.valStatsProc.GetValidatorStatistics()
	if err != nil {
		return nil, err
	}

	return valStats.Statistics, nil
}

// GetAtlasBlockByShardIDAndNonce returns block by shardID and nonce in a BlockAtlas-friendly-format
func (epf *ElrondProxyFacade) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return epf.blockProc.GetAtlasBlockByShardIDAndNonce(shardID, nonce)
}

// GetAddressConverter returns the address converter
func (epf *ElrondProxyFacade) GetAddressConverter() (core.PubkeyConverter, error) {
	return epf.pubKeyConverter, nil
}

// GetLatestFullySynchronizedHyperblockNonce returns the latest fully synchronized hyperblock nonce
func (epf *ElrondProxyFacade) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	return epf.nodeStatusProc.GetLatestFullySynchronizedHyperblockNonce()
}

// ComputeTransactionHash will compute hash of a given transaction
func (epf *ElrondProxyFacade) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	return epf.txProc.ComputeTransactionHash(tx)
}
