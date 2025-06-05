package process

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/alteredAccount"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestHyperblockBuilderWithFinalizedTxs(t *testing.T) {
	builder := &hyperblockBuilder{}

	builder.addMetaBlock(&api.Block{Shard: 4294967295, Nonce: 42, Timestamp: 12345, TimestampMs: 12345678, MiniBlocks: []*api.MiniBlock{
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

	builder.addShardBlock(&shardBlockWithAlteredAccounts{shardBlock: &api.Block{Hash: "hashShard0", Shard: 0, Nonce: 40, Round: 44, StateRootHash: "rootHashShard0",
		MiniBlocks: []*api.MiniBlock{
			{SourceShard: 0, DestinationShard: 0, Hash: "mbSh0Hash0", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "alice", Receiver: "bob"},
			}},
			{SourceShard: 0, DestinationShard: 1, Hash: "mbSh0Hash1", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "alice", Receiver: "carol"},
			}},
		}}})

	builder.addShardBlock(&shardBlockWithAlteredAccounts{shardBlock: &api.Block{Hash: "hashShard1", Shard: 1, Nonce: 41, Round: 45, StateRootHash: "rootHashShard1",
		MiniBlocks: []*api.MiniBlock{
			{SourceShard: 0, DestinationShard: 1, Hash: "mbSh1Hash0", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "alice", Receiver: "carol"},
			}},
			{SourceShard: 1, DestinationShard: 1, Hash: "mbSh1Hash1", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "carol", Receiver: "carol"},
			}},
			{SourceShard: 1, DestinationShard: 1, Hash: "mbSh1Hash2", Type: "PeerBlock", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "foo", Receiver: "bar"},
			}},
		}}})

	hyperblock := builder.build(false)

	require.Equal(t, api.Hyperblock{
		Nonce:  42,
		NumTxs: 4,
		ShardBlocks: []*api.NotarizedBlock{
			{
				Hash:            "hashShard0",
				Nonce:           40,
				Round:           44,
				Shard:           0,
				RootHash:        "rootHashShard0",
				MiniBlockHashes: []string{"mbSh0Hash0", "mbSh0Hash1"},
			},
			{
				Hash:            "hashShard1",
				Nonce:           41,
				Round:           45,
				Shard:           1,
				RootHash:        "rootHashShard1",
				MiniBlockHashes: []string{"mbSh1Hash0", "mbSh1Hash1", "mbSh1Hash2"},
			},
		},
		Timestamp:   12345,
		TimestampMs: 12345678,
		Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "stakingContract"},
			{Sender: "alice", Receiver: "bob"},
			{Sender: "alice", Receiver: "carol"},
			{Sender: "carol", Receiver: "carol"},
		},
	}, hyperblock)
}

func TestHyperblockBuilderWithNotarizedAtSourceTxs(t *testing.T) {
	builder := &hyperblockBuilder{}

	builder.addMetaBlock(&api.Block{Shard: 4294967295, Nonce: 42, Timestamp: 12345, TimestampMs: 12345678, MiniBlocks: []*api.MiniBlock{
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

	builder.addShardBlock(&shardBlockWithAlteredAccounts{shardBlock: &api.Block{Shard: 0, Nonce: 40, MiniBlocks: []*api.MiniBlock{
		{SourceShard: 0, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "bob"},
		}},
		{SourceShard: 1, DestinationShard: 0, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "carol"},
		}},
	}}})

	builder.addShardBlock(&shardBlockWithAlteredAccounts{shardBlock: &api.Block{Shard: 1, Nonce: 41, MiniBlocks: []*api.MiniBlock{
		{SourceShard: 0, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Transactions: []*transaction.ApiTransactionResult{
			{Sender: "carol", Receiver: "carol"},
		}},
		{SourceShard: 1, DestinationShard: 1, Type: "PeerBlock", Transactions: []*transaction.ApiTransactionResult{
			{Sender: "foo", Receiver: "bar"},
		}},
	}}})

	hyperblock := builder.build(true)

	require.Equal(t, 42, int(hyperblock.Nonce))
	require.Equal(t, 4, int(hyperblock.NumTxs))
	require.Equal(t, 2, len(hyperblock.ShardBlocks))
	require.Equal(t, int64(12345), hyperblock.Timestamp)
	require.Equal(t, int64(12345678), hyperblock.TimestampMs)
	require.Equal(t, []*transaction.ApiTransactionResult{
		{Sender: "metachain", Receiver: "alice"},
		{Sender: "staking-contract", Receiver: "carol"},
		{Sender: "alice", Receiver: "bob"},
		{Sender: "carol", Receiver: "carol"},
	}, hyperblock.Transactions)
}

func TestHyperblockBuilderWithAlteredAccounts(t *testing.T) {
	builder := &hyperblockBuilder{}

	builder.addMetaBlock(&api.Block{
		Shard: core.MetachainShardId,
		Nonce: 42,
		NotarizedBlocks: []*api.NotarizedBlock{
			{Shard: 0, Nonce: 40},
			{Shard: 1, Nonce: 41},
		},
	})

	builder.addShardBlock(&shardBlockWithAlteredAccounts{
		shardBlock: &api.Block{
			Shard: 0,
			Nonce: 40,
		},
		alteredAccounts: []*alteredAccount.AlteredAccount{
			{
				Address: "alice",
				Balance: "100",
			},
		},
	})

	builder.addShardBlock(&shardBlockWithAlteredAccounts{
		shardBlock: &api.Block{
			Shard: 1,
			Nonce: 41,
		},
		alteredAccounts: []*alteredAccount.AlteredAccount{
			{
				Address: "bob",
				Balance: "101",
			},
			{
				Address: "carol",
				Balance: "102",
			},
		},
	})

	hyperblock := builder.build(false)
	require.Equal(t, api.Hyperblock{
		Nonce:        42,
		Transactions: make([]*transaction.ApiTransactionResult, 0),
		ShardBlocks: []*api.NotarizedBlock{
			{
				Shard: 0,
				Nonce: 40,
				AlteredAccounts: []*alteredAccount.AlteredAccount{
					{
						Address: "alice",
						Balance: "100",
					},
				},
				MiniBlockHashes: make([]string, 0),
			},
			{
				Shard: 1,
				Nonce: 41,
				AlteredAccounts: []*alteredAccount.AlteredAccount{
					{
						Address: "bob",
						Balance: "101",
					},
					{
						Address: "carol",
						Balance: "102",
					},
				},
				MiniBlockHashes: make([]string, 0),
			},
		},
	}, hyperblock)
}
