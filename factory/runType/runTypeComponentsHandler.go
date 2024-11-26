package runType

import (
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"

	"github.com/multiversx/mx-chain-proxy-go/factory"
	"github.com/multiversx/mx-chain-proxy-go/process"
)

const runTypeComponentsName = "managedRunTypeComponents"

var _ factory.ComponentHandler = (*managedRunTypeComponents)(nil)
var _ factory.RunTypeComponentsHandler = (*managedRunTypeComponents)(nil)
var _ factory.RunTypeComponentsHolder = (*managedRunTypeComponents)(nil)

type managedRunTypeComponents struct {
	*runTypeComponents
	factory                  RunTypeComponentsCreator
	mutRunTypeCoreComponents sync.RWMutex
}

// NewManagedRunTypeComponents returns a news instance of managed runType core components
func NewManagedRunTypeComponents(rtc RunTypeComponentsCreator) (*managedRunTypeComponents, error) {
	if rtc == nil {
		return nil, errNilRunTypeComponents
	}

	return &managedRunTypeComponents{
		runTypeComponents: nil,
		factory:           rtc,
	}, nil
}

// Create will create the managed components
func (mrtc *managedRunTypeComponents) Create() error {
	rtc := mrtc.factory.Create()

	mrtc.mutRunTypeCoreComponents.Lock()
	mrtc.runTypeComponents = rtc
	mrtc.mutRunTypeCoreComponents.Unlock()

	return nil
}

// Close will close all underlying subcomponents
func (mrtc *managedRunTypeComponents) Close() error {
	mrtc.mutRunTypeCoreComponents.Lock()
	defer mrtc.mutRunTypeCoreComponents.Unlock()

	if check.IfNil(mrtc.runTypeComponents) {
		return nil
	}

	err := mrtc.runTypeComponents.Close()
	if err != nil {
		return err
	}
	mrtc.runTypeComponents = nil

	return nil
}

// CheckSubcomponents verifies all subcomponents
func (mrtc *managedRunTypeComponents) CheckSubcomponents() error {
	mrtc.mutRunTypeCoreComponents.RLock()
	defer mrtc.mutRunTypeCoreComponents.RUnlock()

	if check.IfNil(mrtc.runTypeComponents) {
		return errNilRunTypeComponents
	}
	if check.IfNil(mrtc.txNotarizationCheckerHandlerCreator) {
		return process.ErrNilTxNotarizationCheckerHandler
	}
	return nil
}

// TxNotarizationCheckerHandlerCreator returns tx notarization checker handler
func (mrtc *managedRunTypeComponents) TxNotarizationCheckerHandlerCreator() process.TxNotarizationCheckerHandler {
	mrtc.mutRunTypeCoreComponents.RLock()
	defer mrtc.mutRunTypeCoreComponents.RUnlock()

	if check.IfNil(mrtc.runTypeComponents) {
		return nil
	}

	return mrtc.runTypeComponents.txNotarizationCheckerHandlerCreator
}

// IsInterfaceNil returns true if the interface is nil
func (mrtc *managedRunTypeComponents) IsInterfaceNil() bool {
	return mrtc == nil
}

// String returns the name of the component
func (mrtc *managedRunTypeComponents) String() string {
	return runTypeComponentsName
}
