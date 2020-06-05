package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionProcessorStub -
type TransactionProcessorStub struct {
	SendTransactionCalled                      func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsCalled             func(txs []*data.Transaction) (data.ResponseMultipleTransactions, error)
	SendUserFundsCalled                        func(receiver string, value *big.Int) error
	TransactionCostRequestHandler              func(tx *data.Transaction) (string, error)
	GetTransactionStatusHandler                func(txHash string) (string, error)
	GetTransactionCalled                       func(txHash string) (*transaction.ApiTransactionResult, error)
	GetTransactionByHashAndSenderAddressCalled func(txHash string, sndAddr string) (*transaction.ApiTransactionResult, int, error)
}

// SendTransaction -
func (tps *TransactionProcessorStub) SendTransaction(tx *data.Transaction) (int, string, error) {
	return tps.SendTransactionCalled(tx)
}

// SendMultipleTransactions -
func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.Transaction) (data.ResponseMultipleTransactions, error) {
	return tps.SendMultipleTransactionsCalled(txs)
}

// SendUserFunds -
func (tps *TransactionProcessorStub) SendUserFunds(receiver string, value *big.Int) error {
	return tps.SendUserFundsCalled(receiver, value)
}

// GetTransactionStatus -
func (tps *TransactionProcessorStub) GetTransactionStatus(txHash string) (string, error) {
	return tps.GetTransactionStatusHandler(txHash)
}

// GetTransaction -
func (tps *TransactionProcessorStub) GetTransaction(txHash string) (*transaction.ApiTransactionResult, error) {
	return tps.GetTransactionCalled(txHash)
}

// GetTransactionByHashAndSenderAddress -
func (tps *TransactionProcessorStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string) (*transaction.ApiTransactionResult, int, error) {
	return tps.GetTransactionByHashAndSenderAddressCalled(txHash, sndAddr)
}

// TransactionCostRequest --
func (tps *TransactionProcessorStub) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return tps.TransactionCostRequestHandler(tx)
}
