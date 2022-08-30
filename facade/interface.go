package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-go-core/data/vm"
	"github.com/ElrondNetwork/elrond-go-crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ActionsProcessor defines what an actions processor should do
type ActionsProcessor interface {
	ReloadObservers() data.NodesReloadResponse
	ReloadFullHistoryObservers() data.NodesReloadResponse
}

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error)
	GetShardIDForAddress(address string) (uint32, error)
	GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
	GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatus(txHash string, sender string) (string, error)
	GetTransaction(txHash string, withEvents bool) (*transaction.ApiTransactionResult, error)
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
	GetProofCurrentRootHash(address string) (*data.GenericAPIResponse, error)
	VerifyProof(rootHash string, address string, proof []string) (*data.GenericAPIResponse, error)
}

// SCQueryService defines how data should be get from a SC account
type SCQueryService interface {
	ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, error)
}

// HeartbeatProcessor defines what a heartbeat processor should do
type HeartbeatProcessor interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
}

// ValidatorStatisticsProcessor defines what a validator statistics processor should do
type ValidatorStatisticsProcessor interface {
	GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error)
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
}

// BlocksProcessor defines what a blocks processor should do
type BlocksProcessor interface {
	GetBlocksByRound(round uint64, options common.BlockQueryOptions) (*data.BlocksApiResponse, error)
}

// BlockProcessor defines what a block processor should do
type BlockProcessor interface {
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetBlockByHash(shardID uint32, hash string, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)

	GetInternalBlockByHash(shardID uint32, hash string, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
	GetInternalMiniBlockByHash(shardID uint32, hash string, epoch uint32, format common.OutputFormat) (*data.InternalMiniBlockApiResponse, error)
	GetInternalStartOfEpochMetaBlock(epoch uint32, format common.OutputFormat) (*data.InternalBlockApiResponse, error)
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
