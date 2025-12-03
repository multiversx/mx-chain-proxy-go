package process

import (
	"encoding/json"
	"fmt"
)

type cacheableBlock interface {
	Hash() string
	Nonce() uint64
}

// No error checks for this cache.
// These caching errors should never happen, and if they do, they should not be blocking

func (bp *BlockProcessor) cacheObject(obj cacheableBlock, scope string, opts interface{}) {
	objKey, err := makeObjKey(scope, obj.Hash(), opts)
	if err != nil {
		return
	}

	hashLookupKey, err := makeHashCacheKey(scope, obj.Hash(), opts)
	if err != nil {
		return
	}

	nonceLookupKey, err := makeNonceCacheKey(scope, obj.Nonce(), opts)
	if err != nil {
		return
	}

	// Store object
	_ = bp.cache.Put(objKey, obj, 0)

	// Store nonce + hash lookup keys
	_ = bp.cache.Put(hashLookupKey, objKey, 0)
	_ = bp.cache.Put(nonceLookupKey, objKey, 0)
}

func makeObjKey(scope string, hash string, opts interface{}) ([]byte, error) {
	optBytes, err := json.Marshal(opts)
	if err != nil {
		log.Error("makeObjKey", "error", err)
		return nil, err
	}
	return []byte(fmt.Sprintf("%s:%s|opts:%s", scope, hash, string(optBytes))), nil
}

func makeHashCacheKey(scope string, hash string, opts interface{}) ([]byte, error) {
	optBytes, err := json.Marshal(opts)
	if err != nil {
		log.Error("makeHashCacheKey", "error", err)
		return nil, err
	}
	return []byte(fmt.Sprintf("%s:hash:%s|opts:%s", scope, hash, string(optBytes))), nil
}

func makeNonceCacheKey(scope string, nonce uint64, opts interface{}) ([]byte, error) {
	optBytes, err := json.Marshal(opts)
	if err != nil {
		log.Error("makeNonceCacheKey", "error", err)
		return nil, err
	}
	return []byte(fmt.Sprintf("%s:nonce:%d|opts:%s", scope, nonce, string(optBytes))), nil
}

func getObjectFromCacheWithHash[T cacheableBlock](c TimedCache, scope string, hash string, opts interface{}) T {
	var nilRet T
	hashKey, err := makeHashCacheKey(scope, hash, opts)
	if err != nil {
		return nilRet
	}

	return getObjFromCache[T](c, hashKey)
}

func getObjectFromCacheWithNonce[T cacheableBlock](c TimedCache, scope string, nonce uint64, opts interface{}) T {
	var nilRet T
	nonceKey, err := makeNonceCacheKey(scope, nonce, opts)
	if err != nil {
		return nilRet
	}

	return getObjFromCache[T](c, nonceKey)
}

func getObjFromCache[T cacheableBlock](c TimedCache, lookUpKey []byte) T {
	var retObj T

	key, _ := c.Get(lookUpKey)
	if key == nil {
		return retObj
	}

	keyBytes, ok := key.([]byte)
	if !ok {
		return retObj
	}

	val, ok := c.Get(keyBytes)
	if !ok {
		return retObj
	}

	result, ok := val.(T)
	if !ok {
		return retObj
	}
	return result
}
