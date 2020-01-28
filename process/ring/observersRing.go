package ring

import (
	"container/ring"
	"sync"
)

// ObserversRing represents a wrapper over the container/ring which ensures all observers in a shard are called equal times
type ObserversRing struct {
	ring    *ring.Ring
	mutRing sync.Mutex
}

// NewObserversRing will return a new instance of ObserversRing
func NewObserversRing(observers []string) (*ObserversRing, error) {
	numObservers := len(observers)

	if numObservers <= 0 {
		return nil, ErrInvalidObserversSlice
	}

	r := ring.New(numObservers)
	rw := &ObserversRing{
		ring:    r,
		mutRing: sync.Mutex{},
	}

	rw.addObserversToRing(observers)

	return rw, nil
}

// addObserversToRing will add all the observers in the ring
func (rw *ObserversRing) addObserversToRing(observers []string) {
	for i := 0; i < len(observers); i++ {
		rw.ring.Value = observers[i]
		rw.ring = rw.ring.Next()
	}
}

// Len returns the length of the queue
func (rw *ObserversRing) Len() int {
	rw.mutRing.Lock()
	defer rw.mutRing.Unlock()

	return rw.ring.Len()
}

// Next will return the next observer to use and will switch internally to the next observer
func (rw *ObserversRing) Next() string {
	rw.mutRing.Lock()
	observer := rw.ring.Value.(string)
	rw.ring = rw.ring.Next()
	rw.mutRing.Unlock()

	return observer
}

// IsInterfaceNil returns true if there is no value under the interface
func (rw *ObserversRing) IsInterfaceNil() bool {
	return rw == nil
}
