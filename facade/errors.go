package facade

import "github.com/pkg/errors"

// ErrNilAccountProcessor signals that a nil account processor has been provided
var ErrNilAccountProcessor = errors.New("nil account processor provided")

// ErrNilTransactionProcessor signals that a nil transaction processor has been provided
var ErrNilTransactionProcessor = errors.New("nil transaction processor provided")
