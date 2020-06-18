package block

import (
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
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
		shared.RespondWithInvalidAppContext(c)
		return
	}

	shardIDStr := c.Param("shardID")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 32)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrCannotParseShardID.Error(), data.ReturnCodeRequestError)
		return
	}

	nonceStr := c.Param("nonce")
	nonce, err := strconv.ParseUint(nonceStr, 10, 64)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrCannotParseNonce.Error(), data.ReturnCodeRequestError)
		return
	}

	apiBlock, err := ef.GetBlockByShardIDAndNonce(uint32(shardID), nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"block": apiBlock}, "", data.ReturnCodeSuccess)
}
