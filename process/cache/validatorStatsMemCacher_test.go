package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/cache"
	"github.com/stretchr/testify/assert"
)

func TestNewValidatorsStatsMemoryCacher(t *testing.T) {
	t.Parallel()

	mc := cache.NewValidatorsStatsMemoryCacher()
	assert.NotNil(t, mc)
	assert.False(t, mc.IsInterfaceNil())
}

func TestValidatorsStatsMemoryCacher_StoreNilValStatsShouldErr(t *testing.T) {
	t.Parallel()

	mc := cache.NewValidatorsStatsMemoryCacher()

	err := mc.StoreValStats(nil)
	assert.Equal(t, cache.ErrNilValidatorStatsToStoreInCache, err)
}

func TestValidatorsStatsMemoryCacher_StoreShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewValidatorsStatsMemoryCacher()
	valStats := map[string]*data.ValidatorApiResponse{
		"pubk1": {TempRating: 0.5},
	}
	err := mc.StoreValStats(valStats)

	assert.Nil(t, err)
	assert.Equal(t, valStats, mc.GetStoredValStats())
}

func TestValidatorsStatsMemoryCacher_LoadNilStoredValStatsShouldErr(t *testing.T) {
	t.Parallel()

	mc := cache.NewValidatorsStatsMemoryCacher()

	valStats, err := mc.LoadValStats()
	assert.Nil(t, valStats)
	assert.Equal(t, cache.ErrNilValidatorStatsInCache, err)
}

func TestValidatorsStatsMemoryCacher_LoadShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewValidatorsStatsMemoryCacher()
	valStats := map[string]*data.ValidatorApiResponse{
		"pubk1": {TempRating: 50.5},
		"pubk2": {TempRating: 50.6},
	}

	mc.SetStoredValStats(valStats)

	restoredValStatsResp, err := mc.LoadValStats()
	assert.NoError(t, err)
	assert.Equal(t, valStats, restoredValStatsResp)
}

func TestValidatorsStatsMemoryCacher_ConcurrencySafe(t *testing.T) {
	t.Parallel()

	// here we should test if parallel accesses to the cache component leads to a race condition
	// if the mutex from the component is removed then test should fail when run with -race flag
	mc := cache.NewValidatorsStatsMemoryCacher()
	valStatsToStore := map[string]*data.ValidatorApiResponse{
		"pubk1": {TempRating: 50.5},
		"pubk2": {TempRating: 50.6},
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	stopGoRoutinesEvent1 := time.After(1000 * time.Millisecond)
	stopGoRoutinesEvent2 := time.After(1100 * time.Millisecond)

	go func() {
		for {
			select {
			case <-stopGoRoutinesEvent1:
				wg.Done()
				break
			default:
				_ = mc.StoreValStats(valStatsToStore)
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-stopGoRoutinesEvent2:
				wg.Done()
				break
			default:
				_, _ = mc.LoadValStats()
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
}
