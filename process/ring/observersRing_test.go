package ring_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/process/ring"
	"github.com/stretchr/testify/assert"
)

func TestNewRingWrapper_NilSliceShouldErr(t *testing.T) {
	t.Parallel()

	rw, err := ring.NewObserversRing(nil)
	assert.Nil(t, rw)
	assert.Equal(t, ring.ErrInvalidObserversSlice, err)
}

func TestNewRingWrapper_EmptySliceShouldErr(t *testing.T) {
	t.Parallel()

	rw, err := ring.NewObserversRing([]string{})
	assert.Nil(t, rw)
	assert.Equal(t, ring.ErrInvalidObserversSlice, err)
}

func TestNewRingWrapper_ShouldWork(t *testing.T) {
	t.Parallel()

	rw, err := ring.NewObserversRing([]string{"obs"})
	assert.Nil(t, err)
	assert.False(t, check.IfNil(rw))
}

func TestObserversRing_Len(t *testing.T) {
	t.Parallel()

	observers := []string{"obs0", "obs1", "obs2"}
	rw, _ := ring.NewObserversRing(observers)
	assert.Equal(t, len(observers), rw.Len())
}

func TestRingWrapper_Next_AllItemsShouldBeCalled(t *testing.T) {
	t.Parallel()

	mapCalledObservers := make(map[string]bool)
	observers := []string{"obs0", "obs1", "obs2", "obs3", "obs4"}
	rw, _ := ring.NewObserversRing(observers)

	for i := 0; i < len(observers); i++ {
		obs := rw.Next()
		mapCalledObservers[obs] = true
	}

	for _, observer := range observers {
		_, ok := mapCalledObservers[observer]
		if !ok {
			assert.Failf(t, "test failed", "observer %s not called", observer)
		}
	}
}

func TestRingWrapper_Next_ObserversShouldBeCalledAgainAfterQueueIsFinished(t *testing.T) {
	t.Parallel()

	numOfTimesToUseAllObservers := 3
	mapCalledObservers := make(map[string]int)
	observers := []string{"obs0", "obs1", "obs2", "obs3", "obs4"}
	rw, _ := ring.NewObserversRing(observers)

	for i := 0; i < numOfTimesToUseAllObservers*len(observers); i++ {
		obs := rw.Next()
		mapCalledObservers[obs]++
	}

	for _, observer := range observers {
		res, ok := mapCalledObservers[observer]
		if !ok {
			assert.Failf(t, "test failed", "observer %s not called", observer)
		} else {
			if res != numOfTimesToUseAllObservers {
				assert.Failf(t, "test failed",
					"observer %s not called %d times", observer, numOfTimesToUseAllObservers)
			}
		}
	}
}

func TestRingWrapper_NextConcurrentShouldNotFailWithRaceConditionOn(t *testing.T) {
	t.Parallel()

	// if mutex protection is removed from the struct, this test will fail with a race condition

	numOfGoRoutinesToStart := 10
	numOfTimesToCallForEachRoutine := 8
	mapCalledObservers := make(map[string]int)
	mutMap := &sync.RWMutex{}
	observers := []string{"obs0", "obs1", "obs2", "obs3", "obs4"}

	expectedNumOfTimesAnObserverIsCalled := (numOfTimesToCallForEachRoutine * numOfGoRoutinesToStart) / len(observers)

	rw, _ := ring.NewObserversRing(observers)

	for i := 0; i < numOfGoRoutinesToStart; i++ {
		for j := 0; j < numOfTimesToCallForEachRoutine; j++ {
			go func(mutMap *sync.RWMutex, mapCalledObs map[string]int) {
				obs := rw.Next()
				mutMap.Lock()
				mapCalledObs[obs]++
				mutMap.Unlock()
			}(mutMap, mapCalledObservers)
		}
	}
	time.Sleep(500 * time.Millisecond)
	mutMap.RLock()
	for _, res := range mapCalledObservers {
		assert.Equal(t, expectedNumOfTimesAnObserverIsCalled, res)
	}
	mutMap.RUnlock()
}
