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

	notarizedAtSource, err := parseBoolUrlParam(c, common.UrlParameterNotarizedAtSource)
	if err != nil {
		return common.HyperblockQueryOptions{}, err
	}

	withAlteredAccounts, err := parseBoolUrlParam(c, common.UrlParameterWithAlteredAccounts)
	if err != nil {
		return common.HyperblockQueryOptions{}, err
	}

	var alteredAccountsOptions common.GetAlteredAccountsForBlockOptions
	if withAlteredAccounts {
		alteredAccountsOptions, err = parseAlteredAccountOptions(c)
		if err != nil {
			return common.HyperblockQueryOptions{}, err
		}
	}

	return common.HyperblockQueryOptions{
		WithLogs:               withLogs,
		NotarizedAtSource:      notarizedAtSource,
		WithAlteredAccounts:    withAlteredAccounts,
		AlteredAccountsOptions: alteredAccountsOptions,
	}, nil
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
	lastNonce, err := parseBoolUrlParam(c, common.UrlParameterLastNonce)
	if err != nil {
		return common.TransactionsPoolOptions{}, err
	}

	nonceGaps, err := parseBoolUrlParam(c, common.UrlParameterNonceGaps)
	if err != nil {
		return common.TransactionsPoolOptions{}, err
	}

	return common.TransactionsPoolOptions{
		ShardID:   parseStringUrlParam(c, common.UrlParameterShardID),
		Sender:    parseStringUrlParam(c, common.UrlParameterSender),
		Fields:    parseStringUrlParam(c, common.UrlParameterFields),
		LastNonce: lastNonce,
		NonceGaps: nonceGaps,
	}, nil
}

func parseAlteredAccountOptions(c *gin.Context) (common.GetAlteredAccountsForBlockOptions, error) {
	tokensFilter := parseStringUrlParam(c, common.UrlParameterTokensFilter)
	withMetaData, err := parseBoolUrlParam(c, common.UrlParameterWithMetadata)
	if err != nil {
		return common.GetAlteredAccountsForBlockOptions{}, err
	}
	if withMetaData && len(tokensFilter) == 0 {
		return common.GetAlteredAccountsForBlockOptions{}, ErrIncompatibleWithMetadataParam
	}

	return common.GetAlteredAccountsForBlockOptions{
		TokensFilter: tokensFilter,
		WithMetadata: withMetaData,
	}, nil
}
