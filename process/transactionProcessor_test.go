package process_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	hasherFactory "github.com/ElrondNetwork/elrond-go/hashing/factory"
	marshalFactory "github.com/ElrondNetwork/elrond-go/marshal/factory"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hasher, _ = hasherFactory.NewHasher("blake2b")
var marshalizer, _ = marshalFactory.NewMarshalizer("gogo protobuf")

func TestNewTransactionProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(nil, &mock.PubKeyConverterMock{}, hasher, marshalizer)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewTransactionProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, nil, hasher, marshalizer)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewTransactionProcessor_NilHasherShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, nil, marshalizer)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilHasher, err)
}

func TestNewTransactionProcessor_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, nil)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewTransactionProcessor_OkValuesShouldWork(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer)

	require.NotNil(t, tp)
	require.Nil(t, err)
}

//------- SendTransaction

func TestTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer)
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

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer)
	rc, txHash, err := tp.SendTransaction(&data.Transaction{})

	require.Empty(t, txHash)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "no chainID")
	require.Equal(t, http.StatusBadRequest, rc)
}

func TestTransactionProcessor_SendTransactionNoVersionShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer)
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
	)

	response, err := tp.SimulateTransaction(txsToSimulate, true)
	require.Nil(t, err)

	respData := response.Data.(data.TransactionSimulationResponseData)
	require.Equal(t, expectedFailReason, respData.Result.FailReason)
}

func TestTransactionProcessor_SimulateTransactionCrossShardOkOnSenderFailOnReceiverShouldWork(t *testing.T) {
	t.Parallel()

	expectedStatusSh0, expectedStatusSh1 := "ok", "not ok"
	txAddressSh0 := []byte("addr in shard 0")
	txAddressSh1 := []byte("addr in shard 1")
	expectedFailReason := "fail reason"
	txsToSimulate := &data.Transaction{Receiver: hex.EncodeToString(txAddressSh1), Sender: hex.EncodeToString(txAddressSh0), ChainID: "chain", Version: 1}

	obsSh0 := "observer shard 0"
	obsSh1 := "observer shard 1"
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				if bytes.Equal(addressBuff, txAddressSh0) {
					return 0, nil
				}
				return 1, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				if shardId == 0 {
					return []*data.NodeData{{Address: obsSh0, ShardId: 0}}, nil
				}
				return []*data.NodeData{{Address: obsSh1, ShardId: 1}}, nil
			},
			CallPostRestEndPointCalled: func(address string, path string, value interface{}, response interface{}) (int, error) {
				if address == obsSh0 {
					resp := response.(*data.ResponseTransactionSimulation)
					resp.Data.Result.Status = transaction.TxStatus(expectedStatusSh0)
					response = resp
					return http.StatusOK, nil
				}

				resp := response.(*data.ResponseTransactionSimulation)
				resp.Data.Result.FailReason = expectedFailReason
				resp.Data.Result.Status = transaction.TxStatus(expectedStatusSh1)
				response = resp
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
	)

	response, err := tp.SimulateTransaction(txsToSimulate, true)
	require.Nil(t, err)

	respData := response.Data.(data.TransactionSimulationResponseDataCrossShard)
	require.Equal(t, expectedStatusSh0, string(respData.Result["senderShard"].Status))
	require.Equal(t, expectedStatusSh1, string(respData.Result["receiverShard"].Status))
	require.Equal(t, expectedFailReason, respData.Result["receiverShard"].FailReason)
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
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0, 1}
			},
			GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
					}, nil
				}
				if shardId == 1 {
					return []*data.NodeData{
						{Address: addrObs1, ShardId: 1},
					}, nil
				}
				return nil, nil
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
		hasher,
		marshalizer,
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), "")
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}

func TestTransactionProcessor_GetTransactionStatusCrossShardTransaction(t *testing.T) {
	t.Parallel()

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))

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
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0}
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
		hasher,
		marshalizer,
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
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0, 1}
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, err error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
					}, nil
				}
				if shardId == 1 {
					return []*data.NodeData{
						{Address: addrObs1, ShardId: 1},
					}, nil
				}
				return nil, nil
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
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
		hasher,
		marshalizer,
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), sndrShard0)
	assert.NoError(t, err)
	assert.Equal(t, txResponseStatus, txStatus)
}

func TestTransactionProcessor_ComputeTransactionInvalidTransactionValue(t *testing.T) {
	t.Parallel()

	tx := &data.Transaction{
		Nonce:     1,
		Value:     "aaaa",
		Receiver:  "61616161",
		Sender:    "62626262",
		GasPrice:  1,
		GasLimit:  2,
		Data:      []byte("blablabla"),
		Signature: "abcdabcd",
		ChainID:   "1",
		Version:   1,
	}
	marshalizer := marshalizer
	hasher := hasher
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer)

	_, err := tp.ComputeTransactionHash(tx)
	assert.Equal(t, process.ErrInvalidTransactionValueField, err)
}

func TestTransactionProcessor_ComputeTransactionInvalidReceiverAddress(t *testing.T) {
	t.Parallel()

	tx := &data.Transaction{
		Nonce:     1,
		Value:     "1",
		Receiver:  "gfdgfd",
		Sender:    "62626262",
		GasPrice:  1,
		GasLimit:  2,
		Data:      []byte("blablabla"),
		Signature: "abcdabcd",
		ChainID:   "1",
		Version:   1,
	}
	marshalizer := marshalizer
	hasher := hasher
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer)

	_, err := tp.ComputeTransactionHash(tx)
	assert.Equal(t, process.ErrInvalidAddress, err)
}

func TestTransactionProcessor_ComputeTransactionInvalidSenderAddress(t *testing.T) {
	t.Parallel()

	tx := &data.Transaction{
		Nonce:     1,
		Value:     "1",
		Receiver:  "62626262",
		Sender:    "gagasd",
		GasPrice:  1,
		GasLimit:  2,
		Data:      []byte("blablabla"),
		Signature: "abcdabcd",
		ChainID:   "1",
		Version:   1,
	}
	marshalizer := marshalizer
	hasher := hasher
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer)

	_, err := tp.ComputeTransactionHash(tx)
	assert.Equal(t, process.ErrInvalidAddress, err)
}

func TestTransactionProcessor_ComputeTransactionInvalidSignaturesBytes(t *testing.T) {
	t.Parallel()

	tx := &data.Transaction{
		Nonce:     1,
		Value:     "1",
		Receiver:  "62626262",
		Sender:    "62626262",
		GasPrice:  1,
		GasLimit:  2,
		Data:      []byte("blablabla"),
		Signature: "gfgdgfdgfd",
		ChainID:   "1",
		Version:   1,
	}
	marshalizer := marshalizer
	hasher := hasher
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer)

	_, err := tp.ComputeTransactionHash(tx)
	assert.Equal(t, process.ErrInvalidSignatureBytes, err)
}

func TestTransactionProcessor_ComputeTransactionShouldWork1(t *testing.T) {
	t.Parallel()

	tx := &data.Transaction{
		Nonce:     1,
		Value:     "1",
		Receiver:  "61616161",
		Sender:    "62626262",
		GasPrice:  1,
		GasLimit:  2,
		Data:      []byte("blablabla"),
		Signature: "abcdabcd",
		ChainID:   "1",
		Version:   1,
	}
	marshalizer := marshalizer
	hasher := hasher
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer)

	txHashHex := "891694ae6307ee9f17f861816187a6729268397f8fabc055d5b334f552cd3cfb"
	txHash, err := tp.ComputeTransactionHash(tx)
	assert.Nil(t, err)
	assert.Equal(t, txHashHex, txHash)
}

func TestTransactionProcessor_ComputeTransactionShouldWork2(t *testing.T) {
	t.Parallel()

	protoTx := transaction.Transaction{
		Nonce:     1,
		Value:     big.NewInt(1000),
		RcvAddr:   []byte("7c3f38ab6d2f961de7e5ad914cdbd0b6361b5ddb53d504b5297bfa4c901fc1d8"),
		SndAddr:   []byte("7c3f38ab6d2f961de7e5ad914cdbd0b6361b5ddb53d504b5297bfa4c901fc1d8"),
		GasPrice:  12,
		GasLimit:  13,
		Data:      []byte("aGVsbG8="),
		ChainID:   []byte("1"),
		Version:   1,
		Signature: []byte("5e97b3bb223acfe3a152bb8e7fec31909059c90f75b56ffc4edf1695baab561b"),
	}
	protoTxHashBytes, _ := core.CalculateHash(marshalizer, hasher, &protoTx)
	protoTxHash := hex.EncodeToString(protoTxHashBytes)

	marshalizer := marshalizer
	hasher := hasher
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer)

	txHash, err := tp.ComputeTransactionHash(&data.Transaction{
		Nonce:     protoTx.Nonce,
		Value:     protoTx.Value.String(),
		Receiver:  pubKeyConv.Encode(protoTx.RcvAddr),
		Sender:    pubKeyConv.Encode(protoTx.SndAddr),
		GasPrice:  protoTx.GasPrice,
		GasLimit:  protoTx.GasLimit,
		Data:      protoTx.Data,
		Signature: hex.EncodeToString(protoTx.Signature),
		ChainID:   string(protoTx.ChainID),
		Version:   protoTx.Version,
	})
	assert.Nil(t, err)
	assert.Equal(t, protoTxHash, txHash)
}

func TestTransactionProcessor_GetTransactionShouldWork(t *testing.T) {
	t.Parallel()

	expectedNonce := uint64(37)

	sndrShard0 := hex.EncodeToString([]byte("bbbbbb"))
	sndrShard1 := hex.EncodeToString([]byte("cccccc"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"

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
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0, 1}
			},
			GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
					}, nil
				}
				if shardId == 1 {
					return []*data.NodeData{
						{Address: addrObs1, ShardId: 1},
					}, nil
				}
				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs0 {
					responseGetTx := value.(*data.GetTransactionResponse)

					responseGetTx.Data.Transaction = data.FullTransaction{
						Nonce: expectedNonce,
					}
					return http.StatusOK, nil
				}

				return http.StatusBadGateway, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
	)

	tx, err := tp.GetTransaction(string(hash0), false)
	assert.NoError(t, err)
	assert.Equal(t, expectedNonce, tx.Nonce)
}

func TestTransactionProcessor_GetTransactionShouldCallOtherObserverInShardIfHttpError(t *testing.T) {
	t.Parallel()

	addrObs0 := "observer0"
	addrObs1 := "observer1"
	secondObserverWasCalled := false

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(_ []byte) (uint32, error) {
				return 0, nil
			},
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0}
			},
			GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
						{Address: addrObs1, ShardId: 0},
					}, nil
				}
				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs0 {
					return 0, errors.New("rest api error")
				}
				if address == addrObs1 {
					secondObserverWasCalled = true
					return http.StatusOK, nil
				}

				return http.StatusBadGateway, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
	)

	_, _ = tp.GetTransaction(string(hash0), false)
	assert.True(t, secondObserverWasCalled)
}

func TestTransactionProcessor_GetTransactionShouldNotCallOtherObserverInShardIfNoHttpErrorButTxNotFound(t *testing.T) {
	t.Parallel()

	addrObs0 := "observer0"
	addrObs1 := "observer1"

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(_ []byte) (uint32, error) {
				return 0, nil
			},
			GetObserversOnePerShardCalled: func() ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
				}, nil
			},
			GetObserversCalled: func(shardId uint32) ([]*data.NodeData, error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
						{Address: addrObs1, ShardId: 0},
					}, nil
				}
				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs1 {
					require.Fail(t, "second observer should have not been called")
				}

				return http.StatusInternalServerError, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
	)

	_, _ = tp.GetTransaction(string(hash0), false)
}

func TestTransactionProcessor_GetTransactionWithEventsFirstFromDstShardAndAfterSource(t *testing.T) {
	t.Parallel()

	expectedNonce := uint64(37)

	sndrShard0 := hex.EncodeToString([]byte("aaaa"))
	rcvShard1 := hex.EncodeToString([]byte("bbbb"))

	addrObs0 := "observer0"
	addrObs1 := "observer1"

	scHash1 := "scHash1"
	scHash2 := "scHash2"
	scHash3 := "scHash3"

	scRes1 := &transaction.ApiSmartContractResult{
		Hash: scHash1,
	}
	scRes2 := &transaction.ApiSmartContractResult{
		Hash: scHash2,
	}
	scRes3 := &transaction.ApiSmartContractResult{
		Hash: scHash3,
	}

	hash0 := []byte("hash0")
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				if string(addressBuff) == "aaaa" {
					return uint32(0), nil
				}
				if string(addressBuff) == "bbbb" {
					return uint32(1), nil
				}
				return 0, nil
			},
			GetShardIDsCalled: func() []uint32 {
				return []uint32{1, 0}
			},
			GetFullHistoryNodesCalled: func(shardId uint32) ([]*data.NodeData, error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
					}, nil
				}
				if shardId == 1 {
					return []*data.NodeData{
						{Address: addrObs1, ShardId: 1},
					}, nil
				}

				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs1 {
					responseGetTx := value.(*data.GetTransactionResponse)

					responseGetTx.Data.Transaction = data.FullTransaction{
						Sender:           sndrShard0,
						Receiver:         rcvShard1,
						Nonce:            expectedNonce,
						SourceShard:      0,
						DestinationShard: 1,
						ScResults: []*transaction.ApiSmartContractResult{
							scRes1, scRes2,
						},
					}
					return http.StatusOK, nil
				} else if address == addrObs0 {
					responseGetTx := value.(*data.GetTransactionResponse)

					responseGetTx.Data.Transaction = data.FullTransaction{
						Nonce:            expectedNonce,
						SourceShard:      0,
						DestinationShard: 1,
						ScResults: []*transaction.ApiSmartContractResult{
							scRes2, scRes3,
						},
					}
					return http.StatusOK, nil
				}

				return http.StatusBadGateway, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
	)

	tx, err := tp.GetTransaction(string(hash0), true)
	assert.NoError(t, err)
	assert.Equal(t, expectedNonce, tx.Nonce)
	assert.Equal(t, 3, len(tx.ScResults))
}
