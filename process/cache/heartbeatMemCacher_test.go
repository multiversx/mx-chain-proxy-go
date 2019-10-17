package cache_test

import (
	"sync"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process/cache"
	"github.com/stretchr/testify/assert"
)

func TestNewHeartbeatMemoryCacher(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()
	assert.NotNil(t, mc)
	assert.False(t, mc.IsInterfaceNil())
}

func TestHeartbeatMemoryCacher_StoreHeartbeatsShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()
	hbts := []data.PubKeyHeartbeat{
		{
			NodeDisplayName: "node1",
		},
		{
			NodeDisplayName: "node2",
		},
	}
	hbtsResp := data.HeartbeatResponse{Heartbeats: hbts}
	mc.StoreHeartbeats(&hbtsResp)

	assert.Equal(t, hbts, mc.GetStoredHbts())
}

func TestHeartbeatMemoryCacher_LoadHeartbeatsShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()
	hbts := []data.PubKeyHeartbeat{
		{
			NodeDisplayName: "node1",
		},
		{
			NodeDisplayName: "node2",
		},
	}

	mc.SetStoredHbts(hbts)

	restoredHbtsResp := mc.Heartbeats()
	assert.Equal(t, hbts, restoredHbtsResp.Heartbeats)
}

func TestHeartbeatMemoryCacher_ConcurrencySafe(t *testing.T) {
	t.Parallel()

	// here we should test if parallel accesses to the cache component leads to a race condition
	// if the mutex from the component is removed then test should fail when run with -race flag
	mc := cache.NewHeartbeatMemoryCacher()
	hbtsToStore := data.HeartbeatResponse{Heartbeats: []data.PubKeyHeartbeat{{NodeDisplayName: "node1"}}}

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
				mc.StoreHeartbeats(&hbtsToStore)
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
				_ = mc.Heartbeats()
				time.Sleep(5 * time.Millisecond)
			}
		}
	}()

	wg.Wait()
}
