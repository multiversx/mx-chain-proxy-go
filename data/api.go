package data

import (
	"time"

	"github.com/gin-gonic/gin"
)

// GenericAPIResponse defines the structure of all responses on API endpoints
type GenericAPIResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
	Code  ReturnCode  `json:"code"`
}

// NetworkConfig is a dto that will keep information about the network config
type NetworkConfig struct {
	Config struct {
		ChainID               string `json:"erd_chain_id"`
		MinGasLimit           uint64 `json:"erd_min_gas_limit"`
		MinGasPrice           uint64 `json:"erd_min_gas_price"`
		MinTransactionVersion uint32 `json:"erd_min_transaction_version"`
	} `json:"config"`
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
	ApiConfig  ApiRoutesConfig
}

// EndpointHandlerData holds the items needed for creating a new HTTP endpoint
type EndpointHandlerData struct {
	Path    string
	Handler gin.HandlerFunc
	Method  string
}

// GroupHandler defines the actions that an api group handler should be able to do
type GroupHandler interface {
	AddEndpoint(path string, handlerData EndpointHandlerData) error
	UpdateEndpoint(path string, handlerData EndpointHandlerData) error
	RegisterRoutes(ws *gin.RouterGroup, apiConfig ApiRoutesConfig, authenticationFunc gin.HandlerFunc, rateLimiter gin.HandlerFunc, statusMetricExtractor gin.HandlerFunc)
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

// StatusMetricsProvider defines what a status metrics provider should do
type StatusMetricsProvider interface {
	GetAll() map[string]*EndpointMetrics
	GetMetricsForPrometheus() string
	AddRequestData(path string, withError bool, duration time.Duration)
	IsInterfaceNil() bool
}

// ApiRoutesConfig holds the configuration related to Rest API routes
type ApiRoutesConfig struct {
	APIPackages map[string]APIPackageConfig
}

// APIPackageConfig holds the configuration for the routes of each package
type APIPackageConfig struct {
	Routes []RouteConfig
}

// RouteConfig holds the configuration for a single route
type RouteConfig struct {
	Name      string
	Open      bool
	Secured   bool
	RateLimit uint64
}

// Credential holds an username and a password
type Credential struct {
	Username string
	Password string
}
