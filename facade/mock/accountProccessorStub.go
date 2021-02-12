package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessorStub --
type AccountProcessorStub struct {
	GetAccountCalled           func(address string) (*data.Account, error)
	GetValueForKeyCalled       func(address string, key string) (string, error)
	GetShardIDForAddressCalled func(address string) (uint32, error)
	GetTransactionsCalled      func(address string) ([]data.DatabaseTransaction, error)
	ValidatorStatisticsCalled  func() (map[string]*data.ValidatorApiResponse, error)
	GetAllESDTTokensCalled     func(address string) (*data.GenericAPIResponse, error)
	GetESDTTokenDataCalled     func(address string, key string) (*data.GenericAPIResponse, error)
	GetKeyValuePairsCalled     func(address string) (*data.GenericAPIResponse, error)
}

// GetKeyValuePairs -
func (aps *AccountProcessorStub) GetKeyValuePairs(address string) (*data.GenericAPIResponse, error) {
	return aps.GetKeyValuePairsCalled(address)
}

// GetAllESDTTokens -
func (aps *AccountProcessorStub) GetAllESDTTokens(address string) (*data.GenericAPIResponse, error) {
	return aps.GetAllESDTTokensCalled(address)
}

// GetESDTTokenData -
func (aps *AccountProcessorStub) GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error) {
	return aps.GetESDTTokenDataCalled(address, key)
}

// GetAccount --
func (aps *AccountProcessorStub) GetAccount(address string) (*data.Account, error) {
	return aps.GetAccountCalled(address)
}

// GetValueForKey --
func (aps *AccountProcessorStub) GetValueForKey(address string, key string) (string, error) {
	return aps.GetValueForKeyCalled(address, key)
}

// GetShardIDForAddress --
func (aps *AccountProcessorStub) GetShardIDForAddress(address string) (uint32, error) {
	return aps.GetShardIDForAddressCalled(address)
}

// GetTransactions --
func (aps *AccountProcessorStub) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return aps.GetTransactionsCalled(address)
}

// ValidatorStatistics --
func (aps *AccountProcessorStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return aps.ValidatorStatisticsCalled()
}
