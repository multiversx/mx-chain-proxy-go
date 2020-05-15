package address

import (
	"github.com/ElrondNetwork/elrond-go/core/indexer"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetAccount(address string) (*data.Account, error)
	GetTransactions(address string) ([]indexer.Transaction, error)
}
