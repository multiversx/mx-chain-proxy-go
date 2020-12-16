package data

import (
	"github.com/gin-gonic/gin"
)

// GenericAPIResponse defines the structure of all responses on API endpoints
type GenericAPIResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
	Code  ReturnCode  `json:"code"`
}

// ReturnCode defines the type defines to identify return codes
type ReturnCode string

const (
	// ReturnCodeSuccess defines a successful request
	ReturnCodeSuccess ReturnCode = "successful"

	// ReturnCodeInternalError defines a request which hasn't been executed successfully due to an internal error
	ReturnCodeInternalError ReturnCode = "internal_issue"

	// ReturnCodeRequestError defines a request which hasn't been executed successfully due to a bad request received
	ReturnCodeRequestError ReturnCode = "bad_request"
)

// VersionData holds the components specific for each version
type VersionData struct {
	Facade     FacadeHandler
	ApiHandler ApiHandler
}

// EndpointHandlerData holds the items needed for creating a new HTTP endpoint
type EndpointHandlerData struct {
	Handler gin.HandlerFunc
	Method  string
}

// GroupHandler defines the actions that an api group handler should be able to do
type GroupHandler interface {
	AddEndpoint(path string, handlerData EndpointHandlerData) error
	UpdateEndpoint(path string, handlerData EndpointHandlerData) error
	RegisterRoutes(ws *gin.RouterGroup)
	RemoveEndpoint(path string) error
	IsInterfaceNil() bool
}

// ApiHandler defines the actions that an api handler should be able to do
type ApiHandler interface {
	AddGroup(path string, group GroupHandler) error
	UpdateGroup(path string, group GroupHandler) error
	GetGroup(path string) (GroupHandler, error)
	GetAllGroups() map[string]GroupHandler
	RemoveGroup(path string) error
	IsInterfaceNil() bool
}

// FacadeHandler interface defines methods that can be used from facade context variable
type FacadeHandler interface {
}

// VersionsRegistryHandler defines the actions that a versions registry implementation has to do
type VersionsRegistryHandler interface {
	AddVersion(version string, versionData *VersionData) error
	GetAllVersions() (map[string]*VersionData, error)
	IsInterfaceNil() bool
}
