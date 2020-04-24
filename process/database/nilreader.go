package database

import (
	"errors"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

var errNilReaderImplementation = errors.New("database connection is disabled")

type nilReader struct{}

func NewNilReader() *nilReader {
	return new(nilReader)
}

// GetTransactionsByAddress -
func (nr *nilReader) GetTransactionsByAddress(_ string) ([]data.DatabaseTransaction, error) {
	return nil, errNilReaderImplementation
}

// GetLatestBlockHeight -
func (nr *nilReader) GetLatestBlockHeight() (uint64, error) {
	return 0, errNilReaderImplementation
}

// GetBlockByNonce -
func (nr *nilReader) GetBlockByNonce(_ uint64) (data.ApiBlock, error) {
	return data.ApiBlock{}, errNilReaderImplementation
}

// IsInterfaceNil -
func (nr *nilReader) IsInterfaceNil() bool {
	return nr == nil
}
