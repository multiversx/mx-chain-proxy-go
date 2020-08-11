package process

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/require"
)

func TestHyperblockBuilder(t *testing.T) {
	builder := &hyperblockBuilder{}

	builder.addMetaBlock(&data.Block{Shard: 4294967295, Nonce: 42, MiniBlocks: []*data.MiniBlock{
		{SourceShard: 4294967295, DestinationShard: 0, Transactions: []*data.FullTransaction{
			{Sender: "metachain", Receiver: "alice"},
		}},
		{SourceShard: 4294967295, DestinationShard: 1, Transactions: []*data.FullTransaction{
			{Sender: "metachain", Receiver: "carol"},
		}},
		{SourceShard: 0, DestinationShard: 4294967295, Transactions: []*data.FullTransaction{
			{Sender: "alice", Receiver: "stakingContract"},
		}},
	}, NotarizedBlocks: []*data.NotarizedBlock{
		{Shard: 0, Nonce: 40},
		{Shard: 1, Nonce: 41},
	}})

	builder.addShardBlock(&data.Block{Shard: 0, Nonce: 40, MiniBlocks: []*data.MiniBlock{
		{SourceShard: 0, DestinationShard: 0, Transactions: []*data.FullTransaction{
			{Sender: "alice", Receiver: "bob"},
		}},
		{SourceShard: 0, DestinationShard: 1, Transactions: []*data.FullTransaction{
			{Sender: "alice", Receiver: "carol"},
		}},
	}})

	builder.addShardBlock(&data.Block{Shard: 1, Nonce: 41, MiniBlocks: []*data.MiniBlock{
		{SourceShard: 0, DestinationShard: 1, Transactions: []*data.FullTransaction{
			{Sender: "alice", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Transactions: []*data.FullTransaction{
			{Sender: "carol", Receiver: "carol"},
		}},
	}})

	hyperblock := builder.build()

	require.Equal(t, 42, int(hyperblock.Nonce))
	require.Equal(t, 4, int(hyperblock.NumTxs))
	require.Equal(t, 2, len(hyperblock.ShardBlocks))
	require.Equal(t, []*data.FullTransaction{
		{Sender: "alice", Receiver: "stakingContract"},
		{Sender: "alice", Receiver: "bob"},
		{Sender: "alice", Receiver: "carol"},
		{Sender: "carol", Receiver: "carol"},
	}, hyperblock.Transactions)
}
