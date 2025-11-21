package process

import (
	"encoding/json"
	"fmt"

	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type hyperBlockCacheHandler interface {
	cacheHyperblock(resp *data.HyperblockApiResponse, opts common.HyperblockQueryOptions)
	getHyperblockFromCache(hash string, nonce *uint64, opts common.HyperblockQueryOptions) *data.HyperblockApiResponse
}

func (bp *BlockProcessor) cacheHyperblock(resp *data.HyperblockApiResponse, opts common.HyperblockQueryOptions) {
	hashKey := resp.Data.Hyperblock.Hash
	optsStr, _ := json.Marshal(opts)
	objKey := []byte(hashKey + string(optsStr))

	// Store object
	bp.cache.Put(objKey, resp, 0)

	// Store nonce + hash lookup keys
	bp.cache.Put(getHashCacheKey(resp.Data.Hyperblock.Hash, opts), objKey, 0)
	bp.cache.Put(getNonceCacheKey(resp.Data.Hyperblock.Nonce, opts), objKey, 0)
}

func (bp *BlockProcessor) getHyperblockFromCache(hash string, nonce *uint64, opts common.HyperblockQueryOptions) *data.HyperblockApiResponse {
	var key interface{}
	if hash != "" {
		key, _ = bp.cache.Get(getHashCacheKey(hash, opts))
	} else if nonce != nil {
		key, _ = bp.cache.Get(getNonceCacheKey(*nonce, opts))
	}

	if key != nil {
		val, ok := bp.cache.Get(key.([]byte))
		if ok {
			return val.(*data.HyperblockApiResponse)
		}
	}

	return nil
}

func getHashCacheKey(hash string, opts common.HyperblockQueryOptions) []byte {
	optBytes, _ := json.Marshal(opts)
	return []byte(fmt.Sprintf("hash:%s|opts:%s", hash, string(optBytes)))
}

func getNonceCacheKey(nonce uint64, opts common.HyperblockQueryOptions) []byte {
	optBytes, _ := json.Marshal(opts)
	return []byte(fmt.Sprintf("nonce:%d|opts:%s", nonce, string(optBytes)))
}
