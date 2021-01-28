package helpers

import "errors"

// ErrEmptyOwnerAddress signals that an empty owner address has been provided
var ErrEmptyOwnerAddress = errors.New("empty owner address")
