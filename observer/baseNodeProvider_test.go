package observer

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/holder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// path to a configuration file that contains 3 observers for 3 shards (one per shard). the same thing for
// full history observers
const configurationPath = "testdata/config.toml"

func TestBaseNodeProvider_InvalidNodesConfiguration(t *testing.T) {
	t.Parallel()

	nodes := []*data.NodeData{
		{
			Address:        "addr0",
			ShardId:        0,
			IsSnapshotless: false,
		},
		{
			Address:        "addr1",
			ShardId:        0,
			IsSnapshotless: true,
		},
		{
			Address:        "addr2",
			ShardId:        1,
			IsSnapshotless: true,
		},
		{
			Address:        "addr3",
			ShardId:        1,
			IsSnapshotless: true,
		},
	}

	bnp := baseNodeProvider{}
	err := bnp.initNodes(nodes)
	require.Contains(t, err.Error(), "observers for shard 1 must include at least one historical (non-snapshotless) observer")
}

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

	resultSynced, resultFallback, _, _ := initAllNodesSlice(nodesMap)
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

	resultSynced, resultFallback, _, _ := initAllNodesSlice(nodesMap)
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

	result, _, _, _ := initAllNodesSlice(nodesMap)
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

	result, _, _, _ := initAllNodesSlice(nodesMap)
	for i, r := range result {
		assert.Equal(t, expectedOrder[i], r.Address)
	}
}

func createNodesHolder(nodes []*data.NodeData) NodesHolder {
	holderInstance, _ := holder.NewNodesHolder(nodes, []*data.NodeData{}, "")
	return holderInstance
}

func TestBaseNodeProvider_GetNodesShouldWorkAccordingToTheAvailability(t *testing.T) {
	t.Parallel()

	nodes := []*data.NodeData{
		{
			Address:        "addr0",
			ShardId:        1,
			IsSnapshotless: true,
		},
		{
			Address:        "addr1",
			ShardId:        1,
			IsSnapshotless: false,
		},
	}
	syncedNodes, _, syncedSnapshotless, _ := initAllNodesSlice(map[uint32][]*data.NodeData{1: nodes})
	bnp := &baseNodeProvider{
		regularNodes:      createNodesHolder(syncedNodes),
		snapshotlessNodes: createNodesHolder(syncedSnapshotless),
	}

	returnedNodes, err := bnp.getSyncedNodesForShardUnprotected(1, data.AvailabilityRecent)
	require.NoError(t, err)
	require.Equal(t, "addr0", returnedNodes[0].Address)

	returnedNodes, err = bnp.getSyncedNodesForShardUnprotected(1, data.AvailabilityAll)
	require.NoError(t, err)
	require.Equal(t, "addr1", returnedNodes[0].Address)
}

func TestBaseNodeProvider_getSyncedNodesForShardUnprotected(t *testing.T) {
	getInitialNodes := func() []*data.NodeData {
		return []*data.NodeData{
			{
				Address:        "addr0-snapshotless",
				ShardId:        1,
				IsSnapshotless: true,
				IsSynced:       true,
			},
			{
				Address:        "addr1-regular",
				ShardId:        1,
				IsSnapshotless: false,
				IsSynced:       true,
			},
			{
				Address:    "addr2-fallback",
				ShardId:    1,
				IsFallback: true,
				IsSynced:   true,
			},
		}
	}
	initialNodes := getInitialNodes()
	for _, node := range initialNodes {
		node.IsSynced = true
	}
	syncedNodes, _, syncedSnapshotless, _ := initAllNodesSlice(map[uint32][]*data.NodeData{1: initialNodes})
	bnp := &baseNodeProvider{
		regularNodes:      createNodesHolder(syncedNodes),
		snapshotlessNodes: createNodesHolder(syncedSnapshotless),
		shardIds:          []uint32{1},
	}

	nodes, err := bnp.getSyncedNodesForShardUnprotected(1, data.AvailabilityRecent)
	require.NoError(t, err)
	require.Equal(t, "addr0-snapshotless", nodes[0].Address)

	// make the snapshotless node out of sync - it should go to the regular observer
	updatedNodes := getInitialNodes()
	updatedNodes[0].IsSynced = false
	bnp.UpdateNodesBasedOnSyncState(updatedNodes)

	nodes, err = bnp.getSyncedNodesForShardUnprotected(1, data.AvailabilityRecent)
	require.NoError(t, err)
	require.Equal(t, "addr1-regular", nodes[0].Address)

	// make the regular node out of sync - it should go to the fallback observer
	updatedNodes = getInitialNodes()
	updatedNodes[0].IsSynced = false
	updatedNodes[1].IsSynced = false
	bnp.UpdateNodesBasedOnSyncState(updatedNodes)

	nodes, err = bnp.getSyncedNodesForShardUnprotected(1, data.AvailabilityRecent)
	require.NoError(t, err)
	require.Equal(t, "addr2-fallback", nodes[0].Address)

	// make the fallback node out of sync - it should use an out of sync snapshotless node
	updatedNodes = getInitialNodes()
	updatedNodes[0].IsSynced = false
	updatedNodes[1].IsSynced = false
	updatedNodes[2].IsSynced = false
	bnp.UpdateNodesBasedOnSyncState(updatedNodes)

	nodes, err = bnp.getSyncedNodesForShardUnprotected(1, data.AvailabilityRecent)
	require.NoError(t, err)
	require.Equal(t, "addr0-snapshotless", nodes[0].Address)
	require.False(t, nodes[0].IsSynced)
}
