package observer

import (
	"sync"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/assert"
)

func getDummyConfig() config.Config {
	return config.Config{
		Observers: []*data.Observer{
			{
				Address: "dummy1",
				ShardId: 0,
			},
			{
				Address: "dummy2",
				ShardId: 1,
			},
		},
	}
}

func TestNewCircularQueueObserverProvider_EmptyObserversListShouldErr(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.Observers = make([]*data.Observer, 0)
	cqop, err := NewCircularQueueObserversProvider(cfg)
	assert.Nil(t, cqop)
	assert.Equal(t, ErrEmptyObserversList, err)
}

func TestNewCircularQueueObserverProvider_ShouldWork(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cqop, err := NewCircularQueueObserversProvider(cfg)
	assert.Nil(t, err)
	assert.False(t, check.IfNil(cqop))
}

func TestCircularQueueObserversProvider_GetObserversByShardIdShouldWork(t *testing.T) {
	t.Parallel()

	shardId := uint32(0)
	cfg := getDummyConfig()
	cqop, _ := NewCircularQueueObserversProvider(cfg)

	res, err := cqop.GetObserversByShardId(shardId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func TestCircularQueueObserversProvider_GetObserversByShardIdShouldBalanceObservers(t *testing.T) {
	t.Parallel()

	shardId := uint32(0)
	cfg := config.Config{
		Observers: []*data.Observer{
			{
				Address: "addr1",
				ShardId: 0,
			},
			{
				Address: "addr2",
				ShardId: 0,
			},
			{
				Address: "addr3",
				ShardId: 0,
			},
		},
	}
	cqop, _ := NewCircularQueueObserversProvider(cfg)

	res1, _ := cqop.GetObserversByShardId(shardId)
	res2, _ := cqop.GetObserversByShardId(shardId)
	assert.NotEqual(t, res1, res2)

	// there are 3 observers. so after 3 steps, the queue should be the same as the original
	_, _ = cqop.GetObserversByShardId(shardId)

	res4, _ := cqop.GetObserversByShardId(shardId)
	assert.Equal(t, res1, res4)
}

func TestCircularQueueObserversProvider_GetAllObserversShouldWork(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cqop, _ := NewCircularQueueObserversProvider(cfg)

	res := cqop.GetAllObservers()
	assert.Equal(t, 2, len(res))
}

func TestCircularQueueObserversProvider_GetAllObserversShouldWorkAndBalanceObservers(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Observers: []*data.Observer{
			{
				Address: "addr1",
				ShardId: 0,
			},
			{
				Address: "addr2",
				ShardId: 0,
			},
			{
				Address: "addr3",
				ShardId: 0,
			},
		},
	}
	cqop, _ := NewCircularQueueObserversProvider(cfg)

	res1 := cqop.GetAllObservers()
	res2 := cqop.GetAllObservers()
	assert.NotEqual(t, res1, res2)

	// there are 3 observers. so after 3 steps, the queue should be the same as the original
	_ = cqop.GetAllObservers()

	res4 := cqop.GetAllObservers()
	assert.Equal(t, res1, res4)
}

func TestCircularQueueObserversProvider_GetAllObservers_ConcurrentSafe(t *testing.T) {
	numOfGoRoutinesToStart := 10
	numOfTimesToCallForEachRoutine := 8
	mapCalledObservers := make(map[string]int)
	mutMap := &sync.RWMutex{}
	var observers []*data.Observer
	observers = []*data.Observer{
		{
			Address: "addr1",
			ShardId: 0,
		},
		{
			Address: "addr2",
			ShardId: 0,
		},
		{
			Address: "addr3",
			ShardId: 0,
		},
		{
			Address: "addr4",
			ShardId: 0,
		},
		{
			Address: "addr5",
			ShardId: 0,
		},
	}
	cfg := config.Config{
		Observers: observers,
	}

	expectedNumOfTimesAnObserverIsCalled := (numOfTimesToCallForEachRoutine * numOfGoRoutinesToStart) / len(observers)

	cqop, _ := NewCircularQueueObserversProvider(cfg)

	for i := 0; i < numOfGoRoutinesToStart; i++ {
		for j := 0; j < numOfTimesToCallForEachRoutine; j++ {
			go func(mutMap *sync.RWMutex, mapCalledObs map[string]int) {
				obs := cqop.GetAllObservers()
				mutMap.Lock()
				mapCalledObs[obs[0].Address]++
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

func TestCircularQueueObserversProvider_GetObserversByShardId_ConcurrentSafe(t *testing.T) {
	shardId0 := uint32(0)
	shardId1 := uint32(1)
	numOfGoRoutinesToStart := 10
	numOfTimesToCallForEachRoutine := 6
	mapCalledObservers := make(map[string]int)
	mutMap := &sync.RWMutex{}
	var observers []*data.Observer
	observers = []*data.Observer{
		{
			Address: "addr1",
			ShardId: shardId0,
		},
		{
			Address: "addr2",
			ShardId: shardId0,
		},
		{
			Address: "addr3",
			ShardId: shardId0,
		},
		{
			Address: "addr4",
			ShardId: shardId1,
		},
		{
			Address: "addr5",
			ShardId: shardId1,
		},
		{
			Address: "addr6",
			ShardId: shardId1,
		},
	}
	cfg := config.Config{
		Observers: observers,
	}

	expectedNumOfTimesAnObserverIsCalled := 2 * ((numOfTimesToCallForEachRoutine * numOfGoRoutinesToStart) / len(observers))

	cqop, _ := NewCircularQueueObserversProvider(cfg)

	for i := 0; i < numOfGoRoutinesToStart; i++ {
		for j := 0; j < numOfTimesToCallForEachRoutine; j++ {
			go func(mutMap *sync.RWMutex, mapCalledObs map[string]int) {
				obsSh0, _ := cqop.GetObserversByShardId(shardId0)
				obsSh1, _ := cqop.GetObserversByShardId(shardId1)
				mutMap.Lock()
				mapCalledObs[obsSh0[0].Address]++
				mapCalledObs[obsSh1[0].Address]++
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
