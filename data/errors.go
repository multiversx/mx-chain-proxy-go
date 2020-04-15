package data

import "errors"

// ErrNilTransaction signals that a nil transaction has been provided
var ErrNilTransaction = errors.New("nil transaction")

// ErrNilPubKeyConverter signals that a nil pub key converter has been provided
var ErrNilPubKeyConverter = errors.New("nil pub key converter")
