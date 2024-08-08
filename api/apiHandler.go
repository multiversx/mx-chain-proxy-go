package api

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// apiHandler will handle the groups specific to an API version
// This component is not concurrent-safe because the way it is used now doesn't involve any concurrent reads/writes
type apiHandler struct {
	groups map[string]data.GroupHandler
}

// NewApiHandler returns a new instance of commonApiHandler
func NewApiHandler(facade data.FacadeHandler) (*apiHandler, error) {
	if facade == nil {
		return nil, ErrNilFacade
	}

	groupsWithFacade, err := initBaseGroupsWithFacade(facade)
	if err != nil {
		return nil, err
	}

	return &apiHandler{
		groups: groupsWithFacade,
	}, nil
}

func initBaseGroupsWithFacade(facade data.FacadeHandler) (map[string]data.GroupHandler, error) {
	accountsGroup, err := groups.NewAccountsGroup(facade)
	if err != nil {
		return nil, err
	}

	actionsGroup, err := groups.NewActionsGroup(facade)
	if err != nil {
		return nil, err
	}

	blockGroup, err := groups.NewBlockGroup(facade)
	if err != nil {
		return nil, err
	}

	blocksGroup, err := groups.NewBlocksGroup(facade)
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

	statusGroup, err := groups.NewStatusGroup(facade)
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

	proofGroup, err := groups.NewProofGroup(facade)
	if err != nil {
		return nil, err
	}

	internalGroup, err := groups.NewInternalGroup(facade)
	if err != nil {
		return nil, err
	}

	aboutGroup, err := groups.NewAboutGroup(facade)
	if err != nil {
		return nil, err
	}

	return map[string]data.GroupHandler{
		"/actions":     actionsGroup,
		"/address":     accountsGroup,
		"/block":       blockGroup,
		"/blocks":      blocksGroup,
		"/internal":    internalGroup,
		"/hyperblock":  hyperBlocksGroup,
		"/network":     networkGroup,
		"/node":        nodeGroup,
		"/status":      statusGroup,
		"/transaction": transactionsGroup,
		"/validator":   validatorsGroup,
		"/vm-values":   vmValuesGroup,
		"/proof":       proofGroup,
		"/about":       aboutGroup,
	}, nil
}

// AddGroup will add the group at the given path inside the map
func (cah *apiHandler) AddGroup(path string, group data.GroupHandler) error {
	if check.IfNil(group) {
		return ErrNilGroupHandler
	}
	if cah.isGroupRegistered(path) {
		return ErrGroupAlreadyRegistered
	}

	cah.groups[path] = group

	return nil
}

// UpdateGroup updates the group at a given path
func (cah *apiHandler) UpdateGroup(path string, group data.GroupHandler) error {
	if !cah.isGroupRegistered(path) {
		return ErrGroupDoesNotExist
	}
	if check.IfNil(group) {
		return ErrNilGroupHandler
	}

	cah.groups[path] = group

	return nil
}

// GetGroup returns the group at a given path
func (cah *apiHandler) GetGroup(path string) (data.GroupHandler, error) {
	if !cah.isGroupRegistered(path) {
		return nil, ErrGroupDoesNotExist
	}

	return cah.groups[path], nil
}

// GetAllGroups returns the group at a given path
func (cah *apiHandler) GetAllGroups() map[string]data.GroupHandler {
	return cah.groups
}

// RemoveGroup removes the group at a given path
func (cah *apiHandler) RemoveGroup(path string) error {
	if !cah.isGroupRegistered(path) {
		return ErrGroupAlreadyRegistered
	}

	delete(cah.groups, path)

	return nil
}

func (cah *apiHandler) isGroupRegistered(endpoint string) bool {
	_, exists := cah.groups[endpoint]
	return exists
}

// IsInterfaceNil returns true if the value under the interface is nil
func (cah *apiHandler) IsInterfaceNil() bool {
	return cah == nil
}
