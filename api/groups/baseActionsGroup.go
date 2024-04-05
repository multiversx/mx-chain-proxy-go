package groups

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type actionsGroup struct {
	facade ActionsFacadeHandler
	*baseGroup
}

// NewActionsGroup returns a new instance of actionGroup
func NewActionsGroup(facadeHandler data.FacadeHandler) (*actionsGroup, error) {
	facade, ok := facadeHandler.(ActionsFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ng := &actionsGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/reload-observers", Handler: ng.updateObservers, Method: http.MethodPost},
		{Path: "/reload-full-history-observers", Handler: ng.updateFullHistoryObservers, Method: http.MethodPost},
	}
	ng.baseGroup.endpoints = baseRoutesHandlers

	return ng, nil
}

func (group *actionsGroup) updateObservers(c *gin.Context) {
	result := group.facade.ReloadObservers()
	group.handleUpdateResponding(result, c)
}

func (group *actionsGroup) updateFullHistoryObservers(c *gin.Context) {
	result := group.facade.ReloadFullHistoryObservers()
	group.handleUpdateResponding(result, c)
}

func (group *actionsGroup) handleUpdateResponding(result data.NodesReloadResponse, c *gin.Context) {
	if result.Error != "" {
		httpCode := http.StatusInternalServerError
		internalCode := data.ReturnCodeInternalError
		if !result.OkRequest {
			httpCode = http.StatusBadRequest
			internalCode = data.ReturnCodeRequestError
		}

		shared.RespondWith(c, httpCode, result.Description, result.Error, internalCode)
		return
	}

	shared.RespondWith(c, http.StatusOK, result.Description, "", data.ReturnCodeSuccess)
}
