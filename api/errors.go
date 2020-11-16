package api

import "errors"

// ErrNilGroupHandler signals that a nil group handler has been provided
var ErrNilGroupHandler = errors.New("nil group handler")

// ErrGroupAlreadyRegistered signals that the provided group has already been registered
var ErrGroupAlreadyRegistered = errors.New("group already registered")

// ErrGroupDoesNotExist signals that the called group does not exist
var ErrGroupDoesNotExist = errors.New("group does not exist")

// ErrNilFacade signals that a nil facade has been provided
var ErrNilFacade = errors.New("nil facade")
