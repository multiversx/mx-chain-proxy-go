package versions

import (
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/process/v1_2"
)

type ElrondProxyFacadeV1_2 struct {
	AccountsProcessor v1_2.AccountProcessorV1_2
	*facade.ElrondProxyFacade
}

func (epf *ElrondProxyFacadeV1_2) GetShardIDForAddressV1_2(address string, additionalField int) (uint32, error) {
	return epf.AccountsProcessor.GetShardIDForAddress(address, additionalField)
}
