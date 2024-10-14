package facade

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// ActionsProcessor defines what an actions processor should do
type ActionsProcessor interface {
	ReloadObservers() data.NodesReloadResponse
	ReloadFullHistoryObservers() data.NodesReloadResponse
}

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error)
	GetAccounts(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error)
	GetShardIDForAddress(address string) (uint32, error)
	GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error)
	GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatus(txHash string, sender string) (string, error)
	GetTransaction(txHash string, withEvents bool, relayedTxHash string) (*transaction.ApiTransactionResult, error)
	GetProcessedTransactionStatus(txHash string) (*data.ProcessStatusResponse, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error)
	ComputeTransactionHash(tx *data.Transaction) (string, error)
	GetTransactionsPool(fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error)
	GetLastPoolNonceForSender(sender string) (uint64, error)
	GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error)
}

// ProofProcessor defines what a proof request processor should do
type ProofProcessor interface {
	GetProof(rootHash string, address string) (*data.GenericAPIResponse, error)
	GetProofDataTrie(rootHash string, address string, key string) (*data.GenericAPIResponse, error)
	GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error)
	VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error)
}

// SCQueryService defines how data should be get from a SC account
type SCQueryService interface {
	ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, data.BlockInfo, error)
}

// NodeGroupProcessor defines what a node group processor should do
type NodeGroupProcessor interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
	IsOldStorageForToken(tokenID string, nonce uint64) (bool, error)
	GetWaitingEpochsLeftForPublicKey(publicKey string) (*data.WaitingEpochsLeftApiResponse, error)
}

// ValidatorStatisticsProcessor defines what a validator statistics processor should do
type ValidatorStatisticsProcessor interface {
	GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error)
	GetAuctionList() (*data.AuctionListResponse, error)
}

// ESDTSupplyProcessor defines what an esdt supply processor should do
type ESDTSupplyProcessor interface {
	GetESDTSupply(token string) (*data.ESDTSupplyResponse, error)
}

// NodeStatusProcessor defines what a node status processor should do
type NodeStatusProcessor interface {
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error)
	GetEconomicsDataMetrics() (*data.GenericAPIResponse, error)
	GetLatestFullySynchronizedHyperblockNonce() (uint64, error)
	GetAllIssuedESDTs(tokenType string) (*data.GenericAPIResponse, error)
	GetEnableEpochsMetrics() (*data.GenericAPIResponse, error)
	GetDirectStakedInfo() (*data.GenericAPIResponse, error)
	GetDelegatedInfo() (*data.GenericAPIResponse, error)
	GetRatingsConfig() (*data.GenericAPIResponse, error)
	GetGenesisNodesPubKeys() (*data.GenericAPIResponse, error)
	GetGasConfigs() (*data.GenericAPIResponse, error)
	GetTriesStatistics(shardID uint32) (*data.TrieStatisticsAPIResponse, error)
	GetEpochStartData(epoch uint32, shardID uint32) (*data.GenericAPIResponse, error)
}

// BlocksProcessor defines what a blocks processor should do
type BlocksProcessor interface {
	GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error)
}

// BlockProcessor defines what a block processor should do
type BlockProcessor interface {
	GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)

	GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error)
	GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error)

	GetAlteredAccountsByNonce(shardID uint32, nonce uint64, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error)
	GetAlteredAccountsByHash(shardID uint32, hash string, options common.GetAlteredAccountsForBlockOptions) (*data.AlteredAccountsApiResponse, error)
	GetInternalStartOfEpochValidatorsInfo(epoch uint32) (*data.ValidatorsInfoApiResponse, error)
}

// FaucetProcessor defines what a component which will handle faucets should do
type FaucetProcessor interface {
	IsEnabled() bool
	SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error)
	GenerateTxForSendUserFunds(
		senderSk crypto.PrivateKey,
		senderPk string,
		senderNonce uint64,
		receiver string,
		value *big.Int,
		networkConfig *data.NetworkConfig,
	) (*data.Transaction, error)
}

// StatusProcessor defines what a component which will handle status request should do
type StatusProcessor interface {
	GetMetrics() map[string]*data.EndpointMetrics
	GetMetricsForPrometheus() string
}

// AboutInfoProcessor defines the behaviour of about info processor
type AboutInfoProcessor interface {
	GetAboutInfo() *data.GenericAPIResponse
	GetNodesVersions() (*data.GenericAPIResponse, error)
}
