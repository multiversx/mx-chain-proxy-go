package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type accountAPIService struct {
	elrondClient *client.ElrondClient
	config       *configuration.Configuration
}

// NewAccountAPIService
func NewAccountAPIService(elrondClient *client.ElrondClient, cfg *configuration.Configuration) server.AccountAPIServicer {
	return &accountAPIService{
		elrondClient: elrondClient,
		config:       cfg,
	}
}

// AccountBalance implements the /account/balance endpoint.
func (s *accountAPIService) AccountBalance(
	_ context.Context,
	request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	// TODO cannot return balance at a specific nonce right now
	if request.AccountIdentifier.Address == "" {
		return nil, ErrInvalidAccountAddress
	}

	latestBlockData, err := s.elrondClient.GetLatestBlockData()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	account, err := s.elrondClient.GetAccount(request.AccountIdentifier.Address)
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
				Currency: s.config.Currency,
			},
		},
		Metadata: map[string]interface{}{
			"nonce": account.Nonce,
		},
	}

	return response, nil
}
