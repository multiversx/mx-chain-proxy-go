package common

import (
	"net/url"
	"strconv"
)

const (
	// UrlParameterWithTransactions represents the name of an URL parameter
	UrlParameterWithTransactions = "withTxs"
	// UrlParameterWithLogs represents the name of an URL parameter
	UrlParameterWithLogs = "withLogs"
	// UrlParameterNotarizedAtSource represents the name of an URL parameter
	UrlParameterNotarizedAtSource = "notarizedAtSource"
	// UrlParameterOnFinalBlock represents the name of an URL parameter
	UrlParameterOnFinalBlock = "onFinalBlock"
	// UrlParameterOnStartOfEpoch represents the name of an URL parameter
	UrlParameterOnStartOfEpoch = "onStartOfEpoch"
	// UrlParameterCheckSignature represents the name of an URL parameter
	UrlParameterCheckSignature = "checkSignature"
	// UrlParameterWithResults represents the name of an URL parameter
	UrlParameterWithResults = "withResults"
	// UrlParameterShardID represents the name of an URL parameter
	UrlParameterShardID = "shard-id"
	// UrlParameterSender represents the name of an URL parameter
	UrlParameterSender = "by-sender"
	// UrlParameterFields represents the name of an URL parameter
	UrlParameterFields = "fields"
	// UrlParameterLastNonce represents the name of an URL parameter
	UrlParameterLastNonce = "last-nonce"
	// UrlParameterNonceGaps represents the name of an URL parameter
	UrlParameterNonceGaps = "nonce-gaps"
	// UrlParameterWithMetadata represents the name of an URL parameter
	UrlParameterWithMetadata = "withMetadata"
	// UrlParameterTokensFilter represents the name of an URL parameter
	UrlParameterTokensFilter = "tokens"
	// UrlParameterWithAlteredAccounts represents the name of an URL parameter
	UrlParameterWithAlteredAccounts = "withAlteredAccounts"
)

// BlockQueryOptions holds options for block queries
type BlockQueryOptions struct {
	WithTransactions bool
	WithLogs         bool
}

// HyperblockQueryOptions holds options for hyperblock queries
type HyperblockQueryOptions struct {
	WithLogs               bool
	NotarizedAtSource      bool
	WithAlteredAccounts    bool
	AlteredAccountsOptions GetAlteredAccountsForBlockOptions
}

// TransactionQueryOptions holds options for transaction queries
type TransactionQueryOptions struct {
	WithResults bool
}

// TransactionSimulationOptions holds options for transaction simulation requests
type TransactionSimulationOptions struct {
	CheckSignature bool
}

// TransactionsPoolOptions holds options for transactions pool requests
type TransactionsPoolOptions struct {
	ShardID   string
	Sender    string
	Fields    string
	LastNonce bool
	NonceGaps bool
}

// GetAlteredAccountsForBlockOptions specifies the options for returning altered accounts for a given block
type GetAlteredAccountsForBlockOptions struct {
	TokensFilter string
	WithMetadata bool
}

// BuildUrlWithBlockQueryOptions builds an URL with block query parameters
func BuildUrlWithBlockQueryOptions(path string, options BlockQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	if options.WithTransactions {
		query.Set(UrlParameterWithTransactions, "true")
	}
	if options.WithLogs {
		query.Set(UrlParameterWithLogs, "true")
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// AccountQueryOptions holds options for account queries
type AccountQueryOptions struct {
	OnFinalBlock   bool
	OnStartOfEpoch uint32
}

// BuildUrlWithAccountQueryOptions builds an URL with block query parameters
func BuildUrlWithAccountQueryOptions(path string, options AccountQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	if options.OnFinalBlock {
		query.Set(UrlParameterOnFinalBlock, "true")
	}
	if options.OnStartOfEpoch != 0 {
		query.Set(UrlParameterOnStartOfEpoch, strconv.Itoa(int(options.OnStartOfEpoch)))
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// BuildUrlWithAlteredAccountsQueryOptions builds an URL with altered accounts parameters
func BuildUrlWithAlteredAccountsQueryOptions(path string, options GetAlteredAccountsForBlockOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	if len(options.TokensFilter) != 0 {
		query.Set(UrlParameterTokensFilter, options.TokensFilter)
	}
	if options.WithMetadata {
		query.Set(UrlParameterWithMetadata, "true")
	}

	u.RawQuery = query.Encode()
	return u.String()
}
