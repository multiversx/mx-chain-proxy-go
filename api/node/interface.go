package node

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetHeartbeatData() (*data.HeartbeatResponse, error)
	GetEconomicsDataMetrics() (*data.GenericAPIResponse, error)
}
