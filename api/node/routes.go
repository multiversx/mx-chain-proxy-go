package node

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/heartbeatstatus", GetHeartbeatData)
	router.GET("/status/:shard", GetNodeStatus)
	router.GET("/epoch/:shard", GetEpochData)
	router.GET("/config", GetConfigData)
}

// GetHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func GetHeartbeatData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	heartbeatResults, err := ef.GetHeartbeatData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": heartbeatResults.Heartbeats})
}

// GetNodeStatus will expose the node status for the given shard
func GetNodeStatus(c *gin.Context) {
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

	nodeStatusResults, err := ef.GetShardStatus(uint32(shardIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": nodeStatusResults})
}

// GetEpochData will expose the node epoch metrics for the given shard
func GetEpochData(c *gin.Context) {
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

	nodeEpochResults, err := ef.GetEpochMetrics(uint32(shardIDUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": nodeEpochResults})
}

// GetConfigData will expose the node configuration metrics
func GetConfigData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	nodeConfigResults, err := ef.GetConfigMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": nodeConfigResults})
}
