package cache

import (
	"sync"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

// genericApiResponseMemoryCacher will handle caching the ValidatorsStatss response
type genericApiResponseMemoryCacher struct {
	storedResponse        *data.GenericAPIResponse
	mutGenericApiResponse sync.RWMutex
}

// NewGenericApiResponseMemoryCacher will return a new instance of genericApiResponseMemoryCacher
func NewGenericApiResponseMemoryCacher() *genericApiResponseMemoryCacher {
	return &genericApiResponseMemoryCacher{
		storedResponse:        nil,
		mutGenericApiResponse: sync.RWMutex{},
	}
}

// Load will return the generic api response stored in cache (if found)
func (garmc *genericApiResponseMemoryCacher) Load() (*data.GenericAPIResponse, error) {
	garmc.mutGenericApiResponse.RLock()
	defer garmc.mutGenericApiResponse.RUnlock()

	if garmc.storedResponse == nil {
		return nil, ErrNilGenericApiResponseInCache
	}

	return garmc.storedResponse, nil
}

// Store will update the generic api response response in cache
func (garmc *genericApiResponseMemoryCacher) Store(genericApiResponse *data.GenericAPIResponse) {
	garmc.mutGenericApiResponse.Lock()
	garmc.storedResponse = genericApiResponse
	garmc.mutGenericApiResponse.Unlock()
}

// IsInterfaceNil will return true if there is no value under the interface
func (garmc *genericApiResponseMemoryCacher) IsInterfaceNil() bool {
	return garmc == nil
}
