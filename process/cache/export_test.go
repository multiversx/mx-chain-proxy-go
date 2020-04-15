package cache

import "github.com/ElrondNetwork/elrond-proxy-go/data"

func (hmc *HeartbeatMemoryCacher) GetStoredHbts() []data.PubKeyHeartbeat {
	hmc.mutHeartbeats.RLock()
	defer hmc.mutHeartbeats.RUnlock()

	return hmc.storedHeartbeats
}

func (hmc *HeartbeatMemoryCacher) SetStoredHbts(hbts []data.PubKeyHeartbeat) {
	hmc.mutHeartbeats.Lock()
	hmc.storedHeartbeats = hbts
	hmc.mutHeartbeats.Unlock()
}

func (vsmc *validatorsStatsMemoryCacher) GetStoredValStats() map[string]*data.ValidatorApiResponse {
	vsmc.mutValidatorsStatss.RLock()
	defer vsmc.mutValidatorsStatss.RUnlock()

	return vsmc.storedValidatorsStats
}

func (vsmc *validatorsStatsMemoryCacher) SetStoredValStats(valStats map[string]*data.ValidatorApiResponse) {
	vsmc.mutValidatorsStatss.Lock()
	vsmc.storedValidatorsStats = valStats
	vsmc.mutValidatorsStatss.Unlock()
}
