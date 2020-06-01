package validator

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/statistics", Statistics)
}

// Statistics returns the validator statistics
func Statistics(c *gin.Context) {
	epf, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  string(data.ReturnCodeInternalError),
			},
		)
		return
	}

	validatorStatistics, err := epf.ValidatorStatistics()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  string(data.ReturnCodeRequestErrror),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"statistics": validatorStatistics},
			Error: "",
			Code:  string(data.ReturnCodeSuccess),
		},
	)
}
