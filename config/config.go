package config

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GeneralSettingsConfig will hold the general settings for a node
type GeneralSettingsConfig struct {
	ServerPort  int
	FaucetValue string
	ServerPort                        int
	RequestTimeoutSec                 int
	HeartbeatCacheValidityDurationSec int
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings GeneralSettingsConfig
	Observers       []*data.Observer
}
