package groups

import (
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

func NewBaseBlockAtlasGroup() *baseGroup {
	baseEndpointsHandlers := map[string]*data.EndpointHandlerData{
		"/:shard/:nonce": {Handler: GetBlockByShardIDAndNonceFromElastic, Method: http.MethodGet},
	}

	return &baseGroup{
		endpoints: baseEndpointsHandlers,
	}
}

// GetBlockByShardIDAndNonceFromElastic returns the block by shardID and nonce
func GetBlockByShardIDAndNonceFromElastic(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(BlockAtlasFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	shardID, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, apiErrors.ErrCannotParseShardID.Error(), data.ReturnCodeRequestError)
		return
	}

	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, apiErrors.ErrCannotParseNonce.Error(), data.ReturnCodeRequestError)
		return
	}

	apiBlock, err := ef.GetAtlasBlockByShardIDAndNonce(uint32(shardID), nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"block": apiBlock}, "", data.ReturnCodeSuccess)
}
