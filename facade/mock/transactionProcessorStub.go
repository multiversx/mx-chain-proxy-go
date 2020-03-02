package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type TransactionProcessorStub struct {
	SendTransactionCalled             func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsCalled    func(txs []*data.Transaction) (uint64, error)
	SendUserFundsCalled               func(receiver string, value *big.Int) error
	SendTransactionCostRequestHandler func(tx *data.Transaction) (string, error)
}

func (tps *TransactionProcessorStub) SendTransaction(tx *data.Transaction) (int, string, error) {
	return tps.SendTransactionCalled(tx)
}

func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	return tps.SendMultipleTransactionsCalled(txs)
}

func (tps *TransactionProcessorStub) SendUserFunds(receiver string, value *big.Int) error {
	return tps.SendUserFundsCalled(receiver, value)
}

// SendTransactionCostRequest --
func (tps *TransactionProcessorStub) SendTransactionCostRequest(tx *data.Transaction) (string, error) {
	return tps.SendTransactionCostRequestHandler(tx)
}
