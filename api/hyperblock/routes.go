package hyperblock

import (
	"encoding/hex"
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines the HTTP routes
func Routes(router *gin.RouterGroup) {
	router.GET("/by-hash/:hash", ByHashHandler)
	router.GET("/by-nonce/:nonce", ByNonceHandler)
}

// ByHashHandler handles "by-hash" requests
func ByHashHandler(c *gin.Context) {
	epf, ok := c.MustGet("elrondProxyFacade").(facadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	hash := c.Param("hash")
	_, err := hex.DecodeString(hash)
	if err != nil {
		shared.ResponsWithBadRequest(c, apiErrors.ErrInvalidBlockHashParam.Error())
		return
	}

	blockByHashResponse, err := epf.GetHyperBlockByHash(hash)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByHashResponse)
}

// ByNonceHandler handles "by-nonce" requests
func ByNonceHandler(c *gin.Context) {
	epf, ok := c.MustGet("elrondProxyFacade").(facadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.ResponsWithBadRequest(c, apiErrors.ErrCannotParseNonce.Error())
		return
	}

	blockByNonceResponse, err := epf.GetHyperBlockByNonce(nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByNonceResponse)
}
