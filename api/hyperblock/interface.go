package hyperblock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type facadeHandler interface {
	GetHyperBlockByNonce(nonce uint64) (*data.GenericAPIResponse, error)
	GetHyperBlockByHash(hash string) (*data.GenericAPIResponse, error)
}
