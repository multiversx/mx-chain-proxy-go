package database

import "errors"

var errCannotFindBlockInDb = errors.New("cannot find blocks in database")
var errCannotUnmarshalBlock = errors.New("cannot unmarshal block")
var errCannotGetTxsFromBody = errors.New("cannot get transactions from decoded body")
