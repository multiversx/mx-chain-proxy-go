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
func (pf *ProxyFacade) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	return pf.accountProc.GetAccount(address, options)
}

// GetCodeHash returns the code hash for the given address
func (pf *ProxyFacade) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetCodeHash(address, options)
}

// GetKeyValuePairs returns the key-value pairs for the given address
func (pf *ProxyFacade) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetKeyValuePairs(address, options)
}

// GetAccounts returns data about the provided addresses
func (pf *ProxyFacade) GetAccounts(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error) {
	return pf.accountProc.GetAccounts(addresses, options)
}

// GetValueForKey returns the value for the given address and key
func (pf *ProxyFacade) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	return pf.accountProc.GetValueForKey(address, key, options)
}

// GetGuardianData returns the guardian data for the given address
func (pf *ProxyFacade) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetGuardianData(address, options)
}

// GetShardIDForAddress returns the computed shard ID for the given address based on the current proxy's configuration
func (pf *ProxyFacade) GetShardIDForAddress(address string) (uint32, error) {
	return pf.accountProc.GetShardIDForAddress(address)
}

// GetESDTTokenData returns the token data for a given token name
func (pf *ProxyFacade) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetESDTTokenData(address, key, options)
}

// GetESDTNftTokenData returns the token data for a given token name
func (pf *ProxyFacade) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetESDTNftTokenData(address, key, nonce, options)
}

// GetESDTsWithRole returns the tokens where the given address has the assigned role
func (pf *ProxyFacade) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetESDTsWithRole(address, role, options)
}

// GetESDTsRoles returns the tokens and roles for the given address
func (pf *ProxyFacade) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetESDTsRoles(address, options)
}

// GetNFTTokenIDsRegisteredByAddress returns the token identifiers of the NFTs registered by the address
func (pf *ProxyFacade) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetNFTTokenIDsRegisteredByAddress(address, options)
}

// GetAllESDTTokens returns all the ESDT tokens for a given address
func (pf *ProxyFacade) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.GetAllESDTTokens(address, options)
}

// SendTransaction should send the transaction to the correct observer
func (pf *ProxyFacade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return pf.txProc.SendTransaction(tx)
}

// SendMultipleTransactions should send the transactions to the correct observers
func (pf *ProxyFacade) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	return pf.txProc.SendMultipleTransactions(txs)
}

// SimulateTransaction should send the transaction to the correct observer for simulation
func (pf *ProxyFacade) SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
	return pf.txProc.SimulateTransaction(tx, checkSignature)
}

// TransactionCostRequest should return how many gas units a transaction will cost
func (pf *ProxyFacade) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	return pf.txProc.TransactionCostRequest(tx)
}

// GetTransactionStatus should return transaction status
func (pf *ProxyFacade) GetTransactionStatus(txHash string, sender string) (string, error) {
	return pf.txProc.GetTransactionStatus(txHash, sender)
}

// GetProcessedTransactionStatus should return transaction status after internal processing of the transaction results
func (pf *ProxyFacade) GetProcessedTransactionStatus(txHash string) (*data.ProcessStatusResponse, error) {
	return pf.txProc.GetProcessedTransactionStatus(txHash)
}

// GetTransaction should return a transaction by hash
func (pf *ProxyFacade) GetTransaction(txHash string, withResults bool, relayedTxHash string) (*transaction.ApiTransactionResult, error) {
	return pf.txProc.GetTransaction(txHash, withResults, relayedTxHash)
}

// ReloadObservers will try to reload the observers
func (pf *ProxyFacade) ReloadObservers() data.NodesReloadResponse {
	return pf.actionsProc.ReloadObservers()
}

// ReloadFullHistoryObservers will try to reload the full history observers
func (pf *ProxyFacade) ReloadFullHistoryObservers() data.NodesReloadResponse {
	return pf.actionsProc.ReloadFullHistoryObservers()
}

// GetTransactionByHashAndSenderAddress should return a transaction by hash and sender address
func (pf *ProxyFacade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error) {
	return pf.txProc.GetTransactionByHashAndSenderAddress(txHash, sndAddr, withEvents)
}

// IsFaucetEnabled returns true if the faucet mechanism is enabled or false otherwise
func (pf *ProxyFacade) IsFaucetEnabled() bool {
	return pf.faucetProc.IsEnabled()
}

// SendUserFunds should send a transaction to load one user's account with extra funds from an account in the pem file
func (pf *ProxyFacade) SendUserFunds(receiver string, value *big.Int) error {
	senderSk, senderPk, err := pf.faucetProc.SenderDetailsFromPem(receiver)
	if err != nil {
		return err
	}

	senderAccount, err := pf.accountProc.GetAccount(senderPk, common.AccountQueryOptions{})
	if err != nil {
		return err
	}

	networkCfg, err := pf.getNetworkConfig()
	if err != nil {
		return err
	}

	tx, err := pf.faucetProc.GenerateTxForSendUserFunds(
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

	_, _, err = pf.txProc.SendTransaction(tx)
	return err
}

func (pf *ProxyFacade) getNetworkConfig() (*data.NetworkConfig, error) {
	genericResponse, err := pf.nodeStatusProc.GetNetworkConfigMetrics()
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
func (pf *ProxyFacade) ExecuteSCQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error) {
	return pf.scQueryService.ExecuteQuery(query)
}

// GetHeartbeatData retrieves the heartbeat status from one observer
func (pf *ProxyFacade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return pf.nodeGroupProc.GetHeartbeatData()
}

// GetNetworkConfigMetrics retrieves the node's configuration's metrics
func (pf *ProxyFacade) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetNetworkConfigMetrics()
}

// GetNetworkStatusMetrics retrieves the node's network metrics for a given shard
func (pf *ProxyFacade) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetNetworkStatusMetrics(shardID)
}

// GetESDTSupply retrieves the supply for the provided token
func (pf *ProxyFacade) GetESDTSupply(token string) (*data.ESDTSupplyResponse, error) {
	return pf.esdtSuppliesProc.GetESDTSupply(token)
}

// GetEconomicsDataMetrics retrieves the node's network metrics for a given shard
func (pf *ProxyFacade) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetEconomicsDataMetrics()
}

// GetDelegatedInfo retrieves the node's network delegated info
func (pf *ProxyFacade) GetDelegatedInfo() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetDelegatedInfo()
}

// GetDirectStakedInfo retrieves the node's direct staked values
func (pf *ProxyFacade) GetDirectStakedInfo() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetDirectStakedInfo()
}

// GetAllIssuedESDTs retrieves all the issued ESDTs from the node
func (pf *ProxyFacade) GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetAllIssuedESDTs(tokenType)
}

// GetEnableEpochsMetrics retrieves the activation epochs
func (pf *ProxyFacade) GetEnableEpochsMetrics() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetEnableEpochsMetrics()
}

// GetRatingsConfig retrieves the node's configuration's metrics
func (pf *ProxyFacade) GetRatingsConfig() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetRatingsConfig()
}

// GetBlockByHash retrieves the block by hash for a given shard
func (pf *ProxyFacade) GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return pf.blockProc.GetBlockByHash(shardID, hash, options)
}

// GetBlockByNonce retrieves the block by nonce for a given shard
func (pf *ProxyFacade) GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error) {
	return pf.blockProc.GetBlockByNonce(shardID, nonce, options)
}

// GetBlocksByRound retrieves the blocks for a given round
func (pf *ProxyFacade) GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error) {
	return pf.blocksProc.GetBlocksByRound(round, options)
}

// GetInternalBlockByHash retrieves the internal block by hash for a given shard
func (pf *ProxyFacade) GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return pf.blockProc.GetInternalBlockByHash(shardID, hash, format)
}

// GetInternalBlockByNonce retrieves the internal block by nonce for a given shard
func (pf *ProxyFacade) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return pf.blockProc.GetInternalBlockByNonce(shardID, nonce, format)
}

// GetInternalStartOfEpochMetaBlock retrieves the internal block by nonce for a given shard
func (pf *ProxyFacade) GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error) {
	return pf.blockProc.GetInternalStartOfEpochMetaBlock(epoch, format)
}

// GetInternalMiniBlockByHash retrieves the internal miniblock by hash for a given shard
func (pf *ProxyFacade) GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error) {
	return pf.blockProc.GetInternalMiniBlockByHash(shardID, hash, epoch, format)
}

// GetHyperBlockByHash retrieves the hyperblock by hash
func (pf *ProxyFacade) GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return pf.blockProc.GetHyperBlockByHash(hash, options)
}

// GetHyperBlockByNonce retrieves the block by nonce
func (pf *ProxyFacade) GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error) {
	return pf.blockProc.GetHyperBlockByNonce(nonce, options)
}

// ValidatorStatistics will return the statistics from an observer
func (pf *ProxyFacade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	valStats, err := pf.valStatsProc.GetValidatorStatistics()
	if err != nil {
		return nil, err
	}

	return valStats.Statistics, nil
}

// AuctionList will return the auction list
func (epf *ProxyFacade) AuctionList() ([]*data.AuctionListValidatorAPIResponse, error) {
	auctionList, err := epf.valStatsProc.GetAuctionList()
	if err != nil {
		return nil, err
	}

	return auctionList.AuctionListValidators, nil
}

// GetAddressConverter returns the address converter
func (pf *ProxyFacade) GetAddressConverter() (core.PubkeyConverter, error) {
	return pf.pubKeyConverter, nil
}

// GetLatestFullySynchronizedHyperblockNonce returns the latest fully synchronized hyperblock nonce
func (pf *ProxyFacade) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	return pf.nodeStatusProc.GetLatestFullySynchronizedHyperblockNonce()
}

// ComputeTransactionHash will compute hash of a given transaction
func (pf *ProxyFacade) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	return pf.txProc.ComputeTransactionHash(tx)
}

// GetTransactionsPool returns all txs from pool
func (pf *ProxyFacade) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	return pf.txProc.GetTransactionsPool(fields)
}

// GetTransactionsPoolForShard returns all txs from shard's pool
func (pf *ProxyFacade) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	return pf.txProc.GetTransactionsPoolForShard(shardID, fields)
}

// GetTransactionsPoolForSender returns tx pool for sender
func (pf *ProxyFacade) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	return pf.txProc.GetTransactionsPoolForSender(sender, fields)
}

// GetLastPoolNonceForSender returns last nonce from tx pool for sender
func (pf *ProxyFacade) GetLastPoolNonceForSender(sender string) (uint64, error) {
	return pf.txProc.GetLastPoolNonceForSender(sender)
}

// IsOldStorageForToken returns true is the storage for a given token is old
func (pf *ProxyFacade) IsOldStorageForToken(tokenID string, nonce uint64) (bool, error) {
	return pf.nodeGroupProc.IsOldStorageForToken(tokenID, nonce)
}

// GetTransactionsPoolNonceGapsForSender returns all nonce gaps from tx pool for sender
func (pf *ProxyFacade) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	return pf.txProc.GetTransactionsPoolNonceGapsForSender(sender)
}

// GetProof returns the Merkle proof for the given address
func (pf *ProxyFacade) GetProof(rootHash string, address string) (*data.GenericAPIResponse, error) {
	return pf.proofProc.GetProof(rootHash, address)
}

// GetProofDataTrie returns a Merkle proof for the given address and a Merkle proof for the given key
func (pf *ProxyFacade) GetProofDataTrie(rootHash string, address string, key string) (*data.GenericAPIResponse, error) {
	return pf.proofProc.GetProofDataTrie(rootHash, address, key)
}

// GetProofCurrentRootHash returns the Merkle proof for the given address
func (pf *ProxyFacade) GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error) {
	return pf.proofProc.GetProofCurrentRootHash(address)
}

// VerifyProof verifies the given Merkle proof
func (pf *ProxyFacade) VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error) {
	return pf.proofProc.VerifyProof(rootHash, address, proof)
}

// GetMetrics will return the status metrics
func (pf *ProxyFacade) GetMetrics() map[string]*data.EndpointMetrics {
	return pf.statusProc.GetMetrics()
}

// GetMetricsForPrometheus will return the status metrics in a prometheus format
func (pf *ProxyFacade) GetMetricsForPrometheus() string {
	return pf.statusProc.GetMetricsForPrometheus()
}

// GetGenesisNodesPubKeys retrieves the node's configuration public keys
func (pf *ProxyFacade) GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetGenesisNodesPubKeys()
}

// GetGasConfigs retrieves the current gas schedule configs
func (pf *ProxyFacade) GetGasConfigs() (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetGasConfigs()
}

// GetAboutInfo will return the app info
func (pf *ProxyFacade) GetAboutInfo() (*data.GenericAPIResponse, error) {
	return pf.aboutInfoProc.GetAboutInfo(), nil
}

// GetNodesVersions will return the version of the nodes
func (pf *ProxyFacade) GetNodesVersions() (*data.GenericAPIResponse, error) {
	return pf.aboutInfoProc.GetNodesVersions()
}

// GetAlteredAccountsByNonce returns altered accounts by nonce in block
func (pf *ProxyFacade) GetAlteredAccountsByNonce(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	return pf.blockProc.GetAlteredAccountsByNonce(shardID, nonce, options)
}

// GetAlteredAccountsByHash returns altered accounts by hash in block
func (pf *ProxyFacade) GetAlteredAccountsByHash(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error) {
	return pf.blockProc.GetAlteredAccountsByHash(shardID, hash, options)
}

// GetTriesStatistics will return trie statistics
func (pf *ProxyFacade) GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error) {
	return pf.nodeStatusProc.GetTriesStatistics(shardID)
}

// GetEpochStartData retrieves epoch start data for the provides epoch and shard ID
func (pf *ProxyFacade) GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error) {
	return pf.nodeStatusProc.GetEpochStartData(epoch, shardID)
}

// GetInternalStartOfEpochValidatorsInfo retrieves the validators info by epoch
func (pf *ProxyFacade) GetInternalStartOfEpochValidatorsInfo(epoch uint32) (*data.ValidatorsInfoApiResponse, error) {
	return pf.blockProc.GetInternalStartOfEpochValidatorsInfo(epoch)
}

// GetWaitingEpochsLeftForPublicKey returns the number of epochs left for the public key until it becomes eligible
func (epf *ProxyFacade) GetWaitingEpochsLeftForPublicKey(publicKey string) (*data.WaitingEpochsLeftApiResponse, error) {
	return epf.nodeGroupProc.GetWaitingEpochsLeftForPublicKey(publicKey)
}

// IsDataTrieMigrated returns true if the data trie for the given address is migrated
func (pf *ProxyFacade) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return pf.accountProc.IsDataTrieMigrated(address, options)
}
