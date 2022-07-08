package groups

import (
	"strconv"

	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/gin-gonic/gin"
)

func parseBlockQueryOptions(c *gin.Context) (common.BlockQueryOptions, error) {
	withTxs, err := parseBoolUrlParam(c, common.UrlParameterWithTransactions)
	if err != nil {
		return common.BlockQueryOptions{}, err
	}

	withLogs, err := parseBoolUrlParam(c, common.UrlParameterWithLogs)
	if err != nil {
		return common.BlockQueryOptions{}, err
	}

	options := common.BlockQueryOptions{WithTransactions: withTxs, WithLogs: withLogs}
	return options, nil
}

func parseHyperblockQueryOptions(c *gin.Context) (common.HyperblockQueryOptions, error) {
	withLogs, err := parseBoolUrlParam(c, common.UrlParameterWithLogs)
	if err != nil {
		return common.HyperblockQueryOptions{}, err
	}

	options := common.HyperblockQueryOptions{WithLogs: withLogs}
	return options, nil
}

func parseAccountQueryOptions(c *gin.Context) (common.AccountQueryOptions, error) {
	onFinalBlock, err := parseBoolUrlParam(c, common.UrlParameterOnFinalBlock)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	onStartOfEpoch, err := parseUintUrlParam(c, common.UrlParameterOnStartOfEpoch)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	options := common.AccountQueryOptions{OnFinalBlock: onFinalBlock, OnStartOfEpoch: onStartOfEpoch}
	return options, nil
}

func parseTransactionQueryOptions(c *gin.Context) (common.TransactionQueryOptions, error) {
	withResults, err := parseBoolUrlParam(c, common.UrlParameterWithResults)
	if err != nil {
		return common.TransactionQueryOptions{}, err
	}

	options := common.TransactionQueryOptions{WithResults: withResults}
	return options, nil
}

func parseTransactionSimulationOptions(c *gin.Context) (common.TransactionSimulationOptions, error) {
	checkSignature, err := parseBoolUrlParamWithDefault(c, common.UrlParameterCheckSignature, true)
	if err != nil {
		return common.TransactionSimulationOptions{}, err
	}

	options := common.TransactionSimulationOptions{CheckSignature: checkSignature}
	return options, nil
}

func parseBoolUrlParam(c *gin.Context, name string) (bool, error) {
	return parseBoolUrlParamWithDefault(c, name, false)
}

func parseBoolUrlParamWithDefault(c *gin.Context, name string, defaultValue bool) (bool, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return defaultValue, nil
	}

	return strconv.ParseBool(param)
}

func parseStringUrlParam(c *gin.Context, name string) string {
	return c.Request.URL.Query().Get(name)
}

func parseUintUrlParam(c *gin.Context, name string) (uint32, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return 0, nil
	}

	value, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(value), nil
}

func parseTransactionsPoolQueryOptions(c *gin.Context) (common.TransactionsPoolOptions, error) {
	shardId, err := parseUintUrlParam(c, common.UrlParameterShardID)
	if err != nil {
		return common.TransactionsPoolOptions{}, err
	}

	lastNonce, err := parseBoolUrlParam(c, common.UrlParameterLastNonce)
	if err != nil {
		return common.TransactionsPoolOptions{}, err
	}

	nonceGaps, err := parseBoolUrlParam(c, common.UrlParameterNonceGaps)
	if err != nil {
		return common.TransactionsPoolOptions{}, err
	}

	return common.TransactionsPoolOptions{
		ShardID:   shardId,
		Sender:    parseStringUrlParam(c, common.UrlParameterSender),
		Fields:    parseStringUrlParam(c, common.UrlParameterFields),
		LastNonce: lastNonce,
		NonceGaps: nonceGaps,
	}, nil
}
