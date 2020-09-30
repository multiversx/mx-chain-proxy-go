package services

import (
	"context"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/mocks"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestAccountAPIService_AccountBalance(t *testing.T) {
	t.Parallel()

	address := "erd13lx7zldumunqvf74g5z407gwl5r35jha06rjzc32qcujamknzdgsnt2yvn"
	accountBalance := "1234"
	latestBlockNonce := uint64(1)
	lastestBlockHash := "hash-hash-hash"
	elrondClientMock := &mocks.ElrondClientMock{
		GetAccountCalled: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: "erd13lx7zldumunqvf74g5z407gwl5r35jha06rjzc32qcujamknzdgsnt2yvn",
				Nonce:   1,
				Balance: "1234",
			}, nil
		},
		GetLatestBlockDataCalled: func() (*client.BlockData, error) {
			return &client.BlockData{
				Nonce: latestBlockNonce,
				Hash:  lastestBlockHash,
			}, nil
		},
	}
	cfg := &configuration.Configuration{}

	accountAPIService := NewAccountAPIService(elrondClientMock, cfg)
	assert.NotNil(t, accountAPIService)

	_, err := accountAPIService.AccountBalance(context.Background(), &types.AccountBalanceRequest{
		AccountIdentifier: &types.AccountIdentifier{
			Address: "",
		},
	})
	assert.Equal(t, ErrInvalidAccountAddress, err)

	// Get account balance should work
	accountBalanceResponse, err := accountAPIService.AccountBalance(context.Background(), &types.AccountBalanceRequest{
		AccountIdentifier: &types.AccountIdentifier{
			Address: address,
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, accountBalance, accountBalanceResponse.Balances[0].Value)
	assert.Equal(t, lastestBlockHash, accountBalanceResponse.BlockIdentifier.Hash)
	assert.Equal(t, int64(latestBlockNonce), accountBalanceResponse.BlockIdentifier.Index)
}
