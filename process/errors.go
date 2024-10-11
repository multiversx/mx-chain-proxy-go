package process

import "errors"

// ErrMissingObserver signals that no observers have been provided for provided shard ID
var ErrMissingObserver = errors.New("missing observer")

// ErrSendingRequest signals that sending the request failed on all observers
var ErrSendingRequest = errors.New("sending request error")

// ErrNilShardCoordinator signals that a nil shard coordinator has been provided
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrInvalidRequestTimeout signals that the provided number of seconds before timeout is invalid
var ErrInvalidRequestTimeout = errors.New("invalid duration until timeout for requests")

// ErrNilCoreProcessor signals that a nil core processor has been provided
var ErrNilCoreProcessor = errors.New("nil core processor")

// ErrNilPrivateKeysLoader signals that a nil private keys loader has been provided
var ErrNilPrivateKeysLoader = errors.New("nil private keys loader")

// ErrEmptyMapOfAccountsFromPem signals that an empty map of accounts was received
var ErrEmptyMapOfAccountsFromPem = errors.New("empty map of accounts read from the pem file")

// ErrHeartbeatNotAvailable signals that the heartbeat status is not found
var ErrHeartbeatNotAvailable = errors.New("heartbeat status not found at any observer")

// ErrNilHeartbeatCacher signals that the provided heartbeat cacher is nil
var ErrNilHeartbeatCacher = errors.New("nil heartbeat cacher")

// ErrNilValidatorStatisticsCacher signals that the provided validator statistics cacher is nil
var ErrNilValidatorStatisticsCacher = errors.New("nil validator statistics cacher")

// ErrNilEconomicMetricsCacher signals that the provided economic metrics cacher is nil
var ErrNilEconomicMetricsCacher = errors.New("nil economic metrics cacher")

// ErrValidatorStatisticsNotAvailable signals that the validator statistics data is not found
var ErrValidatorStatisticsNotAvailable = errors.New("validator statistics data not found on any observer")

// ErrAuctionListNotAvailable signals that the auction list data is not found
var ErrAuctionListNotAvailable = errors.New("auction list data not found on any observer")

// ErrInvalidCacheValidityDuration signals that the given validity duration for cache data is invalid
var ErrInvalidCacheValidityDuration = errors.New("invalid cache validity duration")

// ErrNilDefaultFaucetValue signals that a nil default faucet value has been provided
var ErrNilDefaultFaucetValue = errors.New("nil default faucet value provided")

// ErrInvalidDefaultFaucetValue signals that the provided faucet value is not strictly positive
var ErrInvalidDefaultFaucetValue = errors.New("default faucet value is not strictly positive")

// ErrNoFaucetAccountForGivenShard signals that no account was found for the shard of the given address
var ErrNoFaucetAccountForGivenShard = errors.New("no faucet account found for the given shard")

// ErrNilNodesProvider signals that a nil observers provider has been provided
var ErrNilNodesProvider = errors.New("nil nodes provider")

// ErrNilPubKeyConverter signals that a nil pub key converter has been provided
var ErrNilPubKeyConverter = errors.New("nil pub key converter provided")

// ErrNoValidTransactionToSend signals that no valid transaction were received
var ErrNoValidTransactionToSend = errors.New("no valid transaction to send")

// ErrCannotParseNodeStatusMetrics signals that the node status metrics cannot be parsed
var ErrCannotParseNodeStatusMetrics = errors.New("cannot parse node status metrics")

// ErrNilHasher is raised when a valid hasher is expected but nil used
var ErrNilHasher = errors.New("hasher is nil")

// ErrNilMarshalizer is raised when a valid marshalizer is expected but nil used
var ErrNilMarshalizer = errors.New("marshalizer is nil")

// ErrNilNewTxCostHandlerFunc is raised when a nil function that creates a new transaction cost handler has been provided
var ErrNilNewTxCostHandlerFunc = errors.New("nil new transaction cost handler function")

// ErrInvalidTransactionValueField signals that field value of transaction is invalid
var ErrInvalidTransactionValueField = errors.New("invalid transaction value field")

// ErrInvalidAddress signals that an invalid address has been provided
var ErrInvalidAddress = errors.New("could not create address from provided param")

// ErrInvalidSignatureBytes signal that an invalid signature hash been provided
var ErrInvalidSignatureBytes = errors.New("invalid signatures bytes")

// ErrNoObserverAvailable signals that no observer could be found
var ErrNoObserverAvailable = errors.New("no observer available")

// ErrInvalidTokenType signals that the provided token type is invalid
var ErrInvalidTokenType = errors.New("invalid token type")

// ErrNilLogsMerger signals that the provided logs merger is nil
var ErrNilLogsMerger = errors.New("nil logs merger")

// ErrNilSCQueryService signals that a nil smart contracts query service has been provided
var ErrNilSCQueryService = errors.New("nil smart contracts query service provided")

// ErrInvalidOutputFormat signals that the output format type is not valid
var ErrInvalidOutputFormat = errors.New("the output format type is invalid")

// ErrNilStatusMetricsProvider signals that a nil status metrics provider has been given
var ErrNilStatusMetricsProvider = errors.New("nil status metrics provider")

// ErrEmptyAppVersionString signals than an empty app version string has been provided
var ErrEmptyAppVersionString = errors.New("empty app version string")

// ErrEmptyCommitString signals than an empty commit id string has been provided
var ErrEmptyCommitString = errors.New("empty commit id string")

// ErrEmptyPubKey signals that an empty public key has been provided
var ErrEmptyPubKey = errors.New("public key is empty")

// ErrNilHttpClient signals that a nil http client has been provided
var ErrNilHttpClient = errors.New("nil http client")

// ErrInvalidHash signals that an invalid hash has been provided
var ErrInvalidHash = errors.New("invalid hash")
