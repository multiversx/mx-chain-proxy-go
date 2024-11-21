package runType

import (
	"github.com/multiversx/mx-chain-proxy-go/process/factory"
)

type runTypeComponentsFactory struct{}

// NewRunTypeComponentsFactory will return a new instance of run type components factory
func NewRunTypeComponentsFactory() *runTypeComponentsFactory {
	return &runTypeComponentsFactory{}
}

// Create will create the run type components
func (rtcf *runTypeComponentsFactory) Create() *runTypeComponents {
	return &runTypeComponents{
		txNotarizationCheckerHandlerCreator: factory.NewTxNotarizationChecker(),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (rtcf *runTypeComponentsFactory) IsInterfaceNil() bool {
	return rtcf == nil
}
