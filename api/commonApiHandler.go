package api

import (
	"sync"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type commonApiHandler struct {
	groups map[string]data.GroupHandler
	sync.RWMutex
}

func NewCommonApiHandler() *commonApiHandler {
	return &commonApiHandler{
		groups: initBaseGroups(),
	}
}

func initBaseGroups() map[string]data.GroupHandler {
	accountsGroup := groups.NewBaseAccountsGroup()
	blocksGroup := groups.NewBaseBlockGroup()
	blockAtlasGroup := groups.NewBaseBlockAtlasGroup()
	hyperBlocksGroup := groups.NewBaseHyperBlockGroup()
	networkGroup := groups.NewBaseNetworkGroup()
	nodeGroup := groups.NewBaseNodeGroup()
	transactionsGroup := groups.NewBaseTransactionsGroup()
	validatorsGroup := groups.NewBaseValidatorGroup()
	vmValuesGroup := groups.NewBaseValidatorGroup()

	return map[string]data.GroupHandler{
		"/address":     accountsGroup,
		"/block":       blocksGroup,
		"/block-atlas": blockAtlasGroup,
		"/hyperblock":  hyperBlocksGroup,
		"/network":     networkGroup,
		"/node":        nodeGroup,
		"/transaction": transactionsGroup,
		"/validator":   validatorsGroup,
		"/vm-values":   vmValuesGroup,
	}
}

// AddGroup will add the group at the given path inside the map
func (cah *commonApiHandler) AddGroup(path string, group data.GroupHandler) error {
	if check.IfNil(group) {
		return ErrNilGroupHandler
	}
	if cah.isGroupRegistered(path) {
		return ErrGroupAlreadyRegistered
	}

	cah.Lock()
	cah.groups[path] = group
	cah.Unlock()

	return nil
}

// UpdateGroup updates the group at a given path
func (cah *commonApiHandler) UpdateGroup(path string, group data.GroupHandler) error {
	if !cah.isGroupRegistered(path) {
		return ErrGroupDoesNotExist
	}
	if check.IfNil(group) {
		return ErrNilGroupHandler
	}

	cah.Lock()
	cah.groups[path] = group
	cah.Unlock()

	return nil
}

// GetGroup returns the group at a given path
func (cah *commonApiHandler) GetGroup(path string) (data.GroupHandler, error) {
	if !cah.isGroupRegistered(path) {
		return nil, ErrGroupDoesNotExist
	}

	cah.RLock()
	defer cah.RUnlock()
	return cah.groups[path], nil
}

// GetAllGroups returns the group at a given path
func (cah *commonApiHandler) GetAllGroups() map[string]data.GroupHandler {
	cah.RLock()
	defer cah.RUnlock()
	return cah.groups
}

// RemoveGroup removes the group at a given path
func (cah *commonApiHandler) RemoveGroup(path string) error {
	if !cah.isGroupRegistered(path) {
		return ErrGroupAlreadyRegistered
	}

	cah.Lock()
	delete(cah.groups, path)
	cah.Unlock()

	return nil
}

func (cah *commonApiHandler) isGroupRegistered(endpoint string) bool {
	cah.RLock()
	defer cah.RUnlock()

	_, exists := cah.groups[endpoint]
	return exists
}

// IsInterfaceNil returns true if the value under the interface is nil
func (cah *commonApiHandler) IsInterfaceNil() bool {
	return cah == nil
}
