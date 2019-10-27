package config

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GeneralSettingsConfig will hold the general settings for a node
type GeneralSettingsConfig struct {
	ServerPort  int
	FaucetValue *big.Int
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings GeneralSettingsConfig
	Observers       []*data.Observer
}
