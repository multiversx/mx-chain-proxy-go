package mock

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AccountProcessorStub --
type AccountProcessorStub struct {
	GetAccountCalled          func(address string) (*data.Account, error)
	GetTransactionsCalled     func(address string) ([]data.DatabaseTransaction, error)
	ValidatorStatisticsCalled func() (map[string]*data.ValidatorApiResponse, error)
}

// GetAccount --
func (aps *AccountProcessorStub) GetAccount(address string) (*data.Account, error) {
	return aps.GetAccountCalled(address)
}

// GetTransactions --
func (aps *AccountProcessorStub) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return aps.GetTransactionsCalled(address)
}

// ValidatorStatistics --
func (aps *AccountProcessorStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return aps.ValidatorStatisticsCalled()
}
