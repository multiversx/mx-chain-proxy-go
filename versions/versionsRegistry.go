package versions

import (
	"sync"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

type versionsRegistry struct {
	versions map[string]*data.VersionData
	sync.RWMutex
}

// NewVersionsRegistry returns a new instance of versionsRegistry
func NewVersionsRegistry() *versionsRegistry {
	return &versionsRegistry{
		versions: make(map[string]*data.VersionData),
	}
}

// AddVersion will add the version and its corresponding handler to the inner map
func (vm *versionsRegistry) AddVersion(version string, versionData *data.VersionData) error {
	if versionData.Facade == nil {
		return ErrNilFacadeHandler
	}
	if versionData.ApiHandler == nil {
		return ErrNilApiHandler
	}

	vm.Lock()
	vm.versions[version] = versionData
	vm.Unlock()

	return nil
}

// GetAllVersions returns a slice containing all the versions in string format
func (vm *versionsRegistry) GetAllVersions() (map[string]*data.VersionData, error) {
	vm.RLock()
	defer vm.RUnlock()
	if len(vm.versions) == 0 {
		return nil, ErrNoVersionIsSet
	}

	return vm.versions, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (vm *versionsRegistry) IsInterfaceNil() bool {
	return vm == nil
}
