package config

import (
	"github.com/multiversx/mx-chain-proxy-go/data"
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
	NumShardsTimeoutInSec                    int
	TimeBetweenNodesRequestsInSec            int
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings        GeneralSettingsConfig
	AddressPubkeyConverter PubkeyConfig
	Marshalizer            TypeConfig
	Hasher                 TypeConfig
	ApiLogging             ApiLoggingConfig
	Observers              []*data.NodeData
	FullHistoryNodes       []*data.NodeData
}

// TypeConfig will map the string type configuration
type TypeConfig struct {
	Type string
}

// PubkeyConfig will map the public key configuration
type PubkeyConfig struct {
	Length          int
	Type            string
	SignatureLength int
}

// ApiLoggingConfig holds the configuration related to API requests logging
type ApiLoggingConfig struct {
	LoggingEnabled          bool
	ThresholdInMicroSeconds int
}

// CredentialsConfig holds the credential pairs
type CredentialsConfig struct {
	Credentials []data.Credential
	Hasher      TypeConfig
}
