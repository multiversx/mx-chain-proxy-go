package holder

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestNodesHolder_GetLastSyncedNodes(t *testing.T) {
	t.Parallel()

	syncedNodes := []*data.NodeData{{Address: "addr0", ShardId: core.MetachainShardId}, {Address: "addr1", ShardId: 0}}
	fallbackNodes := []*data.NodeData{{Address: "fallback-addr0", ShardId: core.MetachainShardId}, {Address: "fallback-addr1", ShardId: 0}}
	shardIds := []uint32{0, core.MetachainShardId}

	nodesHolder, err := NewNodesHolder(syncedNodes, fallbackNodes, shardIds, data.AvailabilityAll)
	require.NoError(t, err)

	require.Equal(t, syncedNodes, nodesHolder.GetSyncedNodes())
	require.Equal(t, fallbackNodes, nodesHolder.GetSyncedFallbackNodes())
	require.Empty(t, nodesHolder.GetOutOfSyncFallbackNodes())
	require.Empty(t, nodesHolder.GetOutOfSyncNodes())
	require.Empty(t, nodesHolder.GetLastSyncedNodes())
}

func TestComputeSyncAndOutOfSyncNodes(t *testing.T) {
	t.Parallel()

	t.Run("all nodes synced", testComputeSyncedAndOutOfSyncNodesAllNodesSynced)
	t.Run("enough synced nodes", testComputeSyncedAndOutOfSyncNodesEnoughSyncedObservers)
	t.Run("all nodes are out of sync", testComputeSyncedAndOutOfSyncNodesAllNodesNotSynced)
	t.Run("invalid config - no node", testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeAtAll)
	t.Run("invalid config - no node in a shard", testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeInAShard)
	t.Run("snapshotless nodes should work with no node in a shard", testSnapshotlessNodesShouldWorkIfNoNodeInShardExists)
	t.Run("edge case - address should not exist in both sync and not-synced lists", testEdgeCaseAddressShouldNotExistInBothLists)
}

func testComputeSyncedAndOutOfSyncNodesAllNodesSynced(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, 1}
	input := []*data.NodeData{
		{Address: "0", ShardId: 0, IsSynced: true},
		{Address: "1", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "2", ShardId: 1, IsSynced: true},
		{Address: "3", ShardId: 1, IsSynced: true, IsFallback: true},
	}

	synced, syncedFb, notSynced, _ := computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityAll)
	require.Equal(t, []*data.NodeData{
		{Address: "0", ShardId: 0, IsSynced: true},
		{Address: "2", ShardId: 1, IsSynced: true},
	}, synced)
	require.Equal(t, []*data.NodeData{
		{Address: "1", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "3", ShardId: 1, IsSynced: true, IsFallback: true},
	}, syncedFb)
	require.Empty(t, notSynced)
}

func testComputeSyncedAndOutOfSyncNodesEnoughSyncedObservers(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, 1}
	input := []*data.NodeData{
		{Address: "0", ShardId: 0, IsSynced: true},
		{Address: "1", ShardId: 0, IsSynced: false},
		{Address: "2", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "3", ShardId: 1, IsSynced: true},
		{Address: "4", ShardId: 1, IsSynced: false},
		{Address: "5", ShardId: 1, IsSynced: true, IsFallback: true},
	}

	synced, syncedFb, notSynced, _ := computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityAll)
	require.Equal(t, []*data.NodeData{
		{Address: "0", ShardId: 0, IsSynced: true},
		{Address: "3", ShardId: 1, IsSynced: true},
	}, synced)
	require.Equal(t, []*data.NodeData{
		{Address: "2", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "5", ShardId: 1, IsSynced: true, IsFallback: true},
	}, syncedFb)
	require.Equal(t, []*data.NodeData{
		{Address: "1", ShardId: 0, IsSynced: false},
		{Address: "4", ShardId: 1, IsSynced: false},
	}, notSynced)
}

func testComputeSyncedAndOutOfSyncNodesAllNodesNotSynced(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, 1}
	input := []*data.NodeData{
		{Address: "0", ShardId: 0, IsSynced: false},
		{Address: "1", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "2", ShardId: 1, IsSynced: false},
		{Address: "3", ShardId: 1, IsSynced: false, IsFallback: true},
	}

	synced, syncedFb, notSynced, _ := computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityAll)
	require.Equal(t, []*data.NodeData{}, synced)
	require.Equal(t, []*data.NodeData{}, syncedFb)
	require.Equal(t, input, notSynced)
}

func testEdgeCaseAddressShouldNotExistInBothLists(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(10)

	nodesMap := nodesSliceToShardedMap(allNodes)
	nh := &nodesHolder{
		shardIDs:    getSortedShardIDsSlice(nodesMap),
		syncedNodes: allNodes,
	}

	setSyncedStateToNodes(allNodes, false, 1, 3, 5, 7, 9)

	nh.UpdateNodes(allNodes)
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr6", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(nh.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: false},
		{Address: "addr7", ShardId: 1, IsSynced: false},
		{Address: "addr9", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(nh.outOfSyncNodes))
	require.False(t, slicesHaveCommonObjects(nh.syncedNodes, nh.outOfSyncNodes))

	allNodes = prepareNodes(10)

	nh.UpdateNodes(allNodes)

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr1", ShardId: 0, IsSynced: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr5", ShardId: 1, IsSynced: true},
		{Address: "addr6", ShardId: 1, IsSynced: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(nh.syncedNodes))
	require.False(t, slicesHaveCommonObjects(nh.syncedNodes, nh.outOfSyncNodes))
}

func testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeAtAll(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, 1}
	var input []*data.NodeData
	synced, syncedFb, notSynced, err := computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityAll)
	require.Error(t, err)
	require.Nil(t, synced)
	require.Nil(t, syncedFb)
	require.Nil(t, notSynced)

	// no node in one shard
	shardIDs = []uint32{0, 1}
	input = []*data.NodeData{
		{
			Address: "0", ShardId: 0, IsSynced: true,
		},
	}
	synced, syncedFb, notSynced, err = computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityAll)
	require.True(t, errors.Is(err, errWrongConfiguration))
	require.Nil(t, synced)
	require.Nil(t, syncedFb)
	require.Nil(t, notSynced)
}

func testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeInAShard(t *testing.T) {
	t.Parallel()

	// no node in one shard
	shardIDs := []uint32{0, 1}
	input := []*data.NodeData{
		{
			Address: "0", ShardId: 0, IsSynced: true,
		},
	}
	synced, syncedFb, notSynced, err := computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityAll)
	require.True(t, errors.Is(err, errWrongConfiguration))
	require.Nil(t, synced)
	require.Nil(t, syncedFb)
	require.Nil(t, notSynced)
}

func testSnapshotlessNodesShouldWorkIfNoNodeInShardExists(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, core.MetachainShardId}
	input := []*data.NodeData{
		{
			Address: "m", ShardId: core.MetachainShardId, IsSynced: true, IsSnapshotless: true,
		},
	}
	synced, syncedFb, notSynced, err := computeSyncedAndOutOfSyncNodes(input, shardIDs, data.AvailabilityRecent)
	require.NoError(t, err)
	require.Empty(t, notSynced)
	require.Empty(t, syncedFb)
	require.Equal(t, input, synced)
}

func slicesHaveCommonObjects(firstSlice []*data.NodeData, secondSlice []*data.NodeData) bool {
	nodeDataToStr := func(nd *data.NodeData) string {
		return fmt.Sprintf("%s%d", nd.Address, nd.ShardId)
	}
	firstSliceItems := make(map[string]struct{})
	for _, el := range firstSlice {
		firstSliceItems[nodeDataToStr(el)] = struct{}{}
	}

	for _, el := range secondSlice {
		nodeDataStr := nodeDataToStr(el)
		_, found := firstSliceItems[nodeDataStr]
		if found {
			return true
		}
	}

	return false
}

func prepareNodes(count int) []*data.NodeData {
	nodes := make([]*data.NodeData, 0, count)
	for i := 0; i < count; i++ {
		shardID := uint32(0)
		if i >= count/2 {
			shardID = 1
		}
		nodes = append(nodes, &data.NodeData{
			ShardId:  shardID,
			Address:  fmt.Sprintf("addr%d", i),
			IsSynced: true,
		})
	}

	return nodes
}

func getSortedShardIDsSlice(nodesOnShards map[uint32][]*data.NodeData) []uint32 {
	shardIDs := make([]uint32, 0)
	for shardID := range nodesOnShards {
		shardIDs = append(shardIDs, shardID)
	}
	sort.SliceStable(shardIDs, func(i, j int) bool {
		return shardIDs[i] < shardIDs[j]
	})

	return shardIDs
}

func setSyncedStateToNodes(nodes []*data.NodeData, state bool, indices ...int) {
	for _, idx := range indices {
		nodes[idx].IsSynced = state
	}
}

func convertAndSortSlice(nodes []*data.NodeData) []data.NodeData {
	newSlice := make([]data.NodeData, 0, len(nodes))
	for _, node := range nodes {
		newSlice = append(newSlice, *node)
	}

	sort.Slice(newSlice, func(i, j int) bool {
		return newSlice[i].Address < newSlice[j].Address
	})

	return newSlice
}
