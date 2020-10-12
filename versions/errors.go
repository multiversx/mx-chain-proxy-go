package versions

import "errors"

// ErrNilFacadeHandler signals that a nil facade handler has been provided
var ErrNilFacadeHandler = errors.New("nil facade handler")

// ErrVersionNotFound signals that the provided version has no setup in the environment
var ErrVersionNotFound = errors.New("version not found")

// ErrNoVersionIsSet signals that no version is provided in the environment
var ErrNoVersionIsSet = errors.New("no version is set")
