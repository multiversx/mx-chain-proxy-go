package groups

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/api/shared"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type transactionGroup struct {
	facade TransactionFacadeHandler
	*baseGroup
}

// NewTransactionGroup returns a new instance of transactionGroup
func NewTransactionGroup(facadeHandler data.FacadeHandler) (*transactionGroup, error) {
	facade, ok := facadeHandler.(TransactionFacadeHandler)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	tg := &transactionGroup{
		facade:    facade,
		baseGroup: &baseGroup{},
	}

	baseRoutesHandlers := []*data.EndpointHandlerData{
		{Path: "/send", Handler: tg.sendTransaction, Method: http.MethodPost},
		{Path: "/simulate", Handler: tg.simulateTransaction, Method: http.MethodPost},
		{Path: "/send-multiple", Handler: tg.sendMultipleTransactions, Method: http.MethodPost},
		{Path: "/send-user-funds", Handler: tg.sendUserFunds, Method: http.MethodPost},
		{Path: "/cost", Handler: tg.requestTransactionCost, Method: http.MethodPost},
		{Path: "/:txhash/status", Handler: tg.getTransactionStatus, Method: http.MethodGet},
		{Path: "/:txhash/process-status", Handler: tg.getProcessedTransactionStatus, Method: http.MethodGet},
		{Path: "/:txhash", Handler: tg.getTransaction, Method: http.MethodGet},
		{Path: "/pool", Handler: tg.getTransactionsPool, Method: http.MethodGet},
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
			data.ReturnCodeRequestError,
		)
		return
	}

	statusCode, txHash, err := group.facade.SendTransaction(&tx)
	if err != nil {
		shared.RespondWith(c, statusCode, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txHash": txHash}, "", data.ReturnCodeSuccess)
}

// sendUserFunds will receive an address from the client and propagate a transaction for sending some ERD to that address
func (group *transactionGroup) sendUserFunds(c *gin.Context) {
	if !group.facade.IsFaucetEnabled() {
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

	err = group.facade.SendUserFunds(gtx.Receiver, gtx.Value)
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
			data.ReturnCodeRequestError,
		)
		return
	}

	response, err := group.facade.SendMultipleTransactions(txs)
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
			data.ReturnCodeRequestError,
		)
		return
	}

	options, err := parseTransactionSimulationOptions(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrValidatorQueryParameterCheckSignature.Error(), data.ReturnCodeRequestError)
		return
	}

	simulationResponse, err := group.facade.SimulateTransaction(&tx, options.CheckSignature)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
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
			data.ReturnCodeInternalError,
		)
		return
	}

	cost, err := group.facade.TransactionCostRequest(&tx)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, cost, "", data.ReturnCodeSuccess)
}

// getTransactionStatus will return the transaction's status
func (group *transactionGroup) getTransactionStatus(c *gin.Context) {
	txHash := c.Param("txhash")
	sender := c.Request.URL.Query().Get("sender")
	txStatus, err := group.facade.GetTransactionStatus(txHash, sender)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"status": txStatus}, "", data.ReturnCodeSuccess)
}

// getTransaction should return a transaction from observer
func (group *transactionGroup) getTransaction(c *gin.Context) {
	txHash := c.Param("txhash")
	if txHash == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrTransactionHashMissing.Error(), data.ReturnCodeRequestError)
		return
	}

	options, err := parseTransactionQueryOptions(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrValidationQueryParameterWithResult.Error(), data.ReturnCodeRequestError)
		return
	}

	sndAddr := c.Request.URL.Query().Get("sender")
	if sndAddr != "" {
		getTransactionByHashAndSenderAddress(c, group.facade, txHash, sndAddr, options.WithResults)
		return
	}

	tx, err := group.facade.GetTransaction(txHash, options.WithResults, options.RelayedTxHash)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"transaction": tx}, "", data.ReturnCodeSuccess)
}

func (group *transactionGroup) getProcessedTransactionStatus(c *gin.Context) {
	txHash := c.Param("txhash")
	if txHash == "" {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrTransactionHashMissing.Error(), data.ReturnCodeRequestError)
		return
	}

	status, err := group.facade.GetProcessedTransactionStatus(txHash)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"status": status.Status, "reason": status.Reason}, "", data.ReturnCodeSuccess)
}

func getTransactionByHashAndSenderAddress(c *gin.Context, ef TransactionFacadeHandler, txHash string, sndAddr string, withEvents bool) {
	tx, statusCode, err := ef.GetTransactionByHashAndSenderAddress(txHash, sndAddr, withEvents)
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

// getTransactionsPool should return transactions from pool
func (group *transactionGroup) getTransactionsPool(c *gin.Context) {
	options, err := parseTransactionsPoolQueryOptions(c)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrBadUrlParams.Error(), data.ReturnCodeRequestError)
		return
	}

	err = validateOptions(options)
	if err != nil {
		shared.RespondWith(c, http.StatusBadRequest, nil, err.Error(), data.ReturnCodeRequestError)
		return
	}

	if options.Sender == "" {
		if options.ShardID == "" {
			getTxPool(c, group.facade, options.Fields)
			return
		}

		shardID, err := strconv.ParseUint(options.ShardID, 10, 32)
		if err != nil {
			shared.RespondWith(c, http.StatusBadRequest, nil, errors.ErrBadUrlParams.Error(), data.ReturnCodeRequestError)
			return
		}
		getTxPoolForShard(c, group.facade, uint32(shardID), options.Fields)
		return
	}

	if options.LastNonce {
		getLastTxPoolNonceForSender(c, group.facade, options.Sender)
		return
	}

	if options.NonceGaps {
		getTxPoolNonceGapsForSender(c, group.facade, options.Sender)
		return
	}

	getTxPoolForSender(c, group.facade, options.Sender, options.Fields)
}

func validateOptions(options common.TransactionsPoolOptions) error {
	if options.Fields != "" && options.LastNonce {
		return errors.ErrFetchingLatestNonceCannotIncludeFields
	}

	if options.Fields != "" && options.NonceGaps {
		return errors.ErrFetchingNonceGapsCannotIncludeFields
	}

	if options.Sender == "" && options.LastNonce {
		return errors.ErrEmptySenderToGetLatestNonce
	}

	if options.Sender == "" && options.NonceGaps {
		return errors.ErrEmptySenderToGetNonceGaps
	}

	if options.Fields == "*" {
		return nil
	}

	if options.Fields != "" {
		return validateFields(options.Fields)
	}

	return nil
}

func validateFields(fields string) error {
	for _, c := range fields {
		if c == ',' {
			continue
		}

		isLowerLetter := c >= 'a' && c <= 'z'
		isUpperLetter := c >= 'A' && c <= 'Z'
		if !isLowerLetter && !isUpperLetter {
			return errors.ErrInvalidFields
		}
	}

	return nil
}

func getTxPool(c *gin.Context, ef TransactionFacadeHandler, fields string) {
	txPool, err := ef.GetTransactionsPool(fields)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txPool": txPool}, "", data.ReturnCodeSuccess)
}

func getTxPoolForShard(c *gin.Context, ef TransactionFacadeHandler, shardID uint32, fields string) {
	txPool, err := ef.GetTransactionsPoolForShard(shardID, fields)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txPool": txPool}, "", data.ReturnCodeSuccess)
}

func getLastTxPoolNonceForSender(c *gin.Context, ef TransactionFacadeHandler, sender string) {
	lastNonce, err := ef.GetLastPoolNonceForSender(sender)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"nonce": lastNonce}, "", data.ReturnCodeSuccess)
}

func getTxPoolNonceGapsForSender(c *gin.Context, ef TransactionFacadeHandler, sender string) {
	nonceGaps, err := ef.GetTransactionsPoolNonceGapsForSender(sender)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"nonceGaps": nonceGaps}, "", data.ReturnCodeSuccess)
}

func getTxPoolForSender(c *gin.Context, ef TransactionFacadeHandler, sender, fields string) {
	txPool, err := ef.GetTransactionsPoolForSender(sender, fields)
	if err != nil {
		shared.RespondWith(c, http.StatusInternalServerError, nil, err.Error(), data.ReturnCodeInternalError)
		return
	}

	shared.RespondWith(c, http.StatusOK, gin.H{"txPool": txPool}, "", data.ReturnCodeSuccess)
}
