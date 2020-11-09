package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type validatorGroup struct {
	facade ValidatorFacadeHandler
	*baseGroup
}

// NewNodeGroup returns a new instance of nodeGroup
func NewValidatorGroup(facadeHandler data.FacadeHandler) (*validatorGroup, error) {
	facade, ok := facadeHandler.(ValidatorFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	vg := &validatorGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/statistics": {Handler: vg.Statistics, Method: http.MethodGet},
	}
	vg.baseGroup.endpoints = baseRoutesHandlers

	return vg, nil
}

// Statistics returns the validator statistics
func (vg *validatorGroup) Statistics(c *gin.Context) {
	validatorStatistics, err := vg.facade.ValidatorStatistics()
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, err.Error(), data.ReturnCodeRequestError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"statistics": validatorStatistics}, "", data.ReturnCodeSuccess)
}
