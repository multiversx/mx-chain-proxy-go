package hyperblock

import (
	"encoding/hex"
	"net/http"
	"strconv"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines full blocks related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/by-nonce/:nonce", ByNonceHandler)
	router.GET("/by-hash/:hash", ByHashHandler)
}

// ByHashHandler will handle the fetching and returning a block based on its hash
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

// ByNonceHandler will handle the fetching and returning a block based on its nonce
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

func getQueryParamWithTxs(c *gin.Context) (bool, error) {
	withTxsStr := c.Request.URL.Query().Get("withTxs")
	if withTxsStr == "" {
		return false, nil
	}

	return strconv.ParseBool(withTxsStr)
}
