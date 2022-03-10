package data

import (
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

// closableComponent defines the behaviour of a component that is closable
type closableComponent interface {
	Close() error
}

var log = logger.GetOrCreate("data")

type ClosableComponentsHandler struct {
	components []closableComponent
	sync.Mutex
}

// NewClosableComponentsHandler returns a new instance of closableComponentsHandler
func NewClosableComponentsHandler() *ClosableComponentsHandler {
	return &ClosableComponentsHandler{
		components: make([]closableComponent, 0),
	}
}

// Add will add a component to the internal closable components slice
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
