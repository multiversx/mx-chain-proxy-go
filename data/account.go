package data

import "github.com/multiversx/mx-chain-core-go/data/validator"

// AccountModel defines an account model (with associated information)
type AccountModel struct {
	Account   Account   `json:"account"`
	BlockInfo BlockInfo `json:"blockInfo"`
}

// AccountsModel defines the model of the accounts response
type AccountsModel struct {
	Accounts map[string]*Account `json:"accounts"`
}

// Account defines the data structure for an account
type Account struct {
	Address         string            `json:"address"`
	Nonce           uint64            `json:"nonce"`
	Balance         string            `json:"balance"`
	Username        string            `json:"username"`
	Code            string            `json:"code"`
	CodeHash        []byte            `json:"codeHash"`
	RootHash        []byte            `json:"rootHash"`
	CodeMetadata    []byte            `json:"codeMetadata"`
	DeveloperReward string            `json:"developerReward"`
	OwnerAddress    string            `json:"ownerAddress"`
	Pairs           map[string]string `json:"pairs,omitempty"`
}

// ValidatorApiResponse represents the data which is fetched from each validator for returning it in API call
type ValidatorApiResponse = validator.ValidatorStatistics

// ValidatorStatisticsResponse respects the format the validator statistics are received from the observers
type ValidatorStatisticsResponse struct {
	Statistics map[string]*ValidatorApiResponse `json:"statistics"`
}

// ValidatorStatisticsApiResponse respects the format the validator statistics are received from the observers
type ValidatorStatisticsApiResponse struct {
	Data  ValidatorStatisticsResponse `json:"data"`
	Error string                      `json:"error"`
	Code  string                      `json:"code"`
}

// AccountApiResponse defines a wrapped account that the node respond with
type AccountApiResponse struct {
	Data  AccountModel `json:"data"`
	Error string       `json:"error"`
	Code  string       `json:"code"`
}

// AccountsApiResponse defines the response that will be returned by the node when requesting multiple accounts
type AccountsApiResponse struct {
	Data  AccountsModel `json:"data"`
	Error string        `json:"error"`
	Code  string        `json:"code"`
}

// AccountKeyValueResponseData follows the format of the data field on an account key-value response
type AccountKeyValueResponseData struct {
	Value string `json:"value"`
}

// AccountKeyValueResponse defines the response for a request for a value of a key for an account
type AccountKeyValueResponse struct {
	Data  AccountKeyValueResponseData `json:"data"`
	Error string                      `json:"error"`
	Code  string                      `json:"code"`
}

// IterateKeysRequest defines the request for iterating keys of an account
type IterateKeysRequest struct {
	Address       string   `json:"address"`
	NumKeys       uint     `json:"numKeys"`
	IteratorState [][]byte `json:"iteratorState"`
}
