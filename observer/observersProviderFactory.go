package observer

import "github.com/ElrondNetwork/elrond-proxy-go/config"

// ObserversProviderFactory handles the creation of an observers provider based on config
type ObserversProviderFactory struct {
	cfg config.Config
}

// NewObserversProviderFactory returns a new instance of ObserversProviderFactory
func NewObserversProviderFactory(cfg config.Config) (*ObserversProviderFactory, error) {
	return &ObserversProviderFactory{
		cfg: cfg,
	}, nil
}

// Create will create and return an object of type ObserversProviderHandler based on a flag
func (opf *ObserversProviderFactory) Create() (ObserversProviderHandler, error) {
	if opf.cfg.GeneralSettings.BalancedObservers {
		return NewCircularQueueObserversProvider(opf.cfg)
	}

	return NewSimpleObserversProvider(opf.cfg)
}
