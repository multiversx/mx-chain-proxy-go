package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type FaucetProcessorStub struct {
	GenerateTxForSendUserFundsCalled func(receiver string, value *big.Int) (*data.Transaction, error)
}

func (fps *FaucetProcessorStub) GenerateTxForSendUserFunds(receiver string, value *big.Int) (*data.Transaction, error) {
	return fps.GenerateTxForSendUserFundsCalled(receiver, value)
}
