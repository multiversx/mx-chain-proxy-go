package facade

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"math/big"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (string, error)
	SendMultipleTransactions(txs []*data.Transaction) (uint64, error)
	SendUserFunds(receiver string, value *big.Int) error
}

// VmValuesProcessor defines what a get value processor should do
type VmValuesProcessor interface {
	GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

// HeartbeatProcessor defines what a heartbeat processor should do
type HeartbeatProcessor interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
}
