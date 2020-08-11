package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
	GetValueForKey(address string, key string) (string, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	TransactionCostRequest(tx *data.Transaction) (string, error)
	GetTransactionStatus(txHash string, sender string) (string, error)
	GetTransaction(txHash string) (*transaction.ApiTransactionResult, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string) (*transaction.ApiTransactionResult, int, error)
}

// SCQueryService defines how data should be get from a SC account
type SCQueryService interface {
	ExecuteQuery(query *data.SCQuery) (*vmcommon.VMOutput, error)
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
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error)
}

// BlockProcessor defines what a block processor should do
type BlockProcessor interface {
	GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error)
	GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.GenericAPIResponse, error)
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.GenericAPIResponse, error)
	GetHyperBlockByHash(hash string, withTxs bool) (*data.GenericAPIResponse, error)
	GetHyperBlockByNonce(nonce uint64, withTxs bool) (*data.GenericAPIResponse, error)
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
		chainID string,
		version uint32,
	) (*data.Transaction, error)
}
