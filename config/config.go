package config

// GeneralSettingsConfig will hold the general settings for a node
type GeneralSettingsConfig struct {
	ServerPort          int
	CfgFileReadInterval int
}

// ObserverConfig will hold data to access an observer
type ObserverConfig struct {
	ShardId uint32
	Address string
}

// Config will hold the whole config file's data
type Config struct {
	GeneralSettings GeneralSettingsConfig
	Observers       []*ObserverConfig
}
