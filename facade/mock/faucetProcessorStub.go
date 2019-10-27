package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type FaucetProcessorStub struct {
	GenerateTxForSendUserFundsCalled func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64,
		receiver string, value *big.Int) (*data.Transaction, error)
	SenderDetailsFromPemCalled func(receiver string) (crypto.PrivateKey, string, error)
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
) (*data.Transaction, error) {

	return fps.GenerateTxForSendUserFundsCalled(senderSk, senderPk, senderNonce, receiver, value)
}
