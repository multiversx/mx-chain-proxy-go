package facade

import "github.com/pkg/errors"

// ErrNilAccountProcessor signals that a nil account processor has been provided
var ErrNilAccountProcessor = errors.New("nil account processor provided")

// ErrNilTransactionProcessor signals that a nil transaction processor has been provided
var ErrNilTransactionProcessor = errors.New("nil transaction processor provided")

// ErrNilSCQueryService signals that a nil smart contracts query service has been provided
var ErrNilSCQueryService = errors.New("nil smart contracts query service provided")

// ErrNilHeartbeatProcessor signals that a nil heartbeat processor has been provided
var ErrNilHeartbeatProcessor = errors.New("nil heartbeat processor provided")

// ErrNilFaucetProcessor signals that a nil faucet processor has been provided
var ErrNilFaucetProcessor = errors.New("nil faucet processor provided")
