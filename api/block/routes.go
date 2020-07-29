package block

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines full blocks related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/:shard/by-nonce/:nonce", ByNonceHandler)
	router.GET("/:shard/by-hash/:hash", ByHashHandler)
}

// ByHashHandler will handle the fetching and returning a block based on its hash
func ByHashHandler(c *gin.Context) {
	epf, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	shardID, err := getShardIdParam(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseShardID.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	hash := c.Param("hash")
	_, err = hex.DecodeString(hash)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrInvalidBlockHashParam.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	withTxs, err := getQueryParamWithTxs(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: withTxs param", apiErrors.ErrValidation),
			data.ReturnCodeInternalError,
		)
		return
	}

	blockByHashResponse, err := epf.GetBlockByHash(shardID, hash, withTxs)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// ByNonceHandler will handle the fetching and returning a block based on its nonce
func ByNonceHandler(c *gin.Context) {
	epf, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	shardID, err := getShardIdParam(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseShardID.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	nonce, err := getQueryParamNonce(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseNonce.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	withTxs, err := getQueryParamWithTxs(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: with txs param", apiErrors.ErrValidation),
			data.ReturnCodeRequestError,
		)
		return
	}

	blockByNonceResponse, err := epf.GetBlockByNonce(shardID, nonce, withTxs)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}

func getQueryParamWithTxs(c *gin.Context) (bool, error) {
	withTxsStr := c.Request.URL.Query().Get("withTxs")
	if withTxsStr == "" {
		return false, nil
	}

	return strconv.ParseBool(withTxsStr)
}

func getQueryParamNonce(c *gin.Context) (uint64, error) {
	nonceStr := c.Param("nonce")
	if nonceStr == "" {
		return 0, apiErrors.ErrInvalidBlockNonceParam
	}

	return strconv.ParseUint(nonceStr, 10, 64)
}

func getShardIdParam(c *gin.Context) (uint32, error) {
	shardStr := c.Param("shard")
	if shardStr == "" {
		return 0, apiErrors.ErrInvalidShardIDParam
	}

	shardID, err := strconv.ParseUint(shardStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(shardID), nil
}
