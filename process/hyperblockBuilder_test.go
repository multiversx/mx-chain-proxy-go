package process

import (
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestHyperblockBuilderWithFinalizedTxs(t *testing.T) {
	builder := &hyperblockBuilder{}

	builder.addMetaBlock(&api.Block{Shard: 4294967295, Nonce: 42, Timestamp: time.Duration(12345), MiniBlocks: []*api.MiniBlock{
		{SourceShard: 4294967295, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "metachain", Receiver: "alice"},
		}},
		{SourceShard: 4294967295, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "metachain", Receiver: "carol"},
		}},
		{SourceShard: 0, DestinationShard: 4294967295, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "stakingContract"},
		}},
	}, NotarizedBlocks: []*api.NotarizedBlock{
		{Shard: 0, Nonce: 40},
		{Shard: 1, Nonce: 41},
	}})

	builder.addShardBlock(&api.Block{Shard: 0, Nonce: 40, MiniBlocks: []*api.MiniBlock{
		{SourceShard: 0, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "bob"},
		}},
		{SourceShard: 0, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "carol"},
		}},
	}})

	builder.addShardBlock(&api.Block{Shard: 1, Nonce: 41, MiniBlocks: []*api.MiniBlock{
		{SourceShard: 0, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "carol", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Type: "PeerBlock", Transactions: []*transaction.ApiTransactionResult{
			{Sender: "foo", Receiver: "bar"},
		}},
	}})

	hyperblock := builder.build(false)

	require.Equal(t, 42, int(hyperblock.Nonce))
	require.Equal(t, 4, int(hyperblock.NumTxs))
	require.Equal(t, 2, len(hyperblock.ShardBlocks))
	require.Equal(t, time.Duration(12345), hyperblock.Timestamp)
	require.Equal(t, []*transaction.ApiTransactionResult{
		{Sender: "alice", Receiver: "stakingContract"},
		{Sender: "alice", Receiver: "bob"},
		{Sender: "alice", Receiver: "carol"},
		{Sender: "carol", Receiver: "carol"},
	}, hyperblock.Transactions)
}

func TestHyperblockBuilderWithNotarizedAtSourceTxs(t *testing.T) {
	builder := &hyperblockBuilder{}

	builder.addMetaBlock(&api.Block{Shard: 4294967295, Nonce: 42, Timestamp: time.Duration(12345), MiniBlocks: []*api.MiniBlock{
		{SourceShard: 4294967295, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "metachain", Receiver: "alice"},
		}},
		{SourceShard: 4294967295, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "staking-contract", Receiver: "carol"},
		}},
		{SourceShard: 0, DestinationShard: 4294967295, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "stakingContract"},
		}},
	}, NotarizedBlocks: []*api.NotarizedBlock{
		{Shard: 0, Nonce: 40},
		{Shard: 1, Nonce: 41},
	}})

	builder.addShardBlock(&api.Block{Shard: 0, Nonce: 40, MiniBlocks: []*api.MiniBlock{
		{SourceShard: 0, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "bob"},
		}},
		{SourceShard: 1, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "carol"},
		}},
	}})

	builder.addShardBlock(&api.Block{Shard: 1, Nonce: 41, MiniBlocks: []*api.MiniBlock{
		{SourceShard: 0, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "carol", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Type: "PeerBlock", Transactions: []*transaction.ApiTransactionResult{
			{Sender: "foo", Receiver: "bar"},
		}},
	}})

	hyperblock := builder.build(true)

	require.Equal(t, 42, int(hyperblock.Nonce))
	require.Equal(t, 4, int(hyperblock.NumTxs))
	require.Equal(t, 2, len(hyperblock.ShardBlocks))
	require.Equal(t, time.Duration(12345), hyperblock.Timestamp)
	require.Equal(t, []*transaction.ApiTransactionResult{
		{Sender: "metachain", Receiver: "alice"},
		{Sender: "staking-contract", Receiver: "carol"},
		{Sender: "alice", Receiver: "bob"},
		{Sender: "carol", Receiver: "carol"},
	}, hyperblock.Transactions)
}
