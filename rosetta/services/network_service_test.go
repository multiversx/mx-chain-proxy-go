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

func TestNetworkAPIService_NetworkList(t *testing.T) {
	t.Parallel()

	elrondClientMock := &mocks.ElrondClientMock{}
	cfg := &configuration.Configuration{
		Network: &types.NetworkIdentifier{
			Blockchain: configuration.BlockchainName,
			Network:    "local_network",
		},
	}

	networkAPIService := NewNetworkAPIService(elrondClientMock, cfg)

	networkListResponse, err := networkAPIService.NetworkList(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, []*types.NetworkIdentifier{{
		Blockchain: configuration.BlockchainName,
		Network:    "local_network",
	}}, networkListResponse.NetworkIdentifiers)
}

func TestNetworkAPIService_NetworkOptions(t *testing.T) {
	t.Parallel()

	clientVersion := "1"
	elrondClientMock := &mocks.ElrondClientMock{
		GetNetworkConfigCalled: func() (*client.NetworkConfig, error) {
			return &client.NetworkConfig{
				ClientVersion: clientVersion,
			}, nil
		},
	}
	cfg := &configuration.Configuration{
		Network: &types.NetworkIdentifier{
			Blockchain: configuration.BlockchainName,
			Network:    "local_network",
		},
	}
	networkAPIService := NewNetworkAPIService(elrondClientMock, cfg)

	networkOptions, err := networkAPIService.NetworkOptions(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: RosettaVersion,
			NodeVersion:    clientVersion,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     OpStatusSuccess,
					Successful: true,
				},
				{
					Status:     OpStatusFailed,
					Successful: false,
				},
			},
			OperationTypes: SupportedOperationTypes,
			Errors:         Errors,
		},
	}, networkOptions)
}

func TestNetworkAPIService_NetworkStatus(t *testing.T) {
	t.Parallel()

	latestBlockNonce := int64(1000)
	latestBlockHash := "hash"
	oldestBlockNonce := int64(800)
	oldestBlockHash := "old"
	elrondClientMock := &mocks.ElrondClientMock{
		GetLatestBlockDataCalled: func() (*client.BlockData, error) {
			return &client.BlockData{
				Hash:  latestBlockHash,
				Nonce: uint64(latestBlockNonce),
			}, nil
		},
		GetBlockByNonceCalled: func(nonce int64) (*data.Hyperblock, error) {
			return &data.Hyperblock{
				Hash:  oldestBlockHash,
				Nonce: uint64(oldestBlockNonce),
			}, nil
		},
	}
	cfg := &configuration.Configuration{
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: 1,
			Hash:  configuration.GenesisBlockHashMainnet,
		},
		Peers: []*types.Peer{
			{
				PeerID: "bla-bla-bla",
			},
		},
	}
	networkAPIService := NewNetworkAPIService(elrondClientMock, cfg)

	networkStatusResponse, err := networkAPIService.NetworkStatus(context.Background(), nil)
	assert.Nil(t, err)
	assert.Equal(t, &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: latestBlockNonce,
			Hash:  latestBlockHash,
		},
		CurrentBlockTimestamp:  0,
		GenesisBlockIdentifier: cfg.GenesisBlockIdentifier,
		OldestBlockIdentifier: &types.BlockIdentifier{
			Index: oldestBlockNonce,
			Hash:  oldestBlockHash,
		},
		SyncStatus: nil,
		Peers:      cfg.Peers,
	}, networkStatusResponse)
}
