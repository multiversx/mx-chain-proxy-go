package groups

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type aboutGroup struct {
	facade AboutFacadeHandler
	*baseGroup
}

// NewAboutGroup returns a new instance of aboutGroup
func NewAboutGroup(facadeHandler data.FacadeHandler) (*aboutGroup, error) {
	facade, ok := facadeHandler.(AboutFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}
	ag := &aboutGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "", Handler: ag.getAboutInfo, Method: http.MethodGet},
	}
	ag.baseGroup.endpoints = baseRoutesHandlers

	return ag, nil
}

func (ag *aboutGroup) getAboutInfo(c *gin.Context) {
	aboutInfo, err := ag.facade.GetAboutInfo()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, aboutInfo)
}
