package runType

import (
	"errors"

	"github.com/multiversx/mx-chain-proxy-go/process"
)

var errNilRunTypeComponents = errors.New("nil run type components")

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
