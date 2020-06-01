package block

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
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

	shardIDStr := c.Param("shardID")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 32)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: "cannot parse shardID",
				Code:  string(data.ReturnCodeRequestErrror),
			},
		)
		return
	}

	nonceStr := c.Param("nonce")
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: "cannot parse nonce",
				Code:  string(data.ReturnCodeRequestErrror),
			},
		)
		return
	}

	apiBlock, err := ef.GetBlockByShardIDAndNonce(uint32(shardID), nonce)
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

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"block": apiBlock},
			Error: "",
			Code:  string(data.ReturnCodeSuccess),
		},
	)
}
