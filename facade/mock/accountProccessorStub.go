package mock

import "github.com/ElrondNetwork/elrond-proxy-go/data"

type AccountProcessorStub struct {
	GetAccountCalled              func(address string) (*data.Account, error)
	PublicKeyFromPrivateKeyCalled func(privateKeyHex string) (string, error)
}

func (aps *AccountProcessorStub) GetAccount(address string) (*data.Account, error) {
	return aps.GetAccountCalled(address)
}

func (aps *AccountProcessorStub) PublicKeyFromPrivateKey(privateKeyHex string) (string, error) {
	return aps.PublicKeyFromPrivateKeyCalled(privateKeyHex)
}
