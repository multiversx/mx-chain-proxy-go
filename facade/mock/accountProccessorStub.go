package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type AccountProcessorStub struct {
	GetAccountCalled          func(address string) (*data.Account, error)
	ValidatorStatisticsCalled func() (map[string]*data.ValidatorApiResponse, error)
}

func (aps *AccountProcessorStub) GetAccount(address string) (*data.Account, error) {
	return aps.GetAccountCalled(address)
}

func (aps *AccountProcessorStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return aps.ValidatorStatisticsCalled()
}
