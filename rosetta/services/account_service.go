package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type accountAPIService struct {
	elrondClient *client.ElrondClient
}

// NewAccountAPIService
func NewAccountAPIService(elrondClient *client.ElrondClient) server.AccountAPIServicer {
	return &accountAPIService{
		elrondClient: elrondClient,
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
		return nil, ErrUnableToGetBlock
	}

	account, err := s.elrondClient.GetAccount(request.AccountIdentifier.Address)
	if err != nil {
		return nil, ErrUnableToGetAccount
	}

	response := &types.AccountBalanceResponse{
		BlockIdentifier: &types.BlockIdentifier{
			Index: int64(latestBlockData.Nonce),
			Hash:  latestBlockData.Hash,
		},
		Balances: []*types.Amount{
			{
				Value:    account.Balance,
				Currency: ElrondCurrency,
			},
		},
		Metadata: map[string]interface{}{
			"nonce": account.Nonce,
		},
	}

	return response, nil
}
