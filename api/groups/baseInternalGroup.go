package groups

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
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
		{Path: "/:shard/json/miniblock/by-hash/:hash/epoch/:epoch", Handler: bg.internalMiniBlockbyHashHandler, Method: http.MethodGet},
		{Path: "/:shard/raw/miniblock/by-hash/:hash/epoch/:epoch", Handler: bg.rawMiniBlockbyHashHandler, Method: http.MethodGet},
		{Path: "/raw/startofepoch/metablock/by-epoch/:epoch", Handler: bg.rawStartOfEpochMetaBlock, Method: http.MethodGet},
		{Path: "/json/startofepoch/metablock/by-epoch/:epoch", Handler: bg.internalStartOfEpochMetaBlock, Method: http.MethodGet},
		{Path: "/json/startofepoch/validators/by-epoch/:epoch", Handler: bg.internalStartOfEpochValidatorsInfo, Method: http.MethodGet},
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

	epoch, err := shared.FetchEpochFromRequest(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseEpoch.Error(),
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

	miniBlockByHashResponse, err := group.facade.GetInternalMiniBlockByHash(shardID, hash, epoch, common.Internal)
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

	epoch, err := shared.FetchEpochFromRequest(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseEpoch.Error(),
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

	miniBlockByHashResponse, err := group.facade.GetInternalMiniBlockByHash(shardID, hash, epoch, common.Proto)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, miniBlockByHashResponse)
}

// internalStartOfEpochMetaBlock will handle the fetching and returning the start of epoch metablock by epoch
func (group *internalGroup) internalStartOfEpochMetaBlock(c *gin.Context) {
	epoch, err := shared.FetchEpochFromRequest(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseEpoch.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	miniBlockByHashResponse, err := group.facade.GetInternalStartOfEpochMetaBlock(epoch, common.Internal)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, miniBlockByHashResponse)
}

// rawStartOfEpochMetaBlock will handle the fetching and returning the start of epoch metablock by epoch
func (group *internalGroup) rawStartOfEpochMetaBlock(c *gin.Context) {
	epoch, err := shared.FetchEpochFromRequest(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseEpoch.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	miniBlockByHashResponse, err := group.facade.GetInternalStartOfEpochMetaBlock(epoch, common.Proto)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, miniBlockByHashResponse)
}

// internalStartOfEpochValidatorsInfo will handle the fetching and returning the start of epoch validators info by epoch
func (group *internalGroup) internalStartOfEpochValidatorsInfo(c *gin.Context) {
	epoch, err := shared.FetchEpochFromRequest(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseEpoch.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	validatorsInfo, err := group.facade.GetInternalStartOfEpochValidatorsInfo(epoch)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, validatorsInfo)
}
