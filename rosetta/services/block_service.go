package services

import (
	"context"
	"github.com/ElrondNetwork/elrond-proxy-go/data"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type blockAPIService struct {
	elrondClient *client.ElrondClient
}

// NewBlockAPIService
func NewBlockAPIService(elrondClient *client.ElrondClient) server.BlockAPIServicer {
	return &blockAPIService{
		elrondClient: elrondClient,
	}
}

// Block implements the /block endpoint.
func (s *blockAPIService) Block(
	_ context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	if request.BlockIdentifier.Index != nil {
		return s.getBlockByNonce(*request.BlockIdentifier.Index)
	}

	if request.BlockIdentifier.Hash != nil {
		return s.getBlockByHash(*request.BlockIdentifier.Hash)
	}

	return nil, ErrMustQueryByIndexOrByHash
}

func (s *blockAPIService) getBlockByNonce(nonce int64) (*types.BlockResponse, *types.Error) {
	hyperBlock, err := s.elrondClient.GetBlockByNonce(nonce)
	if err != nil {
		return nil, ErrUnableToGetBlock
	}

	return s.parseHyperBlock(hyperBlock)
}

func (s *blockAPIService) getBlockByHash(hash string) (*types.BlockResponse, *types.Error) {
	hyperBlock, err := s.elrondClient.GetBlockByHash(hash)
	if err != nil {
		return nil, ErrUnableToGetBlock
	}

	return s.parseHyperBlock(hyperBlock)
}

func (s *blockAPIService) parseHyperBlock(hyperBlock *data.Hyperblock) (*types.BlockResponse, *types.Error) {
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
			Timestamp:             client.CalculateBlockTimestampUnix(hyperBlock.Round),
			Transactions:          parseTxsFromHyperBlock(hyperBlock),
			Metadata: objectsMap{
				"epoch": hyperBlock.Epoch,
				// TODO can add extra data in hyperBlock
			},
		},
	}, nil
}

// BlockTransaction - not implemented
// We dont need this method because all transactions are returned by method Block
func (s *blockAPIService) BlockTransaction(
	_ context.Context,
	_ *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	return nil, ErrNotImplemented
}
