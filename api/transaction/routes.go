package transaction

import (
	"fmt"
	"net/http"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

var log = logger.GetOrCreate("api/transaction")

// Routes defines transaction related routes
func Routes(router *gin.RouterGroup) {
	router.POST("/send", SendTransaction)
	router.POST("/send-multiple", SendMultipleTransactions)
	router.POST("/simulate", SimulateTransaction)
	router.POST("/send-user-funds", SendUserFunds)
	router.POST("/cost", RequestTransactionCost)
	router.GET("/:txhash/status", GetTransactionStatus)
	router.GET("/:txhash", GetTransaction)
}

func respondWithError(c *gin.Context, status int, dataField interface{}, error string, code data.ReturnCode, payload interface{}, apiRoute string) {
	log.Error(apiRoute,
		"status", status,
		"dataField", spew.Sdump(dataField),
		"error", error,
		"code", code,
		"payload", spew.Sdump(payload),
	)
	shared.RespondWith(c, status, nil, error, code)
}

// SendTransaction will receive a transaction from the client and propagate it for processing
func SendTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		respondWithError(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
			nil,
			"SendTransaction",
		)
		return
	}

	statusCode, txHash, err := ef.SendTransaction(&tx)
	if err != nil {
		respondWithError(c, statusCode, nil, err.Error(), data.ReturnCodeInternalError, &tx, "SendTransaction")
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txHash": txHash}, "", data.ReturnCodeSuccess)
}

// SendUserFunds will receive an address from the client and propagate a transaction for sending some ERD to that address
func SendUserFunds(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	if !ef.IsFaucetEnabled() {
		respondWithError(
			c,
			http.StatusBadRequest,
			nil,
			errors.ErrFaucetNotEnabled.Error(),
			data.ReturnCodeRequestError,
			nil,
			"SendUserFunds",
		)
		return
	}

	var gtx = data.FundsRequest{}
	err := c.ShouldBindJSON(&gtx)
	if err != nil {
		respondWithError(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
			nil,
			"SendUserFunds",
		)
		return
	}

	err = ef.SendUserFunds(gtx.Receiver, gtx.Value)
	if err != nil {
		respondWithError(
			c,
			http.StatusInternalServerError,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
			data.ReturnCodeRequestError,
			&gtx,
			"SendUserFunds",
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"message": "ok"}, "", data.ReturnCodeSuccess)
}

// SendMultipleTransactions will send multiple transactions at once
func SendMultipleTransactions(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var txs []*data.Transaction
	err := c.ShouldBindJSON(&txs)
	if err != nil {
		respondWithError(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
			nil,
			"SendMultipleTransactions",
		)
		return
	}

	response, err := ef.SendMultipleTransactions(txs)
	if err != nil {
		respondWithError(
			c,
			http.StatusInternalServerError,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
			data.ReturnCodeInternalError,
			&txs,
			"SendMultipleTransactions",
		)
		return
	}

	shared.RespondWith(
		c,
		http.StatusOK,
		gin.H{
			"numOfSentTxs": response.NumOfTxs,
			"txsHashes":    response.TxsHashes,
		},
		"",
		data.ReturnCodeSuccess,
	)
}

// SimulateTransaction will receive a transaction from the client and will send it for simulation purpose
func SimulateTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		respondWithError(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
			nil,
			"SimulateTransaction",
		)
		return
	}

	simulationResponse, err := ef.SimulateTransaction(&tx)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError,
			&tx, "SimulateTransaction")
		return
	}

	c.JSON(
		http.StatusOK,
		simulationResponse,
	)
}

// RequestTransactionCost will return an estimation of how many gas unit a transaction will cost
func RequestTransactionCost(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		respondWithError(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeInternalError,
			nil,
			"RequestTransactionCost",
		)
		return
	}

	cost, err := ef.TransactionCostRequest(&tx)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError,
			&tx, "RequestTransactionCost")
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txGasUnits": cost}, "", data.ReturnCodeSuccess)
}

// GetTransactionStatus will return the transaction's status
func GetTransactionStatus(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	txHash := c.Param("txhash")
	sender := c.Request.URL.Query().Get("sender")
	txStatus, err := ef.GetTransactionStatus(txHash, sender)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError,
			txHash, "GetTransactionStatus")
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"status": txStatus}, "", data.ReturnCodeSuccess)
}

// GetTransaction should return a transaction from observer
func GetTransaction(c *gin.Context) {
	ef, ok := c.MustGet("elrondProxyFacade").(FacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	txHash := c.Param("txhash")
	if txHash == "" {
		respondWithError(c, http.StatusBadRequest, nil, errors.ErrTransactionHashMissing.Error(), data.ReturnCodeRequestError, nil, "GetTransaction")
		return
	}

	sndAddr := c.Request.URL.Query().Get("sender")
	if sndAddr != "" {
		getTransactionByHashAndSenderAddress(c, ef, txHash, sndAddr)
		return
	}

	tx, err := ef.GetTransaction(txHash)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError, txHash, "GetTransaction")
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "", data.ReturnCodeSuccess)
}

func getTransactionByHashAndSenderAddress(c *gin.Context, ef FacadeHandler, txHash string, sndAddr string) {
	tx, statusCode, err := ef.GetTransactionByHashAndSenderAddress(txHash, sndAddr)
	if err != nil {
		internalCode := data.ReturnCodeInternalError
		if statusCode == http.StatusBadRequest {
			internalCode = data.ReturnCodeRequestError
		}
		respondWithError(c, statusCode, nil, err.Error(), internalCode, txHash, "getTransactionByHashAndSenderAddress")
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "", data.ReturnCodeSuccess)
}
