package address

import (
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	GetAccount(address string) (*data.Account, error)
}

// Routes defines address related routes
func Routes(router *gin.RouterGroup) {
	router.GET("/:address", GetAccount)
	router.GET("/:address/balance", GetBalance)
	router.GET("/:address/nonce", GetNonce)
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

// GetAccount returns an accountResponse containing information
//  about the account correlated with provided address
func GetAccount(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// GetBalance returns the balance for the address parameter
func GetBalance(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

// GetNonce returns the nonce for the address parameter
func GetNonce(c *gin.Context) {
	account, status, err := getAccount(c)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nonce": account.Nonce})
}
