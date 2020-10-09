package services

import "github.com/coinbase/rosetta-sdk-go/types"

var (
	ErrUnableToGetChainID = &types.Error{
		Code:      1,
		Message:   "unable to get chain ID",
		Retriable: true,
	}

	ErrUnableToGetAccount = &types.Error{
		Code:      2,
		Message:   "unable to get account",
		Retriable: true,
	}

	ErrInvalidAccountAddress = &types.Error{
		Code:      3,
		Message:   "invalid account address",
		Retriable: false,
	}

	ErrUnableToGetBlock = &types.Error{
		Code:      4,
		Message:   "unable to get block",
		Retriable: true,
	}

	ErrNotImplemented = &types.Error{
		Code:      5,
		Message:   "operation not implemented",
		Retriable: false,
	}

	ErrUnableToSubmitTransaction = &types.Error{
		Code:      6,
		Message:   "unable to submit transaction",
		Retriable: false,
	}

	ErrMalformedValue = &types.Error{
		Code:      7,
		Message:   "malformed value",
		Retriable: false,
	}

	ErrUnableToGetNodeStatus = &types.Error{
		Code:      8,
		Message:   "unable to get node status",
		Retriable: true,
	}

	ErrUnableToGetClientVersion = &types.Error{
		Code:      9,
		Message:   "unable to get client version",
		Retriable: true,
	}
	ErrMustQueryByIndexOrByHash = &types.Error{
		Code:      10,
		Message:   "must query block by index or by hash",
		Retriable: false,
	}
	ErrConstructionCheck = &types.Error{
		Code:      11,
		Message:   "operation construction check error",
		Retriable: false,
	}
	ErrUnableToGetNetworkConfig = &types.Error{
		Code:      12,
		Message:   "unable to get network config",
		Retriable: true,
	}
	ErrInvalidInputParam = &types.Error{
		Code:      13,
		Message:   "Invalid input param: ",
		Retriable: false,
	}
	ErrUnsupportedCurveType = &types.Error{
		Code:      14,
		Message:   "unsupported curve type",
		Retriable: false,
	}
	ErrInsufficientGasLimit = &types.Error{
		Code:      15,
		Message:   "insufficient gas limit",
		Retriable: false,
	}
	ErrGasPriceTooLow = &types.Error{
		Code:      16,
		Message:   "gas price is to low",
		Retriable: false,
	}
	ErrTransactionIsNotInPool = &types.Error{
		Code:      17,
		Message:   "transaction is not in pool",
		Retriable: false,
	}
	ErrCannotParsePoolTransaction = &types.Error{
		Code:      18,
		Message:   "cannot parse pool transaction",
		Retriable: false,
	}

	Errors = []*types.Error{
		ErrUnableToGetChainID,
		ErrUnableToGetAccount,
		ErrInvalidAccountAddress,
		ErrUnableToGetBlock,
		ErrNotImplemented,
		ErrUnableToSubmitTransaction,
		ErrMalformedValue,
		ErrUnableToGetNodeStatus,
		ErrUnableToGetClientVersion,
		ErrMustQueryByIndexOrByHash,
		ErrConstructionCheck,
		ErrUnableToGetNetworkConfig,
		ErrUnsupportedCurveType,
		ErrInsufficientGasLimit,
		ErrGasPriceTooLow,
		ErrTransactionIsNotInPool,
		ErrCannotParsePoolTransaction,
		ErrInvalidInputParam,
	}
)

// wrapErr adds details to the types.Error provided. We use a function
// to do this so that we don't accidentally overwrite the standard
// errors.
func wrapErr(rErr *types.Error, err error) *types.Error {
	newErr := &types.Error{
		Code:      rErr.Code,
		Message:   rErr.Message,
		Retriable: rErr.Retriable,
	}
	if err != nil {
		newErr.Details = map[string]interface{}{
			"context": err.Error(),
		}
	}

	return newErr
}
