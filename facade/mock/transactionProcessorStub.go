package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type TransactionProcessorStub struct {
	SendTransactionCalled          func(tx *data.Transaction) (string, error)
	SignAndSendTransactionCalled   func(tx *data.Transaction, sk []byte) (string, error)
	SendMultipleTransactionsCalled func(txs []*data.Transaction) (uint64, error)
	SendUserFundsCalled            func(receiver string, value *big.Int) error
}

func (tps *TransactionProcessorStub) SendTransaction(tx *data.Transaction) (string, error) {
	return tps.SendTransactionCalled(tx)
}

func (tps *TransactionProcessorStub) SignAndSendTransaction(tx *data.Transaction, sk []byte) (string, error) {
	return tps.SignAndSendTransactionCalled(tx, sk)
}

func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.Transaction) (uint64, error) {
	return tps.SendMultipleTransactionsCalled(txs)
}

func (tps *TransactionProcessorStub) SendUserFunds(receiver string, value *big.Int) error {
	return tps.SendUserFundsCalled(receiver, value)
}
