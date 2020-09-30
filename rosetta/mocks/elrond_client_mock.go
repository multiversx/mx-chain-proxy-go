package mocks

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
)

type ElrondClientMock struct {
	GetNetworkConfigCalled            func() (*client.NetworkConfig, error)
	GetNetworkStatusCalled            func() (*client.NetworkStatus, error)
	GetLatestBlockDataCalled          func() (*client.BlockData, error)
	GetBlockByNonceCalled             func(nonce int64) (*data.Hyperblock, error)
	GetBlockByHashCalled              func(hash string) (*data.Hyperblock, error)
	GetAccountCalled                  func(address string) (*data.Account, error)
	EncodeAddressCalled               func(address []byte) (string, error)
	SendTxCalled                      func(tx *data.Transaction) (string, error)
	SimulateTxCalled                  func(tx *data.Transaction) (string, error)
	CalculateBlockTimestampUnixCalled func(round uint64) int64
}

func (ecm *ElrondClientMock) GetNetworkConfig() (*client.NetworkConfig, error) {
	if ecm.GetNetworkConfigCalled != nil {
		return ecm.GetNetworkConfigCalled()
	}
	return nil, nil
}
func (ecm *ElrondClientMock) GetNetworkStatus() (*client.NetworkStatus, error) {
	return nil, nil
}
func (ecm *ElrondClientMock) GetLatestBlockData() (*client.BlockData, error) {
	if ecm.GetLatestBlockDataCalled != nil {
		return ecm.GetLatestBlockDataCalled()
	}

	return nil, nil
}
func (ecm *ElrondClientMock) GetBlockByNonce(nonce int64) (*data.Hyperblock, error) {
	if ecm.GetBlockByNonceCalled != nil {
		return ecm.GetBlockByNonceCalled(nonce)
	}
	return nil, nil
}
func (ecm *ElrondClientMock) GetBlockByHash(hash string) (*data.Hyperblock, error) {
	return nil, nil
}
func (ecm *ElrondClientMock) GetAccount(address string) (*data.Account, error) {
	if ecm.GetAccountCalled != nil {
		return ecm.GetAccountCalled(address)
	}
	return nil, nil
}
func (ecm *ElrondClientMock) EncodeAddress(address []byte) (string, error) {
	return "", nil
}
func (ecm *ElrondClientMock) SendTx(tx *data.Transaction) (string, error) {
	return "", nil
}
func (ecm *ElrondClientMock) SimulateTx(tx *data.Transaction) (string, error) {
	return "", nil
}
func (ecm *ElrondClientMock) CalculateBlockTimestampUnix(round uint64) int64 {
	return 0
}
