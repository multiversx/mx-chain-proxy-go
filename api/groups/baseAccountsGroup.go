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

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/:address", Handler: ag.getAccount, Method: http.MethodGet},
		{Path: "/:address/balance", Handler: ag.getBalance, Method: http.MethodGet},
		{Path: "/:address/username", Handler: ag.getUsername, Method: http.MethodGet},
		{Path: "/:address/nonce", Handler: ag.getNonce, Method: http.MethodGet},
		{Path: "/:address/shard", Handler: ag.getShard, Method: http.MethodGet},
		{Path: "/:address/transactions", Handler: ag.getTransactions, Method: http.MethodGet},
		{Path: "/:address/keys", Handler: ag.getKeyValuePairs, Method: http.MethodGet},
		{Path: "/:address/key/:key", Handler: ag.getValueForKey, Method: http.MethodGet},
		{Path: "/:address/esdt", Handler: ag.getESDTTokens, Method: http.MethodGet},
		{Path: "/:address/esdt/:tokenIdentifier", Handler: ag.getESDTTokenData, Method: http.MethodGet},
		{Path: "/:address/esdts-with-role/:role", Handler: ag.getESDTsWithRole, Method: http.MethodGet},
		{Path: "/:address/owned-nfts", Handler: ag.getOwnedNFTs, Method: http.MethodGet},
		{Path: "/:address/nft/:tokenIdentifier/nonce/:nonce", Handler: ag.getESDTNftTokenData, Method: http.MethodGet},
	}
	ag.baseGroup.endpoints = baseRoutesHandlers

	return ag, nil
}

func (group *accountsGroup) getAccountFromFacade(c *gin.Context) (*data.Account, int, error) {
	addr := c.Param("address")
	acc, err := group.facade.GetAccount(addr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return acc, http.StatusOK, nil
}

func (group *accountsGroup) getTransactionsFromFacade(c *gin.Context) ([]data.DatabaseTransaction, int, error) {
	addr := c.Param("address")
	transactions, err := group.facade.GetTransactions(addr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return transactions, http.StatusOK, nil
}

// getAccount returns an accountResponse containing information
// about the account correlated with provided address
func (group *accountsGroup) getAccount(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"account": account}, "", data.ReturnCodeSuccess)
}

// getBalance returns the balance for the address parameter
func (group *accountsGroup) getBalance(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"balance": account.Balance}, "", data.ReturnCodeSuccess)
}

// getUsername returns the username for the address parameter
func (group *accountsGroup) getUsername(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"username": account.Username}, "", data.ReturnCodeSuccess)
}

// getNonce returns the nonce for the address parameter
func (group *accountsGroup) getNonce(c *gin.Context) {
	account, status, err := group.getAccountFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"nonce": account.Nonce}, "", data.ReturnCodeSuccess)
}

// getTransactions returns the transactions for the address parameter
func (group *accountsGroup) getTransactions(c *gin.Context) {
	transactions, status, err := group.getTransactionsFromFacade(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transactions": transactions}, "", data.ReturnCodeSuccess)
}

// getKeyValuePairs returns the key-value pairs for the address parameter
func (group *accountsGroup) getKeyValuePairs(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetKeyValuePairs, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
		)
		return
	}

	keyValuePairs, err := group.facade.GetKeyValuePairs(addr)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(http.StatusOK, keyValuePairs)
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

	value, err := group.facade.GetValueForKey(addr, key)
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

// getShard returns the shard for the given address based on the current proxy's configuration
func (group *accountsGroup) getShard(c *gin.Context) {
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

	shardID, err := group.facade.GetShardIDForAddress(addr)
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

// getESDTTokenData returns the balance for the given address and esdt token
func (group *accountsGroup) getESDTTokenData(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTTokenData, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
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
			data.ReturnCodeRequestError,
		)
		return
	}

	esdtTokenResponse, err := group.facade.GetESDTTokenData(addr, tokenIdentifier)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			err.Error(),
			data.ReturnCodeInternalError,
		)
		return
	}

	c.JSON(http.StatusOK, esdtTokenResponse)
}

// getESDTsWithRole returns the token identifiers of the tokens where  the given address has the given role
func (group *accountsGroup) getESDTsWithRole(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTsWithRole, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
		)
		return
	}

	role := c.Param("role")
	if role == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTsWithRole, errors.ErrEmptyTokenIdentifier),
			data.ReturnCodeRequestError,
		)
		return
	}

	esdtsWithRole, err := group.facade.GetESDTsWithRole(addr, role)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			err.Error(),
			data.ReturnCodeInternalError,
		)
		return
	}

	c.JSON(http.StatusOK, esdtsWithRole)
}

// getOwnedNFTs returns the token identifiers of the NFTs where the given address is the owner
func (group *accountsGroup) getOwnedNFTs(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetOwnedNFTs, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
		)
		return
	}

	tokens, err := group.facade.GetOwnedNFTs(addr)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			err.Error(),
			data.ReturnCodeInternalError,
		)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// getESDTNftTokenData returns the esdt nft data for the given address, esdt token and nonce
func (group *accountsGroup) getESDTNftTokenData(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTTokenData, errors.ErrEmptyAddress),
			data.ReturnCodeRequestError,
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
			data.ReturnCodeRequestError,
		)
		return
	}

	nonce, err := shared.FetchNonceFromRequest(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrCannotParseNonce.Error(), data.ReturnCodeRequestError)
		return
	}

	esdtTokenResponse, err := group.facade.GetESDTNftTokenData(addr, tokenIdentifier, nonce)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			err.Error(),
			data.ReturnCodeInternalError,
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
			data.ReturnCodeRequestError,
		)
		return
	}

	tokens, err := group.facade.GetAllESDTTokens(addr)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			err.Error(),
			data.ReturnCodeInternalError,
		)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
