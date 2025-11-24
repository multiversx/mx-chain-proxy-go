package process

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestBlockProcessorCache(t *testing.T) {
	t.Parallel()

	bp, _ := NewBlockProcessor(&mock.ProcessorStub{})

	nonceBlock1 := uint64(1)
	nonceBlock2 := uint64(2)
	hashBlock1 := "hashBlock1"
	hashBlock2 := "hashBlock2"
	scope1 := "block:shard=1"
	scope2 := "block:shard=2"
	opts1 := common.BlockQueryOptions{WithTransactions: true}
	opts2 := common.BlockQueryOptions{WithLogs: true}

	blockApi1 := &data.BlockApiResponse{
		Data: data.BlockApiResponsePayload{
			Block: api.Block{
				Nonce: nonceBlock1,
				Hash:  hashBlock1,
			},
		},
	}
	blockApi2 := &data.BlockApiResponse{
		Data: data.BlockApiResponsePayload{
			Block: api.Block{
				Nonce: nonceBlock2,
				Hash:  hashBlock2,
			},
		},
	}

	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		hashBlock1,
		nil,
		opts1,
	))
	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		"",
		&nonceBlock1,
		opts1,
	))

	bp.cacheObject(blockApi1, scope1, opts1)
	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope2,
		hashBlock1,
		nil,
		opts1,
	))
	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		hashBlock1,
		nil,
		opts2,
	))

	require.Equal(t, blockApi1, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		hashBlock1,
		nil,
		opts1,
	))
	require.Equal(t, blockApi1, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		"",
		&nonceBlock1,
		opts1,
	))

	bp.cacheObject(blockApi2, scope2, opts2)
	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope2,
		hashBlock2,
		nil,
		opts1,
	))
	require.Equal(t, blockApi2, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope2,
		hashBlock2,
		nil,
		opts2,
	))
	require.Equal(t, blockApi2, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope2,
		"",
		&nonceBlock2,
		opts2,
	))

	scopeHyperBlock := "hyperBlock"
	hyperBlockNonce := uint64(4)
	hyperBlockHash := "hyperBlockHash"
	hyperBlock := &data.HyperblockApiResponse{
		Data: data.HyperblockApiResponsePayload{
			Hyperblock: api.Hyperblock{
				Hash:  hyperBlockHash,
				Nonce: hyperBlockNonce,
			},
		},
	}

	bp.cacheObject(hyperBlock, scopeHyperBlock, opts1)
	bp.cacheObject(hyperBlock, scopeHyperBlock, opts2)
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		hyperBlockHash,
		nil,
		opts2,
	))
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		"",
		&hyperBlockNonce,
		opts2,
	))
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		hyperBlockHash,
		nil,
		opts1,
	))
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		"",
		&hyperBlockNonce,
		opts1,
	))
}
