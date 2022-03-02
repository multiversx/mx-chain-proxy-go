package observer

import (
	"fmt"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
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
		nodes: map[uint32][]*data.NodeData{
			0: {{Address: "addr1", ShardId: 0}},
			1: {{Address: "addr2", ShardId: 1}},
		},
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
		nodes: map[uint32][]*data.NodeData{
			0: {{Address: "addr1", ShardId: 0}},
			1: {{Address: "addr2", ShardId: 1}},
			2: {{Address: "addr3", ShardId: core.MetachainShardId}},
		},
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
		},
		1: {
			{Address: "shard 1 - id 0"},
			{Address: "shard 1 - id 1"},
			{Address: "shard 1 - id 2"},
			{Address: "shard 1 - id 3"},
		},
		2: {
			{Address: "shard 2 - id 0"},
			{Address: "shard 2 - id 1"},
			{Address: "shard 2 - id 2"},
			{Address: "shard 2 - id 3"},
		},
		core.MetachainShardId: {
			{Address: "shard meta - id 0"},
			{Address: "shard meta - id 1"},
			{Address: "shard meta - id 2"},
			{Address: "shard meta - id 3"},
		},
	}

	expectedOrder := []string{
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

	result := initAllNodesSlice(nodesMap)
	for i, r := range result {
		assert.Equal(t, expectedOrder[i], r.Address)
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
		},
	}

	expectedOrder := []string{
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

	result := initAllNodesSlice(nodesMap)
	for i, r := range result {
		assert.Equal(t, expectedOrder[i], r.Address)
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

	result := initAllNodesSlice(nodesMap)
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

	result := initAllNodesSlice(nodesMap)
	for i, r := range result {
		assert.Equal(t, expectedOrder[i], r.Address)
	}
}

func TestBaseNodeProvider_UpdateNodesBasedOnSyncState(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(6)

	nodesMap := nodesSliceToShardedMap(allNodes)
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		nodes:                 nodesMap,
		allNodes:              allNodes,
	}

	setSyncedStateToNodes(allNodes, false, 1, 2, 4, 5)

	bnp.UpdateNodesBasedOnSyncState(allNodes)

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 1, IsSynced: true},
	}, convertSlice(bnp.allNodes))

	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: false},
		{Address: "addr2", ShardId: 0, IsSynced: false},
		{Address: "addr4", ShardId: 1, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: false},
	}, convertSlice(bnp.outOfSyncNodes))
}

func TestBaseNodeProvider_UpdateNodesBasedOnSyncStateShouldNotRemoveNodeIfNotEnoughLeft(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(4)

	nodesMap := nodesSliceToShardedMap(allNodes)
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		nodes:                 nodesMap,
		allNodes:              allNodes,
	}

	setSyncedStateToNodes(allNodes, false, 1, 3)

	bnp.UpdateNodesBasedOnSyncState(allNodes)

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr2", ShardId: 1, IsSynced: true},
	}, convertSlice(bnp.allNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 1, IsSynced: false},
	}, convertSlice(bnp.outOfSyncNodes))
}

func TestBaseNodeProvider_UpdateNodesBasedOnSyncStateShouldWorkAfterMultipleIterations(t *testing.T) {
	t.Parallel()

	allNodes := prepareNodes(10)

	nodesMap := nodesSliceToShardedMap(allNodes)
	bnp := &baseNodeProvider{
		configurationFilePath: configurationPath,
		nodes:                 nodesMap,
		allNodes:              allNodes,
	}

	setSyncedStateToNodes(allNodes, false, 1, 3, 5, 7, 9)

	bnp.UpdateNodesBasedOnSyncState(allNodes)
	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr6", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
	}, convertSlice(bnp.allNodes))
	require.Equal(t, []data.NodeData{
		{Address: "addr1", ShardId: 0, IsSynced: false},
		{Address: "addr3", ShardId: 0, IsSynced: false},
		{Address: "addr5", ShardId: 1, IsSynced: false},
		{Address: "addr7", ShardId: 1, IsSynced: false},
		{Address: "addr9", ShardId: 1, IsSynced: false},
	}, convertSlice(bnp.outOfSyncNodes))

	allNodes = prepareNodes(10)

	bnp.UpdateNodesBasedOnSyncState(allNodes)

	require.Equal(t, []data.NodeData{
		{Address: "addr0", ShardId: 0, IsSynced: true},
		{Address: "addr2", ShardId: 0, IsSynced: true},
		{Address: "addr4", ShardId: 0, IsSynced: true},
		{Address: "addr6", ShardId: 1, IsSynced: true},
		{Address: "addr8", ShardId: 1, IsSynced: true},
		{Address: "addr1", ShardId: 0, IsSynced: true},
		{Address: "addr3", ShardId: 0, IsSynced: true},
		{Address: "addr5", ShardId: 1, IsSynced: true},
		{Address: "addr7", ShardId: 1, IsSynced: true},
		{Address: "addr9", ShardId: 1, IsSynced: true},
	}, convertSlice(bnp.allNodes))
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

func setSyncedStateToNodes(nodes []*data.NodeData, state bool, indices ...int) {
	for _, idx := range indices {
		nodes[idx].IsSynced = state
	}
}

func convertSlice(nodes []*data.NodeData) []data.NodeData {
	newSlice := make([]data.NodeData, 0, len(nodes))
	for _, node := range nodes {
		newSlice = append(newSlice, *node)
	}

	return newSlice
}
