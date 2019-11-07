package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (uint64, error)
}

// VmValuesProcessor defines what a get value processor should do
type VmValuesProcessor interface {
	GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

// HeartbeatProcessor defines what a heartbeat processor should do
type HeartbeatProcessor interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
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
	) (*data.Transaction, error)
}
