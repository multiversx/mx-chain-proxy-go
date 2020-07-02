package validator

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
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
		shared.RespondWithInvalidAppContext(c)
		return
	}

	validatorStatistics, err := epf.ValidatorStatistics()
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, err.Error(), data.ReturnCodeRequestError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"statistics": validatorStatistics}, "", data.ReturnCodeSuccess)
}
