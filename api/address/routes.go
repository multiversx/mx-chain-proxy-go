package address

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/:address", GetAccount)
	router.GET("/:address/balance", GetBalance)
	router.GET("/:address/nonce", GetNonce)
	router.GET("/:address/transactions", GetTransactions)
	router.GET("/:address/storage/:key", GetValueForKey)
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
		c.JSON(
			status,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"account": account},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetBalance returns the balance for the address parameter
func GetBalance(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		c.JSON(
			status,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"balance": account.Balance},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetNonce returns the nonce for the address parameter
func GetNonce(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		c.JSON(
			status,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"nonce": account.Nonce},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetTransactions returns the transactions for the address parameter
func GetTransactions(c *gin.Context) {
	transactions, status, err := getTransactions(c)
	if err != nil {
		c.JSON(
			status,
			data.GenericAPIResponse{
				Data:  nil,
				Error: err.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"transactions": transactions},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetValueForKey returns the value for the given address and key
func GetValueForKey(c *gin.Context) {
	ef, ok := c.MustGet("elrondFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	addr := c.Param("address")
	if addr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v: %v", errors.ErrGetValueForKey, errors.ErrEmptyAddress)})
		return
	}

	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v: %v", errors.ErrGetValueForKey, errors.ErrEmptyKey)})
		return
	}

	value, err := ef.GetValueForKey(addr, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrGetValueForKey.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": value})
}
