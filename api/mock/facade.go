package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node's router handler
type Facade struct {
	GetAccountHandler               func(address string) (*data.Account, error)
	SendTransactionHandler          func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler func(txs []*data.Transaction) (uint64, error)
	SendUserFundsCalled             func(receiver string, value *big.Int) error
	GetVmValueHandler               func(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error)
	GetHeartbeatDataHandler         func() (*data.HeartbeatResponse, error)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SendMultipleTransactions is the mock implementation of a handler's SendMultipleTransactions method
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// SendUserFunds is the mock implementation of a handler's SendUserFunds method
func (f *Facade) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// GetVmValue is the mock implementation of a handler's GetVmValue method
func (f *Facade) GetVmValue(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error) {
	return f.GetVmValueHandler(resType, address, funcName, argsBuff...)
}

// GetHeartbeatData is the mock implementation of a handler's GetHeartbeatData method
func (f *Facade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
