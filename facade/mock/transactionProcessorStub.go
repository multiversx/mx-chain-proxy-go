package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// TransactionProcessorStub -
type TransactionProcessorStub struct {
	SendTransactionCalled                       func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsCalled              func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionCalled                   func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                         func(receiver string, value *big.Int) error
	TransactionCostRequestHandler               func(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatusHandler                 func(txHash string, sender string) (string, error)
	GetTransactionCalled                        func(txHash string, withEvents bool) (*transaction.ApiTransactionResult, error)
	GetTransactionByHashAndSenderAddressCalled  func(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error)
	ComputeTransactionHashCalled                func(tx *data.Transaction) (string, error)
	GetTransactionsPoolCalled                   func(fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForShardCalled           func(shardID uint32, fields string) (*data.TransactionsPool, error)
	GetTransactionsPoolForSenderCalled          func(sender, fields string) (*data.TransactionsPoolForSender, error)
	GetLastPoolNonceForSenderCalled             func(sender string) (uint64, error)
	GetTransactionsPoolNonceGapsForSenderCalled func(sender string) (*data.TransactionsPoolNonceGaps, error)
}

// SimulateTransaction -
func (tps *TransactionProcessorStub) SimulateTransaction(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
	return tps.SimulateTransactionCalled(tx, checkSignature)
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
func (tps *TransactionProcessorStub) GetTransaction(txHash string, withEvents bool) (*transaction.ApiTransactionResult, error) {
	return tps.GetTransactionCalled(txHash, withEvents)
}

// GetTransactionByHashAndSenderAddress -
func (tps *TransactionProcessorStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error) {
	return tps.GetTransactionByHashAndSenderAddressCalled(txHash, sndAddr, withEvents)
}

// TransactionCostRequest -
func (tps *TransactionProcessorStub) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	return tps.TransactionCostRequestHandler(tx)
}

// GetTransactionsPool -
func (tps *TransactionProcessorStub) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	if tps.GetTransactionsPoolCalled != nil {
		return tps.GetTransactionsPoolCalled(fields)
	}

	return nil, nil
}

// GetTransactionsPoolForShard -
func (tps *TransactionProcessorStub) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	if tps.GetTransactionsPoolForShardCalled != nil {
		return tps.GetTransactionsPoolForShardCalled(shardID, fields)
	}

	return nil, nil
}

// GetTransactionsPoolForSender -
func (tps *TransactionProcessorStub) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	if tps.GetTransactionsPoolForSenderCalled != nil {
		return tps.GetTransactionsPoolForSenderCalled(sender, fields)
	}

	return nil, nil
}

// GetLastPoolNonceForSender -
func (tps *TransactionProcessorStub) GetLastPoolNonceForSender(sender string) (uint64, error) {
	if tps.GetLastPoolNonceForSenderCalled != nil {
		return tps.GetLastPoolNonceForSenderCalled(sender)
	}

	return 0, nil
}

// GetTransactionsPoolNonceGapsForSender -
func (tps *TransactionProcessorStub) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	if tps.GetTransactionsPoolNonceGapsForSenderCalled != nil {
		return tps.GetTransactionsPoolNonceGapsForSenderCalled(sender)
	}

	return nil, nil
}
