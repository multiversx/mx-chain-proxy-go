package v_next

import "github.com/multiversx/mx-chain-proxy-go/data"

// AccountsFacadeHandlerV_next interface defines methods that can be used from facade context variable
type AccountsFacadeHandlerV_next interface {
	GetAccount(address string) (*data.AccountModel, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
	GetShardIDForAddressV_next(address string, additional int) (uint32, error)
	GetValueForKey(address string, key string) (string, error)
	NextEndpointHandler() string
}
