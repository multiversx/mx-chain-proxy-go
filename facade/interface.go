package facade

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessor defines what an account request processor should do
type AccountProcessor interface {
	GetAccount(address string) (*data.Account, error)
}

// TransactionProcessor defines what a transaction request processor should do
type TransactionProcessor interface {
	SendTransaction(nonce uint64, sender string, receiver string, value *big.Int,
		data string, signature []byte, gasPrice uint64, gasLimit uint64) (string, error)
	SendUserFunds(receiver string) error
}

// GetValuesProcessor defines what a get value processor should do
type GetValuesProcessor interface {
	GetDataValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}
