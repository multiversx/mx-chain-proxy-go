package transaction

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SendUserFunds(receiver string, value *big.Int) error
	TransactionCostRequest(tx *data.Transaction) (string, error)
	GetTransactionStatus(txHash string, sender string) (string, error)
	GetTransaction(txHash string) (*transaction.ApiTransactionResult, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string) (*transaction.ApiTransactionResult, int, error)
}
