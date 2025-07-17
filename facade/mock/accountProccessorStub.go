package mock

import (
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// AccountProcessorStub -
type AccountProcessorStub struct {
	GetAccountCalled                        func(address string, options common.AccountQueryOptions) (*data.AccountModel, error)
	GetAccountsCalled                       func(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error)
	GetValueForKeyCalled                    func(address string, key string, options common.AccountQueryOptions) (string, error)
	GetShardIDForAddressCalled              func(address string) (uint32, error)
	GetTransactionsCalled                   func(address string) ([]data.DatabaseTransaction, error)
	ValidatorStatisticsCalled               func() (map[string]*data.ValidatorApiResponse, error)
	GetAllESDTTokensCalled                  func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTTokenDataCalled                  func(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTNftTokenDataCalled               func(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsWithRoleCalled                  func(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetNFTTokenIDsRegisteredByAddressCalled func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetKeyValuePairsCalled                  func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetESDTsRolesCalled                     func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetCodeHashCalled                       func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	GetGuardianDataCalled                   func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	IsDataTrieMigratedCalled                func(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
	IterateKeysCalled                       func(address string, numKeys uint, iteratorState [][]byte, options common.AccountQueryOptions) (*data.GenericAPIResponse, error)
}

// GetKeyValuePairs -
func (aps *AccountProcessorStub) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetKeyValuePairsCalled(address, options)
}

// GetAllESDTTokens -
func (aps *AccountProcessorStub) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetAllESDTTokensCalled(address, options)
}

// GetESDTTokenData -
func (aps *AccountProcessorStub) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetESDTTokenDataCalled(address, key, options)
}

// GetESDTNftTokenData -
func (aps *AccountProcessorStub) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetESDTNftTokenDataCalled(address, key, nonce, options)
}

// GetESDTsWithRole -
func (aps *AccountProcessorStub) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetESDTsWithRoleCalled(address, role, options)
}

// GetESDTsRoles -
func (aps *AccountProcessorStub) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if aps.GetESDTsRolesCalled != nil {
		return aps.GetESDTsRolesCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// GetNFTTokenIDsRegisteredByAddress -
func (aps *AccountProcessorStub) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetNFTTokenIDsRegisteredByAddressCalled(address, options)
}

// GetAccount -
func (aps *AccountProcessorStub) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	return aps.GetAccountCalled(address, options)
}

// GetAccounts -
func (aps *AccountProcessorStub) GetAccounts(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error) {
	return aps.GetAccountsCalled(addresses, options)
}

// GetValueForKey -
func (aps *AccountProcessorStub) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	return aps.GetValueForKeyCalled(address, key, options)
}

// GetGuardianData -
func (aps *AccountProcessorStub) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetGuardianDataCalled(address, options)
}

// GetShardIDForAddress -
func (aps *AccountProcessorStub) GetShardIDForAddress(address string) (uint32, error) {
	return aps.GetShardIDForAddressCalled(address)
}

// GetCodeHash -
func (aps *AccountProcessorStub) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	return aps.GetCodeHashCalled(address, options)
}

// ValidatorStatistics -
func (aps *AccountProcessorStub) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return aps.ValidatorStatisticsCalled()
}

// IsDataTrieMigrated --
func (aps *AccountProcessorStub) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if aps.IsDataTrieMigratedCalled != nil {
		return aps.IsDataTrieMigratedCalled(address, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// IterateKeys -
func (aps *AccountProcessorStub) IterateKeys(address string, numKeys uint, iteratorState [][]byte, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	if aps.IterateKeysCalled != nil {
		return aps.IterateKeysCalled(address, numKeys, iteratorState, options)
	}

	return &data.GenericAPIResponse{}, nil
}

// AuctionList -
func (aps *AccountProcessorStub) AuctionList() ([]*data.AuctionListValidatorAPIResponse, error) {
	return nil, nil
}
