package faucet

import "errors"

// ErrNilShardCoordinator signals that the provided shard coordinator is nil
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrFaucetPemFileDoesNotExist signals that the faucet pem file does not exist
var ErrFaucetPemFileDoesNotExist = errors.New("faucet pem file does not exist")

// ErrNilPubKeyConverter signals that the provided pub key converter is nil
var ErrNilPubKeyConverter = errors.New("nil pub key converter")
