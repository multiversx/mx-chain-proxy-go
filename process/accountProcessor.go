package process

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// AddressPath defines the address path at which the nodes answer
const AddressPath = "/address/"

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
func (ap *AccountProcessor) GetAccount(address string) (*data.Account, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		responseAccount := &data.AccountApiResponse{}

		_, err = ap.proc.CallGetRestEndPoint(observer.Address, AddressPath+address, responseAccount)
		if err == nil {
			log.Info("account request", "address", address, "shard ID", observer.ShardId, "observer", observer.Address)
			return &responseAccount.Data.AccountData, nil
		}

		log.Error("account request", "observer", observer.Address, "address", address, "error", err.Error())
	}

	return nil, ErrSendingRequest
}

// GetValueForKey returns the value for the given address and key
func (ap *AccountProcessor) GetValueForKey(address string, key string) (string, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return "", err
	}

	for _, observer := range observers {
		apiResponse := data.AccountKeyValueResponse{}
		apiPath := AddressPath + address + "/key/" + key
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
func (ap *AccountProcessor) GetESDTTokenData(address string, key string) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/esdt/" + key
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
func (ap *AccountProcessor) GetESDTsWithRole(address string, role string) (*data.GenericAPIResponse, error) {
	observers, err := ap.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/esdts-with-role/" + role
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
func (ap *AccountProcessor) GetESDTsRoles(address string) (*data.GenericAPIResponse, error) {
	observers, err := ap.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/esdts/roles"
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
func (ap *AccountProcessor) GetNFTTokenIDsRegisteredByAddress(address string) (*data.GenericAPIResponse, error) {
	//TODO: refactor the entire proxy so endpoints like this which simply forward the response will use a common
	// component, as described in task EN-9857.
	observers, err := ap.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/registered-nfts/"
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
func (ap *AccountProcessor) GetESDTNftTokenData(address string, key string, nonce uint64) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		nonceAsString := fmt.Sprintf("%d", nonce)
		apiPath := AddressPath + address + "/nft/" + key + "/nonce/" + nonceAsString
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
func (ap *AccountProcessor) GetAllESDTTokens(address string) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/esdt"
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
func (ap *AccountProcessor) GetKeyValuePairs(address string) (*data.GenericAPIResponse, error) {
	observers, err := ap.getObserversForAddress(address)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		apiResponse := data.GenericAPIResponse{}
		apiPath := AddressPath + address + "/keys"
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

// GetTransactions resolves the request and returns a slice of transaction for the specific address
func (ap *AccountProcessor) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	if _, err := ap.pubKeyConverter.Decode(address); err != nil {
		return nil, fmt.Errorf("%w, %v", ErrInvalidAddress, err)
	}

	return ap.connector.GetTransactionsByAddress(address)
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
