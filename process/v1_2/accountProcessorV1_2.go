package v1_2

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-proxy-go/process"
)

// AccountProcessorV1_2 is the account processor for the version v1.1
type AccountProcessorV1_2 struct {
	*process.AccountProcessor
}

func (ap *AccountProcessorV1_2) GetShardIDForAddress(address string, additionalField int) (uint32, error) {
	fmt.Println(address, additionalField)
	return 0, nil
}
