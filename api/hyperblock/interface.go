package hyperblock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type facadeHandler interface {
	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error)
}
