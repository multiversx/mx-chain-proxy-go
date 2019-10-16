package cache_test

import (
	"testing"

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

func TestHeartbeatMemoryCacher_StoreHeartbeatsNilHbtsShouldErr(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()

	err := mc.StoreHeartbeats(nil)
	assert.Equal(t, cache.ErrNilHeartbeatsToStoreInCache, err)
}

func TestHeartbeatMemoryCacher_StoreHeartbeatsShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()
	hbts := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node1",
			},
			{
				NodeDisplayName: "node2",
			},
		},
	}

	err := mc.StoreHeartbeats(&hbts)
	assert.Nil(t, err)
}

func TestHeartbeatMemoryCacher_LoadHeartbeatsNilStoredHbtsShouldErr(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()

	hbts, err := mc.LoadHeartbeats()
	assert.Nil(t, hbts)
	assert.Equal(t, cache.ErrNilHeartbeatsInCache, err)
}

func TestHeartbeatMemoryCacher_LoadHeartbeatsShouldWork(t *testing.T) {
	t.Parallel()

	mc := cache.NewHeartbeatMemoryCacher()
	hbts := data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				NodeDisplayName: "node1",
			},
			{
				NodeDisplayName: "node2",
			},
		},
	}

	_ = mc.StoreHeartbeats(&hbts)

	restoredHbts, err := mc.LoadHeartbeats()
	assert.Nil(t, err)
	assert.Equal(t, hbts, *restoredHbts)
}
