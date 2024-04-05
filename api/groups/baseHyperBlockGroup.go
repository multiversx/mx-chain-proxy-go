package groups

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type hyperBlockGroup struct {
	facade HyperBlockFacadeHandler
	*baseGroup
}

// NewHyperBlockGroup returns a new instance of hyperBlockGroup
func NewHyperBlockGroup(facadeHandler data.FacadeHandler) (*hyperBlockGroup, error) {
	facade, ok := facadeHandler.(HyperBlockFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	hbg := &hyperBlockGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/by-hash/:hash", Handler: hbg.hyperBlockByHashHandler, Method: http.MethodGet},
		{Path: "/by-nonce/:nonce", Handler: hbg.hyperBlockByNonceHandler, Method: http.MethodGet},
	}
	hbg.baseGroup.endpoints = baseRoutesHandlers

	return hbg, nil
}

// hyperBlockByHashHandler handles "by-hash" requests
func (group *hyperBlockGroup) hyperBlockByHashHandler(c *gin.Context) {
	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrInvalidBlockHashParam.Error())
		return
	}

	options, err := parseHyperblockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByHashResponse, err := group.facade.GetHyperBlockByHash(hash, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// hyperBlockByNonceHandler handles "by-nonce" requests
func (group *hyperBlockGroup) hyperBlockByNonceHandler(c *gin.Context) {
	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrCannotParseNonce.Error())
		return
	}

	options, err := parseHyperblockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByNonceResponse, err := group.facade.GetHyperBlockByNonce(nonce, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}
