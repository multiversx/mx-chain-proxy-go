package errors

import (
	"errors"
	"fmt"
)

// ErrGetAccount signals an error in fetching an account
var ErrGetAccount = errors.New("cannot get account")

// ErrGetValueForKey signals an error in getting the value of a key for an account
var ErrGetValueForKey = errors.New("get value for key error")

// ErrGetKeyValuePairs signals an error in getting the key-value pairs for a given address
var ErrGetKeyValuePairs = errors.New("get key value pairs error")

// ErrInvalidAddressesArray signals that an invalid input has been provided
var ErrInvalidAddressesArray = errors.New("invalid addresses array")

// ErrCannotGetAddresses signals an error when trying to fetch a bulk of accounts
var ErrCannotGetAddresses = errors.New("error while fetching a bulk of accounts")

// ErrComputeShardForAddress signals an error in computing the shard ID for a given address
var ErrComputeShardForAddress = errors.New("compute shard ID for address error")

// ErrGetESDTTokenData signals an error in fetching an ESDT token data
var ErrGetESDTTokenData = errors.New("cannot get ESDT token data")

// ErrGetGuardianData signals an error in fetching an address guardian data
var ErrGetGuardianData = errors.New("cannot get guardian data")

// ErrGetESDTsWithRole signals an error in fetching an tokens with role for an address
var ErrGetESDTsWithRole = errors.New("cannot get ESDTs with role")

// ErrGetRolesForAccount signals an error in getting esdt tokens and roles for a given address
var ErrGetRolesForAccount = errors.New("get roles for account error")

// ErrGetNFTTokenIDsRegisteredByAddress signals an error in fetching owned NFTs for an address
var ErrGetNFTTokenIDsRegisteredByAddress = errors.New("cannot get owned NFTs for account")

// ErrEmptyAddress signals that an empty address was provided
var ErrEmptyAddress = errors.New("address is empty")

// ErrEmptyKey signals that an empty key was provided
var ErrEmptyKey = errors.New("key is empty")

// ErrEmptyTokenIdentifier signals that an empty token identifier was provided
var ErrEmptyTokenIdentifier = errors.New("token identifier is empty")

// ErrCannotParseShardID signals that the shard ID cannot be parsed
var ErrCannotParseShardID = errors.New("cannot parse shard ID")

// ErrCannotParseNonce signals that the nonce cannot be parsed
var ErrCannotParseNonce = errors.New("cannot parse nonce")

// ErrCannotParseRound signals that the round cannot be parsed
var ErrCannotParseRound = errors.New("cannot parse round")

// ErrCannotParseEpoch signals that the epoch cannot be parsed
var ErrCannotParseEpoch = errors.New("cannot parse epoch")

// ErrInvalidJSONRequest signals an error in json request formatting
var ErrInvalidJSONRequest = errors.New("invalid json request")

// ErrValidation signals an error in validation
var ErrValidation = errors.New("validation error")

// ErrBadUrlParams signals one or more incorrectly provided URL params (generic error)
var ErrBadUrlParams = errors.New("bad url parameter(s)")

// ErrGetCodeHash signals an error in fetching the code hash for an account
var ErrGetCodeHash = errors.New("cannot get code hash")

// ErrValidationQueryParameterWithResult signals that an invalid query parameter has been provided
var ErrValidationQueryParameterWithResult = errors.New("invalid query parameter withResults")

// ErrValidatorQueryParameterCheckSignature signals that an invalid query parameter has been provided
var ErrValidatorQueryParameterCheckSignature = errors.New("invalid query parameter checkSignature")

// ErrInvalidSignatureHex signals a wrong hex value was provided for the signature
var ErrInvalidSignatureHex = errors.New("invalid signature, could not decode hex value")

// ErrInvalidGuardianSignatureHex signals a wrong hex value provided for the guardian signature
var ErrInvalidGuardianSignatureHex = errors.New("invalid guardian signature, could not decode hex value")

// ErrInvalidGuardianAddress signals a wrong format for receiver address was provided
var ErrInvalidGuardianAddress = errors.New("invalid hex receiver address provided")

// ErrTxGenerationFailed signals an error generating a transaction
var ErrTxGenerationFailed = errors.New("transaction generation failed")

// ErrInvalidSenderAddress signals a wrong format for sender address was provided
var ErrInvalidSenderAddress = errors.New("invalid hex sender address provided")

// ErrInvalidReceiverAddress signals a wrong format for receiver address was provided
var ErrInvalidReceiverAddress = errors.New("invalid hex receiver address provided")

// ErrTransactionNotFound signals that a transaction was not found
var ErrTransactionNotFound = errors.New("transaction not found")

// ErrSCRsNoFound signals that smart contract results were not found
var ErrSCRsNoFound = errors.New("smart contract results not found")

// ErrTransactionsNotFoundInPool signals that no transaction was not found in pool
var ErrTransactionsNotFoundInPool = errors.New("transactions not found in pool")

// ErrTransactionHashMissing signals that a transaction was not found
var ErrTransactionHashMissing = errors.New("transaction hash missing")

// ErrFaucetNotEnabled signals that the faucet mechanism is not enabled
var ErrFaucetNotEnabled = errors.New("faucet not enabled")

// ErrInvalidBlockNonceParam signals that an invalid block's nonce parameter has been provided
var ErrInvalidBlockNonceParam = errors.New("invalid block nonce parameter")

// ErrInvalidBlockHashParam signals that an invalid block's hash parameter has been provided
var ErrInvalidBlockHashParam = errors.New("invalid block hash parameter")

// ErrInvalidShardIDParam signals that an invalid shard ID parameter has been provided
var ErrInvalidShardIDParam = errors.New("invalid shard ID parameter")

// ErrInvalidEpochParam signals that an invalid epoch parameter has been provided
var ErrInvalidEpochParam = errors.New("invalid epoch parameter")

// ErrEmptyRootHash signals that an empty root hash has been provided
var ErrEmptyRootHash = errors.New("empty root hash")

// ErrEmptySenderToGetLatestNonce signals that an error happened when trying to fetch latest nonce
var ErrEmptySenderToGetLatestNonce = errors.New("empty sender to get latest nonce")

// ErrEmptySenderToGetNonceGaps signals that an error happened when trying to fetch nonce gaps
var ErrEmptySenderToGetNonceGaps = errors.New("empty sender to get nonce gaps")

// ErrFetchingLatestNonceCannotIncludeFields signals that an error happened when trying to fetch latest nonce
var ErrFetchingLatestNonceCannotIncludeFields = errors.New("fetching latest nonce cannot include fields")

// ErrFetchingNonceGapsCannotIncludeFields signals that an error happened when trying to fetch nonce gaps
var ErrFetchingNonceGapsCannotIncludeFields = errors.New("fetching nonce gaps cannot include fields")

// ErrInvalidFields signals that invalid fields were provided
var ErrInvalidFields = errors.New("invalid fields")

// ErrOperationNotAllowed signals that the operation is not allowed
var ErrOperationNotAllowed = errors.New("operation not allowed")

// ErrIsDataTrieMigrated signals that an error occurred while trying to verify the migration status of the data trie
var ErrIsDataTrieMigrated = errors.New("could not verify the migration status of the data trie")

// ErrInvalidTxFields signals that one or more field of a transaction are invalid
type ErrInvalidTxFields struct {
	Message string
	Reason  string
}

// Error returns the string message of the ErrInvalidTxFields custom error struct
func (eitx *ErrInvalidTxFields) Error() string {
	return fmt.Sprintf("%s : %s", eitx.Message, eitx.Reason)
}
