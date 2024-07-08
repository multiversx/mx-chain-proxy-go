package resultsParser

import (
	"errors"
)

var (
	ErrNoReturnCode           = errors.New("no return code")
	ErrEmptyDataField         = errors.New("empty data field")
	ErrFoundMoreThanOneEvent  = errors.New("found more than one event")
	ErrCannotProcessDataField = errors.New("cannot process data field")
)
