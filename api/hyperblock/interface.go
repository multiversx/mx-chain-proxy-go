package hyperblock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type facadeHandler interface {
	GetHyperBlockByNonce(nonce uint64, withTxs bool) (*data.GenericAPIResponse, error)
	GetHyperBlockByHash(hash string, withTxs bool) (*data.GenericAPIResponse, error)
}
