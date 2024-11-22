package runType

import (
	"github.com/multiversx/mx-chain-proxy-go/process"
)

type runTypeComponents struct {
	txNotarizationCheckerHandlerCreator process.TxNotarizationCheckerHandler
}

// Close does nothing
func (rtc *runTypeComponents) Close() error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (rtc *runTypeComponents) IsInterfaceNil() bool {
	return rtc == nil
}
