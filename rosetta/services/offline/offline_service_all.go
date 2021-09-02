package offline

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/services"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type offlineService struct{}

func NewOfflineService() *offlineService {
	return &offlineService{}
}

// AccountBalance implements the /account/balance endpoint.
func (os *offlineService) AccountBalance(
	_ context.Context,
	_ *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// AccountCoins implements the /account/coins endpoint.
func (os *offlineService) AccountCoins(_ context.Context, _ *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// Block implements the /block endpoint.
func (os *offlineService) Block(
	_ context.Context,
	_ *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// BlockTransaction - not implemented
// We dont need this method because all transactions are returned by method Block
func (os *offlineService) BlockTransaction(
	_ context.Context,
	_ *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// Mempool is not implemented yet
func (os *offlineService) Mempool(context.Context, *types.NetworkRequest) (*types.MempoolResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// MempoolTransaction will return operations for a transaction that is in pool
func (os *offlineService) MempoolTransaction(
	_ context.Context,
	_ *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// NetworkStatus implements the /network/status endpoint.
func (os *offlineService) NetworkStatus(
	_ context.Context,
	_ *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// NetworkOptions implements the /network/options endpoint.
func (os *offlineService) NetworkOptions(
	_ context.Context,
	_ *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}

// NetworkList implements the /network/list endpoint
func (os *offlineService) NetworkList(
	_ context.Context,
	_ *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return nil, services.ErrOfflineMode
}
