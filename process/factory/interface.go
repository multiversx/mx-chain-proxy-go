package factory

import (
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// Processor defines what a processor should be able to do
type Processor interface {
	GetObservers(shardID uint32) ([]*data.Observer, error)
	ComputeShardId(addressBuff []byte) (uint32, error)
	CallGetRestEndPoint(address string, path string, value interface{}) (int, error)
	CallPostRestEndPoint(address string, path string, data interface{}, response interface{}) (int, error)
	GetAllObservers() []*data.Observer
	IsInterfaceNil() bool
}

// PrivateKeysLoaderHandler defines what a component which handles loading of the private keys file should do
type PrivateKeysLoaderHandler interface {
	PrivateKeysByShard() (map[uint32][]crypto.PrivateKey, error)
}
