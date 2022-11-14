package factory

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

var log = logger.GetOrCreate("process/factory")

// CreateFaucetProcessor will return the faucet processor needed for current settings
func CreateFaucetProcessor(
	baseProc Processor,
	shardCoordinator common.Coordinator,
	defaultFaucetValue *big.Int,
	pubKeyConverter core.PubkeyConverter,
	pemFileLocation string,
) (facade.FaucetProcessor, error) {
	if defaultFaucetValue.Cmp(big.NewInt(0)) == 0 {
		log.Info("faucet is disabled")
		return &disabledFaucetProcessor{}, nil
	}

	log.Info("faucet is enabled", "pem file location", pemFileLocation)
	privKeysLoader, err := faucet.NewPrivateKeysLoader(shardCoordinator, pemFileLocation, pubKeyConverter)
	if err != nil {
		return nil, err
	}

	return process.NewFaucetProcessor(baseProc, privKeysLoader, defaultFaucetValue, pubKeyConverter)
}
