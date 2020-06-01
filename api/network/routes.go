package network

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/status/:shard", GetNetworkStatusData)
	router.GET("/config", GetNetworkConfigData)
}

// GetNetworkStatusData will expose the node network metrics for the given shard
func GetNetworkStatusData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  string(data.ReturnCodeInternalError),
			},
		)
		return
	}

	shardID := c.Param("shard")
	shardIDUint, err := strconv.ParseUint(shardID, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: process.ErrInvalidShardId.Error(),
				Code:  string(data.ReturnCodeRequestErrror),
			},
		)
		return
	}

	networkStatusResults, err := ef.GetNetworkStatusMetrics(uint32(shardIDUint))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  string(data.ReturnCodeInternalError),
			},
		)
		return
	}

	c.JSON(http.StatusOK, networkStatusResults)
}

// GetNetworkConfigData will expose the node network metrics for the given shard
func GetNetworkConfigData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  string(data.ReturnCodeInternalError),
			},
		)
		return
	}

	networkConfigResults, err := ef.GetNetworkConfigMetrics()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  string(data.ReturnCodeInternalError),
			},
		)
		return
	}

	c.JSON(http.StatusOK, networkConfigResults)
}
