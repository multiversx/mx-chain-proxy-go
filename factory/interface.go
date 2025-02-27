package factory

import (
	"github.com/multiversx/mx-chain-proxy-go/process"
)

// ComponentHandler defines the actions common to all component handlers
type ComponentHandler interface {
	Create() error
	Close() error
	CheckSubcomponents() error
	String() string
}

// RunTypeComponentsHandler defines the run type components handler actions
type RunTypeComponentsHandler interface {
	ComponentHandler
	RunTypeComponentsHolder
}

// RunTypeComponentsHolder holds the run type components
type RunTypeComponentsHolder interface {
	TxNotarizationCheckerHandlerCreator() process.TxNotarizationCheckerHandler
	Create() error
	Close() error
	CheckSubcomponents() error
	String() string
	IsInterfaceNil() bool
}
