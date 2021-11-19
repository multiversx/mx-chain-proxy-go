package configuration

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/types"
)

const (
	BlockchainName = "Elrond"
	MainnetChainID = "1"

	MainnetElrondSymbol = "eGLD"
	DevnetElrondSymbol  = "XeGLD"
	NumDecimals         = 18

	// GenesisBlockHashMainnet is const that will keep genesis block hash in hex format
	GenesisBlockHashMainnet = "cd229e4ad2753708e4bab01d7f249affe29441829524c9529e84d51b6d12f2a7"
	DevnetGenesisBlock      = "0000000000000000000000000000000000000000000000000000000000000000"
)

// Configuration is structure used for rosetta provider configuration
type Configuration struct {
	ElrondNetworkConfig    *provider.NetworkConfig
	Network                *types.NetworkIdentifier
	Currency               *types.Currency
	GenesisBlockIdentifier *types.BlockIdentifier
	Peers                  []*types.Peer
}

// Settings is the structure used for rosetta offline config
type Settings struct {
	Offline struct {
		ChainID     string
		MinGasPrice uint64
		MinGasLimit uint64
	} `toml:"OfflineSettings"`
}

// LoadConfiguration will load configuration
func LoadConfiguration(networkConfig *provider.NetworkConfig, generalConfig *config.Config) *Configuration {
	return loadConfig(networkConfig, generalConfig)
}

// LoadOfflineConfig will load the offline configuration for the elrond rosetta server
func LoadOfflineConfig(generalConfig *config.Config, pathToOfflineConfig string) (*Configuration, error) {
	settings := &Settings{}
	err := core.LoadTomlFile(settings, pathToOfflineConfig)
	if err != nil {
		return nil, err
	}

	networkConfig := &provider.NetworkConfig{
		ChainID:     settings.Offline.ChainID,
		MinGasPrice: settings.Offline.MinGasPrice,
		MinGasLimit: settings.Offline.MinGasLimit,
	}

	return loadConfig(networkConfig, generalConfig), nil
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
		// other
		return &Configuration{
			Network: &types.NetworkIdentifier{
				Blockchain: BlockchainName,
				Network:    networkConfig.ChainID,
			},
			Currency: &types.Currency{
				Symbol:   DevnetElrondSymbol,
				Decimals: NumDecimals,
			},
			GenesisBlockIdentifier: &types.BlockIdentifier{
				Index: 1,
				Hash:  DevnetGenesisBlock,
			},
			Peers:               peers,
			ElrondNetworkConfig: networkConfig,
		}
	}
}
