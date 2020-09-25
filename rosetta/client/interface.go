package client

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type ElrondProxyClient interface {
	GetNetworkConfigMetrics() (*data.GenericAPIResponse, error)
	GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error)
	GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetAccount(address string) (*data.Account, error)

	GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error)

	SendTransaction(tx *data.Transaction) (int, string, error)
	SimulateTransaction(tx *data.Transaction) (*data.ResponseTransactionSimulation, error)

	GetAddressConverter() (core.PubkeyConverter, error)
}
