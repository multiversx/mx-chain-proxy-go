package config

import (
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// GeneralSettingsConfig will hold the general settings for a node
type GeneralSettingsConfig struct {
	ServerPort                               int
	RequestTimeoutSec                        int
	HeartbeatCacheValidityDurationSec        int
	ValStatsCacheValidityDurationSec         int
	EconomicsMetricsCacheValidityDurationSec int
	FaucetValue                              string
	RateLimitWindowDurationSeconds           int
	BalancedObservers                        bool
	BalancedFullHistoryNodes                 bool
	AllowEntireTxPoolFetch                   bool
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings        GeneralSettingsConfig
	AddressPubkeyConverter config.PubkeyConfig
	Marshalizer            config.TypeConfig
	Hasher                 config.TypeConfig
	ApiLogging             ApiLoggingConfig
	Observers              []*data.NodeData
	FullHistoryNodes       []*data.NodeData
}

// ApiLoggingConfig holds the configuration related to API requests logging
type ApiLoggingConfig struct {
	LoggingEnabled          bool
	ThresholdInMicroSeconds int
}

// CredentialsConfig holds the credential pairs
type CredentialsConfig struct {
	Credentials []data.Credential
	Hasher      config.TypeConfig
}
