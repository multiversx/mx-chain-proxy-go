package network

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/:shard", GetNetworkData)
}

// GetNetworkData will expose the node network metrics for the given shard
func GetNetworkData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	shardID := c.Param("shard")
	shardIDUint, err := strconv.ParseUint(shardID, 10, 32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": process.ErrInvalidShardId})
		return
	}

	networkResults, err := ef.GetNetworkMetrics(uint32(shardIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": networkResults})
}
