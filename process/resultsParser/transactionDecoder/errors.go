package transactionDecoder

import (
	"errors"
)

var (
	// ErrValueSet will signal a parsing error from string to big.Int
	ErrValueSet = errors.New("failed to set Value")

	// ErrNotESDTTransfer will signal that the transaction is not an ESDTTransfer.
	ErrNotESDTTransfer = errors.New("not ESDTTransfer transaction metadata")

	// ErrNotESDTNFTTransfer will signal that the transaction is not an ESDTNFTTransfer.
	ErrNotESDTNFTTransfer = errors.New("not ESDTNFTTransfer transaction metadata")

	// ErrNotMultiESDTNFTTransfer will signal that the transaction is not a MultiESDTNFTTransfer.
	ErrNotMultiESDTNFTTransfer = errors.New("not MultiESDTNFTTransfer transaction metadata")

	// ErrNoArgs will signal that the transaction does not have any function arguments.
	ErrNoArgs = errors.New("no arguments provided")

	// ErrSenderReceiver will signal that the sender and receiver address in the transaction do not match.
	ErrSenderReceiver = errors.New("sender does not match receiver")

	// ErrInvalidAddress will signal that the address in invalid.
	ErrInvalidAddress = errors.New("invalid address")
)
