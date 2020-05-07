package network

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetNetworkStatusMetrics(shardID uint32) (map[string]interface{}, error)
	GetNetworkConfigMetrics() (map[string]interface{}, error)
}
