package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
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

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/heartbeatstatus": {Handler: ng.getHeartbeatData, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func (ng *nodeGroup) getHeartbeatData(c *gin.Context) {
	heartbeatResults, err := ng.facade.GetHeartbeatData()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"heartbeats": heartbeatResults.Heartbeats}, "", data.ReturnCodeSuccess)
}
