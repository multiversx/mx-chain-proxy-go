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

// ErrInvalidRequestTimeout signals that the provided number of seconds before timeout is invalid
var ErrInvalidRequestTimeout = errors.New("invalid duration until timeout for requests")

// ErrNilCoreProcessor signals that a nil core processor has been provided
var ErrNilCoreProcessor = errors.New("nil core processor")

// ErrNilKeyGen signals that a nil key generator has been provided
var ErrNilKeyGen = errors.New("nil keygen")

// ErrNilSingleSigner signals that a nil single signer has been provided
var ErrNilSingleSigner = errors.New("nil single signer")

// ErrNoObserverConnected signals that no observer from the list is online
var ErrNoObserverConnected = errors.New("no observer is online")

// ErrHeartbeatNotAvailable signals that the heartbeat status is not found
var ErrHeartbeatNotAvailable = errors.New("heartbeat status not found at any observer")
