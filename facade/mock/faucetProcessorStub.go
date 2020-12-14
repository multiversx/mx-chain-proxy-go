package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type FaucetProcessorStub struct {
	IsEnabledCalled                  func() bool
	GenerateTxForSendUserFundsCalled func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64,
		receiver string, value *big.Int, chainID string, version uint32) (*data.Transaction, error)
	SenderDetailsFromPemCalled func(receiver string) (crypto.PrivateKey, string, int, error)
}

func (fps *FaucetProcessorStub) IsEnabled() bool {
	if fps.IsEnabledCalled != nil {
		return fps.IsEnabledCalled()
	}

	return true
}

func (fps *FaucetProcessorStub) SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, int, error) {
	return fps.SenderDetailsFromPemCalled(receiver)
}

func (fps *FaucetProcessorStub) GenerateTxForSendUserFunds(
	senderSk crypto.PrivateKey,
	senderPk string,
	senderNonce uint64,
	receiver string,
	value *big.Int,
	chainID string,
	version uint32,
) (*data.Transaction, error) {
	return fps.GenerateTxForSendUserFundsCalled(senderSk, senderPk, senderNonce, receiver, value, chainID, version)
}
