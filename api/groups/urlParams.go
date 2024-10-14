package groups

import (
	"encoding/hex"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/common"
)

// SystemAccountAddressBech is the const for the system account address
const SystemAccountAddressBech = "erd1lllllllllllllllllllllllllllllllllllllllllllllllllllsckry7t"

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

func parseAccountQueryOptions(c *gin.Context, address string) (common.AccountQueryOptions, error) {
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

	shardID, err := parseUint32UrlParam(c, common.UrlParameterForcedShardID)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	withKeys, err := parseBoolUrlParam(c, common.UrlParameterWithKeys)
	if err != nil {
		return common.AccountQueryOptions{}, err
	}

	if shardID.HasValue && address != SystemAccountAddressBech {
		return common.AccountQueryOptions{}, ErrForcedShardIDCannotBeProvided
	}

	options := common.AccountQueryOptions{
		OnFinalBlock:   onFinalBlock,
		OnStartOfEpoch: onStartOfEpoch,
		BlockNonce:     blockNonce,
		BlockHash:      blockHash,
		BlockRootHash:  blockRootHash,
		HintEpoch:      hintEpoch,
		ForcedShardID:  shardID,
		WithKeys:       withKeys,
	}

	return options, nil
}

func parseTransactionQueryOptions(c *gin.Context) (common.TransactionQueryOptions, error) {
	withResults, err := parseBoolUrlParam(c, common.UrlParameterWithResults)
	if err != nil {
		return common.TransactionQueryOptions{}, err
	}

	relayedTxHash := parseStringUrlParam(c, common.UrlParameterRelayedTxHash)

	options := common.TransactionQueryOptions{
		WithResults:   withResults,
		RelayedTxHash: relayedTxHash,
	}
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

func parseAlteredAccountOptions(c *gin.Context) (common.GetAlteredAccountsForBlockOptions, error) {
	tokensFilter := parseStringUrlParam(c, common.UrlParameterTokensFilter)

	return common.GetAlteredAccountsForBlockOptions{
		TokensFilter: tokensFilter,
	}, nil
}
