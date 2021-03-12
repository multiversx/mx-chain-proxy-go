package groups

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type baseGroup struct {
	endpoints []*data.EndpointHandlerData
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
	bg.endpoints = append(bg.endpoints, &handlerData)
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
	for i := 0; i < len(bg.endpoints); i++ {
		if bg.endpoints[i].Path == path {
			bg.endpoints[i] = &handlerData
		}
	}
	bg.Unlock()

	return nil
}

// RemoveEndpoint removes the handler for a given endpoint path
func (bg *baseGroup) RemoveEndpoint(path string) error {
	if !bg.isEndpointRegistered(path) {
		return ErrHandlerDoesNotExist
	}

	bg.Lock()
	for i := 0; i < len(bg.endpoints); i++ {
		if bg.endpoints[i].Path == path {
			bg.endpoints = append(bg.endpoints[:i], bg.endpoints[i+1:]...)
			break
		}
	}
	bg.Unlock()

	return nil
}

// RegisterRoutes will register all the endpoints to the given web server
func (bg *baseGroup) RegisterRoutes(ws *gin.RouterGroup) {
	bg.RLock()
	defer bg.RUnlock()

	for _, handlerData := range bg.endpoints {
		ws.Handle(handlerData.Method, handlerData.Path, handlerData.Handler)
	}
}

func (bg *baseGroup) isEndpointRegistered(endpoint string) bool {
	bg.RLock()
	defer bg.RUnlock()

	for _, end := range bg.endpoints {
		if end.Path == endpoint {
			return true
		}
	}

	return false
}

// IsInterfaceNil returns true if the value under the interface is nil
func (bg *baseGroup) IsInterfaceNil() bool {
	return bg == nil
}
