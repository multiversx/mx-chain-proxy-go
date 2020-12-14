package factory

import (
	"errors"
	"math/big"
	"net/http"

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
func (d *disabledFaucetProcessor) SenderDetailsFromPem(_ string) (crypto.PrivateKey, string, int, error) {
	return nil, "", http.StatusInternalServerError, errNotEnabled
}

// GenerateTxForSendUserFunds will return an error that signals that faucet is not enabled
func (d *disabledFaucetProcessor) GenerateTxForSendUserFunds(
	_ crypto.PrivateKey,
	_ string,
	_ uint64,
	_ string,
	_ *big.Int,
	_ string,
	_ uint32,
) (*data.Transaction, error) {
	return nil, errNotEnabled
}
