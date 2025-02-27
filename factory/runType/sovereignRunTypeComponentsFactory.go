package runType

import (
	"github.com/multiversx/mx-chain-proxy-go/process/factory"
)

type sovereignRunTypeComponentsFactory struct{}

// NewSovereignRunTypeComponentsFactory will return a new instance of sovereign run type components factory
func NewSovereignRunTypeComponentsFactory() *sovereignRunTypeComponentsFactory {
	return &sovereignRunTypeComponentsFactory{}
}

// Create will create the run type components
func (srtcf *sovereignRunTypeComponentsFactory) Create() *runTypeComponents {
	return &runTypeComponents{
		txNotarizationCheckerHandlerCreator: factory.NewSovereignTxNotarizationChecker(),
	}
}

// IsInterfaceNil returns true if there is no value under the interface
func (srtcf *sovereignRunTypeComponentsFactory) IsInterfaceNil() bool {
	return srtcf == nil
}
