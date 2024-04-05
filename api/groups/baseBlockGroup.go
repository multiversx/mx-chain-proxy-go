package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type blockGroup struct {
	facade BlockFacadeHandler
	*baseGroup
}

// NewBlockGroup returns a new instance of blockGroup
func NewBlockGroup(facadeHandler data.FacadeHandler) (*blockGroup, error) {
	facade, ok := facadeHandler.(BlockFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	bg := &blockGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/:shard/by-nonce/:nonce", Handler: bg.byNonceHandler, Method: http.MethodGet},
		{Path: "/:shard/by-hash/:hash", Handler: bg.byHashHandler, Method: http.MethodGet},
		{Path: "/:shard/altered-accounts/by-nonce/:nonce", Handler: bg.alteredAccountsByNonceHandler, Method: http.MethodGet},
		{Path: "/:shard/altered-accounts/by-hash/:hash", Handler: bg.alteredAccountsByHashHandler, Method: http.MethodGet},
	}
	bg.baseGroup.endpoints = baseRoutesHandlers

	return bg, nil
}

// byHashHandler will handle the fetching and returning a block based on its hash
func (group *blockGroup) byHashHandler(c *gin.Context) {
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

	hash, err := shared.FetchHashFromRequest(c)
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

	options, err := parseBlockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByHashResponse, err := group.facade.GetBlockByHash(shardID, hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// byNonceHandler will handle the fetching and returning a block based on its nonce
func (group *blockGroup) byNonceHandler(c *gin.Context) {
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

	options, err := parseBlockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByNonceResponse, err := group.facade.GetBlockByNonce(shardID, nonce, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}

func (group *blockGroup) alteredAccountsByNonceHandler(c *gin.Context) {
	shardID, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrCannotParseShardID.Error())
		return
	}

	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrCannotParseNonce.Error())
		return
	}

	options, err := parseAlteredAccountOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByNonceResponse, err := group.facade.GetAlteredAccountsByNonce(shardID, nonce, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}

func (group *blockGroup) alteredAccountsByHashHandler(c *gin.Context) {
	shardID, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrCannotParseShardID.Error())
		return
	}

	hash, err := shared.FetchHashFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrInvalidBlockHashParam.Error())
		return
	}

	options, err := parseAlteredAccountOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByHashResponse, err := group.facade.GetAlteredAccountsByHash(shardID, hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}
