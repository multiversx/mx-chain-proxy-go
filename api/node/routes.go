package node

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/heartbeatstatus", GetHeartbeatData)
}

// GetHeartbeatData will expose heartbeat status from an observer (if any available) in json format
func GetHeartbeatData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	heartbeatResults, err := ef.GetHeartbeatData()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"heartbeats": heartbeatResults.Heartbeats},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}
