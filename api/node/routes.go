package node

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
	router.GET("/heartbeatstatus", GetHeartbeatData)
	router.GET("/status/:shard", GetNodeStatus)
}

// GetHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func GetHeartbeatData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	heartbeatResults, err := ef.GetHeartbeatData()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"heartbeats": heartbeatResults.Heartbeats},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetNodeStatus will expose the node status for the given shard
func GetNodeStatus(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
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
				Code:  data.ReturnCodeRequestError,
			},
		)
		return
	}

	nodeStatusResults, err := ef.GetShardStatus(uint32(shardIDUint))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  nodeStatusResults,
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}
