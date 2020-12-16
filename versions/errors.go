package versions

import "errors"

// ErrNilFacadeHandler signals that a nil facade handler has been provided
var ErrNilFacadeHandler = errors.New("nil facade handler")

// ErrNilApiHandler signals that the provided api handler is nil
var ErrNilApiHandler = errors.New("nil api handler")

// ErrNoVersionIsSet signals that no version is provided in the environment
var ErrNoVersionIsSet = errors.New("no version is set")

// ErrVersionNotFound signals that a provided version does not exist
var ErrVersionNotFound = errors.New("version not found")
