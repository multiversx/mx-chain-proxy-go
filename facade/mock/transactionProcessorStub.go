package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type TransactionProcessorStub struct {
	SendTransactionCalled          func(tx *data.ApiTransaction) (int, string, error)
	SendMultipleTransactionsCalled func(txs []*data.ApiTransaction) (uint64, error)
	SendUserFundsCalled            func(receiver string, value *big.Int) error
	TransactionCostRequestHandler  func(tx *data.ApiTransaction) (string, error)
}

func (tps *TransactionProcessorStub) SendTransaction(tx *data.ApiTransaction) (int, string, error) {
	return tps.SendTransactionCalled(tx)
}

func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.ApiTransaction) (uint64, error) {
	return tps.SendMultipleTransactionsCalled(txs)
}

func (tps *TransactionProcessorStub) SendUserFunds(receiver string, value *big.Int) error {
	return tps.SendUserFundsCalled(receiver, value)
}

// TransactionCostRequest --
func (tps *TransactionProcessorStub) TransactionCostRequest(tx *data.ApiTransaction) (string, error) {
	return tps.TransactionCostRequestHandler(tx)
}
