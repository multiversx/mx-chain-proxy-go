package observer

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// CircularQueueObserversProvider will handle the providing of observers in a circular queue way, guaranteeing the
// balancing of them
type CircularQueueObserversProvider struct {
	*baseObserverProvider
	countersMap            map[uint32]uint32
	counterForAllObservers uint32
	mutCounters            sync.RWMutex
}

// NewCircularQueueObserversProvider returns a new instance of CircularQueueObserversProvider
func NewCircularQueueObserversProvider(cfg config.Config) (*CircularQueueObserversProvider, error) {
	bop := &baseObserverProvider{}

	err := bop.initObserversMaps(cfg)
	if err != nil {
		return nil, err
	}

	countersMap := make(map[uint32]uint32)
	return &CircularQueueObserversProvider{
		baseObserverProvider:   bop,
		countersMap:            countersMap,
		counterForAllObservers: 0,
	}, nil
}

// GetObserversByShardId will return a slice of observers for the given shard
func (cqop *CircularQueueObserversProvider) GetObserversByShardId(shardId uint32) ([]*data.Observer, error) {
	cqop.mutObservers.Lock()
	defer cqop.mutObservers.Unlock()
	observersForShard := cqop.observers[shardId]

	position := cqop.computeCounterForShard(shardId, uint32(len(observersForShard)))
	sliceToRet := append(observersForShard[position:], observersForShard[:position]...)

	return sliceToRet, nil
}

// GetAllObservers will return a slice containing all observers
func (cqop *CircularQueueObserversProvider) GetAllObservers() []*data.Observer {
	cqop.mutObservers.Lock()
	defer cqop.mutObservers.Unlock()
	allObservers := cqop.allObservers

	position := cqop.computeCounterForAllObservers(uint32(len(allObservers)))
	sliceToRet := append(allObservers[position:], allObservers[:position]...)

	return sliceToRet
}

func (cqop *CircularQueueObserversProvider) computeCounterForShard(shardID uint32, lenObservers uint32) uint32 {
	cqop.mutCounters.Lock()
	defer cqop.mutCounters.Unlock()
	cqop.countersMap[shardID]++
	cqop.countersMap[shardID] %= lenObservers

	return cqop.countersMap[shardID]
}

func (cqop *CircularQueueObserversProvider) computeCounterForAllObservers(lenObservers uint32) uint32 {
	cqop.mutCounters.Lock()
	defer cqop.mutCounters.Unlock()
	cqop.counterForAllObservers++
	cqop.counterForAllObservers %= lenObservers

	return cqop.counterForAllObservers
}

// IsInterfaceNil returns true if there is no value under the interface
func (cqop *CircularQueueObserversProvider) IsInterfaceNil() bool {
	return cqop == nil
}
