package txcost

import "errors"

// ErrNilPubKeyConverter signals that a nil pub key converter has been provided
var ErrNilPubKeyConverter = errors.New("nil pub key converter provided")

// ErrNilCoreProcessor signals that a nil core processor has been provided
var ErrNilCoreProcessor = errors.New("nil core processor")

// ErrSendingRequest signals that sending the request failed on all observers
var ErrSendingRequest = errors.New("sending request error")
