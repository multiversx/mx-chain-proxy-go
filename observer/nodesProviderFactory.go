package observer

import (
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/config"
)

var log = logger.GetOrCreate("observer")

// nodesProviderFactory handles the creation of an nodes provider based on config
type nodesProviderFactory struct {
	cfg                   config.Config
	configurationFilePath string
	numberOfShards        uint32
}

// NewNodesProviderFactory returns a new instance of nodesProviderFactory
func NewNodesProviderFactory(cfg config.Config, configurationFilePath string, numberOfShards uint32) (*nodesProviderFactory, error) {
	return &nodesProviderFactory{
		cfg:                   cfg,
		configurationFilePath: configurationFilePath,
		numberOfShards:        numberOfShards,
	}, nil
}

// CreateObservers will create and return an object of type NodesProviderHandler based on a flag
func (npf *nodesProviderFactory) CreateObservers() (NodesProviderHandler, error) {
	if npf.cfg.GeneralSettings.BalancedObservers {
		return NewCircularQueueNodesProvider(
			npf.cfg.Observers,
			npf.configurationFilePath,
			npf.numberOfShards)
	}

	return NewSimpleNodesProvider(
		npf.cfg.Observers,
		npf.configurationFilePath,
		npf.numberOfShards)
}

// CreateFullHistoryNodes will create and return an object of type NodesProviderHandler based on a flag
func (npf *nodesProviderFactory) CreateFullHistoryNodes() (NodesProviderHandler, error) {
	if npf.cfg.GeneralSettings.BalancedFullHistoryNodes {
		nodesProviderHandler, err := NewCircularQueueNodesProvider(
			npf.cfg.FullHistoryNodes,
			npf.configurationFilePath,
			npf.numberOfShards)
		if err != nil {
			return getDisabledFullHistoryNodesProviderIfNeeded(err)
		}

		return nodesProviderHandler, nil
	}

	nodesProviderHandler, err := NewSimpleNodesProvider(
		npf.cfg.FullHistoryNodes,
		npf.configurationFilePath,
		npf.numberOfShards)
	if err != nil {
		return getDisabledFullHistoryNodesProviderIfNeeded(err)
	}

	return nodesProviderHandler, nil
}

func getDisabledFullHistoryNodesProviderIfNeeded(err error) (NodesProviderHandler, error) {
	if err == ErrEmptyObserversList {
		log.Warn("no configuration found for full history nodes. Calls to endpoints specific to full history nodes " +
			"will return an error")
		return NewDisabledNodesProvider("full history nodes not supported"), nil
	}

	return nil, err
}
