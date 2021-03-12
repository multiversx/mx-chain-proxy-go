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

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{
			Path:    "/by-hash/:hash",
			Handler: hbg.hyperBlockByHashHandler,
			Method:  http.MethodGet,
		},
		{
			Path:    "/by-nonce/:nonce",
			Handler: hbg.hyperBlockByNonceHandler,
			Method:  http.MethodGet,
		},
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

	blockByHashResponse, err := group.facade.GetHyperBlockByHash(hash)
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

	blockByNonceResponse, err := group.facade.GetHyperBlockByNonce(nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}
