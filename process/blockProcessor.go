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

// GetBlockByShardIDAndNonce return the block byte shardID and nonce
func (bp *dbBlockProcessor) GetBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.ApiBlock, error) {
	return bp.dbReader.GetBlockByShardIDAndNonce(shardID, nonce)
}
