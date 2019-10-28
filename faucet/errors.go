package faucet

import "errors"

// ErrNilAddressConverter signals that the provided address converter is nil
var ErrNilAddressConverter = errors.New("nil address converter")

// ErrNilShardCoordinator signals that the provided shard coordinator is nil
var ErrNilShardCoordinator = errors.New("nil shard coordinator")

// ErrInvalidPemFileLocation signals that the provided path for the pem file is invalid
var ErrInvalidPemFileLocation = errors.New("invalid pem file location")
