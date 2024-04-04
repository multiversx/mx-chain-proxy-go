package mapCounters

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMapCounter(t *testing.T) {
	t.Parallel()

	mc := newMapCounter()
	require.NotNil(t, mc)
	require.NotNil(t, mc.positions)
}

func TestMapCounter_ComputeShardPositionShouldWorkWithDifferentNumOfNodes(t *testing.T) {
	t.Parallel()

	mc := newMapCounter()
	computeShardPosAndAssert(t, mc, 3, 1)
	computeShardPosAndAssert(t, mc, 3, 2)
	computeShardPosAndAssert(t, mc, 3, 0)
	computeShardPosAndAssert(t, mc, 3, 1)
	// change num nodes
	computeShardPosAndAssert(t, mc, 2, 0)
	computeShardPosAndAssert(t, mc, 2, 1)
	computeShardPosAndAssert(t, mc, 2, 0)
	// change num nodes again
	computeShardPosAndAssert(t, mc, 5, 1)
	computeShardPosAndAssert(t, mc, 5, 2)
	computeShardPosAndAssert(t, mc, 5, 3)
}

func TestMapCounter_ComputeShardPositionShouldWorkMultiShard(t *testing.T) {
	t.Parallel()

	mc := newMapCounter()
	computeShardPosAndAssertForShard(t, mc, 0, 3, 1)
	computeShardPosAndAssertForShard(t, mc, 1, 4, 1)

	computeShardPosAndAssertForShard(t, mc, 0, 3, 2)
	computeShardPosAndAssertForShard(t, mc, 1, 4, 2)

	computeShardPosAndAssertForShard(t, mc, 0, 3, 0)
	computeShardPosAndAssertForShard(t, mc, 1, 4, 3)

	computeShardPosAndAssertForShard(t, mc, 0, 3, 1)
	computeShardPosAndAssertForShard(t, mc, 1, 4, 0)

}

func computeShardPosAndAssertForShard(t *testing.T, mc *mapCounter, shardID uint32, numNodes uint32, expectedPos uint32) {
	actualPos := mc.computePositionForShard(shardID, numNodes)
	require.Equal(t, expectedPos, actualPos)
}

func computeShardPosAndAssert(t *testing.T, mc *mapCounter, numNodes uint32, expectedPos uint32) {
	computeShardPosAndAssertForShard(t, mc, 0, numNodes, expectedPos)
}

func TestMapCounter_ComputeAllNodesPosition(t *testing.T) {
	t.Parallel()

	mc := newMapCounter()
	computeAllNodesPosAndAssert(t, mc, 3, 1)
	computeAllNodesPosAndAssert(t, mc, 3, 2)
	computeAllNodesPosAndAssert(t, mc, 3, 0)
	computeAllNodesPosAndAssert(t, mc, 3, 1)
	// change num nodes - should reset
	computeAllNodesPosAndAssert(t, mc, 5, 1)
	computeAllNodesPosAndAssert(t, mc, 5, 2)
	computeAllNodesPosAndAssert(t, mc, 5, 3)
	// change num nodes again - should reset
	computeAllNodesPosAndAssert(t, mc, 2, 1)
	computeAllNodesPosAndAssert(t, mc, 2, 0)
	computeAllNodesPosAndAssert(t, mc, 2, 1)
}

func computeAllNodesPosAndAssert(t *testing.T, mc *mapCounter, numNodes uint32, expectedPos uint32) {
	actualPos := mc.computePositionForAllNodes(numNodes)
	require.Equal(t, expectedPos, actualPos)
}

func TestMapCounter_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	mc := newMapCounter()

	numOperations := 10_000
	wg := sync.WaitGroup{}
	wg.Add(numOperations)
	for i := 0; i < numOperations; i++ {
		go func(idx int) {
			switch idx {
			case 0:
				mc.computePositionForShard(uint32(idx), uint32(10+idx))
			case 1:
				mc.computePositionForAllNodes(uint32(10 + idx))
			}
			wg.Done()
		}(i % 2)
	}

	wg.Wait()
}
