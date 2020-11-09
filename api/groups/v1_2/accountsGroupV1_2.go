package v1_2

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-go-logger/check"
	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/groups"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type accountsGroupV1_2 struct {
	baseAccountsGroup data.GroupHandler
	facade            AccountsFacadeHandlerV1_2
}

// NewAccountsGroupV1_2 returns a new instance of accountsGroupV1_2
func NewAccountsGroupV1_2(baseAccountsGroup data.GroupHandler, facadeHandler data.FacadeHandler) (*accountsGroupV1_2, error) {
	if check.IfNil(baseAccountsGroup) {
		return nil, fmt.Errorf("nil base accounts group for v1.2")
	}

	facade, ok := facadeHandler.(AccountsFacadeHandlerV1_2)
	if !ok {
		return nil, groups.ErrWrongTypeAssertion
	}

	ag := &accountsGroupV1_2{
		baseAccountsGroup: baseAccountsGroup,
		facade:            facade,
	}

	err := ag.baseAccountsGroup.UpdateEndpoint("/:address/shard", data.EndpointHandlerData{
		Handler: ag.GetShardForAccountV1_2,
		Method:  http.MethodGet,
	})
	if err != nil {
		return nil, err
	}

	return ag, nil
}

// GetShardForAccountV1_2 handles the request for account's shard in the version v1.2
func (ag *accountsGroupV1_2) GetShardForAccountV1_2(c *gin.Context) {
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

	shardID, err := ag.facade.GetShardIDForAddressV1_2(addr, 0)
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
func (ag *accountsGroupV1_2) Group() data.GroupHandler {
	return ag.baseAccountsGroup
}
