package errors

import (
	"errors"
	"fmt"
)

// ErrInvalidAppContext signals an invalid context passed to the routing system
var ErrInvalidAppContext = errors.New("invalid app context")

// ErrGetValueForKey signals an error in getting the value of a key for an account
var ErrGetValueForKey = errors.New("get value for key error")

// ErrGetKeyValuePairs signals an error in getting the key-value pairs for a given address
var ErrGetKeyValuePairs = errors.New("get key value pairs error")

// ErrComputeShardForAddress signals an error in computing the shard ID for a given address
var ErrComputeShardForAddress = errors.New("compute shard ID for address error")

// ErrGetESDTTokenData signals an error in fetching an ESDT token data
var ErrGetESDTTokenData = errors.New("cannot get ESDT token data")

// ErrGetESDTsWithRole signals an error in fetching an tokens with role for an address
var ErrGetESDTsWithRole = errors.New("cannot get ESDTs with role")

// ErrGetESDTTokenData signals an error in fetching owned NFTs for an address
var ErrGetOwnedNFTs = errors.New("cannot get owned NFTs for account")

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

// ErrInvalidJSONRequest signals an error in json request formatting
var ErrInvalidJSONRequest = errors.New("invalid json request")

// ErrValidation signals an error in validation
var ErrValidation = errors.New("validation error")

// ErrValidationQueryParameterWithResult signals that an invalid query parameter has been provided
var ErrValidationQueryParameterWithResult = errors.New("invalid query parameter withResults")

// ErrValidatorQueryParameterCheckSignature signals that an invalid query parameter has been provided
var ErrValidatorQueryParameterCheckSignature = errors.New("invalid query parameter checkSignature")

// ErrInvalidSignatureHex signals a wrong hex value was provided for the signature
var ErrInvalidSignatureHex = errors.New("invalid signature, could not decode hex value")

// ErrTxGenerationFailed signals an error generating a transaction
var ErrTxGenerationFailed = errors.New("transaction generation failed")

// ErrInvalidSenderAddress signals a wrong format for sender address was provided
var ErrInvalidSenderAddress = errors.New("invalid hex sender address provided")

// ErrInvalidReceiverAddress signals a wrong format for receiver address was provided
var ErrInvalidReceiverAddress = errors.New("invalid hex receiver address provided")

// ErrTransactionNotFound signals that a transaction was not found
var ErrTransactionNotFound = errors.New("transaction not found")

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

// ErrEmptyRootHash signals that an empty root hash has been provided
var ErrEmptyRootHash = errors.New("empty root hash")

// ErrInvalidTxFields signals that one or more field of a transaction are invalid
type ErrInvalidTxFields struct {
	Message string
	Reason  string
}

// Error returns the string message of the ErrInvalidTxFields custom error struct
func (eitx *ErrInvalidTxFields) Error() string {
	return fmt.Sprintf("%s : %s", eitx.Message, eitx.Reason)
}
