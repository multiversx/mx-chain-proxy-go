package holder

import (
	"fmt"
	"sync"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestNewNodesHolder(t *testing.T) {
	t.Parallel()

	t.Run("empty regular nodes slice - should error", func(t *testing.T) {
		t.Parallel()

		nh, err := NewNodesHolder([]*data.NodeData{}, []*data.NodeData{}, data.AvailabilityAll)
		require.Equal(t, errEmptyNodesList, err)
		require.Nil(t, nh)
	})

	t.Run("empty snapshotless nodes slice - should not error", func(t *testing.T) {
		t.Parallel()

		nh, err := NewNodesHolder([]*data.NodeData{}, []*data.NodeData{}, data.AvailabilityRecent)
		require.NoError(t, err)
		require.NotNil(t, nh)
	})

	t.Run("should work for regular nodes", func(t *testing.T) {
		t.Parallel()

		nh, err := NewNodesHolder([]*data.NodeData{{Address: "addr"}}, []*data.NodeData{}, data.AvailabilityAll)
		require.NoError(t, err)
		require.NotNil(t, nh)
	})

	t.Run("should work for snapshotless nodes", func(t *testing.T) {
		t.Parallel()

		nh, err := NewNodesHolder([]*data.NodeData{{Address: "addr"}}, []*data.NodeData{}, data.AvailabilityRecent)
		require.NoError(t, err)
		require.NotNil(t, nh)
	})
}

func TestNodesHolder_Getters(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, 1, core.MetachainShardId}
	syncedNodes := createTestNodes(6)
	fallbackNodes := createTestNodes(6)
	setPropertyToNodes(fallbackNodes, "fallback", true, 0, 1, 2, 3, 4, 5)

	nh, err := NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityAll)
	require.NoError(t, err)
	require.NotNil(t, nh)

	t.Run("test getters before updating the nodes", func(t *testing.T) {
		for _, shardID := range shardIDs {
			indices := getIndicesOfNodesInShard(syncedNodes, shardID)
			compareNodesBasedOnIndices(t, nh.GetSyncedNodes(shardID), syncedNodes, indices)
		}
		for _, shardID := range shardIDs {
			require.Empty(t, nh.GetOutOfSyncNodes(shardID))
		}
		for _, shardID := range shardIDs {
			indices := getIndicesOfNodesInShard(fallbackNodes, shardID)
			compareNodesBasedOnIndices(t, nh.GetSyncedNodes(shardID), fallbackNodes, indices)
		}
		for _, shardID := range shardIDs {
			require.Empty(t, nh.GetOutOfSyncFallbackNodes(shardID))
		}
	})

	t.Run("test getters after updating the nodes", func(t *testing.T) {
		setPropertyToNodes(syncedNodes, "synced", true, 3, 4, 5)
		setPropertyToNodes(syncedNodes, "synced", false, 0, 1, 2)

		setPropertyToNodes(fallbackNodes, "synced", true, 0, 2, 3, 4, 5)
		setPropertyToNodes(fallbackNodes, "synced", false, 1)
		nh.UpdateNodes(append(syncedNodes, fallbackNodes...))

		// check synced regular nodes
		compareNodesBasedOnIndices(t, nh.GetSyncedNodes(0), syncedNodes, []int{3})
		compareNodesBasedOnIndices(t, nh.GetSyncedNodes(1), syncedNodes, []int{4})
		compareNodesBasedOnIndices(t, nh.GetSyncedNodes(core.MetachainShardId), syncedNodes, []int{5})

		// check out of sync regular nodes
		compareNodesBasedOnIndices(t, nh.GetOutOfSyncNodes(0), syncedNodes, []int{0})
		compareNodesBasedOnIndices(t, nh.GetOutOfSyncNodes(1), syncedNodes, []int{1})
		compareNodesBasedOnIndices(t, nh.GetOutOfSyncNodes(core.MetachainShardId), syncedNodes, []int{2})

		// check synced fallback nodes
		compareNodesBasedOnIndices(t, nh.GetSyncedFallbackNodes(0), syncedNodes, []int{0, 3})
		compareNodesBasedOnIndices(t, nh.GetSyncedFallbackNodes(1), syncedNodes, []int{4})
		compareNodesBasedOnIndices(t, nh.GetSyncedFallbackNodes(core.MetachainShardId), syncedNodes, []int{2, 5})

		// check out of sync fallback nodes
		require.Empty(t, nh.GetOutOfSyncFallbackNodes(0))
		compareNodesBasedOnIndices(t, nh.GetOutOfSyncFallbackNodes(1), syncedNodes, []int{1})
		require.Empty(t, nh.GetOutOfSyncFallbackNodes(core.MetachainShardId))
	})
}

func compareNodesBasedOnIndices(t *testing.T, firstSlice []*data.NodeData, secondSlice []*data.NodeData, indices []int) {
	if len(firstSlice) > len(indices) {
		t.Fail()
	}

	if len(firstSlice) == 0 {
		t.Fail()
	}

	for i, node := range firstSlice {
		indexInSecondSlice := indices[i]
		if indexInSecondSlice > len(secondSlice) {
			t.Fail()
		}
		require.Equal(t, node.Address, secondSlice[indexInSecondSlice].Address)
	}
}

func getIndicesOfNodesInShard(nodes []*data.NodeData, shardID uint32) []int {
	intSlice := make([]int, 0)
	for idx, node := range nodes {
		if node.ShardId != shardID {
			continue
		}

		intSlice = append(intSlice, idx)
	}

	return intSlice
}

func TestNodesHolder_Count(t *testing.T) {
	t.Parallel()

	syncedNodes := createTestNodes(3)
	nh, _ := NewNodesHolder(syncedNodes, syncedNodes, data.AvailabilityAll)
	require.Equal(t, 2*len(syncedNodes), nh.Count())
}

func TestNodesHolder_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var nh *nodesHolder
	require.True(t, nh.IsInterfaceNil())

	nh, _ = NewNodesHolder([]*data.NodeData{{Address: "adr"}}, []*data.NodeData{}, data.AvailabilityAll)
	require.False(t, nh.IsInterfaceNil())
}

func TestNodesHolder_UpdateNodesAvailabilityAll(t *testing.T) {
	t.Parallel()

	syncedNodes := createTestNodes(3)
	setPropertyToNodes(syncedNodes, "synced", true, 0, 1, 2)

	fallbackNodes := createTestNodes(3)
	setPropertyToNodes(fallbackNodes, "synced", true, 0, 1, 2)
	setPropertyToNodes(fallbackNodes, "fallback", true, 0, 1, 2)

	nh, err := NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityAll)
	require.NoError(t, err)

	syncedNodes[0].IsSynced = false
	syncedNodes[1].IsSynced = false
	nh.UpdateNodes(append(syncedNodes, fallbackNodes...))

	require.Equal(t, []*data.NodeData{}, nh.GetSyncedNodes(0))
	require.Equal(t, []*data.NodeData{}, nh.GetSyncedNodes(1))
	require.Equal(t, []*data.NodeData{syncedNodes[2]}, nh.GetSyncedNodes(core.MetachainShardId))

	require.Equal(t, []*data.NodeData{fallbackNodes[0]}, nh.GetSyncedFallbackNodes(0))
	require.Equal(t, []*data.NodeData{fallbackNodes[1]}, nh.GetSyncedFallbackNodes(1))
	require.Equal(t, []*data.NodeData{fallbackNodes[2]}, nh.GetSyncedFallbackNodes(core.MetachainShardId))
}

func TestNodesHolder_UpdateNodesAvailabilityRecent(t *testing.T) {
	t.Parallel()

	syncedNodes := createTestNodes(3)
	setPropertyToNodes(syncedNodes, "synced", true, 0, 1, 2)
	setPropertyToNodes(syncedNodes, "snapshotless", true, 0, 1, 2)

	fallbackNodes := createTestNodes(3)
	setPropertyToNodes(fallbackNodes, "synced", true, 0, 1, 2)
	setPropertyToNodes(fallbackNodes, "fallback", true, 0, 1, 2)
	setPropertyToNodes(fallbackNodes, "snapshotless", true, 0, 1, 2)

	nh, err := NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityRecent)
	require.NoError(t, err)

	syncedNodes[0].IsSynced = false
	syncedNodes[2].IsSnapshotless = false // this will force the nodesHolder to remove it

	nh.UpdateNodes(append(syncedNodes, fallbackNodes...))

	require.Equal(t, []*data.NodeData{}, nh.GetSyncedNodes(0))
	require.Equal(t, []*data.NodeData{syncedNodes[1]}, nh.GetSyncedNodes(1))
	require.Equal(t, []*data.NodeData{}, nh.GetSyncedNodes(core.MetachainShardId))

	require.Equal(t, []*data.NodeData{fallbackNodes[0]}, nh.GetSyncedFallbackNodes(0))
	require.Equal(t, []*data.NodeData{fallbackNodes[1]}, nh.GetSyncedFallbackNodes(1))
	require.Equal(t, []*data.NodeData{fallbackNodes[2]}, nh.GetSyncedFallbackNodes(core.MetachainShardId))
}

func TestNodesHolder_GettersShouldUseCachedValues(t *testing.T) {
	t.Parallel()

	syncedNodes := createTestNodes(1)
	setPropertyToNodes(syncedNodes, "synced", true, 0)
	setPropertyToNodes(syncedNodes, "snapshotless", true, 0)

	fallbackNodes := createTestNodes(0)

	nh, err := NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityRecent)
	require.NoError(t, err)

	// warm the cache
	require.Equal(t, []*data.NodeData{syncedNodes[0]}, nh.GetSyncedNodes(0))

	// check the cache
	require.Equal(t, []*data.NodeData{syncedNodes[0]}, nh.cache[getCacheKey(syncedNodesCache, 0)])

	// put something else in the cache and test it
	newValue := []*data.NodeData{{Address: "test-cached-observer"}}
	nh.cache[getCacheKey(syncedNodesCache, 0)] = newValue
	require.Equal(t, newValue, nh.GetSyncedNodes(0))

	// invalid nodes update - should not invalidate the cache
	nh.UpdateNodes([]*data.NodeData{})
	require.Equal(t, newValue, nh.GetSyncedNodes(0))

	// invalidate the cache by updating the nodes
	nh.UpdateNodes(syncedNodes)
	require.Equal(t, []*data.NodeData{syncedNodes[0]}, nh.GetSyncedNodes(0))
}

func TestNodesHolder_ConcurrentOperations(t *testing.T) {
	t.Parallel()

	syncedNodes := createTestNodes(100)
	fallbackNodes := createTestNodes(100)
	nh, _ := NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityRecent)

	numOperations := 1_000
	wg := sync.WaitGroup{}
	wg.Add(numOperations)
	for i := 0; i < numOperations; i++ {
		go func(index int) {
			switch index {
			case 0:
				nh.UpdateNodes(createTestNodes(100))
			case 1:
				_ = nh.Count()
			case 2:
				_ = nh.GetSyncedFallbackNodes(uint32(index % 3))
			case 3:
				_ = nh.GetOutOfSyncFallbackNodes(uint32(index % 3))
			case 4:
				_ = nh.GetSyncedNodes(uint32(index % 3))
			case 5:
				_ = nh.GetOutOfSyncNodes(uint32(index % 3))
			}
			wg.Done()
		}(i % 6)
	}
	wg.Wait()
}

func createTestNodes(numNodes int) []*data.NodeData {
	getShard := func(index int) uint32 {
		switch index % 3 {
		case 1:
			return 1
		case 2:
			return core.MetachainShardId
		default:
			return 0
		}
	}
	nodes := make([]*data.NodeData, 0, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes = append(nodes, &data.NodeData{
			Address: fmt.Sprintf("https://observer-%d:8080", i),
			ShardId: getShard(i),
		})
	}

	return nodes
}

func setPropertyToNodes(nodes []*data.NodeData, property string, propertyVal bool, indices ...int) {
	switch property {
	case "snapshotless":
		for _, i := range indices {
			nodes[i].IsSnapshotless = propertyVal
		}
	case "fallback":
		for _, i := range indices {
			nodes[i].IsFallback = propertyVal
		}
	case "synced":
		for _, i := range indices {
			nodes[i].IsSynced = propertyVal
		}
	}
}
