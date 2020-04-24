package network

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetNetworkMetrics(shardID uint32) (map[string]interface{}, error)
}
