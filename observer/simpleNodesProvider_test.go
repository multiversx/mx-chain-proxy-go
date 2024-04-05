package observer

import (
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/config"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleObserversProvider_EmptyObserversListShouldErr(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.Observers = make([]*data.NodeData, 0)
	sop, err := NewSimpleNodesProvider(cfg.Observers, "path")
	assert.Nil(t, sop)
	assert.Equal(t, ErrEmptyObserversList, err)
}

func TestNewSimpleObserversProvider_ShouldWork(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	sop, err := NewSimpleNodesProvider(cfg.Observers, "path")
	assert.Nil(t, err)
	assert.False(t, check.IfNil(sop))
}

func TestSimpleObserversProvider_GetObserversByShardIdShouldErrBecauseInvalidShardId(t *testing.T) {
	t.Parallel()

	invalidShardId := uint32(37)
	cfg := getDummyConfig()
	cqop, _ := NewSimpleNodesProvider(cfg.Observers, "path")

	res, err := cqop.GetNodesByShardId(invalidShardId, "")
	assert.Nil(t, res)
	assert.Equal(t, ErrShardNotAvailable, err)
}

func TestSimpleObserversProvider_GetObserversByShardIdShouldWork(t *testing.T) {
	t.Parallel()

	shardId := uint32(0)
	cfg := getDummyConfig()
	cqop, _ := NewSimpleNodesProvider(cfg.Observers, "path")

	res, err := cqop.GetNodesByShardId(shardId, "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func TestSimpleObserversProvider_GetAllObserversShouldWork(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cqop, _ := NewSimpleNodesProvider(cfg.Observers, "path")

	res, _ := cqop.GetAllNodes("")
	assert.Equal(t, 2, len(res))
}

func TestSimpleObserversProvider_GetObserversByShardId_ConcurrentSafe(t *testing.T) {
	shardId0 := uint32(0)
	shardId1 := uint32(1)
	numOfGoRoutinesToStart := 10
	numOfTimesToCallForEachRoutine := 6
	mapCalledObservers := make(map[string]int)
	mutMap := &sync.RWMutex{}
	var observers []*data.NodeData
	observers = []*data.NodeData{
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

	// only the first elements in the slice will be used, so we expect that only the observers at index 0 in slices
	// will be called
	expectedNumOfTimesAnObserverIsCalled := numOfTimesToCallForEachRoutine * numOfGoRoutinesToStart

	sop, _ := NewSimpleNodesProvider(cfg.Observers, "path")

	for i := 0; i < numOfGoRoutinesToStart; i++ {
		for j := 0; j < numOfTimesToCallForEachRoutine; j++ {
			go func(mutMap *sync.RWMutex, mapCalledObs map[string]int) {
				obsSh0, _ := sop.GetNodesByShardId(shardId0, "")
				obsSh1, _ := sop.GetNodesByShardId(shardId1, "")
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

func TestSimpleObserversProvider_GetAllObservers_ConcurrentSafe(t *testing.T) {
	shardId0 := uint32(0)
	shardId1 := uint32(1)
	numOfGoRoutinesToStart := 10
	numOfTimesToCallForEachRoutine := 6
	mapCalledObservers := make(map[string]int)
	mutMap := &sync.RWMutex{}
	var observers []*data.NodeData
	observers = []*data.NodeData{
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

	// only the first element in the slice will be used, so we expect that only the observer at index 0 in slice
	// will be called
	expectedNumOfTimesAnObserverIsCalled := numOfTimesToCallForEachRoutine * numOfGoRoutinesToStart

	sop, _ := NewSimpleNodesProvider(cfg.Observers, "path")

	for i := 0; i < numOfGoRoutinesToStart; i++ {
		for j := 0; j < numOfTimesToCallForEachRoutine; j++ {
			go func(mutMap *sync.RWMutex, mapCalledObs map[string]int) {
				obs, _ := sop.GetAllNodes("")
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
