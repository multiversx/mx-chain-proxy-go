package process

import (
	"github.com/ElrondNetwork/elrond-go-sandbox/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
)

type BaseProcessor struct {
	lastConfig       *config.Config
	shardCoordinator sharding.Coordinator
}

func (bp *BaseProcessor) ApplyConfig(cfg *config.Config) error {
	if cfg == nil {
		return ErrNilConfig
	}

}

func (bp *BaseProcessor) GetObserverConfig(shardId uint32) (*config.Config, error) {

}

func (bp *BaseProcessor) ComputeShardId(address []byte) (uint32, error) {

}
