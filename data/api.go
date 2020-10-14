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

type VersionData struct {
	Facade     FacadeHandler
	ApiHandler ApiHandler
}

type EndpointHandlerData struct {
	Handler gin.HandlerFunc
	Method  string
}

type GroupHandler interface {
	AddEndpoint(path string, handlerData EndpointHandlerData) error
	UpdateEndpoint(path string, handlerData EndpointHandlerData) error
	Routes(ws *gin.RouterGroup)
	RemoveEndpoint(path string) error
	IsInterfaceNil() bool
}

type ApiHandler interface {
	AddGroup(path string, group GroupHandler) error
	UpdateGroup(path string, group GroupHandler) error
	GetGroup(path string) (GroupHandler, error)
	GetAllGroups() map[string]GroupHandler
	RemoveGroup(path string) error
	IsInterfaceNil() bool
}

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
}

// VersionManagerHandler defines the actions that a version manager implementation has to do
type VersionManagerHandler interface {
	AddVersion(version string, versionData *VersionData) error
	GetAllVersions() (map[string]*VersionData, error)
	IsInterfaceNil() bool
}
