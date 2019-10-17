package cache

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// HeartbeatMemoryCacher will handle caching the heartbeats response
type HeartbeatMemoryCacher struct {
	storedHeartbeats []data.PubKeyHeartbeat
	mutHeartbeats    *sync.RWMutex
}

// NewHeartbeatMemoryCacher will return a new instance of HeartbeatMemoryCacher
func NewHeartbeatMemoryCacher() *HeartbeatMemoryCacher {
	return &HeartbeatMemoryCacher{
		storedHeartbeats: nil,
		mutHeartbeats:    &sync.RWMutex{},
	}
}

// Heartbeats will return the heartbeats response stored in cache (if found)
func (hmc *HeartbeatMemoryCacher) Heartbeats() *data.HeartbeatResponse {
	hmc.mutHeartbeats.RLock()
	defer hmc.mutHeartbeats.RUnlock()
	return &data.HeartbeatResponse{Heartbeats: hmc.storedHeartbeats}
}

// StoreHeartbeats will update the stored heartbeats response in cache
func (hmc *HeartbeatMemoryCacher) StoreHeartbeats(hbts *data.HeartbeatResponse) {
	hmc.mutHeartbeats.Lock()
	hmc.storedHeartbeats = hbts.Heartbeats
	hmc.mutHeartbeats.Unlock()
}

// IsInterfaceNil will return true if there is no value under the interface
func (hmc *HeartbeatMemoryCacher) IsInterfaceNil() bool {
	return hmc == nil
}
