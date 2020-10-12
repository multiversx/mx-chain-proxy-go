package versions

import "sync"

type versionManager struct {
	versions    map[string]FacadeHandler
	mutVersions sync.RWMutex
}

func NewVersionManager() *versionManager {
	return &versionManager{
		versions: make(map[string]FacadeHandler),
	}
}

// AddVersion will add the version and its corresponding handler to the inner map
func (vm *versionManager) AddVersion(version string, facadeHandler FacadeHandler) error {
	if facadeHandler == nil {
		return ErrNilFacadeHandler
	}
	vm.mutVersions.Lock()
	vm.versions[version] = facadeHandler
	vm.mutVersions.Unlock()

	return nil
}

// GetAllVersions returns a slice containing all the versions in string format
func (vm *versionManager) GetAllVersions() (map[string]FacadeHandler, error) {
	vm.mutVersions.RLock()
	defer vm.mutVersions.RUnlock()
	if len(vm.versions) == 0 {
		return nil, ErrNoVersionIsSet
	}

	return vm.versions, nil
}

// GetFacadeForApiVersion returns the facade for the given version or error if it does not exist
func (vm *versionManager) GetFacadeForApiVersion(version string) (FacadeHandler, error) {
	vm.mutVersions.RLock()
	defer vm.mutVersions.RUnlock()

	facade, ok := vm.versions[version]
	if !ok {
		return nil, ErrVersionNotFound
	}

	return facade, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (vm *versionManager) IsInterfaceNil() bool {
	return vm == nil
}
