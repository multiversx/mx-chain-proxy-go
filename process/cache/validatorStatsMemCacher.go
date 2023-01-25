package cache

import (
	"sync"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

// validatorsStatsMemoryCacher will handle caching the ValidatorsStatss response
type validatorsStatsMemoryCacher struct {
	storedValidatorsStats map[string]*data.ValidatorApiResponse
	mutValidatorsStatss   sync.RWMutex
}

// NewValidatorsStatsMemoryCacher will return a new instance of validatorsStatsMemoryCacher
func NewValidatorsStatsMemoryCacher() *validatorsStatsMemoryCacher {
	return &validatorsStatsMemoryCacher{
		storedValidatorsStats: nil,
		mutValidatorsStatss:   sync.RWMutex{},
	}
}

// LoadValStats will return the ValidatorsStats response stored in cache (if found)
func (vsmc *validatorsStatsMemoryCacher) LoadValStats() (map[string]*data.ValidatorApiResponse, error) {
	vsmc.mutValidatorsStatss.RLock()
	defer vsmc.mutValidatorsStatss.RUnlock()

	if vsmc.storedValidatorsStats == nil {
		return nil, ErrNilValidatorStatsInCache
	}

	return vsmc.storedValidatorsStats, nil
}

// StoreValStats will update the stored ValidatorsStatss response in cache
func (vsmc *validatorsStatsMemoryCacher) StoreValStats(valStats map[string]*data.ValidatorApiResponse) error {
	if valStats == nil {
		return ErrNilValidatorStatsToStoreInCache
	}

	vsmc.mutValidatorsStatss.Lock()
	vsmc.storedValidatorsStats = valStats
	vsmc.mutValidatorsStatss.Unlock()

	return nil
}

// IsInterfaceNil will return true if there is no value under the interface
func (vsmc *validatorsStatsMemoryCacher) IsInterfaceNil() bool {
	return vsmc == nil
}
