package wallet

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	SignAndSendTransaction(tx *data.Transaction, sk []byte) (string, error)
	PublicKeyFromPrivateKey(privateKeyHex string) (string, error)
}
