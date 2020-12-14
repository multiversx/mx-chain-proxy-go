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
		"/:address":                       {Handler: ag.getAccount, Method: http.MethodGet},
		"/:address/balance":               {Handler: ag.getBalance, Method: http.MethodGet},
		"/:address/username":              {Handler: ag.getUsername, Method: http.MethodGet},
		"/:address/nonce":                 {Handler: ag.getNonce, Method: http.MethodGet},
		"/:address/shard":                 {Handler: ag.getShard, Method: http.MethodGet},
		"/:address/transactions":          {Handler: ag.getTransactions, Method: http.MethodGet},
		"/:address/key/:key":              {Handler: ag.getValueForKey, Method: http.MethodGet},
		"/:address/esdt":                  {Handler: ag.getESDTTokens, Method: http.MethodGet},
		"/:address/esdt/:tokenIdentifier": {Handler: ag.getESDTTokenData, Method: http.MethodGet},
	}
	ag.baseGroup.endpoints = baseRoutesHandlers

	return ag, nil
}

func (group *accountsGroup) getAccountFromFacade(c *gin.Context) (*data.Account, int, error) {
	addr := c.Param("address")
	acc, status, err := group.facade.GetAccount(addr)
	if err != nil {
		return nil, status, err
	}

	return acc, http.StatusOK, nil
}

func (group *accountsGroup) getTransactionsFromFacade(c *gin.Context) ([]data.DatabaseTransaction, int, error) {
	addr := c.Param("address")
	transactions, status, err := group.facade.GetTransactions(addr)
	if err != nil {
		return nil, status, err
	}

	return transactions, http.StatusOK, nil
}

// getAccount returns an accountResponse containing information
// about the account correlated with provided address
func (group *accountsGroup) getAccount(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"account": account}, "")
}

// getBalance returns the balance for the address parameter
func (group *accountsGroup) getBalance(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"balance": account.Balance}, "")
}

// getUsername returns the username for the address parameter
func (group *accountsGroup) getUsername(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"username": account.Username}, "")
}

// getNonce returns the nonce for the address parameter
func (group *accountsGroup) getNonce(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"nonce": account.Nonce}, "")
}

// getTransactions returns the transactions for the address parameter
func (group *accountsGroup) getTransactions(c *gin.Context) {
	transactions, status, err := group.getTransactionsFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transactions": transactions}, "")
}

// getValueForKey returns the value for the given address and key
func (group *accountsGroup) getValueForKey(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetValueForKey, errors.ErrEmptyAddress),
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
		)
		return
	}

	value, status, err := group.facade.GetValueForKey(addr, key)
	if err != nil {
		shared.RespondWith(
			c,
			status,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrGetValueForKey.Error(), err.Error()),
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"value": value}, "")
}

// getShard returns the shard for the given address based on the current proxy's configuration
func (group *accountsGroup) getShard(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrComputeShardForAddress, errors.ErrEmptyAddress),
		)
		return
	}

	shardID, status, err := group.facade.GetShardIDForAddress(addr)
	if err != nil {
		shared.RespondWith(
			c,
			status,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrComputeShardForAddress.Error(), err.Error()),
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"shardID": shardID}, "")
}

// GetESDTTokenData returns the balance for the given address and esdt token
func (group *accountsGroup) getESDTTokenData(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTTokenData, errors.ErrEmptyAddress),
		)
		return
	}

	tokenIdentifier := c.Param("tokenIdentifier")
	if tokenIdentifier == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTTokenData, errors.ErrEmptyTokenIdentifier),
		)
		return
	}

	esdtTokenResponse, status, err := group.facade.GetESDTTokenData(addr, tokenIdentifier)
	if err != nil {
		shared.RespondWith(
			c,
			status,
			nil,
			err.Error(),
		)
		return
	}

	c.JSON(http.StatusOK, esdtTokenResponse)
}

// getESDTTokens returns the tokens list from this account
func (group *accountsGroup) getESDTTokens(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTTokenData, errors.ErrEmptyAddress),
		)
		return
	}

	tokens, status, err := group.facade.GetAllESDTTokens(addr)
	if err != nil {
		shared.RespondWith(
			c,
			status,
			nil,
			err.Error(),
		)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
