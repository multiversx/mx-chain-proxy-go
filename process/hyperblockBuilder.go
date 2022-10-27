package process

import (
	"github.com/ElrondNetwork/elrond-go-core/data/api"
	"github.com/ElrondNetwork/elrond-go-core/data/outport"
	"github.com/ElrondNetwork/elrond-go-core/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type shardBlockWithAlteredAccounts struct {
	shardBlock      *api.Block
	alteredAccounts []*outport.AlteredAccount
}

type hyperblockBuilder struct {
	metaBlock                      *api.Block
	shardBlocksWithAlteredAccounts []*shardBlockWithAlteredAccounts
}

func (builder *hyperblockBuilder) addMetaBlock(metablock *api.Block) {
	builder.metaBlock = metablock
}

func (builder *hyperblockBuilder) addShardBlock(shardBlock *shardBlockWithAlteredAccounts) {
	builder.shardBlocksWithAlteredAccounts = append(builder.shardBlocksWithAlteredAccounts, shardBlock)
}

func (builder *hyperblockBuilder) build(notarizedAtSource bool) data.Hyperblock {
	hyperblock := data.Hyperblock{}
	bunch := newBunchOfTxs()

	bunch.collectTxs(builder.metaBlock, notarizedAtSource)
	for _, block := range builder.shardBlocksWithAlteredAccounts {
		bunch.collectTxs(block.shardBlock, notarizedAtSource)
	}

	hyperblock.Nonce = builder.metaBlock.Nonce
	hyperblock.Round = builder.metaBlock.Round
	hyperblock.Hash = builder.metaBlock.Hash
	hyperblock.Timestamp = builder.metaBlock.Timestamp
	hyperblock.PrevBlockHash = builder.metaBlock.PrevBlockHash
	hyperblock.Epoch = builder.metaBlock.Epoch
	hyperblock.ShardBlocks = builder.buildShardBlocks()
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

func (builder *hyperblockBuilder) buildShardBlocks() []*api.NotarizedBlock {
	notarizedBlocks := make([]*api.NotarizedBlock, 0, len(builder.shardBlocksWithAlteredAccounts))
	for _, shardBlock := range builder.shardBlocksWithAlteredAccounts {
		notarizedBlocks = append(notarizedBlocks, &api.NotarizedBlock{
			Hash:            shardBlock.shardBlock.Hash,
			Nonce:           shardBlock.shardBlock.Nonce,
			Round:           shardBlock.shardBlock.Round,
			Shard:           shardBlock.shardBlock.Shard,
			AlteredAccounts: shardBlock.alteredAccounts,
		})
	}

	return notarizedBlocks
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
