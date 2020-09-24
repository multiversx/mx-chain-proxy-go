package services

import "github.com/coinbase/rosetta-sdk-go/types"

var (
	ErrUnableToGetChainID = &types.Error{
		Code:      1,
		Message:   "unable to get chain ID",
		Retriable: true,
	}

	ErrInvalidBlockchain = &types.Error{
		Code:      2,
		Message:   "invalid blockchain specified in network identifier",
		Retriable: false,
	}

	ErrInvalidNetwork = &types.Error{
		Code:      3,
		Message:   "invalid network specified in network identifier",
		Retriable: false,
	}

	ErrUnableToGetLatestBlock = &types.Error{
		Code:      4,
		Message:   "unable to get latest block",
		Retriable: true,
	}

	ErrUnableToGetGenesisBlock = &types.Error{
		Code:      5,
		Message:   "unable to get genesis block",
		Retriable: true,
	}

	ErrUnableToGetAccount = &types.Error{
		Code:      6,
		Message:   "unable to get account",
		Retriable: true,
	}

	ErrInvalidAccountAddress = &types.Error{
		Code:      7,
		Message:   "invalid account address",
		Retriable: false,
	}

	ErrUnableToGetBlock = &types.Error{
		Code:      8,
		Message:   "unable to get block",
		Retriable: true,
	}

	ErrNotImplemented = &types.Error{
		Code:      9,
		Message:   "operation not implemented",
		Retriable: false,
	}

	ErrUnableToGetTransactions = &types.Error{
		Code:      10,
		Message:   "unable to get transactions",
		Retriable: true,
	}

	ErrUnableToSubmitTransaction = &types.Error{
		Code:      11,
		Message:   "unable to submit transaction",
		Retriable: false,
	}

	ErrUnableToGetNextNonce = &types.Error{
		Code:      12,
		Message:   "unable to get next nonce",
		Retriable: true,
	}

	ErrMalformedValue = &types.Error{
		Code:      13,
		Message:   "malformed value",
		Retriable: false,
	}

	ErrUnableToGetNodeStatus = &types.Error{
		Code:      14,
		Message:   "unable to get node status",
		Retriable: true,
	}

	ErrTransactionNotFound = &types.Error{
		Code:      15,
		Message:   "transaction not found",
		Retriable: true,
	}

	ErrUnableToGetClientVersion = &types.Error{
		Code:      16,
		Message:   "unable to get client version",
		Retriable: true,
	}
	ErrMustQueryByIndexOrByHash = &types.Error{
		Code:      17,
		Message:   "must query block by index or by hash",
		Retriable: false,
	}
	ErrConstructionCheck = &types.Error{
		Code:      18,
		Message:   "operation construction check error",
		Retriable: false,
	}
	ErrUnableToGetNetworkConfig = &types.Error{
		Code:      19,
		Message:   "unable to get network config",
		Retriable: false,
	}
	ErrInvalidInputParam = &types.Error{
		Code:      20,
		Message:   "Invalid input param: ",
		Retriable: false,
	}

	Errors = []*types.Error{
		ErrUnableToGetChainID,
		ErrInvalidBlockchain,
		ErrInvalidNetwork,
		ErrUnableToGetLatestBlock,
		ErrUnableToGetGenesisBlock,
		ErrUnableToGetAccount,
		ErrInvalidAccountAddress,
		ErrUnableToGetBlock,
		ErrNotImplemented,
		ErrUnableToGetTransactions,
		ErrUnableToSubmitTransaction,
		ErrUnableToGetNextNonce,
		ErrMalformedValue,
		ErrUnableToGetNodeStatus,
		ErrTransactionNotFound,
		ErrUnableToGetClientVersion,
		ErrMustQueryByIndexOrByHash,
		ErrConstructionCheck,
		ErrUnableToGetNetworkConfig,
	}
)
