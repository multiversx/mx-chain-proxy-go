package versions

import (
	"sync"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type versionManager struct {
	commonApiHandler data.ApiHandler
	versions         map[string]*data.VersionData
	mutVersions      sync.RWMutex
}

func NewVersionManager(commonApiHandler data.ApiHandler) *versionManager {
	return &versionManager{
		commonApiHandler: commonApiHandler,
		versions:         make(map[string]*data.VersionData),
	}
}

// AddVersion will add the version and its corresponding handler to the inner map
func (vm *versionManager) AddVersion(version string, versionData *data.VersionData) error {
	if versionData.Facade == nil {
		return ErrNilFacadeHandler
	}
	if versionData.ApiHandler == nil {
		versionData.ApiHandler = vm.commonApiHandler
	}

	vm.mutVersions.Lock()
	vm.versions[version] = versionData
	vm.mutVersions.Unlock()

	return nil
}

// GetAllVersions returns a slice containing all the versions in string format
func (vm *versionManager) GetAllVersions() (map[string]*data.VersionData, error) {
	vm.mutVersions.RLock()
	defer vm.mutVersions.RUnlock()
	if len(vm.versions) == 0 {
		return nil, ErrNoVersionIsSet
	}

	return vm.versions, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (vm *versionManager) IsInterfaceNil() bool {
	return vm == nil
}
