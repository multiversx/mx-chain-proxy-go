package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type accountAPIService struct {
	elrondClient client.ElrondClientHandler
	config       *configuration.Configuration
}

// NewAccountAPIService will create a new instance of accountAPIService
func NewAccountAPIService(elrondClient client.ElrondClientHandler, cfg *configuration.Configuration) server.AccountAPIServicer {
	return &accountAPIService{
		elrondClient: elrondClient,
		config:       cfg,
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

	latestBlockData, err := aas.elrondClient.GetLatestBlockData()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	account, err := aas.elrondClient.GetAccount(request.AccountIdentifier.Address)
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
