package node

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
	GetShardStatus(shardID uint32) (map[string]interface{}, error)
	GetEpochMetrics(shardID uint32) (map[string]interface{}, error)
}
