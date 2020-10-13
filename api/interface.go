package api

import (
	"github.com/ElrondNetwork/elrond-proxy-go/versions"
	"github.com/gin-gonic/gin"
)

// ElrondProxyHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type ElrondProxyHandler interface {
}

// VersionManagerHandler defines the actions that a version manager implementation has to do
type VersionManagerHandler interface {
	AddVersion(version string, facadeHandler versions.FacadeHandler) error
	GetAllVersions() (map[string]versions.FacadeHandler, error)
	GetFacadeForApiVersion(version string) (versions.FacadeHandler, error)
	IsInterfaceNil() bool
}

type GroupHandler interface {
	Routes(ws *gin.RouterGroup)
	IsInterfaceNil() bool
}
