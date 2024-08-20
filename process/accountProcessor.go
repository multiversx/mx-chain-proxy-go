package process

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/observer/availabilityCommon"
)

// addressPath defines the address path at which the nodes answer
const addressPath = "/address/"

// AccountProcessor is able to process account requests
type AccountProcessor struct {
	proc                 Processor
	pubKeyConverter      core.PubkeyConverter
	availabilityProvider availabilityCommon.AvailabilityProvider
}

// NewAccountProcessor creates a new instance of AccountProcessor
func NewAccountProcessor(proc Processor, pubKeyConverter core.PubkeyConverter) (*AccountProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	return &AccountProcessor{
		proc:                 proc,
		pubKeyConverter:      pubKeyConverter,
		availabilityProvider: availabilityCommon.AvailabilityProvider{},
	}, nil
}

// GetShardIDForAddress resolves the request by returning the shard ID for a given address for the current proxy's configuration
func (ap *AccountProcessor) GetShardIDForAddress(address string) (uint32, error) {
	addressBytes, err := ap.pubKeyConverter.Decode(address)
	if err != nil {
		return 0, err
	}

	return ap.proc.ComputeShardId(addressBytes)
}

// GetAccount resolves the request by sending the request to the right observer and returns the response
func (ap *AccountProcessor) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	responseAccount := data.AccountApiResponse{}
	for _, observer := range observers {

		url := common.BuildUrlWithAccountQueryOptions(addressPath+address, options)
		_, err = ap.proc.CallGetRestEndPoint(observer.Address, url, &responseAccount)
		if err == nil {
			log.Info("account request", "address", address, "shard ID", observer.ShardId, "observer", observer.Address)
			return &responseAccount.Data, nil
		}

		log.Error("account request", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(responseAccount.Error)
}

// GetAccounts will return data about the provided accounts
func (ap *AccountProcessor) GetAccounts(addresses []string, options common.AccountQueryOptions) (*data.AccountsModel, error) {
	addressesInShards := make(map[uint32][]string)
	var shardID uint32
	var err error
	for _, address := range addresses {
		shardID, err = ap.GetShardIDForAddress(address)
		if err != nil {
			return nil, fmt.Errorf("%w while trying to compute shard ID of address %s", err, address)
		}

		addressesInShards[shardID] = append(addressesInShards[shardID], address)
	}

	var wg sync.WaitGroup
	wg.Add(len(addressesInShards))

	var shardErr error
	var mut sync.Mutex // Mutex to protect the shared map and error
	accountsResponse := make(map[string]*data.Account)

	for shID, accounts := range addressesInShards {
		go func(shID uint32, accounts []string) {
			defer wg.Done()
			accountsInShard, errGetAccounts := ap.getAccountsInShard(accounts, shID, options)

			mut.Lock()
			defer mut.Unlock()

			if errGetAccounts != nil {
				shardErr = errGetAccounts
				return
			}

			for address, account := range accountsInShard {
				accountsResponse[address] = account
			}
		}(shID, accounts)
	}

	wg.Wait()

	if shardErr != nil {
		return nil, shardErr
	}

	return &data.AccountsModel{
		Accounts: accountsResponse,
	}, nil
}

func (ap *AccountProcessor) getAccountsInShard(addresses []string, shardID uint32, options common.AccountQueryOptions) (map[string]*data.Account, error) {
	observers, err := ap.proc.GetObservers(shardID, data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	apiResponse := data.AccountsApiResponse{}
	apiPath := addressPath + "bulk"
	apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
	for _, observer := range observers {
		respCode, err := ap.proc.CallPostRestEndPoint(observer.Address, apiPath, addresses, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("bulk accounts request",
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return apiResponse.Data.Accounts, nil
		}

		log.Error("bulk accounts request", "observer", observer.Address, "error", err.Error())
	}

	return nil, ErrSendingRequest
}

// GetValueForKey returns the value for the given address and key
func (ap *AccountProcessor) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return "", err
	}

	apiResponse := data.AccountKeyValueResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/key/" + key
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account value for key request",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return "", errors.New(apiResponse.Error)
			}

			return apiResponse.Data.Value, nil
		}

		log.Error("account value for key request", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return "", WrapObserversError(apiResponse.Error)
}

// GetESDTTokenData returns the token data for a token with the given name
func (ap *AccountProcessor) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/esdt/" + key
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account ESDT token data",
				"address", address,
				"token", key,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get ESDT token data", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetESDTsWithRole returns the token identifiers where the given address has the given role assigned
func (ap *AccountProcessor) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.proc.GetObservers(core.MetachainShardId, availability)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/esdts-with-role/" + role
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account ESDTs with role",
				"address", address,
				"role", role,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get ESDTs with role", "observer", observer.Address, "address", address, "role", role, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetESDTsRoles returns all the tokens and their roles for a given address
func (ap *AccountProcessor) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.proc.GetObservers(core.MetachainShardId, availability)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/esdts/roles"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, errGet := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if errGet == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account ESDTs roles",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get ESDTs roles", "observer", observer.Address, "address", address, "error", errGet.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetNFTTokenIDsRegisteredByAddress returns the token identifiers of the NFTs registered by the address
func (ap *AccountProcessor) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	//TODO: refactor the entire proxy so endpoints like this which simply forward the response will use a common
	// component, as described in task EN-9857.
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.proc.GetObservers(core.MetachainShardId, availability)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/registered-nfts/"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account get owned NFTs",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get owned NFTs", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetESDTNftTokenData returns the nft token data for a token with the given identifier and nonce
func (ap *AccountProcessor) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		nonceAsString := fmt.Sprintf("%d", nonce)
		apiPath := addressPath + address + "/nft/" + key + "/nonce/" + nonceAsString
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account ESDT NFT token data",
				"address", address,
				"token", key,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get ESDT nft token data", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetAllESDTTokens returns all the tokens for a given address
func (ap *AccountProcessor) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/esdt"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account all ESDT tokens",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get all ESDT tokens", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetKeyValuePairs returns all the key-value pairs for a given address
func (ap *AccountProcessor) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/keys"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account get all key-value pairs",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get all key-value pairs error", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetGuardianData returns the guardian data for the given address
func (ap *AccountProcessor) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/guardian-data"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account get guardian data",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get guardian data", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// GetCodeHash returns the code hash for a given address
func (ap *AccountProcessor) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	availability := ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
	observers, err := ap.getObserversForAddress(address, availability, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/code-hash"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("account get code hash",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account get code hash error", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

func (ap *AccountProcessor) getShardIfOdAddress(address string) (uint32, error) {
	addressBytes, err := ap.pubKeyConverter.Decode(address)
	if err != nil {
		return 0, err
	}

	return ap.proc.ComputeShardId(addressBytes)
}

func (ap *AccountProcessor) getObserversForAddress(address string, availability data.ObserverDataAvailabilityType, forcedShardID core.OptionalUint32) ([]*data.NodeData, error) {
	if forcedShardID.HasValue {
		return ap.proc.GetObservers(forcedShardID.Value, availability)
	}

	addressBytes, err := ap.pubKeyConverter.Decode(address)
	if err != nil {
		return nil, err
	}

	shardID, err := ap.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	return ap.proc.GetObservers(shardID, availability)
}

// GetBaseProcessor returns the base processor
func (ap *AccountProcessor) GetBaseProcessor() Processor {
	return ap.proc
}

// IsDataTrieMigrated returns true if the data trie for the given address is migrated
func (ap *AccountProcessor) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address, data.AvailabilityRecent, options.ForcedShardID)
	if err != nil {
		return nil, err
	}

	apiResponse := data.GenericAPIResponse{}
	for _, observer := range observers {
		apiPath := addressPath + address + "/is-data-trie-migrated"
		apiPath = common.BuildUrlWithAccountQueryOptions(apiPath, options)
		respCode, err := ap.proc.CallGetRestEndPoint(observer.Address, apiPath, &apiResponse)
		if err == nil || respCode == http.StatusBadRequest || respCode == http.StatusInternalServerError {
			log.Info("is data trie migrated",
				"address", address,
				"shard ID", observer.ShardId,
				"observer", observer.Address,
				"http code", respCode)
			if apiResponse.Error != "" {
				return nil, errors.New(apiResponse.Error)
			}

			return &apiResponse, nil
		}

		log.Error("account is data trie migrated", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, WrapObserversError(apiResponse.Error)
}

// WrapObserversError wraps the observers error
func WrapObserversError(responseError string) error {
	if len(responseError) == 0 {
		return ErrSendingRequest
	}

	return fmt.Errorf("%w, %s", ErrSendingRequest, responseError)
}

func (ap *AccountProcessor) getAvailabilityBasedOnAccountQueryOptions(options common.AccountQueryOptions) data.ObserverDataAvailabilityType {
	return ap.availabilityProvider.AvailabilityForAccountQueryOptions(options)
}
