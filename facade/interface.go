package facade

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(tx *data.Transaction) (string, error)
	SendUserFunds(receiver string) error
}

// VmValuesProcessor defines what a get value processor should do
type VmValuesProcessor interface {
	GetVmValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}
