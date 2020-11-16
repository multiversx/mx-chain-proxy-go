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
func (bg *baseGroup) AddEndpoint(path string, handlerData data.EndpointHandlerData) error {
	if handlerData.Handler == nil {
		return ErrNilGinHandler
	}

	if bg.isEndpointRegistered(path) {
		return ErrEndpointAlreadyRegistered
	}

	bg.Lock()
	bg.endpoints[path] = &handlerData
	bg.Unlock()

	return nil
}

// UpdateEndpoint updates the handler for a given endpoint path
func (bg *baseGroup) UpdateEndpoint(path string, handlerData data.EndpointHandlerData) error {
	if !bg.isEndpointRegistered(path) {
		return ErrHandlerDoesNotExist
	}
	if handlerData.Handler == nil {
		return ErrNilGinHandler
	}

	bg.Lock()
	bg.endpoints[path] = &handlerData
	bg.Unlock()

	return nil
}

// RemoveEndpoint removes the handler for a given endpoint path
func (bg *baseGroup) RemoveEndpoint(path string) error {
	if !bg.isEndpointRegistered(path) {
		return ErrHandlerDoesNotExist
	}

	bg.Lock()
	delete(bg.endpoints, path)
	bg.Unlock()

	return nil
}

// RegisterRoutes will register all the endpoints to the given web server
func (bg *baseGroup) RegisterRoutes(ws *gin.RouterGroup) {
	bg.RLock()
	defer bg.RUnlock()

	for path, handlerData := range bg.endpoints {
		ws.Handle(handlerData.Method, path, handlerData.Handler)
	}
}

func (bg *baseGroup) isEndpointRegistered(endpoint string) bool {
	bg.RLock()
	defer bg.RUnlock()

	_, exists := bg.endpoints[endpoint]
	return exists
}

// IsInterfaceNil returns true if the value under the interface is nil
func (bg *baseGroup) IsInterfaceNil() bool {
	return bg == nil
}
