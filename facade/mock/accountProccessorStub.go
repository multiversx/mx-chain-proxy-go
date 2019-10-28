package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type AccountProcessorStub struct {
	GetAccountCalled func(address string) (*data.Account, error)
}

func (aps *AccountProcessorStub) GetAccount(address string) (*data.Account, error) {
	return aps.GetAccountCalled(address)
}
