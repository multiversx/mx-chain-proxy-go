package wallet

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines wallet related routes
func Routes(router *gin.RouterGroup) {
	router.POST("/publickey", PublicKeyFromPrivateKey)
	router.POST("/send", SignAndSendTransaction)
}

type publicKeyRequest struct {
	PrivateKey string `json:"privateKey"`
}

// PublicKeyFromPrivateKey will return a public key from a given private key
func PublicKeyFromPrivateKey(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	var pubKeyReq publicKeyRequest
	err := c.ShouldBindJSON(&pubKeyReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error())})
		return
	}

	pubKeyStr, err := ef.PublicKeyFromPrivateKey(pubKeyReq.PrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrPublicKeyFromPrivKey.Error(), err.Error())})
	}

	c.JSON(http.StatusOK, gin.H{"publicKey": pubKeyStr})
}

// SignAndSendTransaction will receive a transaction request and a private key and will create and send a transaction
func SignAndSendTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInvalidAppContext.Error()})
		return
	}

	var txReq = data.TransactionSignRequest{}
	err := c.ShouldBindJSON(&txReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error())})
		return
	}

	tx, skBytes := getTxAndPrivKeyFromTxRequest(txReq)
	txHash, err := ef.SignAndSendTransaction(tx, skBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"txHash": txHash})
}

func getTxAndPrivKeyFromTxRequest(txReq data.TransactionSignRequest) (*data.Transaction, []byte) {
	tx := data.Transaction{
		Nonce:    txReq.Nonce,
		Value:    txReq.Value,
		Receiver: txReq.Receiver,
		Sender:   txReq.Sender,
		GasPrice: txReq.GasPrice,
		GasLimit: txReq.GasLimit,
		Data:     txReq.Data,
	}

	privKeyBytes, err := hex.DecodeString(txReq.PrivateKey)
	if err != nil {
		fmt.Println("error getting bytes from priv key: ", err.Error())
	}

	return &tx, privKeyBytes
}
