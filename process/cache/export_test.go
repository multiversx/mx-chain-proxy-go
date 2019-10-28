package cache

import "github.com/ElrondNetwork/elrond-proxy-go/data"

func (hmc *HeartbeatMemoryCacher) GetStoredHbts() []data.PubKeyHeartbeat {
	return hmc.storedHeartbeats
}

func (hmc *HeartbeatMemoryCacher) SetStoredHbts(hbts []data.PubKeyHeartbeat) {
	hmc.storedHeartbeats = hbts
}
