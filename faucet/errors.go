package faucet

import "errors"

// ErrNilShardCoordinator signals that the provided shard coordinator is nil
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrInvalidPemFileLocation signals that the provided path for the pem file is invalid
var ErrInvalidPemFileLocation = errors.New("invalid pem file location")

// ErrNilPubKeyConverter signals that the provided pub key converter is nil
var ErrNilPubKeyConverter = errors.New("nil pub key converter")
