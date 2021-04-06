package groups

import (
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type blockAtlasGroup struct {
	facade BlockAtlasFacadeHandler
	*baseGroup
}

// NewBlockAtlasGroup returns a new instance of blockAtlasGroup
func NewBlockAtlasGroup(facadeHandler data.FacadeHandler) (*blockAtlasGroup, error) {
	facade, ok := facadeHandler.(BlockAtlasFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	bag := &blockAtlasGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/:shard/:nonce", Handler: bag.getBlockByShardIDAndNonceFromElastic, Method: http.MethodGet},
	}
	bag.baseGroup.endpoints = baseRoutesHandlers

	return bag, nil
}

// getBlockByShardIDAndNonceFromElastic returns the block by shardID and nonce
func (group *blockAtlasGroup) getBlockByShardIDAndNonceFromElastic(c *gin.Context) {
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

	apiBlock, err := group.facade.GetAtlasBlockByShardIDAndNonce(shardID, nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"block": apiBlock}, "", data.ReturnCodeSuccess)
}
