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
func NewCircularQueueObserversProvider(cfg *config.Config) (*CircularQueueObserversProvider, error) {
	bop := &baseObserverProvider{
		mutObservers: sync.RWMutex{},
	}

	err := bop.initObserversMaps(cfg)
	if err != nil {
		return nil, err
	}

	countersMap := getCountersMap(len(bop.observers))
	return &CircularQueueObserversProvider{
		baseObserverProvider:   bop,
		countersMap:            countersMap,
		counterForAllObservers: 0,
	}, nil
}

func getCountersMap(numShards int) map[uint32]uint32 {
	countersMap := make(map[uint32]uint32, numShards)
	for i := 0; i < numShards; i++ {
		countersMap[uint32(i)] = uint32(0)
	}
	return countersMap
}

// GetObserversByShardId will return a slice of observers for the given shard
func (cqop *CircularQueueObserversProvider) GetObserversByShardId(shardId uint32) ([]*data.Observer, error) {
	cqop.mutCounters.RLock()
	defer cqop.mutCounters.RUnlock()
	counterForShard, ok := cqop.countersMap[shardId]
	if !ok {
		return nil, ErrShardNotAvailable
	}

	cqop.mutObservers.RLock()
	observersForShard := cqop.observers[shardId]
	cqop.mutObservers.RUnlock()

	position := int(counterForShard) % len(observersForShard)
	sliceToRet := append(observersForShard[position:], observersForShard[:position]...)
	cqop.countersMap[shardId]++

	return sliceToRet, nil
}

// GetAllObservers will return a slice containing all observers
func (cqop *CircularQueueObserversProvider) GetAllObservers() ([]*data.Observer, error) {
	cqop.mutObservers.Lock()
	defer cqop.mutObservers.Unlock()
	allObservers := cqop.allObservers
	if len(allObservers) == 0 {
		return nil, ErrEmptyObserversList
	}

	cqop.mutCounters.Lock()
	counter := cqop.counterForAllObservers
	cqop.counterForAllObservers++
	cqop.mutCounters.Unlock()

	position := int(counter) % len(allObservers)
	sliceToRet := append(allObservers[position:], allObservers[:position]...)

	return sliceToRet, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (cqop *CircularQueueObserversProvider) IsInterfaceNil() bool {
	return cqop == nil
}
