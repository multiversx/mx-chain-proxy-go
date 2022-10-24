package process

import (
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type hyperblockBuilder struct {
	metaBlock   *api.Block
	shardBlocks []*api.Block
}

func (builder *hyperblockBuilder) addMetaBlock(metablock *api.Block) {
	builder.metaBlock = metablock
}

func (builder *hyperblockBuilder) addShardBlock(block *api.Block) {
	builder.shardBlocks = append(builder.shardBlocks, block)
}

func (builder *hyperblockBuilder) build(notarizedAtSource bool) data.Hyperblock {
	hyperblock := data.Hyperblock{}
	bunch := newBunchOfTxs()

	bunch.collectTxs(builder.metaBlock, notarizedAtSource)
	for _, block := range builder.shardBlocks {
		bunch.collectTxs(block, notarizedAtSource)
	}

	hyperblock.Nonce = builder.metaBlock.Nonce
	hyperblock.Round = builder.metaBlock.Round
	hyperblock.Hash = builder.metaBlock.Hash
	hyperblock.Timestamp = builder.metaBlock.Timestamp
	hyperblock.PrevBlockHash = builder.metaBlock.PrevBlockHash
	hyperblock.Epoch = builder.metaBlock.Epoch
	hyperblock.ShardBlocks = builder.metaBlock.NotarizedBlocks
	hyperblock.NumTxs = uint32(len(bunch.txs))
	hyperblock.Transactions = bunch.txs
	hyperblock.AccumulatedFees = builder.metaBlock.AccumulatedFees
	hyperblock.DeveloperFees = builder.metaBlock.DeveloperFees
	hyperblock.AccumulatedFeesInEpoch = builder.metaBlock.AccumulatedFeesInEpoch
	hyperblock.DeveloperFeesInEpoch = builder.metaBlock.DeveloperFeesInEpoch
	hyperblock.Status = builder.metaBlock.Status
	hyperblock.EpochStartInfo = builder.metaBlock.EpochStartInfo
	hyperblock.EpochStartShardsData = builder.metaBlock.EpochStartShardsData
	hyperblock.StateRootHash = builder.metaBlock.StateRootHash

	return hyperblock
}

type bunchOfTxs struct {
	txs []*transaction.ApiTransactionResult
}

func newBunchOfTxs() *bunchOfTxs {
	return &bunchOfTxs{
		txs: make([]*transaction.ApiTransactionResult, 0),
	}
}

// In a hyperblock we only return transactions that are fully executed (in both shards), if the notarizedAtSource isn't enabled.
// Furthermore, we ignore miniblocks of type "PeerBlock"
func (bunch *bunchOfTxs) collectTxs(block *api.Block, notarizedAtSource bool) {
	if notarizedAtSource {
		bunch.collectNotarizedAtSourceTxs(block)
		return
	}

	bunch.collectFinalizedTxs(block)
}

func (bunch *bunchOfTxs) collectFinalizedTxs(block *api.Block) {
	for _, miniBlock := range block.MiniBlocks {
		isPeerMiniBlock := miniBlock.Type == "PeerBlock"
		isExecutedOnDestination := miniBlock.DestinationShard == block.Shard

		shouldCollect := !isPeerMiniBlock && isExecutedOnDestination
		if shouldCollect {
			bunch.txs = append(bunch.txs, miniBlock.Transactions...)
		}
	}
}

func (bunch *bunchOfTxs) collectNotarizedAtSourceTxs(block *api.Block) {
	for _, miniBlock := range block.MiniBlocks {
		isPeerMiniBlock := miniBlock.Type == "PeerBlock"
		isNotarizedAtSource := miniBlock.SourceShard == block.Shard

		shouldCollect := !isPeerMiniBlock && isNotarizedAtSource
		if shouldCollect {
			bunch.txs = append(bunch.txs, miniBlock.Transactions...)
		}
	}
}
