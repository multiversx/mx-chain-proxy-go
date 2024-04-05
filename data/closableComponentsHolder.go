package data

import (
	"sync"

	logger "github.com/multiversx/mx-chain-logger-go"
)

// closableComponent defines the behaviour of a component that is closable
type closableComponent interface {
	Close() error
}

var log = logger.GetOrCreate("data")

// ClosableComponentsHandler is a structure that holds a list of closable components and closes them when needed
type ClosableComponentsHandler struct {
	components []closableComponent
	sync.Mutex
}

// NewClosableComponentsHandler will return a new instance of closableComponentsHandler
func NewClosableComponentsHandler() *ClosableComponentsHandler {
	return &ClosableComponentsHandler{
		components: make([]closableComponent, 0),
	}
}

// Add will add one or more components to the internal closable components slice
func (cch *ClosableComponentsHandler) Add(components ...closableComponent) {
	cch.Lock()
	cch.components = append(cch.components, components...)
	cch.Unlock()
}

// Close will handle the closing of all the components from the internal slice
func (cch *ClosableComponentsHandler) Close() {
	cch.Lock()
	defer cch.Unlock()

	for _, component := range cch.components {
		log.LogIfError(component.Close())
	}
}
