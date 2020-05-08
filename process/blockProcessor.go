package process

import (
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type blockProcessor struct {
	dbReader DatabaseReader
}

// NewBlockProcessor will create a new block processor
func NewBlockProcessor(dbReader DatabaseReader) (*blockProcessor, error) {
	if check.IfNil(dbReader) {
		return nil, ErrNilDatabaseReader
	}

	return &blockProcessor{
		dbReader: dbReader,
	}, nil
}

func (bp *blockProcessor) GetHighestBlockNonce() (uint64, error) {
	return bp.dbReader.GetLatestBlockHeight()
}

func (bp *blockProcessor) GetBlockByNonce(nonce uint64) (data.ApiBlock, error) {
	return bp.dbReader.GetBlockByNonce(nonce)
}
