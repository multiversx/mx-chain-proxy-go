package observer

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// SimpleObserversProvider will handle the providing of observers in a simple way, in the order in which they were
// provided in the config file.
type SimpleObserversProvider struct {
	*baseObserverProvider
}

// NewSimpleObserversProvider will return a new instance of SimpleObserverProvider
func NewSimpleObserversProvider(cfg *config.Config) (*SimpleObserversProvider, error) {
	bop := &baseObserverProvider{
		mutObservers: sync.RWMutex{},
	}

	err := bop.initObserversMaps(cfg)
	if err != nil {
		return nil, err
	}

	return &SimpleObserversProvider{
		baseObserverProvider: bop,
	}, nil
}

// GetObserversByShardId will return a slice of observers for the given shard
func (sop *SimpleObserversProvider) GetObserversByShardId(shardId uint32) ([]*data.Observer, error) {
	sop.mutObservers.RLock()
	defer sop.mutObservers.RUnlock()

	observersForShard, ok := sop.observers[shardId]
	if !ok {
		return nil, ErrShardNotAvailable
	}

	return observersForShard, nil
}

// GetAllObservers will return a slice containing all observers
func (sop *SimpleObserversProvider) GetAllObservers() ([]*data.Observer, error) {
	sop.mutObservers.RLock()
	defer sop.mutObservers.RUnlock()

	if len(sop.allObservers) == 0 {
		return nil, ErrEmptyObserversList
	}

	return sop.allObservers, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (sop *SimpleObserversProvider) IsInterfaceNil() bool {
	return sop == nil
}
