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

func (bop *baseObserverProvider) initObserversMaps(cfg config.Config) error {
	if len(cfg.Observers) == 0 {
		return ErrEmptyObserversList
	}

	newObservers := make(map[uint32][]*data.Observer)
	for _, observer := range cfg.Observers {
		shardId := observer.ShardId
		newObservers[shardId] = append(newObservers[shardId], observer)
	}

	bop.mutObservers.Lock()
	bop.observers = newObservers
	bop.allObservers = cfg.Observers
	bop.mutObservers.Unlock()

	return nil
}
