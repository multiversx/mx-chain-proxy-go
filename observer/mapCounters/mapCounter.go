package mapCounters

import "sync"

type mapCounter struct {
	positions        map[uint32]uint32
	allNodesCount    uint32
	allNodesPosition uint32
	mut              sync.RWMutex
}

// newMapCounter returns a new instance of a mapCounter
func newMapCounter() *mapCounter {
	return &mapCounter{
		positions:        make(map[uint32]uint32),
		allNodesPosition: 0,
	}
}

func (mc *mapCounter) computePositionForShard(shardID uint32, numNodes uint32) uint32 {
	mc.mut.Lock()
	defer mc.mut.Unlock()

	mc.initShardPositionIfNeededUnprotected(shardID)

	mc.positions[shardID]++
	mc.positions[shardID] %= numNodes

	return mc.positions[shardID]
}

func (mc *mapCounter) computePositionForAllNodes(numNodes uint32) uint32 {
	mc.mut.Lock()
	defer mc.mut.Unlock()

	mc.initAllNodesPositionIfNeededUnprotected(numNodes)

	mc.allNodesPosition++
	mc.allNodesPosition %= numNodes

	return mc.allNodesPosition
}

func (mc *mapCounter) initShardPositionIfNeededUnprotected(shardID uint32) {
	_, shardExists := mc.positions[shardID]
	if !shardExists {
		mc.positions[shardID] = 0
	}
}

func (mc *mapCounter) initAllNodesPositionIfNeededUnprotected(numNodes uint32) {
	if numNodes != mc.allNodesCount {
		mc.allNodesCount = numNodes
		mc.allNodesPosition = 0
	}
}
