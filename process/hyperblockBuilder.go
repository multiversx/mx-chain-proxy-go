package process

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type hyperblockBuilder struct {
	metaBlock   *data.Block
	shardBlocks []*data.Block
}

func (builder *hyperblockBuilder) addMetaBlock(metablock *data.Block) {
	builder.metaBlock = metablock
}

func (builder *hyperblockBuilder) addShardBlock(block *data.Block) {
	builder.shardBlocks = append(builder.shardBlocks, block)
}

func (builder *hyperblockBuilder) build() data.Hyperblock {
	hyperblock := data.Hyperblock{}
	bunch := newBunchOfTxs()

	bunch.collectTxs(builder.metaBlock)
	for _, block := range builder.shardBlocks {
		bunch.collectTxs(block)
	}

	hyperblock.Nonce = builder.metaBlock.Nonce
	hyperblock.Round = builder.metaBlock.Round
	hyperblock.Hash = builder.metaBlock.Hash
	hyperblock.PrevBlockHash = builder.metaBlock.PrevBlockHash
	hyperblock.Epoch = builder.metaBlock.Epoch
	hyperblock.ShardBlocks = builder.metaBlock.NotarizedBlocks
	hyperblock.NumTxs = uint32(len(bunch.txs))
	hyperblock.Transactions = bunch.txs

	return hyperblock
}

type bunchOfTxs struct {
	txs []*data.FullTransaction
}

func newBunchOfTxs() *bunchOfTxs {
	return &bunchOfTxs{
		txs: make([]*data.FullTransaction, 0),
	}
}

// In a hyperblock we only return transactions that are fully executed (in both shards).
// Furthermore, we ignore miniblocks of type "PeerBlock"
func (bunch *bunchOfTxs) collectTxs(block *data.Block) {
	for _, miniBlock := range block.MiniBlocks {
		isPeerMiniBlock := miniBlock.Type == "PeerBlock"
		isExecutedOnDestination := miniBlock.DestinationShard == block.Shard

		shouldCollect := !isPeerMiniBlock && isExecutedOnDestination
		if shouldCollect {
			bunch.txs = append(bunch.txs, miniBlock.Transactions...)
		}
	}
}
