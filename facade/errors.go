package facade

import "github.com/pkg/errors"

// ErrNilAccountProcessor signals that a nil account processor has been provided
var ErrNilAccountProcessor = errors.New("nil account processor provided")

// ErrNilTransactionProcessor signals that a nil transaction processor has been provided
var ErrNilTransactionProcessor = errors.New("nil transaction processor provided")

// ErrNilVmValueProcessor signals that a nil vm value processor has been provided
var ErrNilVmValueProcessor = errors.New("nil vm value processor provided")

// ErrNilHeartbeatProcessor signals that a nil heartbeat processor has been provided
var ErrNilHeartbeatProcessor = errors.New("nil heartbeat processor provided")
