package mock

import "math/big"

type TransactionProcessorStub struct {
	SendTransactionCalled func(nonce uint64, sender string, receiver string,
		value *big.Int, code string, signature []byte, gasPrice uint64, gasLimit uint64) (string, error)
	SendUserFundsCalled func(receiver string) error
}

func (tps *TransactionProcessorStub) SendTransaction(nonce uint64, sender string,
	receiver string, value *big.Int, code string, signature []byte,
	gasPrice uint64, gasLimit uint64) (string, error) {

	return tps.SendTransactionCalled(nonce, sender, receiver, value, code,
		signature, gasPrice, gasLimit)
}

func (tps *TransactionProcessorStub) SendUserFunds(receiver string) error {
	return tps.SendUserFundsCalled(receiver)
}
