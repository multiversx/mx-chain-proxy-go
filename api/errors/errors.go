package errors

import (
	"errors"
	"fmt"
)

// ErrInvalidAppContext signals an invalid context passed to the routing system
var ErrInvalidAppContext = errors.New("invalid app context")

// ErrGetValueForKey signals an error in getting the value of a key for an account
var ErrGetValueForKey = errors.New("get value for key error")

// ErrEmptyAddress signals that an empty address was provided
var ErrEmptyAddress = errors.New("address is empty")

// ErrEmptyKey signals that an empty key was provided
var ErrEmptyKey = errors.New("key is empty")

// ErrCannotParseShardID signals that the shard ID cannot be parsed
var ErrCannotParseShardID = errors.New("cannot parse shard ID")

// ErrCannotParseNonce signals that the nonce cannot be parsed
var ErrCannotParseNonce = errors.New("cannot parse nonce")

// ErrInvalidJSONRequest signals an error in json request formatting
var ErrInvalidJSONRequest = errors.New("invalid json request")

// ErrValidation signals an error in validation
var ErrValidation = errors.New("validation error")

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

// ErrInvalidTxFields signals that one or more field of a transaction are invalid
type ErrInvalidTxFields struct {
	Message string
	Reason  string
}

// Error returns the string message of the ErrInvalidTxFields custom error struct
func (eitx *ErrInvalidTxFields) Error() string {
	return fmt.Sprintf("%s : %s", eitx.Message, eitx.Reason)
}
