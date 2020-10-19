package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type networkAPIService struct {
	elrondProvider provider.ElrondProviderHandler
	config         *configuration.Configuration
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(elrondProvider provider.ElrondProviderHandler, cfg *configuration.Configuration) server.NetworkAPIServicer {
	return &networkAPIService{
		elrondProvider: elrondProvider,
		config:         cfg,
	}
}

// NetworkList implements the /network/list endpoint
func (nas *networkAPIService) NetworkList(
	_ context.Context,
	_ *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			nas.config.Network,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (nas *networkAPIService) NetworkStatus(
	_ context.Context,
	_ *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	latestBlockData, err := nas.elrondProvider.GetLatestBlockData()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetNodeStatus, err)
	}

	networkStatusResponse := &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(latestBlockData.Nonce),
			Hash:  latestBlockData.Hash,
		},
		CurrentBlockTimestamp:  latestBlockData.Timestamp,
		GenesisBlockIdentifier: nas.config.GenesisBlockIdentifier,
		Peers:                  nas.config.Peers,
	}

	oldBlock, err := nas.getOldestBlock(latestBlockData.Nonce)
	if err == nil {
		networkStatusResponse.OldestBlockIdentifier = &types.BlockIdentifier{
			Index: int64(oldBlock.Nonce),
			Hash:  oldBlock.Hash,
		}
	}

	return networkStatusResponse, nil
}

func (nas *networkAPIService) getOldestBlock(latestBlockNonce uint64) (*provider.BlockData, error) {
	oldestBlockNonce := uint64(1)

	if latestBlockNonce > NumBlocksToGet {
		oldestBlockNonce = latestBlockNonce - NumBlocksToGet
	}

	block, err := nas.elrondProvider.GetBlockByNonce(int64(oldestBlockNonce))
	if err != nil {
		return nil, err
	}

	return &provider.BlockData{
		Nonce: block.Nonce,
		Hash:  block.Hash,
	}, nil

}

// NetworkOptions implements the /network/options endpoint.
func (nas *networkAPIService) NetworkOptions(
	_ context.Context,
	_ *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: RosettaVersion,
			NodeVersion:    NodeVersion,
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
	}, nil
}
