package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessorStub --
type AccountProcessorStub struct {
	GetAccountCalled           func(address string) (*data.Account, int, error)
	GetValueForKeyCalled       func(address string, key string) (string, int, error)
	GetShardIDForAddressCalled func(address string) (uint32, int, error)
	GetTransactionsCalled      func(address string) ([]data.DatabaseTransaction, int, error)
	ValidatorStatisticsCalled  func() (map[string]*data.ValidatorApiResponse, error)
	GetAllESDTTokensCalled     func(address string) (*data.GenericAPIResponse, int, error)
	GetESDTTokenDataCalled     func(address string, key string) (*data.GenericAPIResponse, int, error)
}

// GetAllESDTTokens -
func (aps *AccountProcessorStub) GetAllESDTTokens(address string) (*data.GenericAPIResponse, int, error) {
	return aps.GetAllESDTTokensCalled(address)
}

// GetESDTTokenData -
func (aps *AccountProcessorStub) GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, int, error) {
	return aps.GetESDTTokenDataCalled(address, key)
}

// GetAccount --
func (aps *AccountProcessorStub) GetAccount(address string) (*data.Account, int, error) {
	return aps.GetAccountCalled(address)
}

// GetValueForKey --
func (aps *AccountProcessorStub) GetValueForKey(address string, key string) (string, int, error) {
	return aps.GetValueForKeyCalled(address, key)
}

// GetShardIDForAddress --
func (aps *AccountProcessorStub) GetShardIDForAddress(address string) (uint32, int, error) {
	return aps.GetShardIDForAddressCalled(address)
}

// GetTransactions --
func (aps *AccountProcessorStub) GetTransactions(address string) ([]data.DatabaseTransaction, int, error) {
	return aps.GetTransactionsCalled(address)
}

// ValidatorStatistics --
func (aps *AccountProcessorStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return aps.ValidatorStatisticsCalled()
}
