package groups

import (
	"fmt"
	"net/http"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

/*
router.POST("/send", SendTransaction)
	router.POST("/send-multiple", SendMultipleTransactions)
	router.POST("/simulate", SimulateTransaction)
	router.POST("/send-user-funds", SendUserFunds)
	router.POST("/cost", RequestTransactionCost)
	router.GET("/:txhash/status", GetTransactionStatus)
	router.GET("/:txhash", GetTransaction)
*/

func NewBaseTransactionsGroup() *baseGroup {
	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/send":            {Handler: SendTransaction, Method: http.MethodPost},
		"/simulate":        {Handler: SimulateTransaction, Method: http.MethodPost},
		"/send-multiple":   {Handler: SendMultipleTransactions, Method: http.MethodPost},
		"/send-user-funds": {Handler: SendUserFunds, Method: http.MethodPost},
		"/cost":            {Handler: RequestTransactionCost, Method: http.MethodPost},
		"/:txhash/status":  {Handler: GetTransactionStatus, Method: http.MethodGet},
		"/:txhash":         {Handler: GetTransaction, Method: http.MethodGet},
	}

	return &baseGroup{
		endpoints: baseRoutesHandlers,
	}
}

// SendTransaction will receive a transaction from the client and propagate it for processing
func SendTransaction(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	statusCode, txHash, err := ef.SendTransaction(&tx)
	if err != nil {
		shared.RespondWith(c, statusCode, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txHash": txHash}, "", data.ReturnCodeSuccess)
}

// SendUserFunds will receive an address from the client and propagate a transaction for sending some ERD to that address
func SendUserFunds(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	if !ef.IsFaucetEnabled() {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			errors.ErrFaucetNotEnabled.Error(),
			data.ReturnCodeRequestError,
		)
		return
	}

	var gtx = data.FundsRequest{}
	err := c.ShouldBindJSON(&gtx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	err = ef.SendUserFunds(gtx.Receiver, gtx.Value)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"message": "ok"}, "", data.ReturnCodeSuccess)
}

// SendMultipleTransactions will send multiple transactions at once
func SendMultipleTransactions(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var txs []*data.Transaction
	err := c.ShouldBindJSON(&txs)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	response, err := ef.SendMultipleTransactions(txs)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusInternalServerError,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
			data.ReturnCodeInternalError,
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
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeRequestError,
		)
		return
	}

	simulationResponse, err := ef.SimulateTransaction(&tx)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	c.JSON(
		http.StatusOK,
		simulationResponse,
	)
}

// RequestTransactionCost will return an estimation of how many gas unit a transaction will cost
func RequestTransactionCost(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
			data.ReturnCodeInternalError,
		)
		return
	}

	cost, err := ef.TransactionCostRequest(&tx)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txGasUnits": cost}, "", data.ReturnCodeSuccess)
}

// GetTransactionStatus will return the transaction's status
func GetTransactionStatus(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	txHash := c.Param("txhash")
	sender := c.Request.URL.Query().Get("sender")
	txStatus, err := ef.GetTransactionStatus(txHash, sender)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"status": txStatus}, "", data.ReturnCodeSuccess)
}

// GetTransaction should return a transaction from observer
func GetTransaction(c *gin.Context) {
	ef, ok := c.MustGet(shared.GetFacadeVersion(c)).(TransactionFacadeHandler)
	if !ok {
		shared.RespondWithInvalidAppContext(c)
		return
	}

	txHash := c.Param("txhash")
	if txHash == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrTransactionHashMissing.Error(), data.ReturnCodeRequestError)
		return
	}

	sndAddr := c.Request.URL.Query().Get("sender")
	if sndAddr != "" {
		getTransactionByHashAndSenderAddress(c, ef, txHash, sndAddr)
		return
	}

	tx, err := ef.GetTransaction(txHash)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "", data.ReturnCodeSuccess)
}

func getTransactionByHashAndSenderAddress(c *gin.Context, ef TransactionFacadeHandler, txHash string, sndAddr string) {
	tx, statusCode, err := ef.GetTransactionByHashAndSenderAddress(txHash, sndAddr)
	if err != nil {
		internalCode := data.ReturnCodeInternalError
		if statusCode == http.StatusBadRequest {
			internalCode = data.ReturnCodeRequestError
		}
		shared.RespondWith(c, statusCode, nil, err.Error(), internalCode)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "", data.ReturnCodeSuccess)
}
