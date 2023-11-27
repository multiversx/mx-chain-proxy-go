package facade

import (
	"encoding/json"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// interfaces assertions. verifies that all API endpoint have their corresponding methods in the facade
var _ groups.ActionsFacadeHandler = (*ProxyFacade)(nil)
var _ groups.AccountsFacadeHandler = (*ProxyFacade)(nil)
var _ groups.BlockFacadeHandler = (*ProxyFacade)(nil)
var _ groups.BlocksFacadeHandler = (*ProxyFacade)(nil)
var _ groups.BlockAtlasFacadeHandler = (*ProxyFacade)(nil)
var _ groups.HyperBlockFacadeHandler = (*ProxyFacade)(nil)
var _ groups.NetworkFacadeHandler = (*ProxyFacade)(nil)
var _ groups.NodeFacadeHandler = (*ProxyFacade)(nil)
var _ groups.TransactionFacadeHandler = (*ProxyFacade)(nil)
var _ groups.ValidatorFacadeHandler = (*ProxyFacade)(nil)
var _ groups.VmValuesFacadeHandler = (*ProxyFacade)(nil)
var _ groups.ProofFacadeHandler = (*ProxyFacade)(nil)

// ProxyFacade implements the facade used in api calls
type ProxyFacade struct {
	actionsProc      ActionsProcessor
	accountProc      AccountProcessor
	txProc           TransactionProcessor
	scQueryService   SCQueryService
	nodeGroupProc    NodeGroupProcessor
	valStatsProc     ValidatorStatisticsProcessor
	faucetProc       FaucetProcessor
	nodeStatusProc   NodeStatusProcessor
	blockProc        BlockProcessor
	blocksProc       BlocksProcessor
	proofProc        ProofProcessor
	esdtSuppliesProc ESDTSupplyProcessor
	statusProc       StatusProcessor

	pubKeyConverter core.PubkeyConverter
	aboutInfoProc   AboutInfoProcessor
}

// NewProxyFacade creates a new ProxyFacade instance
func NewProxyFacade(
	actionsProc ActionsProcessor,
	accountProc AccountProcessor,
	txProc TransactionProcessor,
	scQueryService SCQueryService,
	nodeGroupProc NodeGroupProcessor,
	valStatsProc ValidatorStatisticsProcessor,
	faucetProc FaucetProcessor,
	nodeStatusProc NodeStatusProcessor,
	blockProc BlockProcessor,
	blocksProc BlocksProcessor,
	proofProc ProofProcessor,
	pubKeyConverter core.PubkeyConverter,
	esdtSuppliesProc ESDTSupplyProcessor,
	statusProc StatusProcessor,
	aboutInfoProc AboutInfoProcessor,
) (*ProxyFacade, error) {
	if actionsProc == nil {
		return nil, ErrNilActionsProcessor
	}
	if accountProc == nil {
		return nil, ErrNilAccountProcessor
	}
	if txProc == nil {
		return nil, ErrNilTransactionProcessor
	}
	if scQueryService == nil {
		return nil, ErrNilSCQueryService
	}
	if nodeGroupProc == nil {
		return nil, ErrNilNodeGroupProcessor
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
	if blocksProc == nil {
		return nil, ErrNilBlocksProcessor
	}
	if proofProc == nil {
		return nil, ErrNilProofProcessor
	}
	if esdtSuppliesProc == nil {
		return nil, ErrNilESDTSuppliesProcessor
	}
	if statusProc == nil {
		return nil, ErrNilStatusProcessor
	}
	if aboutInfoProc == nil {
		return nil, ErrNilAboutInfoProcessor
	}

	return &ProxyFacade{
		actionsProc:      actionsProc,
		accountProc:      accountProc,
		txProc:           txProc,
		scQueryService:   scQueryService,
		nodeGroupProc:    nodeGroupProc,
		valStatsProc:     valStatsProc,
		faucetProc:       faucetProc,
		nodeStatusProc:   nodeStatusProc,
		blockProc:        blockProc,
		blocksProc:       blocksProc,
		proofProc:        proofProc,
		pubKeyConverter:  pubKeyConverter,
		esdtSuppliesProc: esdtSuppliesProc,
		statusProc:       statusProc,
		aboutInfoProc:    aboutInfoProc,
	}, nil
}

// GetAccount returns an account based on the input address
func (epf *ProxyFacade) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	return epf.accountProc.GetAccount(address, options)
}

// GetCodeHash returns the code hash for the given address
func (epf *ProxyFacade) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetCodeHash(address, options)
}

// GetKeyValuePairs returns the key-value pairs for the given address
func (epf *ProxyFacade) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetKeyValuePairs(address, options)
}

// GetValueForKey returns the value for the given address and key
func (epf *ProxyFacade) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	return epf.accountProc.GetValueForKey(address, key, options)
}

// GetGuardianData returns the guardian data for the given address
func (epf *ProxyFacade) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetGuardianData(address, options)
}

// GetShardIDForAddress returns the computed shard ID for the given address based on the current proxy's configuration
func (epf *ProxyFacade) GetShardIDForAddress(address string) (uint32, error) {
	return epf.accountProc.GetShardIDForAddress(address)
}

// GetTransactions returns transactions by address
func (epf *ProxyFacade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return epf.accountProc.GetTransactions(address)
}

// GetESDTTokenData returns the token data for a given token name
func (epf *ProxyFacade) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetESDTTokenData(address, key, options)
}

// GetESDTNftTokenData returns the token data for a given token name
func (epf *ProxyFacade) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetESDTNftTokenData(address, key, nonce, options)
}

// GetESDTsWithRole returns the tokens where the given address has the assigned role
func (epf *ProxyFacade) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetESDTsWithRole(address, role, options)
}

// GetESDTsRoles returns the tokens and roles for the given address
func (epf *ProxyFacade) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetESDTsRoles(address, options)
}

// GetNFTTokenIDsRegisteredByAddress returns the token identifiers of the NFTs registered by the address
func (epf *ProxyFacade) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetNFTTokenIDsRegisteredByAddress(address, options)
}

// GetAllESDTTokens returns all the ESDT tokens for a given address
func (epf *ProxyFacade) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.GetAllESDTTokens(address, options)
}

// SendTransaction should send the transaction to the correct observer
func (epf *ProxyFacade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return epf.txProc.SendTransaction(tx)
}

// SendMultipleTransactions should send the transactions to the correct observers
func (epf *ProxyFacade) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	return epf.txProc.SendMultipleTransactions(txs)
}

// SimulateTransaction should send the transaction to the correct observer for simulation
func (epf *ProxyFacade) SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
	return epf.txProc.SimulateTransaction(tx, checkSignature)
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (epf *ProxyFacade) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	return epf.txProc.TransactionCostRequest(tx)
}

// GetTransactionStatus should return transaction status
func (epf *ProxyFacade) GetTransactionStatus(txHash string, sender string) (string, error) {
	return epf.txProc.GetTransactionStatus(txHash, sender)
}

// GetProcessedTransactionStatus should return transaction status after internal processing of the transaction results
func (epf *ProxyFacade) GetProcessedTransactionStatus(txHash string) (string, error) {
	return epf.txProc.GetProcessedTransactionStatus(txHash)
}

// GetTransaction should return a transaction by hash
func (epf *ProxyFacade) GetTransaction(txHash string, withResults bool) (*transaction.ApiTransactionResult, error) {
	return epf.txProc.GetTransaction(txHash, withResults)
}

// ReloadObservers will try to reload the observers
func (epf *ProxyFacade) ReloadObservers() data.NodesReloadResponse {
	return epf.actionsProc.ReloadObservers()
}

// ReloadFullHistoryObservers will try to reload the full history observers
func (epf *ProxyFacade) ReloadFullHistoryObservers() data.NodesReloadResponse {
	return epf.actionsProc.ReloadFullHistoryObservers()
}

// GetTransactionByHashAndSenderAddress should return a transaction by hash and sender address
func (epf *ProxyFacade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error) {
	return epf.txProc.GetTransactionByHashAndSenderAddress(txHash, sndAddr, withEvents)
}

// IsFaucetEnabled returns true if the faucet mechanism is enabled or false otherwise
func (epf *ProxyFacade) IsFaucetEnabled() bool {
	return epf.faucetProc.IsEnabled()
}

// SendUserFunds should send a transaction to load one user's account with extra funds from an account in the pem file
func (epf *ProxyFacade) SendUserFunds(receiver string, value *big.Int) error {
	senderSk, senderPk, err := epf.faucetProc.SenderDetailsFromPem(receiver)
	if err != nil {
		return err
	}

	senderAccount, err := epf.accountProc.GetAccount(senderPk, common.AccountQueryOptions{})
	if err != nil {
		return err
	}

	networkCfg, err := epf.getNetworkConfig()
	if err != nil {
		return err
	}

	tx, err := epf.faucetProc.GenerateTxForSendUserFunds(
		senderSk,
		senderPk,
		senderAccount.Account.Nonce,
		receiver,
		value,
		networkCfg,
	)
	if err != nil {
		return err
	}

	_, _, err = epf.txProc.SendTransaction(tx)
	return err
}

func (epf *ProxyFacade) getNetworkConfig() (*data.NetworkConfig, error) {
	genericResponse, err := epf.nodeStatusProc.GetNetworkConfigMetrics()
	if err != nil {
		return nil, err
	}

	networkConfigBytes, err := json.Marshal(&genericResponse.Data)
	if err != nil {
		return nil, err
	}

	networkCfg := &data.NetworkConfig{}
	err = json.Unmarshal(networkConfigBytes, networkCfg)

	return networkCfg, err
}

// ExecuteSCQuery retrieves data from existing SC trie through the use of a VM
func (epf *ProxyFacade) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	return epf.scQueryService.ExecuteQuery(query)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (epf *ProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return epf.nodeGroupProc.GetHeartbeatData()
}

// GetNetworkConfigMetrics retrieves the node's configuration's metrics
func (epf *ProxyFacade) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetNetworkConfigMetrics()
}

// GetNetworkStatusMetrics retrieves the node's network metrics for a given shard
func (epf *ProxyFacade) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetNetworkStatusMetrics(shardID)
}

// GetESDTSupply retrieves the supply for the provided token
func (epf *ProxyFacade) GetESDTSupply(token string) (*data.ESDTSupplyResponse, error) {
	return epf.esdtSuppliesProc.GetESDTSupply(token)
}

// GetEconomicsDataMetrics retrieves the node's network metrics for a given shard
func (epf *ProxyFacade) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetEconomicsDataMetrics()
}

// GetDelegatedInfo retrieves the node's network delegated info
func (epf *ProxyFacade) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetDelegatedInfo()
}

// GetDirectStakedInfo retrieves the node's direct staked values
func (epf *ProxyFacade) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetDirectStakedInfo()
}

// GetAllIssuedESDTs retrieves all the issued ESDTs from the node
func (epf *ProxyFacade) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetAllIssuedESDTs(tokenType)
}

// GetEnableEpochsMetrics retrieves the activation epochs
func (epf *ProxyFacade) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetEnableEpochsMetrics()
}

// GetRatingsConfig retrieves the node's configuration's metrics
func (epf *ProxyFacade) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetRatingsConfig()
}

// GetBlockByHash retrieves the block by hash for a given shard
func (epf *ProxyFacade) GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return epf.blockProc.GetBlockByHash(shardID, hash, options)
}

// GetBlockByNonce retrieves the block by nonce for a given shard
func (epf *ProxyFacade) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return epf.blockProc.GetBlockByNonce(shardID, nonce, options)
}

// GetBlocksByRound retrieves the blocks for a given round
func (epf *ProxyFacade) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	return epf.blocksProc.GetBlocksByRound(round, options)
}

// GetInternalBlockByHash retrieves the internal block by hash for a given shard
func (epf *ProxyFacade) GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return epf.blockProc.GetInternalBlockByHash(shardID, hash, format)
}

// GetInternalBlockByNonce retrieves the internal block by nonce for a given shard
func (epf *ProxyFacade) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return epf.blockProc.GetInternalBlockByNonce(shardID, nonce, format)
}

// GetInternalStartOfEpochMetaBlock retrieves the internal block by nonce for a given shard
func (epf *ProxyFacade) GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return epf.blockProc.GetInternalStartOfEpochMetaBlock(epoch, format)
}

// GetInternalMiniBlockByHash retrieves the internal miniblock by hash for a given shard
func (epf *ProxyFacade) GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
	return epf.blockProc.GetInternalMiniBlockByHash(shardID, hash, epoch, format)
}

// GetHyperBlockByHash retrieves the hyperblock by hash
func (epf *ProxyFacade) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return epf.blockProc.GetHyperBlockByHash(hash, options)
}

// GetHyperBlockByNonce retrieves the block by nonce
func (epf *ProxyFacade) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return epf.blockProc.GetHyperBlockByNonce(nonce, options)
}

// ValidatorStatistics will return the statistics from an observer
func (epf *ProxyFacade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	valStats, err := epf.valStatsProc.GetValidatorStatistics()
	if err != nil {
		return nil, err
	}

	return valStats.Statistics, nil
}

// GetAtlasBlockByShardIDAndNonce returns block by shardID and nonce in a BlockAtlas-friendly-format
func (epf *ProxyFacade) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return epf.blockProc.GetAtlasBlockByShardIDAndNonce(shardID, nonce)
}

// GetAddressConverter returns the address converter
func (epf *ProxyFacade) GetAddressConverter() (core.PubkeyConverter, error) {
	return epf.pubKeyConverter, nil
}

// GetLatestFullySynchronizedHyperblockNonce returns the latest fully synchronized hyperblock nonce
func (epf *ProxyFacade) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	return epf.nodeStatusProc.GetLatestFullySynchronizedHyperblockNonce()
}

// ComputeTransactionHash will compute hash of a given transaction
func (epf *ProxyFacade) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	return epf.txProc.ComputeTransactionHash(tx)
}

// GetTransactionsPool returns all txs from pool
func (epf *ProxyFacade) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	return epf.txProc.GetTransactionsPool(fields)
}

// GetTransactionsPoolForShard returns all txs from shard's pool
func (epf *ProxyFacade) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	return epf.txProc.GetTransactionsPoolForShard(shardID, fields)
}

// GetTransactionsPoolForSender returns tx pool for sender
func (epf *ProxyFacade) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	return epf.txProc.GetTransactionsPoolForSender(sender, fields)
}

// GetLastPoolNonceForSender returns last nonce from tx pool for sender
func (epf *ProxyFacade) GetLastPoolNonceForSender(sender string) (uint64, error) {
	return epf.txProc.GetLastPoolNonceForSender(sender)
}

// IsOldStorageForToken returns true is the storage for a given token is old
func (epf *ProxyFacade) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	return epf.nodeGroupProc.IsOldStorageForToken(tokenID, nonce)
}

// GetTransactionsPoolNonceGapsForSender returns all nonce gaps from tx pool for sender
func (epf *ProxyFacade) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	return epf.txProc.GetTransactionsPoolNonceGapsForSender(sender)
}

// GetProof returns the Merkle proof for the given address
func (epf *ProxyFacade) GetProof(rootHash string, address string) (*data.GenericAPIResponse, error) {
	return epf.proofProc.GetProof(rootHash, address)
}

// GetProofDataTrie returns a Merkle proof for the given address and a Merkle proof for the given key
func (epf *ProxyFacade) GetProofDataTrie(rootHash string, address string, key string) (*data.GenericAPIResponse, error) {
	return epf.proofProc.GetProofDataTrie(rootHash, address, key)
}

// GetProofCurrentRootHash returns the Merkle proof for the given address
func (epf *ProxyFacade) GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error) {
	return epf.proofProc.GetProofCurrentRootHash(address)
}

// VerifyProof verifies the given Merkle proof
func (epf *ProxyFacade) VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error) {
	return epf.proofProc.VerifyProof(rootHash, address, proof)
}

// GetMetrics will return the status metrics
func (epf *ProxyFacade) GetMetrics() map[string]*data.EndpointMetrics {
	return epf.statusProc.GetMetrics()
}

// GetMetricsForPrometheus will return the status metrics in a prometheus format
func (epf *ProxyFacade) GetMetricsForPrometheus() string {
	return epf.statusProc.GetMetricsForPrometheus()
}

// GetGenesisNodesPubKeys retrieves the node's configuration public keys
func (epf *ProxyFacade) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetGenesisNodesPubKeys()
}

// GetGasConfigs retrieves the current gas schedule configs
func (epf *ProxyFacade) GetGasConfigs() (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetGasConfigs()
}

// GetAboutInfo will return the app info
func (epf *ProxyFacade) GetAboutInfo() (*data.GenericAPIResponse, error) {
	return epf.aboutInfoProc.GetAboutInfo(), nil
}

// GetNodesVersions will return the version of the nodes
func (epf *ProxyFacade) GetNodesVersions() (*data.GenericAPIResponse, error) {
	return epf.aboutInfoProc.GetNodesVersions()
}

// GetAlteredAccountsByNonce returns altered accounts by nonce in block
func (epf *ProxyFacade) GetAlteredAccountsByNonce(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	return epf.blockProc.GetAlteredAccountsByNonce(shardID, nonce, options)
}

// GetAlteredAccountsByHash returns altered accounts by hash in block
func (epf *ProxyFacade) GetAlteredAccountsByHash(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	return epf.blockProc.GetAlteredAccountsByHash(shardID, hash, options)
}

// GetTriesStatistics will return trie statistics
func (epf *ProxyFacade) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	return epf.nodeStatusProc.GetTriesStatistics(shardID)
}

// GetEpochStartData retrieves epoch start data for the provides epoch and shard ID
func (epf *ProxyFacade) GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error) {
	return epf.nodeStatusProc.GetEpochStartData(epoch, shardID)
}

// GetInternalStartOfEpochValidatorsInfo retrieves the validators info by epoch
func (epf *ProxyFacade) GetInternalStartOfEpochValidatorsInfo(epoch uint32) (*data.ValidatorsInfoApiResponse, error) {
	return epf.blockProc.GetInternalStartOfEpochValidatorsInfo(epoch)
}

// GetWaitingEpochsLeftForPublicKey returns the number of epochs left for the public key until it becomes eligible
func (epf *ProxyFacade) GetWaitingEpochsLeftForPublicKey(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
	return epf.nodeGroupProc.GetWaitingEpochsLeftForPublicKey(publicKey)
}

// IsDataTrieMigrated returns true if the data trie for the given address is migrated
func (epf *ProxyFacade) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return epf.accountProc.IsDataTrieMigrated(address, options)
}
