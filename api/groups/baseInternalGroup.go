package groups

import (
	"encoding/hex"
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type internalGroup struct {
	facade InternalFacadeHandler
	*baseGroup
}

// NewInternalGroup returns a new instance of blockGroup
func NewInternalGroup(facadeHandler data.FacadeHandler) (*internalGroup, error) {
	facade, ok := facadeHandler.(InternalFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	bg := &internalGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/:shard/raw/block/by-nonce/:nonce", Handler: bg.rawBlockbyNonceHandler, Method: http.MethodGet},
		{Path: "/:shard/raw/block/by-hash/:hash", Handler: bg.rawBlockbyHashHandler, Method: http.MethodGet},
		{Path: "/:shard/json/block/by-nonce/:nonce", Handler: bg.internalBlockbyNonceHandler, Method: http.MethodGet},
		{Path: "/:shard/json/block/by-hash/:hash", Handler: bg.internalBlockbyHashHandler, Method: http.MethodGet},
		{Path: "/:shard/json/miniblock/by-hash/:hash", Handler: bg.internalMiniBlockbyHashHandler, Method: http.MethodGet},
		{Path: "/:shard/raw/miniblock/by-hash/:hash", Handler: bg.rawMiniBlockbyHashHandler, Method: http.MethodGet},
	}
	bg.baseGroup.endpoints = baseRoutesHandlers

	return bg, nil
}

// internalBlockbyHashHandler will handle the fetching and returning a block based on its hash
func (group *internalGroup) internalBlockbyHashHandler(c *gin.Context) {
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

	blockByHashResponse, err := group.facade.GetInternalBlockByHash(shardID, hash, common.Internal)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// internalBlockbyNonceHandler will handle the fetching and returning a block based on its hash
func (group *internalGroup) internalBlockbyNonceHandler(c *gin.Context) {
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

	blockByNonceResponse, err := group.facade.GetInternalBlockByNonce(shardID, nonce, common.Internal)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}

// rawBlockbyHashHandler will handle the fetching and returning a raw block based on its hash
func (group *internalGroup) rawBlockbyHashHandler(c *gin.Context) {
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

	blockByHashResponse, err := group.facade.GetInternalBlockByHash(shardID, hash, common.Proto)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// rawBlockbyNonceHandler will handle the fetching and returning a raw block based on its hash
func (group *internalGroup) rawBlockbyNonceHandler(c *gin.Context) {
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

	blockByNonceResponse, err := group.facade.GetInternalBlockByNonce(shardID, nonce, common.Proto)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}

// internalMiniBlockbyHashHandler will handle the fetching and returning a miniblock based on its hash
func (group *internalGroup) internalMiniBlockbyHashHandler(c *gin.Context) {
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

	miniBlockByHashResponse, err := group.facade.GetInternalMiniBlockByHash(shardID, hash, common.Internal)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, miniBlockByHashResponse)
}

// rawMiniBlockbyHashHandler will handle the fetching and returning a miniblock based on its hash
func (group *internalGroup) rawMiniBlockbyHashHandler(c *gin.Context) {
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

	miniBlockByHashResponse, err := group.facade.GetInternalMiniBlockByHash(shardID, hash, common.Proto)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, miniBlockByHashResponse)
}
