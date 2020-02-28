package transaction

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	SendTransaction(tx *data.Transaction) (int, string, error)
	SendMultipleTransactions(txs []*data.Transaction) (uint64, error)
	SendUserFunds(receiver string, value *big.Int) error
	GetTransactionCost(tx *data.Transaction) (string, error)
}
