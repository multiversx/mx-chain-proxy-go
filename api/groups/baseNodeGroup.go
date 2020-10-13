package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

func NewBaseNodeGroup() *baseGroup {
	baseEndpointsHandlers := map[string]*shared.EndpointHandlerData{
		"/heartbeatstatus": {Handler: GetHeartbeatData, Method: http.MethodGet},
	}

	return &baseGroup{
		endpoints: baseEndpointsHandlers,
	}
}

// GetHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func GetHeartbeatData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(NodeFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	heartbeatResults, err := ef.GetHeartbeatData()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"heartbeats": heartbeatResults.Heartbeats}, "", data.ReturnCodeSuccess)
}
