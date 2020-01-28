package ring

import "errors"

// ErrInvalidObserversSlice signals that an invalid observers slice has been provided
var ErrInvalidObserversSlice = errors.New("invalid slice of observers provided")
