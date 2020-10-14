package groups

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type baseGroup struct {
	endpoints map[string]*data.EndpointHandlerData
	sync.RWMutex
}

// AddEndpoint will add the handler data for the given path inside the map
func (bag *baseGroup) AddEndpoint(path string, handlerData data.EndpointHandlerData) error {
	if handlerData.Handler == nil {
		return ErrNilGinHandler
	}

	if bag.isEndpointRegistered(path) {
		return ErrEndpointAlreadyRegistered
	}

	bag.Lock()
	bag.endpoints[path] = &handlerData
	bag.Unlock()

	return nil
}

// UpdateEndpoint updates the handler for a given endpoint path
func (bag *baseGroup) UpdateEndpoint(path string, handlerData data.EndpointHandlerData) error {
	if !bag.isEndpointRegistered(path) {
		return ErrHandlerDoesNotExist
	}
	if handlerData.Handler == nil {
		return ErrNilGinHandler
	}

	bag.Lock()
	bag.endpoints[path] = &handlerData
	bag.Unlock()

	return nil
}

// RemoveEndpoint removes the handler for a given endpoint path
func (bag *baseGroup) RemoveEndpoint(path string) error {
	if !bag.isEndpointRegistered(path) {
		return ErrHandlerDoesNotExist
	}

	bag.Lock()
	delete(bag.endpoints, path)
	bag.Unlock()

	return nil
}

func (bag *baseGroup) Routes(ws *gin.RouterGroup) {
	bag.RLock()
	defer bag.RUnlock()

	for path, handlerData := range bag.endpoints {
		ws.Handle(handlerData.Method, path, handlerData.Handler)
	}
}

func (bag *baseGroup) isEndpointRegistered(endpoint string) bool {
	bag.RLock()
	defer bag.RUnlock()

	_, exists := bag.endpoints[endpoint]
	return exists
}

// IsInterfaceNil returns true if the value under the interface is nil
func (bag *baseGroup) IsInterfaceNil() bool {
	return bag == nil
}
