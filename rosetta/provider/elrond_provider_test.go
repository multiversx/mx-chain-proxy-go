package provider

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/state/factory"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/provider/mock"
	"github.com/stretchr/testify/assert"
)

func TestInitializeElrondProvider(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	count := 0
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxy := &mock.ElrondProxyClientMock{}
	elrondProxy.GetNetworkConfigMetricsCalled = func() (*data.GenericAPIResponse, error) {
		if count == 2 {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		}
		count++
		return nil, localErr
	}

	elrondProxyClient, err := NewElrondProvider(elrondProxy)
	assert.Nil(t, err)
	assert.Equal(t, roundDuration, elrondProxyClient.roundDurationMilliseconds)
	assert.Equal(t, startTime, elrondProxyClient.genesisTime)
}

func TestNewElrondProvider_InvalidHandlerShouldErr(t *testing.T) {
	t.Parallel()

	elrondProvider, err := NewElrondProvider(nil)

	assert.Nil(t, elrondProvider)
	assert.Equal(t, ErrInvalidElrondProxyHandler, err)
}

func TestElrondProvider_GetLatestBlockData(t *testing.T) {
	t.Parallel()

	blockNonce := uint64(10)
	blockHash := "hash"
	preBlockHash := "prevBlockHash"
	round := uint64(11)
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetLatestFullySynchronizedHyperblockNonceCalled: func() (uint64, error) {
			return blockNonce, nil
		},
		GetBlockByNonceCalled: func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
			return &data.BlockApiResponse{
				Data: data.BlockApiResponsePayload{
					Block: data.Block{
						Nonce:         blockNonce,
						Round:         round,
						Hash:          blockHash,
						PrevBlockHash: preBlockHash,
					},
				},
			}, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	blockData, err := elrondProvider.GetLatestBlockData()
	assert.Nil(t, err)
	assert.Equal(t, &BlockData{
		Nonce:         blockNonce,
		Hash:          blockHash,
		PrevBlockHash: preBlockHash,
		Timestamp:     1044000,
	}, blockData)
}

func TestElrondProvider_GetBlockByNonce(t *testing.T) {
	t.Parallel()

	blockNonce := uint64(10)
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetHyperBlockByNonceCalled: func(nonce uint64) (*data.HyperblockApiResponse, error) {
			return &data.HyperblockApiResponse{
				Data: data.HyperblockApiResponsePayload{
					Hyperblock: data.Hyperblock{
						Nonce: blockNonce,
					},
				},
			}, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	hyperBlock, err := elrondProvider.GetBlockByNonce(int64(blockNonce))
	assert.Nil(t, err)
	assert.Equal(t, &data.Hyperblock{Nonce: blockNonce}, hyperBlock)
}

func TestElrondProvider_GetBlockByHash(t *testing.T) {
	t.Parallel()

	blockHash := "hash-hash"
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetHyperBlockByHashCalled: func(hash string) (*data.HyperblockApiResponse, error) {
			return &data.HyperblockApiResponse{
				Data: data.HyperblockApiResponsePayload{
					Hyperblock: data.Hyperblock{
						Hash: blockHash,
					},
				},
			}, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	hyperBlock, err := elrondProvider.GetBlockByHash(blockHash)
	assert.Nil(t, err)
	assert.Equal(t, &data.Hyperblock{Hash: blockHash}, hyperBlock)
}

func TestElrondProvider_GetAccount(t *testing.T) {
	t.Parallel()

	accountAddr := "addr-addr"
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetAccountCalled: func(address string) (*data.Account, error) {
			return &data.Account{
				Address: accountAddr,
			}, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	accountRet, err := elrondProvider.GetAccount(accountAddr)
	assert.Nil(t, err)
	assert.Equal(t, &data.Account{Address: accountAddr}, accountRet)
}

func TestElrondProvider_ComputeTransactionHash(t *testing.T) {
	t.Parallel()

	transactionHash := "hash-hash"
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		ComputeTransactionHashCalled: func(tx *data.Transaction) (string, error) {
			return transactionHash, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	hash, err := elrondProvider.ComputeTransactionHash(&data.Transaction{})
	assert.Nil(t, err)
	assert.Equal(t, transactionHash, hash)
}

func TestElrondProvider_EncodeAddress(t *testing.T) {
	t.Parallel()

	addrBytes, _ := hex.DecodeString("7c3f38ab6d2f961de7e5ad914cdbd0b6361b5ddb53d504b5297bfa4c901fc1d8")
	expectedAddr := "erd10sln32md97tpmel94kg5ek7skcmpkhwm202sfdff00ayeyqlc8vqpajkz5"
	pubKeyConverter, _ := factory.NewPubkeyConverter(config.PubkeyConfig{
		Length: 32,
		Type:   "bech32",
	})
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetAddressConverterCalled: func() (core.PubkeyConverter, error) {
			return pubKeyConverter, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	bech32Addr, err := elrondProvider.EncodeAddress(addrBytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedAddr, bech32Addr)
}

func TestElrondProvider_SendTx(t *testing.T) {
	t.Parallel()

	transactionHash := "hash-hash"
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
			return 0, transactionHash, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	hash, err := elrondProvider.SendTx(&data.Transaction{})
	assert.Nil(t, err)
	assert.Equal(t, transactionHash, hash)
}

func TestElrondProvider_GetTransactionByHashFromPool_TxNotInPool(t *testing.T) {
	t.Parallel()

	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetTransactionByHashAndSenderAddressCalled: func(hash string, sndAddr string) (*data.FullTransaction, int, error) {
			return &data.FullTransaction{
				Status: transaction.TxStatusSuccess,
			}, 0, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	tx, isInPool := elrondProvider.GetTransactionByHashFromPool("hash")
	assert.Nil(t, tx)
	assert.False(t, isInPool)
}

func TestElrondProvider_GetTransactionByHashFromPool_TxInPool(t *testing.T) {
	t.Parallel()

	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxyMock := &mock.ElrondProxyClientMock{
		GetNetworkConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		},
		GetTransactionByHashAndSenderAddressCalled: func(hash string, sndAddr string) (*data.FullTransaction, int, error) {
			return &data.FullTransaction{
				Status: transaction.TxStatusPending,
			}, 0, nil
		},
	}

	elrondProvider, _ := NewElrondProvider(elrondProxyMock)

	tx, isInPool := elrondProvider.GetTransactionByHashFromPool("hash")
	assert.Equal(t, &data.FullTransaction{Status: transaction.TxStatusPending}, tx)
	assert.True(t, isInPool)
}
