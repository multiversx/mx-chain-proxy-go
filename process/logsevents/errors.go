package logsevents

import "errors"

// ErrNilHasher is raised when a valid hasher is expected but nil used
var ErrNilHasher = errors.New("hasher is nil")

// ErrNilMarshalizer is raised when a valid marshalizer is expected but nil used
var ErrNilMarshalizer = errors.New("marshalizer is nil")
