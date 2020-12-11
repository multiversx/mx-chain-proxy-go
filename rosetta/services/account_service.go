package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type accountAPIService struct {
	elrondProvider provider.ElrondProviderHandler
	config         *configuration.Configuration
}

// NewAccountAPIService will create a new instance of accountAPIService
func NewAccountAPIService(elrondProvider provider.ElrondProviderHandler, cfg *configuration.Configuration) server.AccountAPIServicer {
	return &accountAPIService{
		elrondProvider: elrondProvider,
		config:         cfg,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (aas *accountAPIService) AccountBalance(
	_ context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	// TODO cannot return balance at a specific nonce right now
	if request.AccountIdentifier.Address == "" {
		return nil, ErrInvalidAccountAddress
	}

	latestBlockData, err := aas.elrondProvider.GetLatestBlockData()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	account, err := aas.elrondProvider.GetAccount(request.AccountIdentifier.Address)
	if err != nil {
		return nil, wrapErr(ErrUnableToGetAccount, err)
	}

	response := &types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(latestBlockData.Nonce),
			Hash:  latestBlockData.Hash,
		},
		Balances: []*types.Amount{
			{
				Value:    account.Balance,
				Currency: aas.config.Currency,
			},
		},
		Metadata: map[string]interface{}{
			"nonce": account.Nonce,
		},
	}

	return response, nil
}

// AccountCoins implements the /account/coins endpoint.
func (aas *accountAPIService) AccountCoins(_ context.Context, _ *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	return nil, ErrNotImplemented
}
