package versions

import (
	"github.com/multiversx/mx-chain-proxy-go/facade"
	"github.com/multiversx/mx-chain-proxy-go/process/v_next"
)

// ElrondProxyFacadeV_next is the facade that corresponds to the version v_next
type ElrondProxyFacadeV_next struct {
	AccountsProcessor v_next.AccountProcessorV_next
	*facade.ElrondProxyFacade
}

// GetShardIDForAddressV_next is an example function that demonstrates how to add a new custom handler for a modified api endpoint
func (epf *ElrondProxyFacadeV_next) GetShardIDForAddressV_next(address string, additionalField int) (uint32, error) {
	return epf.AccountsProcessor.GetShardIDForAddress(address, additionalField)
}

// NextEndpointHandler is an example function that demonstrates how to add a new custom handler for a new API endpoint
func (epf *ElrondProxyFacadeV_next) NextEndpointHandler() string {
	return epf.AccountsProcessor.NextEndpointHandler()
}
