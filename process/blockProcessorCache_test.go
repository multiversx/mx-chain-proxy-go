package process

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	facadeMock "github.com/multiversx/mx-chain-proxy-go/facade/mock"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func getObjectFromCache[T cacheableBlock](c TimedCache, scope string, hash string, nonce *uint64, opts interface{}) T {
	if hash != "" {
		hashKey, _ := makeHashCacheKey(scope, hash, opts)
		return getObjFromCache[T](c, hashKey)
	}

	nonceKey, _ := makeNonceCacheKey(scope, *nonce, opts)
	return getObjFromCache[T](c, nonceKey)
}

func TestBlockProcessorCache(t *testing.T) {
	t.Parallel()

	mockCache := facadeMock.NewTimedCacheMock()
	bp, _ := NewBlockProcessor(&mock.ProcessorStub{}, mockCache)

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

	// Some basic checks that the cache is empty
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

	// Cache scope1:blockApi1:opts1
	bp.cacheObject(blockApi1, scope1, opts1)
	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope2, // wrong scope
		hashBlock1,
		nil,
		opts1,
	))
	require.Nil(t, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		hashBlock1,
		nil,
		opts2, // wrong opts
	))

	require.Equal(t, blockApi1, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		hashBlock1, // found cached blockApi1 object by hash
		nil,
		opts1,
	))
	require.Equal(t, blockApi1, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope1,
		"",
		&nonceBlock1, // found cached blockApi1 object by nonce
		opts1,
	))

	// Cache scope2:blockApi2:opts1
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
		hashBlock2, // found cached blockApi2 object by hash
		nil,
		opts2,
	))
	require.Equal(t, blockApi2, getObjectFromCache[*data.BlockApiResponse](
		bp.cache,
		scope2,
		"",
		&nonceBlock2, // found cached blockApi2 object by nonce
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

	// Cache same hyperBlock object from two different routes:
	// 1. scopeHyperBlock:hyperBlock:opts1
	// 2. scopeHyperBlock:hyperBlock:opts2
	bp.cacheObject(hyperBlock, scopeHyperBlock, opts1)
	bp.cacheObject(hyperBlock, scopeHyperBlock, opts2)

	// found cached hyperBlock object by hyperBlockHash + opts2
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		hyperBlockHash,
		nil,
		opts2,
	))
	// found cached hyperBlock object by hyperBlockNonce + opts2
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		"",
		&hyperBlockNonce,
		opts2,
	))
	// found cached hyperBlock object by hyperBlockHash + opts1
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		hyperBlockHash,
		nil,
		opts1,
	))
	// found cached hyperBlock object by hyperBlockNonce + opts1
	require.Equal(t, hyperBlock, getObjectFromCache[*data.HyperblockApiResponse](
		bp.cache,
		scopeHyperBlock,
		"",
		&hyperBlockNonce,
		opts1,
	))

	opts1Str, _ := json.Marshal(&opts1)
	opts2Str, _ := json.Marshal(&opts2)
	require.Len(t, mockCache.Cache, 12)

	expectedObjKeys := []string{
		fmt.Sprintf("%s:%s|opts:%s", scope1, hashBlock1, opts1Str),              // blockApi1
		fmt.Sprintf("%s:%s|opts:%s", scope2, hashBlock2, opts2Str),              // blockApi2
		fmt.Sprintf("%s:%s|opts:%s", scopeHyperBlock, hyperBlockHash, opts1Str), // hyperBlock
		fmt.Sprintf("%s:%s|opts:%s", scopeHyperBlock, hyperBlockHash, opts2Str), // hyperBlock
	}

	require.Equal(t, mockCache.Cache[expectedObjKeys[0]], blockApi1)
	require.Equal(t, mockCache.Cache[expectedObjKeys[1]], blockApi2)
	require.Equal(t, mockCache.Cache[expectedObjKeys[2]], hyperBlock)
	require.Equal(t, mockCache.Cache[expectedObjKeys[3]], hyperBlock)

	expectedLookUpKeys := []string{
		fmt.Sprintf("%s:nonce:%d|opts:%s", scope1, nonceBlock1, string(opts1Str)), // blockApi1
		fmt.Sprintf("%s:hash:%s|opts:%s", scope1, hashBlock1, string(opts1Str)),   // blockApi1

		fmt.Sprintf("%s:nonce:%d|opts:%s", scope2, nonceBlock2, string(opts2Str)), // blockApi2
		fmt.Sprintf("%s:hash:%s|opts:%s", scope2, hashBlock2, string(opts2Str)),   // blockApi2

		fmt.Sprintf("%s:nonce:%d|opts:%s", scopeHyperBlock, hyperBlockNonce, string(opts1Str)), // hyperBlock
		fmt.Sprintf("%s:hash:%s|opts:%s", scopeHyperBlock, hyperBlockHash, string(opts1Str)),   // hyperBlock

		fmt.Sprintf("%s:nonce:%d|opts:%s", scopeHyperBlock, hyperBlockNonce, string(opts2Str)), // hyperBlock
		fmt.Sprintf("%s:hash:%s|opts:%s", scopeHyperBlock, hyperBlockHash, string(opts2Str)),   // hyperBlock
	}

	require.Equal(t, mockCache.Cache[expectedLookUpKeys[0]], []byte(expectedObjKeys[0]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[1]], []byte(expectedObjKeys[0]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[2]], []byte(expectedObjKeys[1]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[3]], []byte(expectedObjKeys[1]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[4]], []byte(expectedObjKeys[2]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[5]], []byte(expectedObjKeys[2]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[6]], []byte(expectedObjKeys[3]))
	require.Equal(t, mockCache.Cache[expectedLookUpKeys[7]], []byte(expectedObjKeys[3]))
}
