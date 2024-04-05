package groups

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

var log = logger.GetOrCreate("api/groups")

type baseGroup struct {
	endpoints []*data.EndpointHandlerData
	sync.RWMutex
}

type endpointProperties struct {
	isOpen           bool
	isSecured        bool
	isFoundInConfig  bool
	rateLimiterPerIP uint64
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
func (bg *baseGroup) RegisterRoutes(
	ws *gin.RouterGroup,
	apiConfig data.ApiRoutesConfig,
	authenticationFunc gin.HandlerFunc,
	rateLimiter gin.HandlerFunc,
	statusMetricsExtractor gin.HandlerFunc,
) {
	bg.RLock()
	defer bg.RUnlock()

	for _, handlerData := range bg.endpoints {
		properties := getEndpointProperties(ws, handlerData.Path, apiConfig)
		if !properties.isFoundInConfig {
			log.Warn("endpoint not found in config", "path", handlerData.Path)
			ws.Handle(handlerData.Method, handlerData.Path, handlerData.Handler)
			continue
		}

		if !properties.isOpen {
			log.Debug("endpoint is not opened", "path", handlerData.Path)
			continue
		}

		middlewares := make([]gin.HandlerFunc, 0)
		if properties.isSecured {
			middlewares = append(middlewares, authenticationFunc)
		}

		if properties.rateLimiterPerIP > 0 {
			middlewares = append(middlewares, rateLimiter)
		}

		middlewares = append(middlewares, statusMetricsExtractor)
		middlewares = append(middlewares, handlerData.Handler)

		ws.Handle(handlerData.Method, handlerData.Path, middlewares...)
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
				isOpen:           route.Open,
				isSecured:        route.Secured,
				isFoundInConfig:  true,
				rateLimiterPerIP: route.RateLimit,
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
