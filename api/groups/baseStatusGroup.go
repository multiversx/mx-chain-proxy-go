package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type statusGroup struct {
	facade StatusFacadeHandler
	*baseGroup
}

// NewStatusGroup returns a new instance of statusGroup
func NewStatusGroup(facadeHandler data.FacadeHandler) (*statusGroup, error) {
	facade, ok := facadeHandler.(StatusFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &statusGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/metrics", Handler: ng.getMetrics, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func (group *statusGroup) getMetrics(c *gin.Context) {
	metricsResults, err := group.facade.GetMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"metrics": metricsResults}, "", data.ReturnCodeSuccess)
}
