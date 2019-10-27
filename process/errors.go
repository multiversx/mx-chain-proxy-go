package process

import "errors"

// ErrNilConfig signals that a nil config has been provided
var ErrNilConfig = errors.New("nil configuration provided")

// ErrEmptyObserversList signals that an empty list of observers has been provided
var ErrEmptyObserversList = errors.New("empty observers list provided")

// ErrMissingObserver signals that no observers have been provided for provided shard ID
var ErrMissingObserver = errors.New("missing observer")

// ErrSendingRequest signals that sending the request failed on all observers
var ErrSendingRequest = errors.New("sending request error")

// ErrNilAddressConverter signals that a nil address converter has been provided
var ErrNilAddressConverter = errors.New("nil address converter")

// ErrNilShardCoordinator signals that a nil shard coordinator has been provided
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrNilCoreProcessor signals that a nil core processor has been provided
var ErrNilCoreProcessor = errors.New("nil core processor")

// ErrNilPrivateKeysLoader signals that a nil private keys loader has been provided
var ErrNilPrivateKeysLoader = errors.New("nil private keys loader")

// ErrEmptyMapOfAccountsFromPem signals that an empty map of accounts was received
var ErrEmptyMapOfAccountsFromPem = errors.New("empty map of accounts read from the pem file")

// ErrNoObserverConnected signals that no observer from the list is online
var ErrNoObserverConnected = errors.New("no observer is online")

// ErrHeartbeatNotAvailable signals that the heartbeat status is not found
var ErrHeartbeatNotAvailable = errors.New("heartbeat status not found at any observer")

// ErrNilDefaultFaucetValue signals that a nil default faucet value has been provided
var ErrNilDefaultFaucetValue = errors.New("nil default faucet value provided")

// ErrInvalidDefaultFaucetValue signals that the provided faucet value is not strictly positive
var ErrInvalidDefaultFaucetValue = errors.New("default faucet value is not strictly positive")
