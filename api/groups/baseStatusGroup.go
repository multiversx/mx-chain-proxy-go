package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
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
		{Path: "/prometheus-metrics", Handler: ng.getPrometheusMetrics, Method: http.MethodGet},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

// getMetrics will expose endpoints statistics in json format
func (group *statusGroup) getMetrics(c *gin.Context) {
	metricsResults := group.facade.GetMetrics()

	shared.RespondWith(c, http.StatusOK, gin.H{"metrics": metricsResults}, "", data.ReturnCodeSuccess)
}

// getPrometheusMetrics will expose proxy metrics in prometheus format
func (group *statusGroup) getPrometheusMetrics(c *gin.Context) {
	metricsResults := group.facade.GetMetricsForPrometheus()

	c.String(http.StatusOK, metricsResults)
}
