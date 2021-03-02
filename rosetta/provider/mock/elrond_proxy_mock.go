package mock

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// ElrondProxyClientMock -
type ElrondProxyClientMock struct {
	GetNetworkConfigMetricsCalled                   func() (*data.GenericAPIResponse, error)
	GetBlockByNonceCalled                           func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetAccountCalled                                func(address string) (*data.Account, error)
	GetHyperBlockByNonceCalled                      func(nonce uint64) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHashCalled                       func(hash string) (*data.HyperblockApiResponse, error)
	SendTransactionCalled                           func(tx *data.Transaction) (int, string, error)
	SimulateTransactionCalled                       func(tx *data.Transaction) (*data.ResponseTransactionSimulation, error)
	GetAddressConverterCalled                       func() (core.PubkeyConverter, error)
	GetLatestFullySynchronizedHyperblockNonceCalled func() (uint64, error)
	ComputeTransactionHashCalled                    func(tx *data.Transaction) (string, error)
	GetTransactionByHashAndSenderAddressCalled      func(hash string, sndAddr string) (*data.FullTransaction, int, error)
}

// GetNetworkConfigMetrics -
func (epcm *ElrondProxyClientMock) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if epcm.GetNetworkConfigMetricsCalled != nil {
		return epcm.GetNetworkConfigMetricsCalled()
	}
	return nil, nil
}

// GetNetworkStatusMetrics -
func (epcm *ElrondProxyClientMock) GetNetworkStatusMetrics(_ uint32) (*data.GenericAPIResponse, error) {
	return nil, nil
}

// GetBlockByNonce -
func (epcm *ElrondProxyClientMock) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	if epcm.GetBlockByNonceCalled != nil {
		return epcm.GetBlockByNonceCalled(shardID, nonce, withTxs)
	}
	return nil, nil
}

// GetAccount -
func (epcm *ElrondProxyClientMock) GetAccount(address string) (*data.Account, error) {
	if epcm.GetAccountCalled != nil {
		return epcm.GetAccountCalled(address)
	}
	return nil, nil
}

// GetHyperBlockByNonce -
func (epcm *ElrondProxyClientMock) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	if epcm.GetHyperBlockByNonceCalled != nil {
		return epcm.GetHyperBlockByNonceCalled(nonce)
	}
	return nil, nil
}

// GetHyperBlockByHash -
func (epcm *ElrondProxyClientMock) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	if epcm.GetHyperBlockByHashCalled != nil {
		return epcm.GetHyperBlockByHashCalled(hash)
	}
	return nil, nil
}

// SendTransaction -
func (epcm *ElrondProxyClientMock) SendTransaction(tx *data.Transaction) (int, string, error) {
	if epcm.SendTransactionCalled != nil {
		return epcm.SendTransactionCalled(tx)
	}
	return 0, "", nil
}

// ComputeTransactionHash -
func (epcm *ElrondProxyClientMock) ComputeTransactionHash(hash *data.Transaction) (string, error) {
	if epcm.ComputeTransactionHashCalled != nil {
		return epcm.ComputeTransactionHashCalled(hash)
	}
	return "", nil
}

// GetAddressConverter -
func (epcm *ElrondProxyClientMock) GetAddressConverter() (core.PubkeyConverter, error) {
	if epcm.GetAddressConverterCalled != nil {
		return epcm.GetAddressConverterCalled()
	}
	return nil, nil
}

// GetLatestBlockNonce -
func (epcm *ElrondProxyClientMock) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	if epcm.GetLatestFullySynchronizedHyperblockNonceCalled != nil {
		return epcm.GetLatestFullySynchronizedHyperblockNonceCalled()
	}
	return 0, nil
}

// GetTransactionByHashAndSenderAddress -
func (epcm *ElrondProxyClientMock) GetTransactionByHashAndSenderAddress(
	hash string,
	sndAddr string,
	_ bool,
) (*data.FullTransaction, int, error) {
	if epcm.GetTransactionByHashAndSenderAddressCalled != nil {
		return epcm.GetTransactionByHashAndSenderAddressCalled(hash, sndAddr)
	}
	return nil, 0, nil
}
