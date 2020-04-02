package config

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GeneralSettingsConfig will hold the general settings for a node
type GeneralSettingsConfig struct {
	ServerPort                        int
	RequestTimeoutSec                 int
	HeartbeatCacheValidityDurationSec int
	ValStatsCacheValidityDurationSec  int
	FaucetValue                       string
	BalancedObservers                 bool
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings GeneralSettingsConfig
	Observers       []*data.Observer
}
