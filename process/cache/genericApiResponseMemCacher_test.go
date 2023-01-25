package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/cache"
	"github.com/stretchr/testify/assert"
)

func TestNewGenericApiResponseMemoryCacher(t *testing.T) {
	t.Parallel()

	mc := cache.NewGenericApiResponseMemoryCacher()
	assert.NotNil(t, mc)
	assert.False(t, mc.IsInterfaceNil())
}

func TestGenericApiResponseMemoryCacher_StoreNilValStatsShouldNotPanic(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		assert.Nil(t, r)
	}()
	mc := cache.NewGenericApiResponseMemoryCacher()

	mc.Store(nil)
}

func TestGenericApiResponseMemoryCacher_StoreShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewGenericApiResponseMemoryCacher()
	apiResp := &data.GenericAPIResponse{
		Data: "test data",
	}

	mc.Store(apiResp)
	assert.Equal(t, apiResp, mc.GetGenericApiResponse())
}

func TestGenericApiResponseMemoryCacher_LoadNilStoredShouldErr(t *testing.T) {
	t.Parallel()

	mc := cache.NewGenericApiResponseMemoryCacher()

	apiResp, err := mc.Load()
	assert.Nil(t, apiResp)
	assert.Equal(t, cache.ErrNilGenericApiResponseInCache, err)
}

func TestGenericApiResponseMemoryCacher_LoadShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewGenericApiResponseMemoryCacher()
	apiResp := &data.GenericAPIResponse{
		Data: "test data 2",
	}

	mc.SetGenericApiResponse(apiResp)

	restoredApiResp, err := mc.Load()
	assert.NoError(t, err)
	assert.Equal(t, apiResp, restoredApiResp)
}

func TestGenericApiResponseMemoryCacher_ConcurrencySafe(t *testing.T) {
	t.Parallel()

	// here we should test if parallel accesses to the cache component leads to a race condition
	// if the mutex from the component is removed then test should fail when run with -race flag
	mc := cache.NewGenericApiResponseMemoryCacher()
	genericApiRespToStore := &data.GenericAPIResponse{
		Data: "test data concurrent test",
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
				mc.Store(genericApiRespToStore)
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
				_, _ = mc.Load()
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
}
