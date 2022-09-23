package observer

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// path to a configuration file that contains 3 observers for 3 shards (one per shard). the same thing for
// full history observers
const configurationPath = "testdata/config.toml"

func TestBaseNodeProvider_ReloadNodesDifferentNumberOfNewShard(t *testing.T) {
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		shardIds:              []uint32{0, 1},
	}

	response := bnp.ReloadNodes(data.Observer)
	require.False(t, response.OkRequest)
	require.Contains(t, response.Error, "different number of shards")
}

func TestBaseNodeProvider_ReloadNodesConfigurationFileNotFound(t *testing.T) {
	bnp := &baseNodeProvider{
		configurationFilePath: "wrong config path",
	}

	response := bnp.ReloadNodes(data.Observer)
	require.Contains(t, response.Error, "path")
}

func TestBaseNodeProvider_ReloadNodesShouldWork(t *testing.T) {
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		shardIds:              []uint32{0, 1, core.MetachainShardId},
	}

	response := bnp.ReloadNodes(data.Observer)
	require.True(t, response.OkRequest)
	require.Empty(t, response.Error)
}

func TestBaseNodeProvider_prepareReloadResponseMessage(t *testing.T) {
	addr0, addr1, addr2 := "addr0", "addr1", "addr2"
	newNodes := map[uint32][]*data.NodeData{
		0: {
			{Address: addr0},
		},
		1: {
			{Address: addr1},
		},
		37: {
			{Address: addr2},
		},
	}

	response := prepareReloadResponseMessage(newNodes)
	require.Contains(t, response, addr0)
	require.Contains(t, response, addr1)
	require.Contains(t, response, addr2)
}

func TestInitAllNodesSlice_BalancesNumObserversDistribution(t *testing.T) {
	t.Parallel()

	nodesMap := map[uint32][]*data.NodeData{
		0: {
			{Address: "shard 0 - id 0"},
			{Address: "shard 0 - id 1"},
			{Address: "shard 0 - id 2"},
			{Address: "shard 0 - id 3"},
			{Address: "shard 0 - id 4", IsFallback: true},
		},
		1: {
			{Address: "shard 1 - id 0"},
			{Address: "shard 1 - id 1"},
			{Address: "shard 1 - id 2"},
			{Address: "shard 1 - id 3"},
			{Address: "shard 1 - id 4", IsFallback: true},
		},
		2: {
			{Address: "shard 2 - id 0"},
			{Address: "shard 2 - id 1"},
			{Address: "shard 2 - id 2"},
			{Address: "shard 2 - id 3"},
			{Address: "shard 2 - id 4", IsFallback: true},
		},
		core.MetachainShardId: {
			{Address: "shard meta - id 0"},
			{Address: "shard meta - id 1"},
			{Address: "shard meta - id 2"},
			{Address: "shard meta - id 3"},
			{Address: "shard meta - id 4", IsFallback: true},
		},
	}

	expectedSyncedOrder := []string{
		"shard 0 - id 0",
		"shard 1 - id 0",
		"shard 2 - id 0",
		"shard meta - id 0",
		"shard 0 - id 1",
		"shard 1 - id 1",
		"shard 2 - id 1",
		"shard meta - id 1",
		"shard 0 - id 2",
		"shard 1 - id 2",
		"shard 2 - id 2",
		"shard meta - id 2",
		"shard 0 - id 3",
		"shard 1 - id 3",
		"shard 2 - id 3",
		"shard meta - id 3",
	}

	resultSynced, resultFallback := initAllNodesSlice(nodesMap)
	for i, r := range resultSynced {
		assert.Equal(t, expectedSyncedOrder[i], r.Address)
	}

	expectedFallbackOrder := []string{
		"shard 0 - id 4",
		"shard 1 - id 4",
		"shard 2 - id 4",
		"shard meta - id 4",
	}

	for i, r := range resultFallback {
		assert.Equal(t, expectedFallbackOrder[i], r.Address)
	}
}

func TestInitAllNodesSlice_UnbalancedNumObserversDistribution(t *testing.T) {
	t.Parallel()

	nodesMap := map[uint32][]*data.NodeData{
		0: {
			{Address: "shard 0 - id 0"},
			{Address: "shard 0 - id 1"},
			{Address: "shard 0 - id 2"},
		},
		1: {
			{Address: "shard 1 - id 0"},
			{Address: "shard 1 - id 1"},
			{Address: "shard 1 - id 2"},
			{Address: "shard 1 - id 3"},
		},
		2: {
			{Address: "shard 2 - id 0"},
		},
		core.MetachainShardId: {
			{Address: "shard meta - id 0"},
			{Address: "shard meta - id 1"},
			{Address: "shard meta - id 2"},
			{Address: "shard meta - id 3"},
			{Address: "shard meta - id 4"},
			{Address: "shard meta - id 5", IsFallback: true},
		},
	}

	expectedSyncedOrder := []string{
		"shard 0 - id 0",
		"shard 1 - id 0",
		"shard 2 - id 0",
		"shard meta - id 0",
		"shard 0 - id 1",
		"shard 1 - id 1",
		"shard meta - id 1",
		"shard 0 - id 2",
		"shard 1 - id 2",
		"shard meta - id 2",
		"shard 1 - id 3",
		"shard meta - id 3",
		"shard meta - id 4",
	}

	resultSynced, resultFallback := initAllNodesSlice(nodesMap)
	for i, r := range resultSynced {
		assert.Equal(t, expectedSyncedOrder[i], r.Address)
	}

	expectedFallbackOrder := []string{
		"shard meta - id 5",
	}
	for i, r := range resultFallback {
		assert.Equal(t, expectedFallbackOrder[i], r.Address)
	}
}

func TestInitAllNodesSlice_EmptyObserversSliceForAShardShouldStillWork(t *testing.T) {
	t.Parallel()

	nodesMap := map[uint32][]*data.NodeData{
		0: {
			{Address: "shard 0 - id 0"},
		},
		1: {}, // empty - possible after a config error
		2: {
			{Address: "shard 2 - id 0"},
		},
		core.MetachainShardId: {
			{Address: "shard meta - id 0"},
			{Address: "shard meta - id 1"},
		},
	}

	expectedOrder := []string{
		"shard 0 - id 0",
		"shard 2 - id 0",
		"shard meta - id 0",
		"shard meta - id 1",
	}

	result, _ := initAllNodesSlice(nodesMap)
	for i, r := range result {
		assert.Equal(t, expectedOrder[i], r.Address)
	}
}

func TestInitAllNodesSlice_SingleShardShouldWork(t *testing.T) {
	t.Parallel()

	nodesMap := map[uint32][]*data.NodeData{
		0: {
			{Address: "shard 0 - id 0"},
		},
	}

	expectedOrder := []string{
		"shard 0 - id 0",
	}

	result, _ := initAllNodesSlice(nodesMap)
	for i, r := range result {
		assert.Equal(t, expectedOrder[i], r.Address)
	}
}

func TestBaseNodeProvider_UpdateNodesBasedOnSyncState(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(8)
	setFallbackNodes(allNodes, 0, 1, 4, 5)

	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
	}
	_ = bnp.initNodes(allNodes)

	nodesCopy := copyNodes(allNodes)
	setSyncedStateToNodes(nodesCopy, false, 1, 2, 5, 6)

	bnp.UpdateNodesBasedOnSyncState(nodesCopy)

	require.Equal(t, []data.NodeData{
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "addr4", ShardId: 1, IsSynced: true, IsFallback: true},
	}, convertAndSortSlice(bnp.syncedFallbackNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr6", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
	}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
}

func TestBaseNodeProvider_UpdateNodesBasedOnSyncStateShouldNotRemoveNodeIfNotEnoughLeft(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(4)

	nodesMap := nodesSliceToShardedMap(allNodes)
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		shardIds:              getSortedShardIDsSlice(nodesMap),
		syncedNodes:           allNodes,
		lastSyncedNodes:       map[uint32]*data.NodeData{},
	}

	nodesCopy := copyNodes(allNodes)
	setSyncedStateToNodes(nodesCopy, false, 0, 2)

	bnp.UpdateNodesBasedOnSyncState(nodesCopy)

	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false},
		{Address: "addr2", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))

	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 1, 3)

	bnp.UpdateNodesBasedOnSyncState(nodesCopy)

	require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false},
		{Address: "addr1", ShardId: 0, IsSynced: false},
		{Address: "addr2", ShardId: 1, IsSynced: false},
		{Address: "addr3", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.Equal(t, data.NodeData{Address: "addr1", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0])
	require.Equal(t, data.NodeData{Address: "addr3", ShardId: 1, IsSynced: false}, *bnp.lastSyncedNodes[1])
}

func TestBaseNodeProvider_UpdateNodesBasedOnSyncStateShouldWorkAfterMultipleIterations(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(10)
	setFallbackNodes(allNodes, 0, 1, 5, 6)

	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
	}
	_ = bnp.initNodes(allNodes)

	// sync all nodes, call update 3 times same state
	nodesCopy := copyNodes(allNodes)
	setFallbackNodes(nodesCopy, 0, 1, 5, 6)
	for i := 0; i < 3; i++ {
		syncAllNodesAndCheck(t, nodesCopy, bnp)
	}

	// unsync some nodes, call Update 3 times same state
	nodesCopy = copyNodes(allNodes)
	setSyncedStateToNodes(nodesCopy, false, 1, 3, 5, 7, 9)
	for i := 0; i < 3; i++ {
		bnp.UpdateNodesBasedOnSyncState(nodesCopy)
		require.Equal(t, []data.NodeData{
			{Address: "addr2", ShardId: 0, IsSynced: true},
			{Address: "addr4", ShardId: 0, IsSynced: true},
			{Address: "addr8", ShardId: 1, IsSynced: true},
		}, convertAndSortSlice(bnp.syncedNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		}, convertAndSortSlice(bnp.syncedFallbackNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr7", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.outOfSyncNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
		}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
		require.Equal(t, 0, len(bnp.lastSyncedNodes))

		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr2", ShardId: 0, IsSynced: true},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: true},
			{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr7", ShardId: 1, IsSynced: false},
			{Address: "addr8", ShardId: 1, IsSynced: true},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
	}

	// sync all nodes, call update 3 times same state
	nodesCopy = prepareNodes(10)
	setFallbackNodes(nodesCopy, 0, 1, 5, 6)
	for i := 0; i < 3; i++ {
		syncAllNodesAndCheck(t, nodesCopy, bnp)
	}

	// unsync all regular observers, call update 3 times same state
	nodesCopy = prepareNodes(10)
	setFallbackNodes(nodesCopy, 0, 1, 5, 6)
	setSyncedStateToNodes(nodesCopy, false, 2, 3, 4, 7, 8, 9)
	for i := 0; i < 3; i++ {
		bnp.UpdateNodesBasedOnSyncState(nodesCopy)
		require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.syncedNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		}, convertAndSortSlice(bnp.syncedFallbackNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr2", ShardId: 0, IsSynced: false},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr7", ShardId: 1, IsSynced: false},
			{Address: "addr8", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.outOfSyncNodes))
		require.Equal(t, data.NodeData{Address: "addr4", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0])
		require.Equal(t, data.NodeData{Address: "addr9", ShardId: 1, IsSynced: false}, *bnp.lastSyncedNodes[1])

		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr2", ShardId: 0, IsSynced: false},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr7", ShardId: 1, IsSynced: false},
			{Address: "addr8", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
	}

	// unsync fallbacks as well -> should keep the last regular ones, call update 3 times same state
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 0, 1, 5, 6)
	for i := 0; i < 3; i++ {
		bnp.UpdateNodesBasedOnSyncState(nodesCopy)
		require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.syncedNodes))
		require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.syncedFallbackNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr2", ShardId: 0, IsSynced: false},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr7", ShardId: 1, IsSynced: false},
			{Address: "addr8", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.outOfSyncNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: false, IsFallback: true},
		}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
		require.Equal(t, data.NodeData{Address: "addr4", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0])
		require.Equal(t, data.NodeData{Address: "addr9", ShardId: 1, IsSynced: false}, *bnp.lastSyncedNodes[1])

		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr2", ShardId: 0, IsSynced: false},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr7", ShardId: 1, IsSynced: false},
			{Address: "addr8", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
	}

	// bring one node on each shard back, call update 3 times same state
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, true, 2, 7)
	for i := 0; i < 3; i++ {
		bnp.UpdateNodesBasedOnSyncState(nodesCopy)
		require.Equal(t, []data.NodeData{
			{Address: "addr2", ShardId: 0, IsSynced: true},
			{Address: "addr7", ShardId: 1, IsSynced: true},
		}, convertAndSortSlice(bnp.syncedNodes))
		require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.syncedFallbackNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr8", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.outOfSyncNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: false, IsFallback: true},
		}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
		require.Equal(t, 0, len(bnp.lastSyncedNodes))

		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
			{Address: "addr2", ShardId: 0, IsSynced: true},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr5", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: false, IsFallback: true},
			{Address: "addr7", ShardId: 1, IsSynced: true},
			{Address: "addr8", ShardId: 1, IsSynced: false},
			{Address: "addr9", ShardId: 1, IsSynced: false},
		}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
	}

	// sync all nodes
	nodesCopy = prepareNodes(10)
	setFallbackNodes(nodesCopy, 0, 1, 5, 6)
	syncAllNodesAndCheck(t, nodesCopy, bnp)

	// unsync fallbacks from shard 0
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 0, 1)
	bnp.UpdateNodesBasedOnSyncState(nodesCopy)
	require.Equal(t, []data.NodeData{
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
	}, convertAndSortSlice(bnp.syncedFallbackNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
	}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
	require.Equal(t, 0, len(bnp.lastSyncedNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))

	// unsync all regular observers
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 2, 3, 4, 7, 8, 9)
	bnp.UpdateNodesBasedOnSyncState(nodesCopy)
	require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
	}, convertAndSortSlice(bnp.syncedFallbackNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 0, IsSynced: false},
		{Address: "addr7", ShardId: 1, IsSynced: false},
		{Address: "addr8", ShardId: 1, IsSynced: false},
		{Address: "addr9", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
	}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
	require.Equal(t, data.NodeData{Address: "addr4", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0])
	require.Equal(t, data.NodeData{Address: "addr9", ShardId: 1, IsSynced: false}, *bnp.lastSyncedNodes[1])

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 0, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr7", ShardId: 1, IsSynced: false},
		{Address: "addr8", ShardId: 1, IsSynced: false},
		{Address: "addr9", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))

	// sync all nodes
	nodesCopy = prepareNodes(10)
	setFallbackNodes(nodesCopy, 0, 1, 5, 6)
	syncAllNodesAndCheck(t, nodesCopy, bnp)

	// unsync all observers from shard 0
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 0, 1, 2, 3, 4)
	bnp.UpdateNodesBasedOnSyncState(nodesCopy)
	require.Equal(t, []data.NodeData{
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
	}, convertAndSortSlice(bnp.syncedFallbackNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 0, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
	}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
	require.Equal(t, data.NodeData{Address: "addr4", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0])
	require.Equal(t, 1, len(bnp.lastSyncedNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 0, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))

	// sync all nodes
	nodesCopy = prepareNodes(10)
	setFallbackNodes(nodesCopy, 0, 1, 5, 6)
	syncAllNodesAndCheck(t, nodesCopy, bnp)

	// unsync last regular observer from shard 0, call update 3 times
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 4)
	for i := 0; i < 3; i++ {
		bnp.UpdateNodesBasedOnSyncState(nodesCopy)
		require.Equal(t, []data.NodeData{
			{Address: "addr2", ShardId: 0, IsSynced: true},
			{Address: "addr3", ShardId: 0, IsSynced: true},
			{Address: "addr7", ShardId: 1, IsSynced: true},
			{Address: "addr8", ShardId: 1, IsSynced: true},
			{Address: "addr9", ShardId: 1, IsSynced: true},
		}, convertAndSortSlice(bnp.syncedNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		}, convertAndSortSlice(bnp.syncedFallbackNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr4", ShardId: 0, IsSynced: false},
		}, convertAndSortSlice(bnp.outOfSyncNodes))
		require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
		require.Equal(t, 0, len(bnp.lastSyncedNodes))

		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr2", ShardId: 0, IsSynced: true},
			{Address: "addr3", ShardId: 0, IsSynced: true},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr7", ShardId: 1, IsSynced: true},
			{Address: "addr8", ShardId: 1, IsSynced: true},
			{Address: "addr9", ShardId: 1, IsSynced: true},
		}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
	}

	// unsync all regular observers from shard 0, call update 3 times
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 2, 3)
	for i := 0; i < 3; i++ {
		bnp.UpdateNodesBasedOnSyncState(nodesCopy)
		require.Equal(t, []data.NodeData{
			{Address: "addr7", ShardId: 1, IsSynced: true},
			{Address: "addr8", ShardId: 1, IsSynced: true},
			{Address: "addr9", ShardId: 1, IsSynced: true},
		}, convertAndSortSlice(bnp.syncedNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		}, convertAndSortSlice(bnp.syncedFallbackNodes))
		require.Equal(t, []data.NodeData{
			{Address: "addr2", ShardId: 0, IsSynced: false},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
		}, convertAndSortSlice(bnp.outOfSyncNodes))
		require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
		require.Equal(t, data.NodeData{Address: "addr3", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0]) // last synced
		require.Equal(t, 1, len(bnp.lastSyncedNodes))

		require.Equal(t, []data.NodeData{
			{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
			{Address: "addr2", ShardId: 0, IsSynced: false},
			{Address: "addr3", ShardId: 0, IsSynced: false},
			{Address: "addr4", ShardId: 0, IsSynced: false},
			{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
			{Address: "addr7", ShardId: 1, IsSynced: true},
			{Address: "addr8", ShardId: 1, IsSynced: true},
			{Address: "addr9", ShardId: 1, IsSynced: true},
		}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
	}

	// unsync all fallback observers from shard 0
	nodesCopy = copyNodes(nodesCopy)
	setSyncedStateToNodes(nodesCopy, false, 0, 1)
	bnp.UpdateNodesBasedOnSyncState(nodesCopy)
	require.Equal(t, []data.NodeData{
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
	}, convertAndSortSlice(bnp.syncedFallbackNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 0, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
	}, convertAndSortSlice(bnp.outOfSyncFallbackNodes))
	require.Equal(t, data.NodeData{Address: "addr3", ShardId: 0, IsSynced: false}, *bnp.lastSyncedNodes[0])
	require.Equal(t, 1, len(bnp.lastSyncedNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: false, IsFallback: true},
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 0, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
}

func syncAllNodesAndCheck(t *testing.T, nodes []*data.NodeData, bnp *baseNodeProvider) {
	bnp.UpdateNodesBasedOnSyncState(nodes)
	require.Equal(t, []data.NodeData{
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
	}, convertAndSortSlice(bnp.syncedFallbackNodes))
	require.Equal(t, 0, len(bnp.lastSyncedNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "addr1", ShardId: 0, IsSynced: true, IsFallback: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr5", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr6", ShardId: 1, IsSynced: true, IsFallback: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.GetAllNodesWithSyncState()))
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

func copyNodes(nodes []*data.NodeData) []*data.NodeData {
	nodesCopy := make([]*data.NodeData, len(nodes))
	for i, node := range nodes {
		nodesCopy[i] = &data.NodeData{
			ShardId:    node.ShardId,
			Address:    node.Address,
			IsSynced:   node.IsSynced,
			IsFallback: node.IsFallback,
		}
	}

	return nodesCopy
}

func setSyncedStateToNodes(nodes []*data.NodeData, state bool, indices ...int) {
	for _, idx := range indices {
		nodes[idx].IsSynced = state
	}
}

func setFallbackNodes(nodes []*data.NodeData, indices ...int) {
	for _, idx := range indices {
		nodes[idx].IsFallback = true
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

func TestComputeSyncAndOutOfSyncNodes(t *testing.T) {
	t.Parallel()

	t.Run("all nodes synced", testComputeSyncedAndOutOfSyncNodesAllNodesSynced)
	t.Run("enough synced nodes", testComputeSyncedAndOutOfSyncNodesEnoughSyncedObservers)
	t.Run("all nodes are out of sync", testComputeSyncedAndOutOfSyncNodesAllNodesNotSynced)
	t.Run("invalid config - no node", testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeAtAll)
	t.Run("invalid config - no node in a shard", testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeInAShard)
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

	synced, syncedFb, notSynced, _ := computeSyncedAndOutOfSyncNodes(input, shardIDs)
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

	synced, syncedFb, notSynced, _ := computeSyncedAndOutOfSyncNodes(input, shardIDs)
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

	synced, syncedFb, notSynced, _ := computeSyncedAndOutOfSyncNodes(input, shardIDs)
	require.Equal(t, []*data.NodeData{}, synced)
	require.Equal(t, []*data.NodeData{}, syncedFb)
	require.Equal(t, input, notSynced)
}

func testEdgeCaseAddressShouldNotExistInBothLists(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(10)

	nodesMap := nodesSliceToShardedMap(allNodes)
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		shardIds:              getSortedShardIDsSlice(nodesMap),
		syncedNodes:           allNodes,
	}

	setSyncedStateToNodes(allNodes, false, 1, 3, 5, 7, 9)

	bnp.UpdateNodesBasedOnSyncState(allNodes)
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr6", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
	}, convertAndSortSlice(bnp.syncedNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: false},
		{Address: "addr7", ShardId: 1, IsSynced: false},
		{Address: "addr9", ShardId: 1, IsSynced: false},
	}, convertAndSortSlice(bnp.outOfSyncNodes))
	require.False(t, slicesHaveCommonObjects(bnp.syncedNodes, bnp.outOfSyncNodes))

	allNodes = prepareNodes(10)

	bnp.UpdateNodesBasedOnSyncState(allNodes)

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
	}, convertAndSortSlice(bnp.syncedNodes))
	require.False(t, slicesHaveCommonObjects(bnp.syncedNodes, bnp.outOfSyncNodes))
}

func testComputeSyncedAndOutOfSyncNodesInvalidConfigurationNoNodeAtAll(t *testing.T) {
	t.Parallel()

	shardIDs := []uint32{0, 1}
	var input []*data.NodeData
	synced, syncedFb, notSynced, err := computeSyncedAndOutOfSyncNodes(input, shardIDs)
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
	synced, syncedFb, notSynced, err = computeSyncedAndOutOfSyncNodes(input, shardIDs)
	require.True(t, errors.Is(err, ErrWrongObserversConfiguration))
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
	synced, syncedFb, notSynced, err := computeSyncedAndOutOfSyncNodes(input, shardIDs)
	require.True(t, errors.Is(err, ErrWrongObserversConfiguration))
	require.Nil(t, synced)
	require.Nil(t, syncedFb)
	require.Nil(t, notSynced)
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
