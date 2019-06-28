package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node router handler
type Facade struct {
	GetAccountHandler      func(address string) (*data.Account, error)
	SendTransactionHandler func(nonce uint64, sender string, receiver string,
		value *big.Int, code string, signature []byte, gasPrice uint64,
		gasLimit uint64) (string, error)
	SendUserFundsCalled func(receiver string) error
	GetDataValueHandler func(address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(
	nonce uint64,
	sender string,
	receiver string,
	value *big.Int,
	code string,
	signature []byte,
	gasPrice uint64,
	gasLimit uint64) (string, error) {

	return f.SendTransactionHandler(
		nonce,
		sender,
		receiver,
		value,
		code,
		signature,
		gasPrice,
		gasLimit)
}

// SendUserFunds is the mock implementation of a handler's SendUserFunds method
func (f *Facade) SendUserFunds(receiver string) error {
	return f.SendUserFundsCalled(receiver)
}

func (f *Facade) GetDataValue(address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return f.GetDataValueHandler(address, funcName, argsBuff...)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
