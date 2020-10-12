package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type mempoolAPIService struct {
	elrondProvider provider.ElrondProviderHandler
	txsParser      *transactionsParser
}

// NewMempoolApiService will create a new instance of mempoolAPIService
func NewMempoolApiService(
	elrondProvider provider.ElrondProviderHandler,
	cfg *configuration.Configuration,
	networkConfig *provider.NetworkConfig,
) server.MempoolAPIServicer {
	return &mempoolAPIService{
		elrondProvider: elrondProvider,
		txsParser:      newTransactionParser(elrondProvider, cfg, networkConfig),
	}
}

// Mempool is not implemented yet
func (mas *mempoolAPIService) Mempool(context.Context, *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	return nil, ErrNotImplemented
}

// MempoolTransaction will return operations for a transaction that is in pool
func (mas *mempoolAPIService) MempoolTransaction(
	_ context.Context,
	request *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	tx, ok := mas.elrondProvider.GetTransactionByHashFromPool(request.TransactionIdentifier.Hash)
	if !ok {
		return nil, ErrTransactionIsNotInPool
	}

	rosettaTx, ok := mas.txsParser.parseTx(tx, true)
	if !ok {
		return nil, ErrCannotParsePoolTransaction
	}

	return &types.MempoolTransactionResponse{
		Transaction: rosettaTx,
	}, nil

}
