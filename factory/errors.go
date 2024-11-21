package factory

import (
	"errors"
)

// ErrNilRunTypeComponents signals that nil run type components were provided
var ErrNilRunTypeComponents = errors.New("nil run type components")
