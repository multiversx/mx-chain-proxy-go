package v_next

import (
	"github.com/multiversx/mx-chain-proxy-go/process"
)

// AccountProcessorV_next is the account processor for the version v_next
type AccountProcessorV_next struct {
	*process.AccountProcessor
}

// GetShardIDForAddress is an example of an updated endpoint the version v_next
func (ap *AccountProcessorV_next) GetShardIDForAddress(address string, additionalField int) (uint32, error) {
	return 37, nil
}

// NextEndpointHandler is an example of a new endpoint in the version v_next
func (ap *AccountProcessorV_next) NextEndpointHandler() string {
	return "test"
}
