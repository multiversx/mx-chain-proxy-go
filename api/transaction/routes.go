package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

const FaucetDefaultValue = 10000
const FaucetMaxValue     = 1000000

// Routes defines transaction related routes
func Routes(router *gin.RouterGroup) {
	router.POST("/send", SendTransaction)
	router.POST("/send-user-funds", SendUserFunds)
}

// SendTransaction will receive a transaction from the client and propagate it for processing
func SendTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error())})
		return
	}

	_, err = hex.DecodeString(tx.Sender)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrInvalidSenderAddress.Error(), err.Error())})
		return
	}

	_, err = hex.DecodeString(tx.Receiver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrInvalidReceiverAddress.Error(), err.Error())})
		return
	}

	_, err = hex.DecodeString(tx.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrInvalidSignatureHex.Error(), err.Error())})
		return
	}

	txHash, err := ef.SendTransaction(&tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"txHash": txHash})
}

// SendUserFunds will receive an address from the client and propagate a transaction for sending some ERD to that address
func SendUserFunds(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	var gtx = data.FundsRequest{}
	err := c.ShouldBindJSON(&gtx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error())})
		return
	}

	err = ef.SendUserFunds(gtx.Receiver, validateAndSetFaucetValue(gtx.Value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func validateAndSetFaucetValue(providedVal *big.Int) *big.Int {
	faucetDefault := big.NewInt(0).SetUint64(uint64(FaucetDefaultValue))
	faucetMax     := big.NewInt(0).SetUint64(uint64(FaucetMaxValue))

	if providedVal == nil {
		return faucetDefault
	}

	if faucetMax.Cmp(providedVal) == -1 {
		return faucetDefault
	}

	return providedVal
}
