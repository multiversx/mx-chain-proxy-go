package groups

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

type blockGroup struct {
	facade BlocksFacadeHandler
	*baseGroup
}

// NewBlockGroup returns a new instance of blockGroup
func NewBlockGroup(facadeHandler data.FacadeHandler) (*blockGroup, error) {
	facade, ok := facadeHandler.(BlocksFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	bg := &blockGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/:shard/by-nonce/:nonce": {Handler: bg.ByNonceHandler, Method: http.MethodGet},
		"/:shard/by-hash/:hash":   {Handler: bg.ByHashHandler, Method: http.MethodGet},
	}
	bg.baseGroup.endpoints = baseRoutesHandlers

	return bg, nil
}

// ByHashHandler will handle the fetching and returning a block based on its hash
func (bg *blockGroup) ByHashHandler(c *gin.Context) {
	shardID, err := shared.FetchShardIDFromRequest(c)
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

	blockByHashResponse, err := bg.facade.GetBlockByHash(shardID, hash, withTxs)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// ByNonceHandler will handle the fetching and returning a block based on its nonce
func (bg *blockGroup) ByNonceHandler(c *gin.Context) {
	shardID, err := shared.FetchShardIDFromRequest(c)
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

	nonce, err := shared.FetchNonceFromRequest(c)
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

	blockByNonceResponse, err := bg.facade.GetBlockByNonce(shardID, nonce, withTxs)
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
