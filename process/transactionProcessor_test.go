package process_test

import (
	"encoding/hex"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransactionProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(nil, &mock.PubKeyConverterMock{})

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewTransactionProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, nil)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewTransactionProcessor_OkValuesShouldWork(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})

	require.NotNil(t, tp)
	require.Nil(t, err)
}

//------- SendTransaction

func TestTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	rc, txHash, err := tp.SendTransaction(&data.Transaction{
		Sender: "invalid hex number",
	})

	require.Empty(t, txHash)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid byte")
	require.Equal(t, http.StatusBadRequest, rc)
}

func TestTransactionProcessor_SendTransactionNoChainIDShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	rc, txHash, err := tp.SendTransaction(&data.Transaction{})

	require.Empty(t, txHash)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "no chainID")
	require.Equal(t, http.StatusBadRequest, rc)
}

func TestTransactionProcessor_SendTransactionNoVersionShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	rc, txHash, err := tp.SendTransaction(&data.Transaction{
		ChainID: "chainID",
	})

	require.Empty(t, txHash)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "no version")
	require.Equal(t, http.StatusBadRequest, rc)
}

func TestTransactionProcessor_SendTransactionComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	rc, txHash, err := tp.SendTransaction(&data.Transaction{
		ChainID: "chain",
		Version: 1,
	})

	require.Empty(t, txHash)
	require.Equal(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rc)
}

func TestTransactionProcessor_SendTransactionGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return nil, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	rc, txHash, err := tp.SendTransaction(&data.Transaction{
		Sender:  address,
		ChainID: "chain",
		Version: 1,
	})

	require.Empty(t, txHash)
	require.Equal(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rc)
}

func TestTransactionProcessor_SendTransactionSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address1", ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, data interface{}, response interface{}) (int, error) {
				return http.StatusInternalServerError, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	rc, txHash, err := tp.SendTransaction(&data.Transaction{
		Sender:  address,
		ChainID: "chain",
		Version: 1,
	})

	require.Empty(t, txHash)
	require.Equal(t, errExpected, err)
	require.Equal(t, http.StatusInternalServerError, rc)
}

func TestTransactionProcessor_SendTransactionSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	txHash := "DEADBEEF01234567890"
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: addressFail, ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				txResponse := response.(*data.ResponseTransaction)
				txResponse.Data.TxHash = txHash
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	rc, resultedTxHash, err := tp.SendTransaction(&data.Transaction{
		Sender:  address,
		ChainID: "chain",
		Version: 1,
	})

	require.Equal(t, resultedTxHash, txHash)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, rc)
}

////------- SendMultipleTransactions

func TestTransactionProcessor_SendMultipleTransactionsShouldWork(t *testing.T) {
	t.Parallel()

	var txsToSend []*data.Transaction
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "aaaaaa", Sender: hex.EncodeToString([]byte("cccccc")), ChainID: "chain", Version: 1})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "bbbbbb", Sender: hex.EncodeToString([]byte("dddddd")), ChainID: "chain", Version: 1})

	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "observer1", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				receivedTxs, ok := value.([]*data.Transaction)
				require.True(t, ok)
				resp := response.(*data.ResponseMultipleTransactions)
				resp.Data.NumOfTxs = uint64(len(receivedTxs))
				resp.Data.TxsHashes = map[int]string{
					0: "hash1",
					1: "hash2",
				}
				response = resp
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	response, err := tp.SendMultipleTransactions(txsToSend)
	require.Nil(t, err)
	require.Equal(t, len(response.TxsHashes), len(txsToSend))
	require.Equal(t, uint64(len(txsToSend)), response.NumOfTxs)
}

func TestTransactionProcessor_SendMultipleTransactionsShouldWorkAndSendTxsByShard(t *testing.T) {
	t.Parallel()

	var txsToSend []*data.Transaction
	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "aaaaaa", Sender: sndrShard0, ChainID: "chain", Version: 1})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "aaaaaa", Sender: sndrShard0, ChainID: "chain", Version: 1})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "aaaaaa", Sender: sndrShard1, ChainID: "chain", Version: 1})
	txsToSend = append(txsToSend, &data.Transaction{Receiver: "aaaaaa", Sender: sndrShard1, ChainID: "chain", Version: 1})
	numOfTimesPostEndpointWasCalled := uint32(0)

	addrObs0 := "observer0"
	addrObs1 := "observer1"

	hash0, hash1, hash2, hash3 := "hash0", "hash1", "hash2", "hash3"

	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				sndrHex := hex.EncodeToString(addressBuff)
				if sndrHex == sndrShard0 {
					return uint32(0), nil
				}
				if sndrHex == sndrShard1 {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetObserversCalled: func(shardID uint32) (observers []*data.NodeData, e error) {
				if shardID == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
					}, nil
				}
				return []*data.NodeData{
					{Address: addrObs1, ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				atomic.AddUint32(&numOfTimesPostEndpointWasCalled, 1)
				resp := response.(*data.ResponseMultipleTransactions)
				resp.Data.NumOfTxs = uint64(2)
				if address == addrObs0 {
					resp.Data.TxsHashes = map[int]string{
						0: hash0,
						1: hash1,
					}
				} else {
					resp.Data.TxsHashes = map[int]string{
						0: hash2,
						1: hash3,
					}
				}

				response = resp
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	response, err := tp.SendMultipleTransactions(txsToSend)
	require.Nil(t, err)
	require.Equal(t, uint64(len(txsToSend)), response.NumOfTxs)
	require.Equal(t, uint32(2), atomic.LoadUint32(&numOfTimesPostEndpointWasCalled))

	require.Equal(t, len(txsToSend), len(response.TxsHashes))
	require.Equal(
		t,
		map[int]string{0: hash0, 1: hash1, 2: hash2, 3: hash3},
		response.TxsHashes,
	)
}

func TestTransactionProcessor_SimulateTransactionShouldWork(t *testing.T) {
	t.Parallel()

	expectedFailReason := "fail reason"
	txsToSimulate := &data.Transaction{Receiver: "aaaaaa", Sender: hex.EncodeToString([]byte("cccccc")), ChainID: "chain", Version: 1}

	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "observer1", ShardId: 0},
				}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				resp := response.(*data.ResponseTransactionSimulation)
				resp.Data.Result.FailReason = expectedFailReason
				response = resp
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	response, err := tp.SimulateTransaction(txsToSimulate)
	require.Nil(t, err)
	require.Equal(t, expectedFailReason, response.Data.Result.FailReason)
}

func TestTransactionProcessor_GetTransactionStatusIntraShardTransaction(t *testing.T) {
	t.Parallel()

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"

	txResponseStatus := "executed"

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				sndrHex := hex.EncodeToString(addressBuff)
				if sndrHex == sndrShard0 {
					return uint32(0), nil
				}
				if sndrHex == sndrShard1 {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetAllObserversCalled: func() ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
					{Address: addrObs1, ShardId: 1},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs0 {
					responseGetTx := value.(*data.GetTransactionResponse)

					responseGetTx.Data.Transaction = data.FullTransaction{
						Status: transaction.TxStatus(txResponseStatus),
					}
					return http.StatusOK, nil
				}

				return http.StatusBadGateway, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), "")
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}

func TestTransactionProcessor_GetTransactionStatusCrossShardTransaction(t *testing.T) {
	t.Parallel()

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"

	txResponseStatus := "executed"

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				sndrHex := hex.EncodeToString(addressBuff)
				if sndrHex == sndrShard0 {
					return uint32(0), nil
				}
				if sndrHex == sndrShard1 {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetAllObserversCalled: func() ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
				}, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
				return []*data.NodeData{
					{Address: addrObs1, ShardId: 1},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				responseGetTx := value.(*data.GetTransactionResponse)

				responseGetTx.Data.Transaction = data.FullTransaction{
					Receiver: sndrShard1,
					Sender:   sndrShard0,
					Status:   transaction.TxStatus(txResponseStatus),
				}
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), "")
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}

func TestTransactionProcessor_GetTransactionStatusCrossShardTransactionDestinationNotAnswer(t *testing.T) {
	t.Parallel()

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"

	txResponseStatus := "partially-executed"

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				sndrHex := hex.EncodeToString(addressBuff)
				if sndrHex == sndrShard0 {
					return uint32(0), nil
				}
				if sndrHex == sndrShard1 {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetAllObserversCalled: func() ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
				}, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
				return []*data.NodeData{
					{Address: addrObs1, ShardId: 1},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if addrObs1 == address {
					return http.StatusBadRequest, nil
				}

				responseGetTx := value.(*data.GetTransactionResponse)

				responseGetTx.Data.Transaction = data.FullTransaction{
					Receiver: sndrShard1,
					Sender:   sndrShard0,
					Status:   transaction.TxStatus(txResponseStatus),
				}
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), "")
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}

func TestTransactionProcessor_GetTransactionStatusWithSenderAddressCrossShard(t *testing.T) {
	t.Parallel()

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	rcvShard1 := hex.EncodeToString([]byte("cccccc"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"
	addrObs2 := "observer2"
	addrObs3 := "observer3"

	txResponseStatus := "executed"

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				sndrHex := hex.EncodeToString(addressBuff)
				if sndrHex == sndrShard0 {
					return uint32(0), nil
				}
				if sndrHex == rcvShard1 {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetAllObserversCalled: func() ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
				}, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
				return []*data.NodeData{
					{Address: addrObs1, ShardId: 1},
					{Address: addrObs2, ShardId: 1},
					{Address: addrObs3, ShardId: 1},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if addrObs1 == address {
					return 0, errors.New("local error")
				}
				if addrObs2 == address {
					return http.StatusBadRequest, nil
				}

				responseGetTx := value.(*data.GetTransactionResponse)

				responseGetTx.Data.Transaction = data.FullTransaction{
					Receiver: rcvShard1,
					Sender:   sndrShard0,
					Status:   transaction.TxStatus(txResponseStatus),
				}
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), sndrShard0)
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}

func TestTransactionProcessor_GetTransactionStatusWithSenderInvaidSender(t *testing.T) {
	t.Parallel()

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return 0, errors.New("local error")
			},
		},
		&mock.PubKeyConverterMock{},
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), "blablabla")
	assert.Error(t, err)
	assert.Equal(t, process.UnknownStatusTx, txStatus)
}

func TestTransactionProcessor_GetTransactionStatusWithSenderAddressIntraShard(t *testing.T) {
	t.Parallel()

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	rcvShard0 := hex.EncodeToString([]byte("cccccc"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"
	addrObs2 := "observer2"

	txResponseStatus := "executed"

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
					{Address: addrObs1, ShardId: 0},
					{Address: addrObs2, ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs0 {
					return http.StatusBadRequest, nil
				}
				if address == addrObs1 {
					return 0, errors.New("local error")
				}

				responseGetTx := value.(*data.GetTransactionResponse)

				responseGetTx.Data.Transaction = data.FullTransaction{
					Receiver: rcvShard0,
					Sender:   sndrShard0,
					Status:   transaction.TxStatus(txResponseStatus),
				}
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), sndrShard0)
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}
