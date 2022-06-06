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
	// UrlParameterOnFinalBlock represents the name of an URL parameter
	UrlParameterOnFinalBlock = "onFinalBlock"
	// UrlParameterOnStartOfEpoch represents the name of an URL parameter
	UrlParameterOnStartOfEpoch = "onStartOfEpoch"
)

// BlockQueryOptions holds options for block queries
type BlockQueryOptions struct {
	WithTransactions bool
	WithLogs         bool
}

// HyperblockQueryOptions holds options for hyperblock queries
type HyperblockQueryOptions struct {
	WithLogs bool
}

// BuildUrlWithBlockQueryOptions builds an URL with block query parameters
func BuildUrlWithBlockQueryOptions(path string, options BlockQueryOptions) string {
	url := url.URL{Path: path}
	query := url.Query()

	if options.WithTransactions {
		query.Set(UrlParameterWithTransactions, "true")
	}
	if options.WithLogs {
		query.Set(UrlParameterWithLogs, "true")
	}

	url.RawQuery = query.Encode()
	return url.String()
}

// AccountQueryOptions holds options for account queries
type AccountQueryOptions struct {
	OnFinalBlock   bool
	OnStartOfEpoch uint32
}

// BuildUrlWithAccountQueryOptions builds an URL with block query parameters
func BuildUrlWithAccountQueryOptions(path string, options AccountQueryOptions) string {
	url := url.URL{Path: path}
	query := url.Query()

	if options.OnFinalBlock {
		query.Set(UrlParameterOnFinalBlock, "true")
	}
	if options.OnStartOfEpoch != 0 {
		query.Set(UrlParameterOnStartOfEpoch, strconv.Itoa(int(options.OnStartOfEpoch)))
	}

	url.RawQuery = query.Encode()
	return url.String()
}
