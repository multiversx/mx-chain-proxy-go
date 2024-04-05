package mock

import (
	"math/big"

	crypto "github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type FaucetProcessorStub struct {
	IsEnabledCalled                  func() bool
	GenerateTxForSendUserFundsCalled func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64,
		receiver string, value *big.Int, networkConfig *data.NetworkConfig) (*data.Transaction, error)
	SenderDetailsFromPemCalled func(receiver string) (crypto.PrivateKey, string, error)
}

func (fps *FaucetProcessorStub) IsEnabled() bool {
	if fps.IsEnabledCalled != nil {
		return fps.IsEnabledCalled()
	}

	return true
}

func (fps *FaucetProcessorStub) SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error) {
	return fps.SenderDetailsFromPemCalled(receiver)
}

func (fps *FaucetProcessorStub) GenerateTxForSendUserFunds(
	senderSk crypto.PrivateKey,
	senderPk string,
	senderNonce uint64,
	receiver string,
	value *big.Int,
	networkConfig *data.NetworkConfig,
) (*data.Transaction, error) {
	return fps.GenerateTxForSendUserFundsCalled(senderSk, senderPk, senderNonce, receiver, value, networkConfig)
}
