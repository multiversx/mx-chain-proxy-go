package block

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/gin-gonic/gin"
)

// Routes defines blocks-related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/:shardID/:nonce", GetBlockByShardIDAndNonce)
}

// GetBlockByShardIDAndNonce returns the block by shardID and nonce
func GetBlockByShardIDAndNonce(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	shardIDStr := c.Param("shardID")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse shardID"})
		return
	}

	nonceStr := c.Param("nonce")
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot parse nonce"})
		return
	}

	apiBlock, err := ef.GetBlockByShardIDAndNonce(uint32(shardID), nonce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"block": apiBlock})
}
