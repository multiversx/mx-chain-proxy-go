package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/gin-gonic/gin"
)

type networkGroup struct {
	facade NetworkFacadeHandler
	*baseGroup
}

// NewNetworkGroup returns a new instance of networkGroup
func NewNetworkGroup(facadeHandler data.FacadeHandler) (*networkGroup, error) {
	facade, ok := facadeHandler.(NetworkFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &networkGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/status/:shard": {Handler: ng.GetNetworkStatusData, Method: http.MethodGet},
		"/config":        {Handler: ng.GetNetworkConfigData, Method: http.MethodGet},
		"/economics":     {Handler: ng.GetEconomicsData, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// GetNetworkStatusData will expose the node network metrics for the given shard
func (ng *networkGroup) GetNetworkStatusData(c *gin.Context) {
	shardIDUint, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, process.ErrInvalidShardId.Error(), data.ReturnCodeRequestError)
		return
	}

	networkStatusResults, err := ng.facade.GetNetworkStatusMetrics(shardIDUint)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkStatusResults)
}

// GetNetworkConfigData will expose the node network metrics for the given shard
func (ng *networkGroup) GetNetworkConfigData(c *gin.Context) {
	networkConfigResults, err := ng.facade.GetNetworkConfigMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkConfigResults)
}

// GetEconomicsData will expose the economics data metrics from an observer (if any available) in json format
func (ng *networkGroup) GetEconomicsData(c *gin.Context) {
	economicsData, err := ng.facade.GetEconomicsDataMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, economicsData)
}
