package factory

import "errors"

// ErrNoDirectoryAtPath signals that the file is not a directory
var ErrNoDirectoryAtPath = errors.New("no directory at the given path")
