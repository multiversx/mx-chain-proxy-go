package configuration

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	BlockchainName = "Elrond"
	MainnetChainID = "1"

	MainnetElrondSymbol = "eGLD"
	TestnetElrondSymbol = "XeGLD"
	NumDecimals         = 18

	// GenesisBlockHashMainnet is const that will keep genesis block hash in hex format
	GenesisBlockHashMainnet = "cd229e4ad2753708e4bab01d7f249affe29441829524c9529e84d51b6d12f2a7"
	TestnetGenesisBlock     = "0000000000000000000000000000000000000000000000000000000000000000"

	MinGasPrice = uint64(1000000000)
	MinGasLimit = uint64(50000)
)

// Configuration is structure used for rosetta provider configuration
type Configuration struct {
	ElrondNetworkConfig    *provider.NetworkConfig
	Network                *types.NetworkIdentifier
	Currency               *types.Currency
	GenesisBlockIdentifier *types.BlockIdentifier
	Peers                  []*types.Peer
}

// LoadConfiguration will load configuration
func LoadConfiguration(networkConfig *provider.NetworkConfig, generalConfig *config.Config) *Configuration {
	return loadConfig(networkConfig, generalConfig)
}

func LoadOfflineConfig(generalConfig *config.Config) *Configuration {
	networkConfig := &provider.NetworkConfig{
		ChainID:     MainnetChainID,
		MinGasPrice: MinGasPrice,
		MinGasLimit: MinGasLimit,
	}

	return loadConfig(networkConfig, generalConfig)
}

func loadConfig(networkConfig *provider.NetworkConfig, generalConfig *config.Config) *Configuration {
	peers := make([]*types.Peer, len(generalConfig.Observers))
	for idx, observer := range generalConfig.Observers {
		peer := &types.Peer{
			PeerID: hex.EncodeToString([]byte(observer.Address)),
			Metadata: map[string]interface{}{
				"address": observer.Address,
				"shardID": observer.ShardId,
			},
		}
		peers[idx] = peer
	}

	switch networkConfig.ChainID {
	case MainnetChainID:
		return &Configuration{
			Network: &types.NetworkIdentifier{
				Blockchain: BlockchainName,
				Network:    networkConfig.ChainID,
			},
			Currency: &types.Currency{
				Symbol:   MainnetElrondSymbol,
				Decimals: NumDecimals,
			},
			GenesisBlockIdentifier: &types.BlockIdentifier{
				Index: 1,
				Hash:  GenesisBlockHashMainnet,
			},
			Peers:               peers,
			ElrondNetworkConfig: networkConfig,
		}
	default:
		// other testnets
		return &Configuration{
			Network: &types.NetworkIdentifier{
				Blockchain: BlockchainName,
				Network:    networkConfig.ChainID,
			},
			Currency: &types.Currency{
				Symbol:   TestnetElrondSymbol,
				Decimals: NumDecimals,
			},
			GenesisBlockIdentifier: &types.BlockIdentifier{
				Index: 1,
				Hash:  TestnetGenesisBlock,
			},
			Peers:               peers,
			ElrondNetworkConfig: networkConfig,
		}
	}
}
