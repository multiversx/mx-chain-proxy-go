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

// NewApiHandler returns a new instance of commonApiHandler
func NewApiHandler(facade data.FacadeHandler) (*commonApiHandler, error) {
	if facade == nil {
		return nil, ErrNilFacade
	}

	groupsWithFacade, err := initBaseGroupsWithFacade(facade)
	if err != nil {
		return nil, err
	}

	return &commonApiHandler{
		groups: groupsWithFacade,
	}, nil
}

func initBaseGroupsWithFacade(facade data.FacadeHandler) (map[string]data.GroupHandler, error) {
	accountsGroup, err := groups.NewAccountsGroup(facade)
	if err != nil {
		return nil, err
	}

	blocksGroup, err := groups.NewBlockGroup(facade)
	if err != nil {
		return nil, err
	}

	blockAtlasGroup, err := groups.NewBlockAtlasGroup(facade)
	if err != nil {
		return nil, err
	}

	hyperBlocksGroup, err := groups.NewHyperBlockGroup(facade)
	if err != nil {
		return nil, err
	}

	networkGroup, err := groups.NewNetworkGroup(facade)
	if err != nil {
		return nil, err
	}

	nodeGroup, err := groups.NewNodeGroup(facade)
	if err != nil {
		return nil, err
	}

	transactionsGroup, err := groups.NewTransactionGroup(facade)
	if err != nil {
		return nil, err
	}

	validatorsGroup, err := groups.NewValidatorGroup(facade)
	if err != nil {
		return nil, err
	}

	vmValuesGroup, err := groups.NewVmValuesGroup(facade)
	if err != nil {
		return nil, err
	}

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
	}, nil
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
