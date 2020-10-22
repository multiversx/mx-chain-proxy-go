package mocks

import (
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider"
)

// ElrondProviderMock -
type ElrondProviderMock struct {
	GetNetworkConfigCalled             func() (*provider.NetworkConfig, error)
	GetLatestBlockDataCalled           func() (*provider.BlockData, error)
	GetBlockByNonceCalled              func(nonce int64) (*data.Hyperblock, error)
	GetBlockByHashCalled               func(hash string) (*data.Hyperblock, error)
	GetAccountCalled                   func(address string) (*data.Account, error)
	EncodeAddressCalled                func(address []byte) (string, error)
	SendTxCalled                       func(tx *data.Transaction) (string, error)
	ComputeTransactionHashCalled       func(tx *data.Transaction) (string, error)
	CalculateBlockTimestampUnixCalled  func(round uint64) int64
	GetTransactionByHashFromPoolCalled func(txHash string) (*data.FullTransaction, bool)
	DecodeAddressCalled                func(address string) ([]byte, error)
}

// GetNetworkConfig -
func (epm *ElrondProviderMock) GetNetworkConfig() (*provider.NetworkConfig, error) {
	if epm.GetNetworkConfigCalled != nil {
		return epm.GetNetworkConfigCalled()
	}
	return nil, nil
}

// GetLatestBlockData -
func (epm *ElrondProviderMock) GetLatestBlockData() (*provider.BlockData, error) {
	if epm.GetLatestBlockDataCalled != nil {
		return epm.GetLatestBlockDataCalled()
	}

	return nil, nil
}

// GetBlockByNonce -
func (epm *ElrondProviderMock) GetBlockByNonce(nonce int64) (*data.Hyperblock, error) {
	if epm.GetBlockByNonceCalled != nil {
		return epm.GetBlockByNonceCalled(nonce)
	}
	return nil, nil
}

// GetBlockByHash -
func (epm *ElrondProviderMock) GetBlockByHash(_ string) (*data.Hyperblock, error) {
	return nil, nil
}

// GetAccount -
func (epm *ElrondProviderMock) GetAccount(address string) (*data.Account, error) {
	if epm.GetAccountCalled != nil {
		return epm.GetAccountCalled(address)
	}
	return nil, nil
}

// EncodeAddress -
func (epm *ElrondProviderMock) EncodeAddress(pubkey []byte) (string, error) {
	if epm.EncodeAddressCalled != nil {
		return epm.EncodeAddressCalled(pubkey)
	}
	return "", nil
}

// SendTx -
func (epm *ElrondProviderMock) SendTx(tx *data.Transaction) (string, error) {
	if epm.SendTxCalled != nil {
		return epm.SendTxCalled(tx)
	}
	return "", nil
}

// ComputeTransactionHash -
func (epm *ElrondProviderMock) ComputeTransactionHash(tx *data.Transaction) (string, error) {
	if epm.ComputeTransactionHashCalled != nil {
		return epm.ComputeTransactionHashCalled(tx)
	}
	return "", nil
}

// CalculateBlockTimestampUnix -
func (epm *ElrondProviderMock) CalculateBlockTimestampUnix(_ uint64) int64 {
	return 0
}

// GetTransactionByHashFromPool -
func (epm *ElrondProviderMock) GetTransactionByHashFromPool(txHash string) (*data.FullTransaction, bool) {
	if epm.GetTransactionByHashFromPoolCalled != nil {
		return epm.GetTransactionByHashFromPoolCalled(txHash)
	}
	return nil, false
}

// DecodeAddress -
func (epm *ElrondProviderMock) DecodeAddress(address string) ([]byte, error) {
	if epm.DecodeAddressCalled != nil {
		return epm.DecodeAddressCalled(address)
	}
	return nil, nil
}
