package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type blockAPIService struct {
	elrondProvider provider.ElrondProviderHandler
	txsParser      *transactionsParser
}

// NewBlockAPIService will create a new instance of blockAPIService
func NewBlockAPIService(
	elrondProvider provider.ElrondProviderHandler,
	cfg *configuration.Configuration,
	networkConfig *provider.NetworkConfig,
) server.BlockAPIServicer {
	return &blockAPIService{
		elrondProvider: elrondProvider,
		txsParser:      newTransactionParser(cfg, networkConfig),
	}
}

// Block implements the /block endpoint.
func (bas *blockAPIService) Block(
	_ context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	if request.BlockIdentifier.Index != nil {
		return bas.getBlockByNonce(*request.BlockIdentifier.Index)
	}

	if request.BlockIdentifier.Hash != nil {
		return bas.getBlockByHash(*request.BlockIdentifier.Hash)
	}

	return nil, ErrMustQueryByIndexOrByHash
}

func (bas *blockAPIService) getBlockByNonce(nonce int64) (*types.BlockResponse, *types.Error) {
	hyperBlock, err := bas.elrondProvider.GetBlockByNonce(nonce)
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	return bas.parseHyperBlock(hyperBlock)
}

func (bas *blockAPIService) getBlockByHash(hash string) (*types.BlockResponse, *types.Error) {
	hyperBlock, err := bas.elrondProvider.GetBlockByHash(hash)
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	return bas.parseHyperBlock(hyperBlock)
}

func (bas *blockAPIService) parseHyperBlock(hyperBlock *data.Hyperblock) (*types.BlockResponse, *types.Error) {
	var parentBlockIdentifier *types.BlockIdentifier
	if hyperBlock.Nonce != 0 {
		parentBlockIdentifier = &types.BlockIdentifier{
			Index: int64(hyperBlock.Nonce - 1),
			Hash:  hyperBlock.PrevBlockHash,
		}
	}

	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Index: int64(hyperBlock.Nonce),
				Hash:  hyperBlock.Hash,
			},
			ParentBlockIdentifier: parentBlockIdentifier,
			Timestamp:             bas.elrondProvider.CalculateBlockTimestampUnix(hyperBlock.Round),
			Transactions:          bas.txsParser.parseTxsFromHyperBlock(hyperBlock),
			Metadata: objectsMap{
				"epoch": hyperBlock.Epoch,
				"round": hyperBlock.Round,
			},
		},
	}, nil
}

// BlockTransaction - not implemented
// We dont need this method because all transactions are returned by method Block
func (bas *blockAPIService) BlockTransaction(
	_ context.Context,
	_ *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	return nil, ErrNotImplemented
}
