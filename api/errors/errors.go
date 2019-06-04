package errors

import "errors"

// ErrInvalidAppContext signals an invalid context passed to the routing system
var ErrInvalidAppContext = errors.New("invalid app context")

// ErrValidation signals an error in validation
var ErrValidation = errors.New("validation error")

// ErrInvalidSignatureHex signals a wrong hex value was provided for the signature
var ErrInvalidSignatureHex = errors.New("invalid signature, could not decode hex value")

// ErrTxGenerationFailed signals an error generating a transaction
var ErrTxGenerationFailed = errors.New("transaction generation failed")
