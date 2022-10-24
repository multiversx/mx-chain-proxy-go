package common

import (
	"encoding/hex"
	"net/url"
	"strconv"

	"github.com/ElrondNetwork/elrond-go-core/core"
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
	// UrlParameterBlockNonce represents the name of an URL parameter
	UrlParameterBlockNonce = "blockNonce"
	// UrlParameterBlockHash represents the name of an URL parameter
	UrlParameterBlockHash = "blockHash"
	// UrlParameterBlockRootHash represents the name of an URL parameter
	UrlParameterBlockRootHash = "blockRootHash"
	// UrlParameterHintEpoch represents the name of an URL parameter
	UrlParameterHintEpoch = "hintEpoch"
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
	OnStartOfEpoch core.OptionalUint32
	BlockNonce     core.OptionalUint64
	BlockHash      []byte
	BlockRootHash  []byte
	HintEpoch      core.OptionalUint32
}

// BuildUrlWithAccountQueryOptions builds an URL with block query parameters
func BuildUrlWithAccountQueryOptions(path string, options AccountQueryOptions) string {
	u := url.URL{Path: path}
	query := u.Query()

	if options.OnFinalBlock {
		query.Set(UrlParameterOnFinalBlock, "true")
	}
	if options.OnStartOfEpoch.HasValue {
		query.Set(UrlParameterOnStartOfEpoch, strconv.Itoa(int(options.OnStartOfEpoch.Value)))
	}
	if options.BlockNonce.HasValue {
		query.Set(UrlParameterBlockNonce, strconv.FormatUint(options.BlockNonce.Value, 10))
	}
	if len(options.BlockHash) > 0 {
		query.Set(UrlParameterBlockHash, hex.EncodeToString(options.BlockHash))
	}
	if len(options.BlockRootHash) > 0 {
		query.Set(UrlParameterBlockRootHash, hex.EncodeToString(options.BlockRootHash))
	}
	if options.HintEpoch.HasValue {
		query.Set(UrlParameterHintEpoch, strconv.Itoa(int(options.HintEpoch.Value)))
	}

	u.RawQuery = query.Encode()
	return u.String()
}
