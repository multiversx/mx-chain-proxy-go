package groups

import (
	"encoding/hex"
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

func NewBaseHyperBlockGroup() *baseGroup {
	baseEndpointsHandlers := map[string]*data.EndpointHandlerData{
		"/by-hash/:hash":   {Handler: HyperBlockByHashHandler, Method: http.MethodGet},
		"/by-nonce/:nonce": {Handler: HyperBlockByNonceHandler, Method: http.MethodGet},
	}

	return &baseGroup{
		endpoints: baseEndpointsHandlers,
	}
}

// HyperBlockByHashHandler handles "by-hash" requests
func HyperBlockByHashHandler(c *gin.Context) {
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
func HyperBlockByNonceHandler(c *gin.Context) {
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
