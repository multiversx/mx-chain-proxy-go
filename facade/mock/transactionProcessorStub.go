package mock

import (
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

var errNotImplemented = errors.New("not implemented")

// TransactionProcessorStub -
type TransactionProcessorStub struct {
	SendTransactionCalled                       func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsCalled              func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SimulateTransactionCalled                   func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error)
	SendUserFundsCalled                         func(receiver string, value *big.Int) error
	TransactionCostRequestCalled                func(tx *data.Transaction) (*data.TxCostResponseData, error)
	GetTransactionStatusCalled                  func(txHash string, sender string) (string, error)
	GetProcessedTransactionStatusCalled         func(txHash string) (*data.ProcessStatusResponse, error)
	GetTransactionCalled                        func(txHash string, withEvents bool, relayedTxHash string) (*transaction.ApiTransactionResult, error)
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
	if tps.SimulateTransactionCalled != nil {
		return tps.SimulateTransactionCalled(tx, checkSignature)
	}

	return nil, errNotImplemented
}

// SendTransaction -
func (tps *TransactionProcessorStub) SendTransaction(tx *data.Transaction) (int, string, error) {
	if tps.SendTransactionCalled != nil {
		return tps.SendTransactionCalled(tx)
	}

	return 0, "", errNotImplemented
}

// SendMultipleTransactions -
func (tps *TransactionProcessorStub) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	if tps.SendMultipleTransactionsCalled != nil {
		return tps.SendMultipleTransactionsCalled(txs)
	}

	return data.MultipleTransactionsResponseData{}, errNotImplemented
}

// ComputeTransactionHash -
func (tps *TransactionProcessorStub) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	if tps.ComputeTransactionHashCalled != nil {
		return tps.ComputeTransactionHashCalled(tx)
	}

	return "", errNotImplemented
}

// SendUserFunds -
func (tps *TransactionProcessorStub) SendUserFunds(receiver string, value *big.Int) error {
	if tps.SendUserFundsCalled != nil {
		return tps.SendUserFundsCalled(receiver, value)
	}

	return errNotImplemented
}

// GetTransactionStatus -
func (tps *TransactionProcessorStub) GetTransactionStatus(txHash string, sender string) (string, error) {
	if tps.GetTransactionStatusCalled != nil {
		return tps.GetTransactionStatusCalled(txHash, sender)
	}

	return "", errNotImplemented
}

// GetProcessedTransactionStatus -
func (tps *TransactionProcessorStub) GetProcessedTransactionStatus(txHash string) (*data.ProcessStatusResponse, error) {
	if tps.GetProcessedTransactionStatusCalled != nil {
		return tps.GetProcessedTransactionStatusCalled(txHash)
	}

	return &data.ProcessStatusResponse{}, errNotImplemented
}

// GetTransaction -
func (tps *TransactionProcessorStub) GetTransaction(txHash string, withEvents bool, relayedTxHash string) (*transaction.ApiTransactionResult, error) {
	if tps.GetTransactionCalled != nil {
		return tps.GetTransactionCalled(txHash, withEvents, relayedTxHash)
	}

	return nil, errNotImplemented
}

// GetTransactionByHashAndSenderAddress -
func (tps *TransactionProcessorStub) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withEvents bool) (*transaction.ApiTransactionResult, int, error) {
	if tps.GetTransactionByHashAndSenderAddressCalled != nil {
		return tps.GetTransactionByHashAndSenderAddressCalled(txHash, sndAddr, withEvents)
	}

	return nil, 0, errNotImplemented
}

// TransactionCostRequest -
func (tps *TransactionProcessorStub) TransactionCostRequest(tx *data.Transaction) (*data.TxCostResponseData, error) {
	if tps.TransactionCostRequestCalled != nil {
		return tps.TransactionCostRequestCalled(tx)
	}

	return nil, errNotImplemented
}

// GetTransactionsPool -
func (tps *TransactionProcessorStub) GetTransactionsPool(fields string) (*data.TransactionsPool, error) {
	if tps.GetTransactionsPoolCalled != nil {
		return tps.GetTransactionsPoolCalled(fields)
	}

	return nil, errNotImplemented
}

// GetTransactionsPoolForShard -
func (tps *TransactionProcessorStub) GetTransactionsPoolForShard(shardID uint32, fields string) (*data.TransactionsPool, error) {
	if tps.GetTransactionsPoolForShardCalled != nil {
		return tps.GetTransactionsPoolForShardCalled(shardID, fields)
	}

	return nil, errNotImplemented
}

// GetTransactionsPoolForSender -
func (tps *TransactionProcessorStub) GetTransactionsPoolForSender(sender, fields string) (*data.TransactionsPoolForSender, error) {
	if tps.GetTransactionsPoolForSenderCalled != nil {
		return tps.GetTransactionsPoolForSenderCalled(sender, fields)
	}

	return nil, errNotImplemented
}

// GetLastPoolNonceForSender -
func (tps *TransactionProcessorStub) GetLastPoolNonceForSender(sender string) (uint64, error) {
	if tps.GetLastPoolNonceForSenderCalled != nil {
		return tps.GetLastPoolNonceForSenderCalled(sender)
	}

	return 0, errNotImplemented
}

// GetTransactionsPoolNonceGapsForSender -
func (tps *TransactionProcessorStub) GetTransactionsPoolNonceGapsForSender(sender string) (*data.TransactionsPoolNonceGaps, error) {
	if tps.GetTransactionsPoolNonceGapsForSenderCalled != nil {
		return tps.GetTransactionsPoolNonceGapsForSenderCalled(sender)
	}

	return nil, errNotImplemented
}
