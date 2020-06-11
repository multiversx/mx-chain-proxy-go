package transaction

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

// Routes defines transaction related routes
func Routes(router *gin.RouterGroup) {
	router.POST("/send", SendTransaction)
	router.POST("/send-multiple", SendMultipleTransactions)
	router.POST("/send-user-funds", SendUserFunds)
	router.POST("/cost", RequestTransactionCost)
	router.GET("/:txhash/status", GetTransactionStatus)
	router.GET("/:txhash", GetTransaction)
}

// SendTransaction will receive a transaction from the client and propagate it for processing
func SendTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
				Code:  data.ReturnCodeRequestError,
			},
		)
		return
	}

	statusCode, txHash, err := ef.SendTransaction(&tx)
	if err != nil {
		c.JSON(
			statusCode,
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
			Data:  gin.H{"txHash": txHash},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// SendUserFunds will receive an address from the client and propagate a transaction for sending some ERD to that address
func SendUserFunds(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	var gtx = data.FundsRequest{}
	err := c.ShouldBindJSON(&gtx)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
				Code:  data.ReturnCodeRequestError,
			},
		)
		return
	}

	err = ef.SendUserFunds(gtx.Receiver, gtx.Value)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"message": "ok"},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// SendMultipleTransactions will send multiple transactions at once
func SendMultipleTransactions(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	var txs []*data.Transaction
	err := c.ShouldBindJSON(&txs)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
				Code:  data.ReturnCodeRequestError,
			},
		)
		return
	}

	response, err := ef.SendMultipleTransactions(txs)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data: gin.H{
				"numOfSentTxs": response.NumOfTxs,
				"txsHashes":    response.TxsHashes,
			},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// RequestTransactionCost will return an estimation of how many gas unit a transaction will cost
func RequestTransactionCost(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
				Code:  data.ReturnCodeRequestError,
			},
		)
		return
	}

	cost, err := ef.TransactionCostRequest(&tx)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
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
			Data:  gin.H{"txGasUnits": cost},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetTransactionStatus will return the transaction's status
func GetTransactionStatus(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	txHash := c.Param("txhash")
	txStatus, err := ef.GetTransactionStatus(txHash)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
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
			Data:  gin.H{"status": txStatus},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

// GetTransaction should return a transaction from observer
func GetTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrInvalidAppContext.Error(),
				Code:  data.ReturnCodeInternalError,
			},
		)
		return
	}

	txHash := c.Param("txhash")
	if txHash == "" {
		c.JSON(
			http.StatusBadRequest,
			data.GenericAPIResponse{
				Data:  nil,
				Error: errors.ErrTransactionHashMissing.Error(),
				Code:  data.ReturnCodeRequestError,
			},
		)
		return
	}

	sndAddr := c.Request.URL.Query().Get("sender")
	if sndAddr != "" {
		getTransactionByHashAndSenderAddress(c, ef, txHash, sndAddr)
		return
	}

	tx, err := ef.GetTransaction(txHash)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
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
			Data:  gin.H{"transaction": tx},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}

func getTransactionByHashAndSenderAddress(c *gin.Context, ef FacadeHandler, txHash string, sndAddr string) {
	tx, statusCode, err := ef.GetTransactionByHashAndSenderAddress(txHash, sndAddr)
	if err != nil {
		internalCode := data.ReturnCodeInternalError
		if statusCode == http.StatusBadRequest {
			internalCode = data.ReturnCodeRequestError
		}
		c.JSON(
			statusCode,
			data.GenericAPIResponse{
				Data:  nil,
				Error: "",
				Code:  internalCode,
			},
		)
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		data.GenericAPIResponse{
			Data:  gin.H{"transaction": tx},
			Error: "",
			Code:  data.ReturnCodeSuccess,
		},
	)
}
