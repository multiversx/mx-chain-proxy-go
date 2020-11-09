package groups

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type accountsGroup struct {
	facade AccountsFacadeHandler
	*baseGroup
}

// NewAccountsGroup returns a new instance of accountsGroup
func NewAccountsGroup(facadeHandler data.FacadeHandler) (*accountsGroup, error) {
	facade, ok := facadeHandler.(AccountsFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	ag := &accountsGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/:address":              {Handler: ag.GetAccount, Method: http.MethodGet},
		"/:address/balance":      {Handler: ag.GetBalance, Method: http.MethodGet},
		"/:address/username":     {Handler: ag.GetUsername, Method: http.MethodGet},
		"/:address/nonce":        {Handler: ag.GetNonce, Method: http.MethodGet},
		"/:address/shard":        {Handler: ag.GetShard, Method: http.MethodGet},
		"/:address/transactions": {Handler: ag.GetTransactions, Method: http.MethodGet},
		"/:address/key/:key":     {Handler: ag.GetValueForKey, Method: http.MethodGet},
	}
	ag.baseGroup.endpoints = baseRoutesHandlers

	return ag, nil
}

func (ag *accountsGroup) getAccount(c *gin.Context) (*data.Account, int, error) {
	addr := c.Param("address")
	acc, err := ag.facade.GetAccount(addr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return acc, http.StatusOK, nil
}

func (ag *accountsGroup) getTransactions(c *gin.Context) ([]data.DatabaseTransaction, int, error) {
	addr := c.Param("address")
	transactions, err := ag.facade.GetTransactions(addr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return transactions, http.StatusOK, nil
}

// GetAccount returns an accountResponse containing information
// about the account correlated with provided address
func (ag *accountsGroup) GetAccount(c *gin.Context) {
	account, status, err := ag.getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"account": account}, "", data.ReturnCodeSuccess)
}

// GetBalance returns the balance for the address parameter
func (ag *accountsGroup) GetBalance(c *gin.Context) {
	account, status, err := ag.getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"balance": account.Balance}, "", data.ReturnCodeSuccess)
}

// GetUsername returns the username for the address parameter
func (ag *accountsGroup) GetUsername(c *gin.Context) {
	account, status, err := ag.getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"username": account.Username}, "", data.ReturnCodeSuccess)
}

// GetNonce returns the nonce for the address parameter
func (ag *accountsGroup) GetNonce(c *gin.Context) {
	account, status, err := ag.getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"nonce": account.Nonce}, "", data.ReturnCodeSuccess)
}

// GetTransactions returns the transactions for the address parameter
func (ag *accountsGroup) GetTransactions(c *gin.Context) {
	transactions, status, err := ag.getTransactions(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transactions": transactions}, "", data.ReturnCodeSuccess)
}

// GetValueForKey returns the value for the given address and key
func (ag *accountsGroup) GetValueForKey(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetValueForKey, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
		)
		return
	}

	key := c.Param("key")
	if key == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetValueForKey, errors.ErrEmptyKey),
			data.ReturnCodeRequestError,
		)
		return
	}

	value, err := ag.facade.GetValueForKey(addr, key)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrGetValueForKey.Error(), err.Error()),
			data.ReturnCodeInternalError,
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"value": value}, "", data.ReturnCodeSuccess)
}

// GetShard returns the shard for the given address based on the current proxy's configuration
func (ag *accountsGroup) GetShard(c *gin.Context) {
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

	shardID, err := ag.facade.GetShardIDForAddress(addr)
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
