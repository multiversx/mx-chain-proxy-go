package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/gin-gonic/gin"
)

func NewBaseNetworkGroup() *baseGroup {
	baseEndpointsHandlers := map[string]*data.EndpointHandlerData{
		"/status/:shard": {Handler: GetNetworkStatusData, Method: http.MethodGet},
		"/config":        {Handler: GetNetworkConfigData, Method: http.MethodGet},
		"/economics":     {Handler: GetEconomicsData, Method: http.MethodGet},
	}

	return &baseGroup{
		endpoints: baseEndpointsHandlers,
	}
}

// GetNetworkStatusData will expose the node network metrics for the given shard
func GetNetworkStatusData(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(NetworkFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	shardIDUint, err := shared.FetchShardIDFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, process.ErrInvalidShardId.Error(), data.ReturnCodeRequestError)
		return
	}

	networkStatusResults, err := ef.GetNetworkStatusMetrics(uint32(shardIDUint))
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkStatusResults)
}

// GetNetworkConfigData will expose the node network metrics for the given shard
func GetNetworkConfigData(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(NetworkFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	networkConfigResults, err := ef.GetNetworkConfigMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, networkConfigResults)
}

// GetEconomicsData will expose the economics data metrics from an observer (if any available) in json format
func GetEconomicsData(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(NetworkFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	economicsData, err := ef.GetEconomicsDataMetrics()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, economicsData)
}
