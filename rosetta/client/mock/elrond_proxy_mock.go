package mock

import (
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type ElrondProxyClientMock struct {
	GetNetworkConfigMetricsCalled func() (*data.GenericAPIResponse, error)
	GetNetworkStatusMetricsCalled func(shardID uint32) (*data.GenericAPIResponse, error)
	GetBlockByNonceCalled         func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetAccountCalled              func(address string) (*data.Account, error)
	GetHyperBlockByNonceCalled    func(nonce uint64) (*data.HyperblockApiResponse, error)
	GetHyperBlockByHashCalled     func(hash string) (*data.HyperblockApiResponse, error)
	SendTransactionCalled         func(tx *data.Transaction) (int, string, error)
	SimulateTransactionCalled     func(tx *data.Transaction) (*data.ResponseTransactionSimulation, error)
	GetAddressConverterCalled     func() (core.PubkeyConverter, error)
}

func (epcm *ElrondProxyClientMock) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if epcm.GetNetworkConfigMetricsCalled != nil {
		return epcm.GetNetworkConfigMetricsCalled()
	}
	return nil, nil
}
func (epcm *ElrondProxyClientMock) GetNetworkStatusMetrics(_ uint32) (*data.GenericAPIResponse, error) {
	return nil, nil
}
func (epcm *ElrondProxyClientMock) GetBlockByNonce(_ uint32, _ uint64, _ bool) (*data.BlockApiResponse, error) {
	return nil, nil
}
func (epcm *ElrondProxyClientMock) GetAccount(_ string) (*data.Account, error) {
	return nil, nil
}

func (epcm *ElrondProxyClientMock) GetHyperBlockByNonce(_ uint64) (*data.HyperblockApiResponse, error) {
	return nil, nil
}
func (epcm *ElrondProxyClientMock) GetHyperBlockByHash(_ string) (*data.HyperblockApiResponse, error) {
	return nil, nil
}

func (epcm *ElrondProxyClientMock) SendTransaction(_ *data.Transaction) (int, string, error) {
	return 0, "", nil
}
func (epcm *ElrondProxyClientMock) SimulateTransaction(_ *data.Transaction) (*data.ResponseTransactionSimulation, error) {
	return nil, nil
}

func (epcm *ElrondProxyClientMock) GetAddressConverter() (core.PubkeyConverter, error) {
	return nil, nil
}

func (epcm *ElrondProxyClientMock) GetLatestBlockNonce() (uint64, error) {
	return 0, nil
}
