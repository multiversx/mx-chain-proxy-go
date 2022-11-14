package groups

import (
	"encoding/hex"
	"strconv"

	"github.com/ElrondNetwork/elrond-go-core/core"
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

	options := common.HyperblockQueryOptions{WithLogs: withLogs, NotarizedAtSource: notarizedAtSource}
	return options, nil
}

func parseAccountQueryOptions(c *gin.Context) (common.AccountQueryOptions, error) {
	onFinalBlock, err := parseBoolUrlParam(c, common.UrlParameterOnFinalBlock)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	onStartOfEpoch, err := parseUint32UrlParam(c, common.UrlParameterOnStartOfEpoch)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	blockNonce, err := parseUint64UrlParam(c, common.UrlParameterBlockNonce)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	blockHash, err := parseHexBytesUrlParam(c, common.UrlParameterBlockHash)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	blockRootHash, err := parseHexBytesUrlParam(c, common.UrlParameterBlockRootHash)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	hintEpoch, err := parseUint32UrlParam(c, common.UrlParameterOnStartOfEpoch)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	options := common.AccountQueryOptions{
		OnFinalBlock:   onFinalBlock,
		OnStartOfEpoch: onStartOfEpoch,
		BlockNonce:     blockNonce,
		BlockHash:      blockHash,
		BlockRootHash:  blockRootHash,
		HintEpoch:      hintEpoch,
	}

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

func parseUint32UrlParam(c *gin.Context, name string) (core.OptionalUint32, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return core.OptionalUint32{}, nil
	}

	value, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return core.OptionalUint32{}, err
	}

	return core.OptionalUint32{
		Value:    uint32(value),
		HasValue: true,
	}, nil
}

func parseUint64UrlParam(c *gin.Context, name string) (core.OptionalUint64, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return core.OptionalUint64{}, nil
	}

	value, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return core.OptionalUint64{}, err
	}

	return core.OptionalUint64{
		Value:    value,
		HasValue: true,
	}, nil
}

func parseHexBytesUrlParam(c *gin.Context, name string) ([]byte, error) {
	param := c.Request.URL.Query().Get(name)
	if param == "" {
		return nil, nil
	}

	decoded, err := hex.DecodeString(param)
	if err != nil {
		return nil, err
	}

	return decoded, nil
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
