package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.ApiTransaction) (int, string, error)
	SendMultipleTransactions(txs []*data.ApiTransaction) (uint64, error)
	TransactionCostRequest(tx *data.ApiTransaction) (string, error)
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
	GetShardStatus(shardID uint32) (map[string]interface{}, error)
	GetEpochMetrics(shardID uint32) (map[string]interface{}, error)
}

// FaucetProcessor defines what a component which will handle faucets should do
type FaucetProcessor interface {
	SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error)
	GenerateTxForSendUserFunds(
		senderSk crypto.PrivateKey,
		senderPk string,
		senderNonce uint64,
		receiver string,
		value *big.Int,
	) (*data.ApiTransaction, error)
}
