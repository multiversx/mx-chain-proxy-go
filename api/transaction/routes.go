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

// FacadeHandler interface defines methods that can be used from `elrondProxyFacade` context variable
type FacadeHandler interface {
	SendTransaction(nonce uint64, sender string, receiver string, value *big.Int, code string, signature []byte) (*data.Transaction, error)
}

// Routes defines transaction related routes
func Routes(router *gin.RouterGroup) {
	router.POST("/send", SendTransaction)
}

// SendTransaction will receive a transaction from the client and propagate it for processing
func SendTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	var gtx = data.Transaction{}
	err := c.ShouldBindJSON(&gtx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error())})
		return
	}

	signature, err := hex.DecodeString(gtx.Signature)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrInvalidSignatureHex.Error(), err.Error())})
		return
	}

	tx, err := ef.SendTransaction(gtx.Nonce, gtx.Sender, gtx.Receiver, gtx.Value, gtx.Data, signature)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": tx})
}
