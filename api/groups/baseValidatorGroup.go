package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

func NewBaseValidatorGroup() *baseGroup {
	baseEndpointsHandlers := map[string]*data.EndpointHandlerData{
		"/statistics": {Handler: Statistics, Method: http.MethodGet},
	}

	return &baseGroup{
		endpoints: baseEndpointsHandlers,
	}
}

// Statistics returns the validator statistics
func Statistics(c *gin.Context) {
	epf, ok := c.MustGet(shared.GetFacadeVersion(c)).(ValidatorFacadeHandler)
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
