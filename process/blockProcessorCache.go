package process

import (
	"encoding/json"
	"fmt"
)

type cacheableBlock interface {
	ID() string
	Hash() string
	Nonce() uint64
}

func makeHashCacheKey(scope string, hash string, opts interface{}) []byte {
	optBytes, _ := json.Marshal(opts)
	return []byte(fmt.Sprintf("%s:hash:%s|opts:%s", scope, hash, string(optBytes)))
}

func makeNonceCacheKey(scope string, nonce uint64, opts interface{}) []byte {
	optBytes, _ := json.Marshal(opts)
	return []byte(fmt.Sprintf("%s:nonce:%d|opts:%s", scope, nonce, string(optBytes)))
}

func makeObjKey(id string, opts interface{}) []byte {
	optBytes, _ := json.Marshal(opts)
	return []byte(id + string(optBytes))
}

func (bp *BlockProcessor) cacheObject(obj cacheableBlock, scope string, opts interface{}) {
	objKey := makeObjKey(obj.ID(), opts)

	// Store object
	bp.cache.Put(objKey, obj, 0)

	// Store nonce + hash lookup keys
	bp.cache.Put(makeHashCacheKey(scope, obj.Hash(), opts), objKey, 0)
	bp.cache.Put(makeNonceCacheKey(scope, obj.Nonce(), opts), objKey, 0)
}

func getObjectFromCache[T cacheableBlock](c TimedCache, scope string, hash string, nonce *uint64, opts interface{}) *T {
	var key interface{}
	if hash != "" {
		key, _ = c.Get(makeHashCacheKey(scope, hash, opts))
	} else if nonce != nil {
		key, _ = c.Get(makeNonceCacheKey(scope, *nonce, opts))
	}

	if key != nil {
		val, ok := c.Get(key.([]byte))
		if ok {
			return val.(*T)
		}
	}
	return nil
}
