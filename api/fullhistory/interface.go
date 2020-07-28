package fullhistory

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.GenericAPIResponse, error)
	GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.GenericAPIResponse, error)
}
