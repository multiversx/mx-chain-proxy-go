package observer

import "github.com/ElrondNetwork/elrond-proxy-go/data"

// ObserversProviderHandler defines what an observer provider should be able to do
type ObserversProviderHandler interface {
	GetObserversByShardId(shardId uint32) ([]*data.Observer, error)
	GetAllObservers() ([]*data.Observer, error)
	IsInterfaceNil() bool
}
