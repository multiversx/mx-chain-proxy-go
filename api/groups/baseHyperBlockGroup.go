package groups

import (
	"encoding/hex"
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
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

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/by-hash/:hash":   {Handler: hbg.HyperBlockByHashHandler, Method: http.MethodGet},
		"/by-nonce/:nonce": {Handler: hbg.HyperBlockByNonceHandler, Method: http.MethodGet},
	}
	hbg.baseGroup.endpoints = baseRoutesHandlers

	return hbg, nil
}

// HyperBlockByHashHandler handles "by-hash" requests
func (hbg *hyperBlockGroup) HyperBlockByHashHandler(c *gin.Context) {
	epf, ok := c.MustGet(shared.GetFacadeVersion(c)).(HyperBlockFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrInvalidBlockHashParam.Error())
		return
	}

	blockByHashResponse, err := epf.GetHyperBlockByHash(hash)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// HyperBlockByNonceHandler handles "by-nonce" requests
func (hbg *hyperBlockGroup) HyperBlockByNonceHandler(c *gin.Context) {
	epf, ok := c.MustGet(shared.GetFacadeVersion(c)).(HyperBlockFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrCannotParseNonce.Error())
		return
	}

	blockByNonceResponse, err := epf.GetHyperBlockByNonce(nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}
