package process

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// addressPath defines the address path at which the nodes answer
const addressPath = "/address/"

// AccountProcessor is able to process account requests
type AccountProcessor struct {
	connector       ExternalStorageConnector
	proc            Processor
	pubKeyConverter core.PubkeyConverter
}

// NewAccountProcessor creates a new instance of AccountProcessor
func NewAccountProcessor(proc Processor, pubKeyConverter core.PubkeyConverter, connector ExternalStorageConnector) (*AccountProcessor, error) {
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}
	if check.IfNil(connector) {
		return nil, ErrNilDatabaseConnector
	}

	return &AccountProcessor{
		proc:            proc,
		pubKeyConverter: pubKeyConverter,
		connector:       connector,
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

// GetAccount resolves the request by sending the request to the right observer and replies back the answer
func (ap *AccountProcessor) GetAccount(address string, options common.AccountQueryOptions) (*data.AccountModel, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		responseAccount := &data.AccountApiResponse{}

		url := common.BuildUrlWithAccountQueryOptions(addressPath+address, options)
		_, err = ap.proc.CallGetRestEndPoint(observer.Address, url, responseAccount)
		if err == nil {
			log.Info("account request", "address", address, "shard ID", observer.ShardId, "observer", observer.Address)
			return &responseAccount.Data, nil
		}

		log.Error("account request", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, ErrSendingRequest
}

// GetValueForKey returns the value for the given address and key
func (ap *AccountProcessor) GetValueForKey(address string, key string, options common.AccountQueryOptions) (string, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return "", err
	}

	for _, observer := range observers {
		apiResponse := data.AccountKeyValueResponse{}
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

	return "", ErrSendingRequest
}

// GetESDTTokenData returns the token data for a token with the given name
func (ap *AccountProcessor) GetESDTTokenData(address string, key string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetESDTsWithRole returns the token identifiers where the given address has the given role assigned
func (ap *AccountProcessor) GetESDTsWithRole(address string, role string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetESDTsRoles returns all the tokens and their roles for a given address
func (ap *AccountProcessor) GetESDTsRoles(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetNFTTokenIDsRegisteredByAddress returns the token identifiers of the NFTs registered by the address
func (ap *AccountProcessor) GetNFTTokenIDsRegisteredByAddress(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	//TODO: refactor the entire proxy so endpoints like this which simply forward the response will use a common
	// component, as described in task EN-9857.
	observers, err := ap.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetESDTNftTokenData returns the nft token data for a token with the given identifier and nonce
func (ap *AccountProcessor) GetESDTNftTokenData(address string, key string, nonce uint64, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetAllESDTTokens returns all the tokens for a given address
func (ap *AccountProcessor) GetAllESDTTokens(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetKeyValuePairs returns all the key-value pairs for a given address
func (ap *AccountProcessor) GetKeyValuePairs(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetGuardianData returns the guardian data for the given address
func (ap *AccountProcessor) GetGuardianData(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

// GetTransactions resolves the request and returns a slice of transaction for the specific address
func (ap *AccountProcessor) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	if _, err := ap.pubKeyConverter.Decode(address); err != nil {
		return nil, fmt.Errorf("%w, %v", ErrInvalidAddress, err)
	}

	return ap.connector.GetTransactionsByAddress(address)
}

// GetCodeHash returns the code hash for a given address
func (ap *AccountProcessor) GetCodeHash(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
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

	return nil, ErrSendingRequest
}

func (ap *AccountProcessor) getObserversForAddress(address string) ([]*data.NodeData, error) {
	addressBytes, err := ap.pubKeyConverter.Decode(address)
	if err != nil {
		return nil, err
	}

	shardID, err := ap.proc.ComputeShardId(addressBytes)
	if err != nil {
		return nil, err
	}

	observers, err := ap.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	return observers, nil
}

// GetBaseProcessor returns the base processor
func (ap *AccountProcessor) GetBaseProcessor() Processor {
	return ap.proc
}

// IsDataTrieMigrated returns true if the data trie for the given address is migrated
func (ap *AccountProcessor) IsDataTrieMigrated(address string, options common.AccountQueryOptions) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/is-data-trie-migrated"
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

	return nil, ErrSendingRequest
}
