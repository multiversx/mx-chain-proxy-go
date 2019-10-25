package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Facade is the mock implementation of a node router handler
type Facade struct {
	GetAccountHandler               func(address string) (*data.Account, error)
	SendTransactionHandler          func(tx *data.Transaction) (string, error)
	SendMultipleTransactionsHandler func(txs []*data.Transaction) (uint64, error)
	SendUserFundsCalled             func(receiver string, value *big.Int) error
	SignAndSendTransactionCalled    func(tx *data.Transaction, sk []byte) (string, error)
	PublicKeyFromPrivateKeyCalled   func(privateKeyHex string) (string, error)
	GetVmValueHandler               func(resType string, address string, funcName string, argsBuff ...[]byte) ([]byte, error)
	GetHeartbeatDataHandler         func() (*data.HeartbeatResponse, error)
}

// getSignedTransaction is the mock implementation of a handler's getSignedTransaction method
func (f *Facade) SignAndSendTransaction(tx *data.Transaction, sk []byte) (string, error) {
	return f.SignAndSendTransactionCalled(tx, sk)
}

// PublicKeyFromPrivateKey is the mock implementation of a handler's PublicKeyFromPrivateKey method
func (f *Facade) PublicKeyFromPrivateKey(privateKeyHex string) (string, error) {
	return f.PublicKeyFromPrivateKeyCalled(privateKeyHex)
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(tx *data.Transaction) (string, error) {
	return f.SendTransactionHandler(tx)
}

// SendMultipleTransactions is the mock implementation of a handler's SendMultipleTransactions method
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// GenerateTxForSendUserFunds is the mock implementation of a handler's GenerateTxForSendUserFunds method
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
