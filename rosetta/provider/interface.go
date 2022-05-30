package provider

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ElrondProxyClient defines what a real elrond proxy client should do
type ElrondProxyClient interface {
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetBlockByNonce(shardID uint32, nonce uint64, options common.BlockQueryOptions) (*data.BlockApiResponse, error)
	GetAccount(address string) (*data.Account, error)

	GetHyperBlockByNonce(nonce uint64, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHash(hash string, options common.HyperblockQueryOptions) (*data.HyperblockApiResponse, error)

	SendTransaction(tx *data.Transaction) (int, string, error)
	ComputeTransactionHash(tx *data.Transaction) (string, error)
	GetTransactionByHashAndSenderAddress(txHash string, sndAddr string, withResults bool) (*data.FullTransaction, int, error)

	GetLatestFullySynchronizedHyperblockNonce() (uint64, error)
	GetAddressConverter() (core.PubkeyConverter, error)
}

// ElrondProviderHandler defines what a real elrond provider should do
type ElrondProviderHandler interface {
	GetNetworkConfig() (*NetworkConfig, error)
	GetLatestBlockData() (*BlockData, error)
	GetBlockByNonce(nonce int64) (*data.Hyperblock, error)
	GetBlockByHash(hash string) (*data.Hyperblock, error)
	GetAccount(address string) (*data.Account, error)
	EncodeAddress(address []byte) (string, error)
	DecodeAddress(address string) ([]byte, error)
	SendTx(tx *data.Transaction) (string, error)
	CalculateBlockTimestampUnix(round uint64) int64
	ComputeTransactionHash(tx *data.Transaction) (string, error)
	GetTransactionByHashFromPool(txHash string) (*data.FullTransaction, bool)
}
