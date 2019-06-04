package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node router handler
type Facade struct {
	GetAccountHandler      func(address string) (*data.Account, error)
	SendTransactionHandler func(nonce uint64, sender string, receiver string, value *big.Int, code string, signature []byte) (*data.Transaction, error)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(nonce uint64, sender string, receiver string, value *big.Int, code string, signature []byte) (*data.Transaction, error) {
	return f.SendTransactionHandler(nonce, sender, receiver, value, code, signature)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
