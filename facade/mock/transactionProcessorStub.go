package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionProcessorStub -
type TransactionProcessorStub struct {
	SendTransactionCalled                      func(tx *data.Transaction) (string, int, error)
	SendMultipleTransactionsCalled             func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, int, error)
	SimulateTransactionCalled                  func(tx *data.Transaction) (*data.GenericAPIResponse, int, error)
	SendUserFundsCalled                        func(receiver string, value *big.Int) error
	TransactionCostRequestHandler              func(tx *data.Transaction) (string, int, error)
	GetTransactionStatusHandler                func(txHash string, sender string) (string, int, error)
	GetTransactionCalled                       func(txHash string, withEvents bool) (*data.FullTransaction, int, error)
	GetTransactionByHashAndSenderAddressCalled func(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error)
	ComputeTransactionHashCalled               func(tx *data.Transaction) (string, error)
}

// SimulateTransaction -
func (tps *TransactionProcessorStub) SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, int, error) {
	return tps.SimulateTransactionCalled(tx)
}

// SendTransaction -
func (tps *TransactionProcessorStub) SendTransaction(tx *data.Transaction) (string, int, error) {
	return tps.SendTransactionCalled(tx)
}

// SendMultipleTransactions -
func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, int, error) {
	return tps.SendMultipleTransactionsCalled(txs)
}

// ComputeTransactionHash -
func (tps *TransactionProcessorStub) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	return tps.ComputeTransactionHashCalled(tx)
}

// SendUserFunds -
func (tps *TransactionProcessorStub) SendUserFunds(receiver string, value *big.Int) error {
	return tps.SendUserFundsCalled(receiver, value)
}

// GetTransactionStatus -
func (tps *TransactionProcessorStub) GetTransactionStatus(txHash string, sender string) (string, int, error) {
	return tps.GetTransactionStatusHandler(txHash, sender)
}

// GetTransaction -
func (tps *TransactionProcessorStub) GetTransaction(txHash string, withEvents bool) (*data.FullTransaction, int, error) {
	return tps.GetTransactionCalled(txHash, withEvents)
}

// GetTransactionByHashAndSenderAddress -
func (tps *TransactionProcessorStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error) {
	return tps.GetTransactionByHashAndSenderAddressCalled(txHash, sndAddr, withEvents)
}

// TransactionCostRequest --
func (tps *TransactionProcessorStub) TransactionCostRequest(tx *data.Transaction) (string, int, error) {
	return tps.TransactionCostRequestHandler(tx)
}
