package middleware

import "errors"

// ErrNilLimitsMapForEndpoints signals that a nil limits map has been provided
var ErrNilLimitsMapForEndpoints = errors.New("nil limits map")
