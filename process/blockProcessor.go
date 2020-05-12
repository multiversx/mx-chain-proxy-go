package process

import (
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type dbBlockProcessor struct {
	dbReader ExternalStorageConnector
}

// NewBlockProcessor will create a new block processor
func NewBlockProcessor(dbReader ExternalStorageConnector) (*dbBlockProcessor, error) {
	if check.IfNil(dbReader) {
		return nil, ErrNilDatabaseConnector
	}

	return &dbBlockProcessor{
		dbReader: dbReader,
	}, nil
}

func (bp *dbBlockProcessor) GetHighestBlockNonce() (uint64, error) {
	return bp.dbReader.GetLatestBlockHeight()
}

func (bp *dbBlockProcessor) GetBlockByNonce(nonce uint64) (data.ApiBlock, error) {
	return bp.dbReader.GetBlockByNonce(nonce)
}
