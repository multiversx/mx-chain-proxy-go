package transaction

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	SendTransaction(tx *data.Transaction) (string, error)
	SendUserFunds(receiver string) error
}
