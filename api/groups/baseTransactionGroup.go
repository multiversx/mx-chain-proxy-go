package groups

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/api/errors"
	"github.com/ElrondNetwork/elrond-proxy-go/api/shared"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/gin-gonic/gin"
)

type transactionGroup struct {
	facade TransactionFacadeHandler
	*baseGroup
}

// NewNodeGroup returns a new instance of nodeGroup
func NewTransactionGroup(facadeHandler data.FacadeHandler) (*transactionGroup, error) {
	facade, ok := facadeHandler.(TransactionFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	tg := &transactionGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := map[string]*data.EndpointHandlerData{
		"/send":            {Handler: tg.sendTransaction, Method: http.MethodPost},
		"/simulate":        {Handler: tg.simulateTransaction, Method: http.MethodPost},
		"/send-multiple":   {Handler: tg.sendMultipleTransactions, Method: http.MethodPost},
		"/send-user-funds": {Handler: tg.sendUserFunds, Method: http.MethodPost},
		"/cost":            {Handler: tg.requestTransactionCost, Method: http.MethodPost},
		"/:txhash/status":  {Handler: tg.getTransactionStatus, Method: http.MethodGet},
		"/:txhash":         {Handler: tg.getTransaction, Method: http.MethodGet},
	}
	tg.baseGroup.endpoints = baseRoutesHandlers

	return tg, nil
}

// sendTransaction will receive a transaction from the client and propagate it for processing
func (group *transactionGroup) sendTransaction(c *gin.Context) {
	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
		)
		return
	}

	txHash, statusCode, err := group.facade.SendTransaction(&tx)
	if err != nil {
		shared.RespondWith(c, statusCode, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txHash": txHash}, "")
}

// sendUserFunds will receive an address from the client and propagate a transaction for sending some ERD to that address
func (group *transactionGroup) sendUserFunds(c *gin.Context) {
	if !group.facade.IsFaucetEnabled() {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			errors.ErrFaucetNotEnabled.Error(),
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
		)
		return
	}

	status, err := group.facade.SendUserFunds(gtx.Receiver, gtx.Value)
	if err != nil {
		shared.RespondWith(
			c,
			status,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
		)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"message": "ok"}, "")
}

// sendMultipleTransactions will send multiple transactions at once
func (group *transactionGroup) sendMultipleTransactions(c *gin.Context) {
	var txs []*data.Transaction
	err := c.ShouldBindJSON(&txs)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
		)
		return
	}

	response, status, err := group.facade.SendMultipleTransactions(txs)
	if err != nil {
		shared.RespondWith(
			c,
			status,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrTxGenerationFailed.Error(), err.Error()),
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
	)
}

// simulateTransaction will receive a transaction from the client and will send it for simulation purpose
func (group *transactionGroup) simulateTransaction(c *gin.Context) {
	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
		)
		return
	}

	simulationResponse, status, err := group.facade.SimulateTransaction(&tx)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	c.JSON(
		http.StatusOK,
		simulationResponse,
	)
}

// requestTransactionCost will return an estimation of how many gas unit a transaction will cost
func (group *transactionGroup) requestTransactionCost(c *gin.Context) {
	var tx = data.Transaction{}
	err := c.ShouldBindJSON(&tx)
	if err != nil {
		shared.RespondWith(
			c,
			http.StatusBadRequest,
			nil,
			fmt.Sprintf("%s: %s", errors.ErrValidation.Error(), err.Error()),
		)
		return
	}

	cost, status, err := group.facade.TransactionCostRequest(&tx)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txGasUnits": cost}, "")
}

// getTransactionStatus will return the transaction's status
func (group *transactionGroup) getTransactionStatus(c *gin.Context) {
	txHash := c.Param("txhash")
	sender := c.Request.URL.Query().Get("sender")
	txStatus, status, err := group.facade.GetTransactionStatus(txHash, sender)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"status": txStatus}, "")
}

// getTransaction should return a transaction from observer
func (group *transactionGroup) getTransaction(c *gin.Context) {
	txHash := c.Param("txhash")
	if txHash == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrTransactionHashMissing.Error())
		return
	}

	withResults, err := getQueryParamWithResults(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrValidationQueryParameterWithResult.Error())
		return
	}

	sndAddr := c.Request.URL.Query().Get("sender")
	if sndAddr != "" {
		getTransactionByHashAndSenderAddress(c, group.facade, txHash, sndAddr, withResults)
		return
	}

	tx, status, err := group.facade.GetTransaction(txHash, withResults)
	if err != nil {
		shared.RespondWith(c, status, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "")
}

func getTransactionByHashAndSenderAddress(c *gin.Context, ef TransactionFacadeHandler, txHash string, sndAddr string, withEvents bool) {
	tx, statusCode, err := ef.GetTransactionByHashAndSenderAddress(txHash, sndAddr, withEvents)
	if err != nil {
		shared.RespondWith(c, statusCode, nil, err.Error())
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "")
}

func getQueryParamWithResults(c *gin.Context) (bool, error) {
	withResultsStr := c.Request.URL.Query().Get("withResults")
	if withResultsStr == "" {
		return false, nil
	}

	return strconv.ParseBool(withResultsStr)
}
