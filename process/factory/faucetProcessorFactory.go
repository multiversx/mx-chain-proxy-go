package factory

import (
	"math/big"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/faucet"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

var log = logger.GetOrCreate("process/factory")

// CreateFaucetProcessor will return the faucet processor needed for current settings
func CreateFaucetProcessor(
	ecConf *config.EconomicsConfig,
	baseProc Processor,
	shardCoordinator sharding.Coordinator,
	defaultFaucetValue *big.Int,
	pubKeyConverter core.PubkeyConverter,
	pemFileLocation string,
) (facade.FaucetProcessor, error) {
	if defaultFaucetValue.Cmp(big.NewInt(0)) == 0 {
		log.Info("faucet is offline")
		return &disabledFaucetProcessor{}, nil
	}

	log.Info("faucet is enabled", "pem file location", pemFileLocation)
	privKeysLoader, err := faucet.NewPrivateKeysLoader(shardCoordinator, pemFileLocation, pubKeyConverter)
	if err != nil {
		return nil, err
	}

	return process.NewFaucetProcessor(ecConf, baseProc, privKeysLoader, defaultFaucetValue, pubKeyConverter)
}
