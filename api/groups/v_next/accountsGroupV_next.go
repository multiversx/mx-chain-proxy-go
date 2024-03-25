package v_next

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/groups"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type accountsGroupV_next struct {
	baseAccountsGroup data.GroupHandler
	facade            AccountsFacadeHandlerV_next
}

// NewAccountsGroupV_next returns a new instance of accountsGroupV_next
func NewAccountsGroupV_next(baseAccountsGroup data.GroupHandler, facadeHandler data.FacadeHandler) (*accountsGroupV_next, error) {
	if check.IfNil(baseAccountsGroup) {
		return nil, fmt.Errorf("nil base accounts group for v_next")
	}

	facade, ok := facadeHandler.(AccountsFacadeHandlerV_next)
	if !ok {
		return nil, groups.ErrWrongTypeAssertion
	}

	ag := &accountsGroupV_next{
		baseAccountsGroup: baseAccountsGroup,
		facade:            facade,
	}

	err := ag.baseAccountsGroup.UpdateEndpoint("/:address/shard", data.EndpointHandlerData{
		Handler: ag.GetShardForAccountV_next,
		Method:  http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	err = ag.baseAccountsGroup.RemoveEndpoint("/:address/nonce")
	if err != nil {
		return nil, err
	}

	err = ag.baseAccountsGroup.AddEndpoint("/:address/new-endpoint", data.EndpointHandlerData{
		Handler: ag.NewEndpoint,
		Method:  http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	return ag, nil
}

// NewEndpoint is an example of a new endpoint added in the version v_next
func (ag *accountsGroupV_next) NewEndpoint(c *gin.Context) {
	res := ag.facade.NextEndpointHandler()
	c.JSON(http.StatusOK, &data.GenericAPIResponse{
		Data:  res,
		Error: "",
		Code:  data.ReturnCodeSuccess,
	})
}

// GetShardForAccountV_next is an example of an updated endpoint in the version v_next
func (ag *accountsGroupV_next) GetShardForAccountV_next(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrComputeShardForAddress, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
		)
		return
	}

	shardID, err := ag.facade.GetShardIDForAddressV_next(addr, 0)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrComputeShardForAddress.Error(), err.Error()),
			data.ReturnCodeInternalError,
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"shardID": shardID}, "", data.ReturnCodeSuccess)
}

// Group returns the base accounts group
func (ag *accountsGroupV_next) Group() data.GroupHandler {
	return ag.baseAccountsGroup
}
