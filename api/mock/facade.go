package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node router handler
type Facade struct {
	GetAccountHandler      func(address string) (*data.Account, error)
	SendTransactionHandler func(tx *data.Transaction) (string, error)
	SendUserFundsCalled    func(receiver string) error
	GetVmValueHandler      func(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(tx *data.Transaction) (string, error) {
	return f.SendTransactionHandler(tx)
}

// SendUserFunds is the mock implementation of a handler's SendUserFunds method
func (f *Facade) SendUserFunds(receiver string) error {
	return f.SendUserFundsCalled(receiver)
}

func (f *Facade) GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return f.GetVmValueHandler(resType, address, funcName, argsBuff...)
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
