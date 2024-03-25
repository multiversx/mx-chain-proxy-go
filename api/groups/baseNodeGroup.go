package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type nodeGroup struct {
	facade NodeFacadeHandler
	*baseGroup
}

// NewNodeGroup returns a new instance of nodeGroup
func NewNodeGroup(facadeHandler data.FacadeHandler) (*nodeGroup, error) {
	facade, ok := facadeHandler.(NodeFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &nodeGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/heartbeatstatus", Handler: ng.getHeartbeatData, Method: http.MethodGet},
		{Path: "/old-storage-token/:token/nonce/:nonce", Handler: ng.isOldStorageForToken, Method: http.MethodGet},
		{Path: "/waiting-epochs-left/:key", Handler: ng.waitingEpochsLeft, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func (group *nodeGroup) getHeartbeatData(c *gin.Context) {
	heartbeatResults, err := group.facade.GetHeartbeatData()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"heartbeats": heartbeatResults.Heartbeats}, "", data.ReturnCodeSuccess)
}

func (group *nodeGroup) isOldStorageForToken(c *gin.Context) {
	// TODO: when the old storage tokens liquidity issue is solved on the protocol, mark this endpoint as deprecated
	// and remove the processing code
	token := c.Param("token")
	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			apiErrors.ErrCannotParseNonce.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}
	isOldStorage, err := group.facade.IsOldStorageForToken(token, nonce)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"isOldStorage": isOldStorage}, "", data.ReturnCodeSuccess)
}

func (group *nodeGroup) waitingEpochsLeft(c *gin.Context) {
	publicKey := c.Param("key")
	response, err := group.facade.GetWaitingEpochsLeftForPublicKey(publicKey)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, response.Data, "", data.ReturnCodeSuccess)
}
