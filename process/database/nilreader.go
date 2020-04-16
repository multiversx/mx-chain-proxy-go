package database

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type nilReader struct{}

func NewNilReader() *nilReader {
	return new(nilReader)
}

func (nr *nilReader) GetTransactionsByAddress(_ string) ([]data.ApiTransaction, error) {
	return nil, errors.New("database connection is disabled")
}
