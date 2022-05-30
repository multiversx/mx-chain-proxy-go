package common

import "net/url"

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
		query.Set("withTxs", "true")
	}
	if options.WithLogs {
		query.Set("withLogs", "true")
	}

	url.RawQuery = query.Encode()
	return url.String()
}
