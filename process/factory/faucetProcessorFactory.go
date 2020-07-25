package factory

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

// CreateFaucetProcessor will return the faucet processor needed for current settings
func CreateFaucetProcessor(
	ecConf *config.EconomicsConfig,
	baseProc Processor,
	privKeysLoader PrivateKeysLoaderHandler,
	defaultFaucetValue *big.Int,
	pubKeyConverter core.PubkeyConverter,
) (facade.FaucetProcessor, error) {
	if defaultFaucetValue.Cmp(big.NewInt(0)) == 0 {
		return &disabledFaucetProcessor{}, nil
	}

	return process.NewFaucetProcessor(ecConf, baseProc, privKeysLoader, defaultFaucetValue, pubKeyConverter)
}
