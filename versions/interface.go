package versions

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
}

// VersionManagerHandler defines the actions that a version manager implementation has to do
type VersionManagerHandler interface {
	AddVersion(version string, facadeHandler FacadeHandler) error
	GetAllVersions() (map[string]FacadeHandler, error)
	GetFacadeForApiVersion(version string) (FacadeHandler, error)
	IsInterfaceNil() bool
}
