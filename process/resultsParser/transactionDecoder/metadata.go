package transactionDecoder

import (
	"math/big"
)

type (
	// TransactionToDecode is the transaction whose Data field will be later decoded into metadata.
	TransactionToDecode struct {
		Sender   string
		Receiver string
		Data     string
		Value    string
	}

	// TransactionMetadata is the result of the decoded data field.
	TransactionMetadata struct {
		Sender   string
		Receiver string
		Value    *big.Int

		FunctionName string
		FunctionArgs []string

		Transfers []TransactionMetadataTransfer
	}

	// TransactionMetadataTransfer contains information about the transfers within the transaction.
	TransactionMetadataTransfer struct {
		Properties TokenTransferProperties
		Value      *big.Int
	}

	// TokenTransferProperties holds information about the token that has been transferred in the transaction.
	TokenTransferProperties struct {
		Token      string
		Collection string
		Identifier string
	}
)
