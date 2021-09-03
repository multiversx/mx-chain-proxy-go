package configuration

// ChainConfig will hold the chain configuration
type ChainConfig struct {
	ChainID     string
	MinGasPrice uint64
	MinGasLimit uint64
}
