package groups

import (
	"strings"
	"sync"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

var log = logger.GetOrCreate("api/groups")

type baseGroup struct {
	endpoints map[string]*data.EndpointHandlerData
	sync.RWMutex
}

type endpointProperties struct {
	isOpen          bool
	isSecured       bool
	isFoundInConfig bool
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
func (bg *baseGroup) RegisterRoutes(ws *gin.RouterGroup, apiConfig data.ApiRoutesConfig, authenticationFunc gin.HandlerFunc) {
	bg.RLock()
	defer bg.RUnlock()

	for path, handlerData := range bg.endpoints {
		properties := getEndpointProperties(ws, path, apiConfig)
		if !properties.isFoundInConfig {
			log.Warn("endpoint not found in config", "path", path)
			ws.Handle(handlerData.Method, path, handlerData.Handler)
			continue
		}

		if !properties.isOpen {
			log.Debug("endpoint is not opened", "path", path)
			continue
		}

		if properties.isSecured {
			ws.Handle(handlerData.Method, path, authenticationFunc, handlerData.Handler)
			continue
		}

		ws.Handle(handlerData.Method, path, handlerData.Handler)
	}
}

func getEndpointProperties(ws *gin.RouterGroup, path string, apiConfig data.ApiRoutesConfig) endpointProperties {
	basePath := ws.BasePath()

	// ws.BasePath will return paths like /group or /v1.0/group so we need the last token after splitting by /
	splitPath := strings.Split(basePath, "/")
	basePath = splitPath[len(splitPath)-1]

	group, ok := apiConfig.APIPackages[basePath]
	if !ok {
		return endpointProperties{
			isFoundInConfig: false,
		}
	}

	for _, route := range group.Routes {
		if route.Name == path {
			return endpointProperties{
				isOpen:          route.Open,
				isSecured:       route.Secured,
				isFoundInConfig: true,
			}
		}
	}

	return endpointProperties{
		isFoundInConfig: false,
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
