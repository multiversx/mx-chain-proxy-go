package process

import (
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/stretchr/testify/require"
)

func TestHyperblockBuilder(t *testing.T) {
	builder := &HyperblockBuilder{}

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

	builder.addShardBlock(&api.Block{Hash: "hashShard0", Shard: 0, Nonce: 40, Round: 44, StateRootHash: "rootHashShard0",
		MiniBlocks: []*api.MiniBlock{
			{SourceShard: 0, DestinationShard: 0, Hash: "mbSh0Hash0", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "alice", Receiver: "bob"},
			}},
			{SourceShard: 0, DestinationShard: 1, Hash: "mbSh0Hash1", Transactions: []*transaction.ApiTransactionResult{
				{Sender: "alice", Receiver: "carol"},
			}},
		}})

	builder.addShardBlock(&api.Block{Hash: "hashShard1", Shard: 1, Nonce: 41, Round: 45, StateRootHash: "rootHashShard1",
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
		}})

	hyperblock := builder.build()

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
		Timestamp: time.Duration(12345),
		Transactions: []*transaction.ApiTransactionResult{
			{Sender: "alice", Receiver: "stakingContract"},
			{Sender: "alice", Receiver: "bob"},
			{Sender: "alice", Receiver: "carol"},
			{Sender: "carol", Receiver: "carol"},
		},
	}, hyperblock)
}
