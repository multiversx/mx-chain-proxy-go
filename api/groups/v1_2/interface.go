package v1_2

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// AccountsFacadeHandlerV1_2 interface defines methods that can be used from facade context variable
type AccountsFacadeHandlerV1_2 interface {
	GetAccount(address string) (*data.Account, error)
	GetTransactions(address string) ([]data.DatabaseTransaction, error)
	GetShardIDForAddressV1_2(address string, additional int) (uint32, error)
	GetValueForKey(address string, key string) (string, error)
}
