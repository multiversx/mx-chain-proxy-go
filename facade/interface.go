package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, int, error)
	GetShardIDForAddress(address string) (uint32, int, error)
	GetValueForKey(address string, key string) (string, int, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, int, error)
	GetAllESDTTokens(address string) (*data.GenericAPIResponse, int, error)
	GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, int, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (string, int, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, int, error)
	SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, int, error)
	TransactionCostRequest(tx *data.Transaction) (string, int, error)
	GetTransactionStatus(txHash string, sender string) (string, int, error)
	GetTransaction(txHash string, withEvents bool) (*data.FullTransaction, int, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error)
	ComputeTransactionHash(tx *data.Transaction) (string, error)
}

// SCQueryService defines how data should be get from a SC account
type SCQueryService interface {
	ExecuteQuery(query *data.SCQuery) (*vm.VMOutputApi, int, error)
}

// HeartbeatProcessor defines what a heartbeat processor should do
type HeartbeatProcessor interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
}

// ValidatorStatisticsProcessor defines what a validator statistics processor should do
type ValidatorStatisticsProcessor interface {
	GetValidatorStatistics() (*data.ValidatorStatisticsResponse, error)
}

// NodeStatusProcessor defines what a node status processor should do
type NodeStatusProcessor interface {
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, int, error)
	GetEconomicsDataMetrics() (*data.GenericAPIResponse, error)
	GetLatestFullySynchronizedHyperblockNonce() (uint64, error)
}

// BlockProcessor defines what a block processor should do
type BlockProcessor interface {
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, int, error)
	GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, int, error)
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, int, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, int, error)
	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, int, error)
}

// FaucetProcessor defines what a component which will handle faucets should do
type FaucetProcessor interface {
	IsEnabled() bool
	SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, int, error)
	GenerateTxForSendUserFunds(
		senderSk crypto.PrivateKey,
		senderPk string,
		senderNonce uint64,
		receiver string,
		value *big.Int,
		chainID string,
		version uint32,
	) (*data.Transaction, error)
}
