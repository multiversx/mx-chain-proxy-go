package factory

import "errors"

// ErrFileNotFound signals that the provided path is invalid and does not belong to a file or a directory
var ErrFileNotFound = errors.New("file not found")

// ErrFileIsNotADirectory signals that the file is not a directory
var ErrFileIsNotADirectory = errors.New("file is not a directory")
