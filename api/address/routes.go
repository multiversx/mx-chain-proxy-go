package address

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/:address", GetAccount)
	router.GET("/:address/balance", GetBalance)
	router.GET("/:address/username", GetUsername)
	router.GET("/:address/nonce", GetNonce)
	router.GET("/:address/shard", GetShard)
	router.GET("/:address/transactions", GetTransactions)
	router.GET("/:address/key/:key", GetValueForKey)
	router.GET("/:address/esdt", GetESDTTokens)
	router.GET("/:address/esdt/:tokenName", GetESDTTokenData)
}

func getAccount(c *gin.Context) (*data.Account, int, error) {
	epf, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		return nil, http.StatusInternalServerError, errors.ErrInvalidAppContext
	}

	addr := c.Param("address")
	acc, err := epf.GetAccount(addr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return acc, http.StatusOK, nil
}

func getTransactions(c *gin.Context) ([]data.DatabaseTransaction, int, error) {
	epf, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		return nil, http.StatusInternalServerError, errors.ErrInvalidAppContext
	}

	addr := c.Param("address")
	transactions, err := epf.GetTransactions(addr)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return transactions, http.StatusOK, nil
}

// GetAccount returns an accountResponse containing information
// about the account correlated with provided address
func GetAccount(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"account": account}, "", data.ReturnCodeSuccess)
}

// GetBalance returns the balance for the address parameter
func GetBalance(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"balance": account.Balance}, "", data.ReturnCodeSuccess)
}

// GetUsername returns the username for the address parameter
func GetUsername(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"username": account.Username}, "", data.ReturnCodeSuccess)
}

// GetNonce returns the nonce for the address parameter
func GetNonce(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"nonce": account.Nonce}, "", data.ReturnCodeSuccess)
}

// GetTransactions returns the transactions for the address parameter
func GetTransactions(c *gin.Context) {
	transactions, status, err := getTransactions(c)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transactions": transactions}, "", data.ReturnCodeSuccess)
}

// GetValueForKey returns the value for the given address and key
func GetValueForKey(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

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

	value, err := ef.GetValueForKey(addr, key)
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
func GetShard(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

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

	shardID, err := ef.GetShardIDForAddress(addr)
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

// GetESDTTokenData returns the balance for the given address and esdt token
func GetESDTTokenData(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

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

	tokenName := c.Param("tokenName")
	if tokenName == "" {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%v: %v", errors.ErrGetESDTTokenData, errors.ErrEmptyTokenName),
			data.ReturnCodeRequestError,
		)
		return
	}

	esdtTokenResponse, err := ef.GetESDTTokenData(addr, tokenName)
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

// GetESDTTokens returns the tokens list from this account
func GetESDTTokens(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

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

	tokens, err := ef.GetAllESDTTokens(addr)
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
