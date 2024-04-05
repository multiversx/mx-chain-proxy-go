package cache

import (
	"sync"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

// HeartbeatMemoryCacher will handle caching the heartbeats response
type HeartbeatMemoryCacher struct {
	storedHeartbeats []data.PubKeyHeartbeat
	mutHeartbeats    sync.RWMutex
}

// NewHeartbeatMemoryCacher will return a new instance of HeartbeatMemoryCacher
func NewHeartbeatMemoryCacher() *HeartbeatMemoryCacher {
	return &HeartbeatMemoryCacher{
		storedHeartbeats: nil,
		mutHeartbeats:    sync.RWMutex{},
	}
}

// LoadHeartbeats will return the heartbeats response stored in cache (if found)
func (hmc *HeartbeatMemoryCacher) LoadHeartbeats() (*data.HeartbeatResponse, error) {
	hmc.mutHeartbeats.RLock()
	defer hmc.mutHeartbeats.RUnlock()

	if hmc.storedHeartbeats == nil {
		return nil, ErrNilHeartbeatsInCache
	}

	return &data.HeartbeatResponse{Heartbeats: hmc.storedHeartbeats}, nil
}

// StoreHeartbeats will update the stored heartbeats response in cache
func (hmc *HeartbeatMemoryCacher) StoreHeartbeats(hbts *data.HeartbeatResponse) error {
	if hbts == nil {
		return ErrNilHeartbeatsToStoreInCache
	}

	hmc.mutHeartbeats.Lock()
	hmc.storedHeartbeats = hbts.Heartbeats
	hmc.mutHeartbeats.Unlock()

	return nil
}

// IsInterfaceNil will return true if there is no value under the interface
func (hmc *HeartbeatMemoryCacher) IsInterfaceNil() bool {
	return hmc == nil
}
