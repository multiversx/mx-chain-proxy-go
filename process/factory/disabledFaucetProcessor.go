package factory

import (
	"errors"
	"math/big"

	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

var errNotEnabled = errors.New("faucet not enabled")

type disabledFaucetProcessor struct {
}

// IsEnabled will return false
func (d *disabledFaucetProcessor) IsEnabled() bool {
	return false
}

// SenderDetailsFromPem will return an error that signals that faucet is not enabled
func (d *disabledFaucetProcessor) SenderDetailsFromPem(receiver string) (crypto.PrivateKey, string, error) {
	return nil, "", errNotEnabled
}

// GenerateTxForSendUserFunds will return an error that signals that faucet is not enabled
func (d *disabledFaucetProcessor) GenerateTxForSendUserFunds(
	senderSk crypto.PrivateKey,
	senderPk string,
	senderNonce uint64,
	receiver string,
	value *big.Int,
	chainID string,
	version uint32,
) (*data.Transaction, error) {
	return nil, errNotEnabled
}
