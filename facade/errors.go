package facade

import "github.com/pkg/errors"

// ErrNilAccountProccessor signals that a nil account processor has been provided
var ErrNilAccountProccessor = errors.New("nil account processor provided")
