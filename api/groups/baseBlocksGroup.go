package groups

import (
	"net/http"

	apiErrors "github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type blocksGroup struct {
	facade BlocksFacadeHandler
	*baseGroup
}

func NewBlocksGroup(facadeHandler data.FacadeHandler) (*blocksGroup, error) {
	facade, ok := facadeHandler.(BlocksFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	bbg := &blocksGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}
	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/by-round/:round", Handler: bbg.byRoundHandler, Method: http.MethodGet},
	}
	bbg.baseGroup.endpoints = baseRoutesHandlers

	return bbg, nil
}

func (bbp *blocksGroup) byRoundHandler(c *gin.Context) {
	round, err := shared.FetchRoundFromRequest(c)
	if err != nil {
		shared.RespondWithBadRequest(c, apiErrors.ErrCannotParseRound.Error())
		return
	}

	options, err := parseBlockQueryOptions(c)
	if err != nil {
		shared.RespondWithValidationError(c, apiErrors.ErrBadUrlParams, err)
		return
	}

	blockByRoundResponse, err := bbp.facade.GetBlocksByRound(round, options)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, blockByRoundResponse)
}
