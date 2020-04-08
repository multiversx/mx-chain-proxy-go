package transaction

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	SendTransaction(tx *data.ApiTransaction) (int, string, error)
	SendMultipleTransactions(txs []*data.ApiTransaction) (uint64, error)
	SendUserFunds(receiver string, value *big.Int) error
	TransactionCostRequest(tx *data.ApiTransaction) (string, error)
}
