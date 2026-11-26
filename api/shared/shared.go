package shared

import (
	"encoding/hex"
	stderrors "errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

// RespondWith will respond with the generic API response
func RespondWith(c *gin.Context, status int, dataField interface{}, error string, code data.ReturnCode) {
	c.JSON(
		status,
		data.GenericAPIResponse{
			Data:  dataField,
			Error: error,
			Code:  code,
		},
	)
}

// FetchNonceFromRequest will try to fetch the nonce from the request
func FetchNonceFromRequest(c *gin.Context) (uint64, error) {
	nonceStr := c.Param("nonce")
	if nonceStr == "" {
		return 0, errors.ErrInvalidBlockNonceParam
	}

	return strconv.ParseUint(nonceStr, 10, 64)
}

// FetchRoundFromRequest will try to fetch the round from the request
func FetchRoundFromRequest(c *gin.Context) (uint64, error) {
	roundStr := c.Param("round")
	if roundStr == "" {
		return 0, errors.ErrInvalidBlockNonceParam
	}

	return strconv.ParseUint(roundStr, 10, 64)
}

// FetchEpochFromRequest will try to fetch the epoch from the request
func FetchEpochFromRequest(c *gin.Context) (uint32, error) {
	epochStr := c.Param("epoch")
	if epochStr == "" {
		return 0, errors.ErrInvalidEpochParam
	}

	epoch, err := strconv.ParseUint(epochStr, 10, 32)
	return uint32(epoch), err
}

// FetchShardIDFromRequest will try to fetch the shard ID from the request
func FetchShardIDFromRequest(c *gin.Context) (uint32, error) {
	shardStr := c.Param("shard")
	if shardStr == "" {
		return 0, errors.ErrInvalidShardIDParam
	}

	shardID, err := strconv.ParseUint(shardStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(shardID), nil
}

// FetchHashFromRequest will try to fetch the hash from the request
func FetchHashFromRequest(c *gin.Context) (string, error) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		return "", fmt.Errorf("%w:%s", errors.ErrInvalidBlockHashParam, hash)
	}

	return hash, nil
}

// RespondWithBadRequest creates a generic response for bad request
func RespondWithBadRequest(c *gin.Context, errorMessage string) {
	RespondWith(c, http.StatusBadRequest, nil, errorMessage, data.ReturnCodeRequestError)
}

// RespondWithValidationError should be called when the request cannot be satisfied due to a (request) validation error
func RespondWithValidationError(c *gin.Context, err error, innerErr error) {
	errMessage := fmt.Sprintf("%s: %s", err.Error(), innerErr.Error())

	RespondWith(
		c,
		http.StatusBadRequest,
		nil,
		errMessage,
		data.ReturnCodeRequestError,
	)
}

// RespondWithInternalError should be called when the request cannot be satisfied due to an internal error
func RespondWithInternalError(c *gin.Context, err error, innerErr error) {
	errMessage := fmt.Sprintf("%s: %s", err.Error(), innerErr.Error())

	RespondWith(
		c,
		http.StatusInternalServerError,
		nil,
		errMessage,
		data.ReturnCodeInternalError,
	)
}

// ObserverErrorStatus maps paging errors without changing unrelated responders.
func ObserverErrorStatus(err error) (int, data.ReturnCode) {
	var observerErr *data.ObserverAPIError
	if !stderrors.As(err, &observerErr) {
		return http.StatusInternalServerError, data.ReturnCodeInternalError
	}
	switch observerErr.StatusCode {
	case http.StatusBadRequest, http.StatusRequestEntityTooLarge:
		return observerErr.StatusCode, data.ReturnCodeRequestError
	case http.StatusTooManyRequests:
		return observerErr.StatusCode, data.ReturnCodeSystemBusy
	default:
		return observerErr.StatusCode, data.ReturnCodeInternalError
	}
}

// RespondWithObserverError retains the address operation's error prefix.
func RespondWithObserverError(c *gin.Context, baseErr, err error) {
	if c.Request.Context().Err() != nil {
		return
	}
	status, code := ObserverErrorStatus(err)
	RespondWith(c, status, nil, fmt.Sprintf("%s: %s", baseErr, err), code)
}
