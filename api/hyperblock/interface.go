package hyperblock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler defines the actions needed for fetching the hyperblocks from the nodes
type FacadeHandler interface {
	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error)
}
