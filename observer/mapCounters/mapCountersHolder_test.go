package mapCounters

import (
	"sync"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestNewMapCountersHolder(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()
	require.NotNil(t, mch)
	require.Len(t, mch.countersMap, 2)
}

func TestMapCountersHolder_ComputeShardPositionShouldFailDueToInvalidAvailability(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()

	pos, err := mch.ComputeShardPosition("invalid", 0, 10)
	require.Equal(t, errInvalidAvailability, err)
	require.Empty(t, pos)
}

func TestMapCountersHolder_ComputeShardPositionShouldFailDueToZeroNumNodes(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()

	pos, err := mch.ComputeShardPosition(data.AvailabilityAll, 0, 0)
	require.Equal(t, errNumNodesMustBeGreaterThanZero, err)
	require.Empty(t, pos)
}

func TestMapCountersHolder_ComputeShardPositionShouldWorkWhileChangingNumNodes(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()

	calculatePosAndAssert(t, mch, 0, 3, 1)
	calculatePosAndAssert(t, mch, 1, 3, 1)
	calculatePosAndAssert(t, mch, 2, 3, 1)
	calculatePosAndAssert(t, mch, core.MetachainShardId, 3, 1)

	calculatePosAndAssert(t, mch, 0, 2, 0)
	calculatePosAndAssert(t, mch, 1, 2, 0)
	calculatePosAndAssert(t, mch, 2, 2, 0)
	calculatePosAndAssert(t, mch, core.MetachainShardId, 2, 0)

	calculatePosAndAssert(t, mch, 0, 2, 1)
	calculatePosAndAssert(t, mch, 1, 2, 1)
	calculatePosAndAssert(t, mch, 2, 2, 1)
	calculatePosAndAssert(t, mch, core.MetachainShardId, 2, 1)

	calculatePosAndAssert(t, mch, 0, 5, 2)
	calculatePosAndAssert(t, mch, 1, 4, 2)
	calculatePosAndAssert(t, mch, 2, 3, 2)
	calculatePosAndAssert(t, mch, core.MetachainShardId, 2, 0)
}

func TestMapCountersHolder_ComputeShardPositionShouldWorkForMultipleAvailabilities(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()

	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 0, 3, 1)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 0, 3, 1)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 1, 3, 1)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 1, 3, 1)

	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 0, 2, 0)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 0, 3, 2)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 1, 2, 0)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 1, 5, 2)

	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 0, 3, 1)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 0, 3, 0)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 1, 3, 1)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 1, 3, 0)

	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 0, 3, 2)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 0, 3, 1)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, 1, 3, 2)
	calculatePosAndAssertForShard(t, mch, data.AvailabilityAll, 1, 3, 1)
}

func calculatePosAndAssertForShard(
	t *testing.T,
	mch *mapCountersHolder,
	availability data.ObserverDataAvailabilityType,
	shardID uint32,
	numNodes uint32,
	expectedPos uint32,
) {
	pos, err := mch.ComputeShardPosition(availability, shardID, numNodes)
	require.NoError(t, err)
	require.Equal(t, expectedPos, pos)
}

func calculatePosAndAssert(t *testing.T, mch *mapCountersHolder, shardID uint32, numNodes uint32, expectedPos uint32) {
	calculatePosAndAssertForShard(t, mch, data.AvailabilityRecent, shardID, numNodes, expectedPos)
}

func TestMapCountersHolder_ComputeAllNodesPositionShouldFailDueToInvalidAvailaility(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()
	pos, err := mch.ComputeAllNodesPosition("invalid", 10)
	require.Equal(t, errInvalidAvailability, err)
	require.Empty(t, pos)
}

func TestMapCountersHolder_ComputeAllNodesPositionShouldFailDueToZeroNumNodes(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()

	pos, err := mch.ComputeAllNodesPosition(data.AvailabilityAll, 0)
	require.Equal(t, errNumNodesMustBeGreaterThanZero, err)
	require.Empty(t, pos)
}

func TestMapCountersHolder_ComputeAllNodesPositionShouldWork(t *testing.T) {
	t.Parallel()

	mch := NewMapCountersHolder()

	calculateAllNodesPosAndAssert(t, mch, 3, 1)
	calculateAllNodesPosAndAssert(t, mch, 3, 2)
	calculateAllNodesPosAndAssert(t, mch, 3, 0)
	calculateAllNodesPosAndAssert(t, mch, 3, 1)

	calculateAllNodesPosAndAssert(t, mch, 5, 1)
	calculateAllNodesPosAndAssert(t, mch, 5, 2)
	calculateAllNodesPosAndAssert(t, mch, 5, 3)
	calculateAllNodesPosAndAssert(t, mch, 5, 4)
	calculateAllNodesPosAndAssert(t, mch, 5, 0)
	calculateAllNodesPosAndAssert(t, mch, 5, 1)

	calculateAllNodesPosAndAssert(t, mch, 2, 1)
	calculateAllNodesPosAndAssert(t, mch, 2, 0)
}

func calculateAllNodesPosAndAssert(t *testing.T, mch *mapCountersHolder, numNodes uint32, expectedPos uint32) {
	pos, err := mch.ComputeAllNodesPosition(data.AvailabilityRecent, numNodes)
	require.NoError(t, err)
	require.Equal(t, expectedPos, pos)
}

func TestMapCountersHolder_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	numOperations := 10_000
	mch := NewMapCountersHolder()

	wg := sync.WaitGroup{}
	wg.Add(numOperations)

	for i := 0; i < numOperations; i++ {
		go func(idx int) {
			switch idx {
			case 0:
				_, _ = mch.ComputeShardPosition(data.AvailabilityRecent, uint32(idx), uint32(10+idx))
			case 1:
				_, _ = mch.ComputeShardPosition(data.AvailabilityAll, uint32(idx), uint32(10+idx))
			case 2:
				_, _ = mch.ComputeAllNodesPosition(data.AvailabilityRecent, uint32(10+idx))
			case 3:
				_, _ = mch.ComputeAllNodesPosition(data.AvailabilityAll, uint32(10+idx))
			}
		}(i % 2)
	}
}

func TestMapCountersHolder_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var mch *mapCountersHolder
	require.True(t, mch.IsInterfaceNil())

	mch = NewMapCountersHolder()
	require.False(t, mch.IsInterfaceNil())
}
