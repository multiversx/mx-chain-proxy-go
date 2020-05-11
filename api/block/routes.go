package block

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/gin-gonic/gin"
)

// Routes defines blocks-related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/meta/latest-nonce", GetHighestBlockNonce)
	router.GET("/meta/:nonce", GetBlockByNonce)
}

// GetHighestBlockNonce returns the highest metablock nonce
func GetHighestBlockNonce(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	highestNonce, err := ef.GetHighestBlockNonce()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nonce": highestNonce})
}

// GetBlockByNonce returns the metablock by nonce
func GetBlockByNonce(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	nonceStr := c.Param("nonce")
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiBlock, err := ef.GetBlockByNonce(nonce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"block": apiBlock})
}
