package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
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
		{Path: "/nodes-versions", Handler: ag.getNodesVersions, Method: http.MethodGet},
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

func (ag *aboutGroup) getNodesVersions(c *gin.Context) {
	nodesVersions, err := ag.facade.GetNodesVersions()
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, nodesVersions)
}
