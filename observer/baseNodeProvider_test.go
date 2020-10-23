package observer

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/assert"
)

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
