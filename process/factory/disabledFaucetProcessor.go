package factory

import (
	"errors"
	"math/big"

	"github.com/multiversx/mx-chain-crypto-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

var errNotEnabled = errors.New("faucet not enabled")

type disabledFaucetProcessor struct {
}

// IsEnabled will return false
func (d *disabledFaucetProcessor) IsEnabled() bool {
	return false
}

// SenderDetailsFromPem will return an error that signals that faucet is not enabled
func (d *disabledFaucetProcessor) SenderDetailsFromPem(_ string) (crypto.PrivateKey, string, error) {
	return nil, "", errNotEnabled
}

// GenerateTxForSendUserFunds will return an error that signals that faucet is not enabled
func (d *disabledFaucetProcessor) GenerateTxForSendUserFunds(
	_ crypto.PrivateKey,
	_ string,
	_ uint64,
	_ string,
	_ *big.Int,
	_ *data.NetworkConfig,
) (*data.Transaction, error) {
	return nil, errNotEnabled
}
