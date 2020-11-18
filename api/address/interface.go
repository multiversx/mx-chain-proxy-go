package address

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetAccount(address string) (*data.Account, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
	GetShardIDForAddress(address string) (uint32, error)
	GetValueForKey(address string, key string) (string, error)
	GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error)
	GetAllESDTTokens(address string) (*data.GenericAPIResponse, error)
}
