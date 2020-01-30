package observer

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type baseObserverProvider struct {
	mutObservers sync.RWMutex
	observers    map[uint32][]*data.Observer
	allObservers []*data.Observer
}

func (bop *baseObserverProvider) initObserversMaps(cfg *config.Config) error {
	if cfg == nil {
		return ErrNilConfig
	}
	if len(cfg.Observers) == 0 {
		return ErrEmptyObserversList
	}

	newObservers := make(map[uint32][]*data.Observer)
	newAllObservers := make([]*data.Observer, 0)
	for _, observer := range cfg.Observers {
		shardId := observer.ShardId
		newObservers[shardId] = append(newObservers[shardId], observer)
		newAllObservers = append(newAllObservers, observer)
	}

	bop.mutObservers.Lock()
	bop.observers = newObservers
	bop.allObservers = newAllObservers
	bop.mutObservers.Unlock()

	return nil
}
