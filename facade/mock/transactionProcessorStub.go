package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionProcessorStub -
type TransactionProcessorStub struct {
	SendTransactionCalled                      func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsCalled             func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionCalled                  func(tx *data.Transaction) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                        func(receiver string, value *big.Int) error
	TransactionCostRequestHandler              func(tx *data.Transaction) (string, error)
	GetTransactionStatusHandler                func(txHash string, sender string) (string, error)
	GetTransactionCalled                       func(txHash string, withEvents bool) (*data.FullTransaction, error)
	GetTransactionByHashAndSenderAddressCalled func(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error)
	ComputeTransactionHashCalled               func(tx *data.Transaction) (string, error)
}

// SimulateTransaction -
func (tps *TransactionProcessorStub) SimulateTransaction(tx *data.Transaction) (*data.GenericAPIResponse, error) {
	return tps.SimulateTransactionCalled(tx)
}

// SendTransaction -
func (tps *TransactionProcessorStub) SendTransaction(tx *data.Transaction) (int, string, error) {
	return tps.SendTransactionCalled(tx)
}

// SendMultipleTransactions -
func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
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
func (tps *TransactionProcessorStub) GetTransactionStatus(txHash string, sender string) (string, error) {
	return tps.GetTransactionStatusHandler(txHash, sender)
}

// GetTransaction -
func (tps *TransactionProcessorStub) GetTransaction(txHash string, withEvents bool) (*data.FullTransaction, error) {
	return tps.GetTransactionCalled(txHash, withEvents)
}

// GetTransactionByHashAndSenderAddress -
func (tps *TransactionProcessorStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*data.FullTransaction, int, error) {
	return tps.GetTransactionByHashAndSenderAddressCalled(txHash, sndAddr, withEvents)
}

// TransactionCostRequest --
func (tps *TransactionProcessorStub) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return tps.TransactionCostRequestHandler(tx)
}
