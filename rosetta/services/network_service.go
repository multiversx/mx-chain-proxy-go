package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type networkAPIService struct {
	elrondClient client.ElrondClientHandler
	config       *configuration.Configuration
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(elrondClient client.ElrondClientHandler, cfg *configuration.Configuration) server.NetworkAPIServicer {
	return &networkAPIService{
		elrondClient: elrondClient,
		config:       cfg,
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
	latestBlockData, err := nas.elrondClient.GetLatestBlockData()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetNodeStatus, err)
	}

	oldBlock, err := nas.getOldestBlock(latestBlockData.Nonce)
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: int64(latestBlockData.Nonce),
			Hash:  latestBlockData.Hash,
		},
		CurrentBlockTimestamp:  latestBlockData.Timestamp,
		GenesisBlockIdentifier: nas.config.GenesisBlockIdentifier,
		OldestBlockIdentifier: &types.BlockIdentifier{
			Index: int64(oldBlock.Nonce),
			Hash:  oldBlock.Hash,
		},
		Peers: nas.config.Peers,
	}, nil
}

func (nas *networkAPIService) getOldestBlock(latestBlockNonce uint64) (*client.BlockData, error) {
	oldestBlockNonce := uint64(1)

	if latestBlockNonce > NumBlocksToGet {
		oldestBlockNonce = latestBlockNonce - NumBlocksToGet
	}

	block, err := nas.elrondClient.GetBlockByNonce(int64(oldestBlockNonce))
	if err != nil {
		return nil, err
	}

	return &client.BlockData{
		Nonce: block.Nonce,
		Hash:  block.Hash,
	}, nil

}

// NetworkOptions implements the /network/options endpoint.
func (nas *networkAPIService) NetworkOptions(
	_ context.Context,
	_ *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	networkConfig, err := nas.elrondClient.GetNetworkConfig()
	if err != nil {
		return nil, wrapErr(ErrUnableToGetClientVersion, err)
	}

	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: RosettaVersion,
			NodeVersion:    networkConfig.ClientVersion,
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
