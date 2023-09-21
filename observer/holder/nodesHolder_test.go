package holder

import (
	"fmt"
	"sync"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestNodesHolder_ConstructorAndGetters(t *testing.T) {
	nh, err := NewNodesHolder([]*data.NodeData{}, []*data.NodeData{}, data.AvailabilityAll)
	require.Equal(t, errEmptyNodesList, err)
	require.True(t, check.IfNil(nh))

	syncedNodes := createTestNodes(3)
	setPropertyToNodes(syncedNodes, "synced", true, 0, 1, 2)

	fallbackNodes := createTestNodes(3)
	setPropertyToNodes(fallbackNodes, "synced", true, 0, 1, 2)
	setPropertyToNodes(fallbackNodes, "fallback", true, 0, 1, 2)

	nh, err = NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityAll)
	require.NoError(t, err)
	require.False(t, nh.IsInterfaceNil())

	require.Equal(t, []*data.NodeData{syncedNodes[0]}, nh.GetSyncedNodes(0))
	require.Equal(t, []*data.NodeData{syncedNodes[1]}, nh.GetSyncedNodes(1))
	require.Equal(t, []*data.NodeData{syncedNodes[2]}, nh.GetSyncedNodes(core.MetachainShardId))

	require.Equal(t, []*data.NodeData{fallbackNodes[0]}, nh.GetSyncedFallbackNodes(0))
	require.Equal(t, []*data.NodeData{fallbackNodes[1]}, nh.GetSyncedFallbackNodes(1))
	require.Equal(t, []*data.NodeData{fallbackNodes[2]}, nh.GetSyncedFallbackNodes(core.MetachainShardId))

	setPropertyToNodes(syncedNodes, "synced", false, 0, 2)
	setPropertyToNodes(fallbackNodes, "synced", false, 1)
	nh.UpdateNodes(append(syncedNodes, fallbackNodes...))
	require.Equal(t, []*data.NodeData{syncedNodes[0]}, nh.GetOutOfSyncNodes(0))
	require.Equal(t, []*data.NodeData{}, nh.GetOutOfSyncNodes(1))
	require.Equal(t, []*data.NodeData{syncedNodes[2]}, nh.GetOutOfSyncNodes(core.MetachainShardId))

	require.Equal(t, []*data.NodeData{}, nh.GetOutOfSyncFallbackNodes(0))
	require.Equal(t, []*data.NodeData{fallbackNodes[1]}, nh.GetOutOfSyncFallbackNodes(1))
	require.Equal(t, []*data.NodeData{}, nh.GetOutOfSyncFallbackNodes(core.MetachainShardId))
}

func TestNodesHolder_UpdateNodesAvailabilityAll(t *testing.T) {
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
	syncedNodes := createTestNodes(100)
	fallbackNodes := createTestNodes(100)
	nh, _ := NewNodesHolder(syncedNodes, fallbackNodes, data.AvailabilityRecent)

	numOperations := 100_000
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
		}(i)
	}
	wg.Wait()
}

func createTestNodes(numNodes int) []*data.NodeData {
	getShard := func(index int) uint32 {
		switch index % 3 {
		case 0:
			return 0
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
