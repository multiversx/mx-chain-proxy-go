package process_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	hasherFactory "github.com/multiversx/mx-chain-core-go/hashing/factory"
	"github.com/multiversx/mx-chain-core-go/marshal"
	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	apiErrors "github.com/multiversx/mx-chain-proxy-go/api/errors"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/logsevents"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var hasher, _ = hasherFactory.NewHasher("blake2b")
var marshalizer, _ = marshalFactory.NewMarshalizer("gogo protobuf")
var funcNewTxCostHandler = func() (process.TransactionCostHandler, error) {
	return &mock.TransactionCostHandlerStub{}, nil
}

var logsMerger, _ = logsevents.NewLogsMerger(hasher, &marshal.JsonMarshalizer{})
var testPubkeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32, "erd")
var testLogger = logger.GetOrCreate("process_test")

type scenarioData struct {
	Transaction *transaction.ApiTransactionResult   `json:"transaction"`
	SCRs        []*transaction.ApiTransactionResult `json:"scrs"`
}

func loadJsonIntoTxAndScrs(tb testing.TB, path string) *scenarioData {
	scenarioDataInstance := &scenarioData{}
	buff, err := os.ReadFile(path)
	require.Nil(tb, err)

	err = json.Unmarshal(buff, scenarioDataInstance)
	require.Nil(tb, err)

	return scenarioDataInstance
}

func createTestProcessorFromScenarioData(testData *scenarioData) *process.TransactionProcessor {
	processorStub := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0} // force everything intra-shard for test setup simplicity
		},
		ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					Address: "test",
					ShardId: 0,
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			for _, scr := range testData.SCRs {
				if strings.Contains(path, scr.Hash) {
					response := value.(*data.GetTransactionResponse)
					response.Data.Transaction = *scr
					return http.StatusOK, nil
				}
			}

			return http.StatusInternalServerError, fmt.Errorf("not found")
		},
	}

	tp, _ := process.NewTransactionProcessor(
		processorStub,
		testPubkeyConverter,
		hasher,
		marshalizer,
		funcNewTxCostHandler,
		logsMerger,
		false,
	)

	return tp
}

func TestNewTransactionProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(nil, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewTransactionProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, nil, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewTransactionProcessor_NilHasherShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, nil, marshalizer, funcNewTxCostHandler, logsMerger, true)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilHasher, err)
}

func TestNewTransactionProcessor_NilMarshalizerShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, nil, funcNewTxCostHandler, logsMerger, true)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilMarshalizer, err)
}

func TestNewTransactionProcessor_NilLogsMergerShouldErr(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, nil, true)

	require.Nil(t, tp)
	require.Equal(t, process.ErrNilLogsMerger, err)
}

func TestNewTransactionProcessor_OkValuesShouldWork(t *testing.T) {
	t.Parallel()

	tp, err := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

	require.NotNil(t, tp)
	require.Nil(t, err)
}

// ------- SendTransaction

func TestTransactionProcessor_SendTransactionInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
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

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
	rc, txHash, err := tp.SendTransaction(&data.Transaction{})

	require.Empty(t, txHash)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "no chainID")
	require.Equal(t, http.StatusBadRequest, rc)
}

func TestTransactionProcessor_SendTransactionNoVersionShouldErr(t *testing.T) {
	t.Parallel()

	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
				return nil, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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

// //------- SendMultipleTransactions

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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardID uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) (observers []*data.NodeData, e error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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

					responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, err error) {
				return []*data.NodeData{
					{Address: addrObs1, ShardId: 1},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				responseGetTx, ok := value.(*data.GetTransactionResponse)
				if !ok {
					return http.StatusOK, nil
				}

				responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, err error) {
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

				responseGetTx, ok := value.(*data.GetTransactionResponse)
				if !ok {
					return http.StatusOK, nil
				}

				responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetAllObserversCalled: func(dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
				}, nil
			},
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, err error) {
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

				responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
		marshalizer, funcNewTxCostHandler,
		logsMerger,
		true,
	)

	txStatus, err := tp.GetTransactionStatus(string(hash0), "blablabla")
	assert.Error(t, err)
	assert.Equal(t, string(data.TxStatusUnknown), txStatus)
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
			GetObserversCalled: func(shardId uint32, _ data.ObserverDataAvailabilityType) (observers []*data.NodeData, err error) {
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

				responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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

	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

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

	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

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
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

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
	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

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

	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

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

	pubKeyConv := &mock.PubKeyConverterMock{}
	tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, pubKeyConv, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)

	txHash, err := tp.ComputeTransactionHash(&data.Transaction{
		Nonce:     protoTx.Nonce,
		Value:     protoTx.Value.String(),
		Receiver:  pubKeyConv.SilentEncode(protoTx.RcvAddr, testLogger),
		Sender:    pubKeyConv.SilentEncode(protoTx.SndAddr, testLogger),
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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

					responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetObserversOnePerShardCalled: func(_ data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				return []*data.NodeData{
					{Address: addrObs0, ShardId: 0},
				}, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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
		funcNewTxCostHandler,
		logsMerger,
		true,
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
			GetFullHistoryNodesCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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
				responseGetTx, ok := value.(*data.GetTransactionResponse)
				if !ok {
					return http.StatusOK, nil
				}

				if strings.Contains(path, scHash1) {
					responseGetTx.Data.Transaction.Hash = scHash1
					return http.StatusOK, nil
				}
				if strings.Contains(path, scHash2) {
					responseGetTx.Data.Transaction.Hash = scHash2
					return http.StatusOK, nil
				}
				if strings.Contains(path, scHash3) {
					responseGetTx.Data.Transaction.Hash = scHash3
					return http.StatusOK, nil
				}

				if address == addrObs1 {
					responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
						Sender:           sndrShard0,
						Receiver:         rcvShard1,
						Nonce:            expectedNonce,
						SourceShard:      0,
						DestinationShard: 1,
						SmartContractResults: []*transaction.ApiSmartContractResult{
							scRes1, scRes2,
						},
						Status: transaction.TxStatusSuccess,
					}
					return http.StatusOK, nil
				} else if address == addrObs0 {
					responseGetTx.Data.Transaction = transaction.ApiTransactionResult{
						Nonce:            expectedNonce,
						SourceShard:      0,
						DestinationShard: 1,
						SmartContractResults: []*transaction.ApiSmartContractResult{
							scRes2, scRes3,
						},
						Status: transaction.TxStatusSuccess,
					}
					return http.StatusOK, nil
				}

				return http.StatusBadGateway, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
		funcNewTxCostHandler,
		logsMerger,
		true,
	)

	tx, err := tp.GetTransaction(string(hash0), true)
	assert.NoError(t, err)
	assert.Equal(t, expectedNonce, tx.Nonce)
	assert.Equal(t, 3, len(tx.SmartContractResults))
}

func TestTransactionProcessor_GetTransactionPool(t *testing.T) {
	t.Parallel()

	// GetTransactionsPool
	t.Run("GetTransactionsPool, flag not enabled", func(t *testing.T) {
		t.Parallel()

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, false)
		require.NotNil(t, tp)

		txs, err := tp.GetTransactionsPool("")
		assert.Nil(t, txs)
		assert.Equal(t, apiErrors.ErrOperationNotAllowed, err)
	})
	t.Run("GetTransactionsPool, no txs in pools", func(t *testing.T) {
		t.Parallel()

		addrObs0 := "observer0"
		addrObs1 := "observer1"

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0, 1}
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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
				response := value.(*data.TransactionsPoolApiResponse)
				response.Data.Transactions = data.TransactionsPool{
					RegularTransactions:  []data.WrappedTransaction{},
					SmartContractResults: []data.WrappedTransaction{},
					Rewards:              []data.WrappedTransaction{},
				}

				return http.StatusOK, nil
			},
		}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
		require.NotNil(t, tp)

		txs, err := tp.GetTransactionsPool("sender,nonce")
		require.NotNil(t, txs)
		assert.NoError(t, err)
	})
	t.Run("GetTransactionsPool, txs in 2 shards, but none in 3rd", func(t *testing.T) {
		t.Parallel()

		sndrShard0 := hex.EncodeToString([]byte("aaaa"))
		sndrShard1 := hex.EncodeToString([]byte("bbbb"))

		addrObs0 := "observer0"
		addrObs1 := "observer1"
		addrObs2 := "observer2"

		regularTxSh0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndrShard0,
				"nonce":  101,
				"hash":   "hashRegularTxSh0",
			},
		}
		rewardsTxSh0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndrShard0,
				"nonce":  102,
				"hash":   "hashRewardsTxSh0",
			},
		}
		scrTxSh0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndrShard0,
				"nonce":  103,
				"hash":   "hashSCRTxSh0",
			},
		}
		regularTxSh1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndrShard1,
				"nonce":  111,
				"hash":   "hashRegularTxSh1",
			},
		}
		rewardsTxSh1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndrShard1,
				"nonce":  112,
				"hash":   "hashRewardsTxSh1",
			},
		}
		scrTxSh1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndrShard0,
				"nonce":  113,
				"hash":   "hashSCRTxSh1",
			},
		}

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
			GetShardIDsCalled: func() []uint32 {
				return []uint32{0, 1, 2}
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
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
				if shardId == 2 {
					return []*data.NodeData{
						{Address: addrObs2, ShardId: 2},
					}, nil
				}

				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs0 {
					response := value.(*data.TransactionsPoolApiResponse)
					response.Data.Transactions = data.TransactionsPool{
						RegularTransactions:  []data.WrappedTransaction{regularTxSh0},
						SmartContractResults: []data.WrappedTransaction{scrTxSh0},
						Rewards:              []data.WrappedTransaction{rewardsTxSh0},
					}

					return http.StatusOK, nil
				} else if address == addrObs1 {
					response := value.(*data.TransactionsPoolApiResponse)
					response.Data.Transactions = data.TransactionsPool{
						RegularTransactions:  []data.WrappedTransaction{regularTxSh1},
						SmartContractResults: []data.WrappedTransaction{scrTxSh1},
						Rewards:              []data.WrappedTransaction{rewardsTxSh1},
					}

					return http.StatusOK, nil
				} else if address == addrObs2 {
					response := value.(*data.TransactionsPoolApiResponse)
					response.Data.Transactions = data.TransactionsPool{
						RegularTransactions:  []data.WrappedTransaction{},
						SmartContractResults: []data.WrappedTransaction{},
						Rewards:              []data.WrappedTransaction{},
					}
				}

				return http.StatusBadGateway, nil
			},
		}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
		require.NotNil(t, tp)

		expectedResponse := &data.TransactionsPool{
			RegularTransactions:  []data.WrappedTransaction{regularTxSh0, regularTxSh1},
			SmartContractResults: []data.WrappedTransaction{scrTxSh0, scrTxSh1},
			Rewards:              []data.WrappedTransaction{rewardsTxSh0, rewardsTxSh1},
		}
		txs, err := tp.GetTransactionsPool("sender,nonce")
		require.Nil(t, err)
		assert.Equal(t, expectedResponse, txs)
	})

	// GetTransactionsPoolForShard
	t.Run("GetTransactionsPoolForShard, flag not enabled", func(t *testing.T) {
		t.Parallel()

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, false)
		require.NotNil(t, tp)

		txs, err := tp.GetTransactionsPoolForShard(0, "")
		assert.Nil(t, txs)
		assert.Equal(t, apiErrors.ErrOperationNotAllowed, err)
	})
	t.Run("GetTransactionsPoolForShard, no txs in pool", func(t *testing.T) {
		t.Parallel()

		addrObs0 := "observer0"

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, uint32(0), shardId)
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
					}, nil
				}

				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				response := value.(*data.TransactionsPoolApiResponse)
				response.Data.Transactions = data.TransactionsPool{
					RegularTransactions:  []data.WrappedTransaction{},
					SmartContractResults: []data.WrappedTransaction{},
					Rewards:              []data.WrappedTransaction{},
				}

				return http.StatusOK, nil
			},
		}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
		require.NotNil(t, tp)

		txs, err := tp.GetTransactionsPoolForShard(0, "sender,nonce")
		require.NotNil(t, txs)
		assert.NoError(t, err)
	})
	t.Run("GetTransactionsPoolForShard, txs in pool", func(t *testing.T) {
		t.Parallel()

		sndr0 := hex.EncodeToString([]byte("aaaa"))
		sndr1 := hex.EncodeToString([]byte("bbbb"))

		addrObs0 := "observer0"
		addrObs1 := "observer1"
		addrObs2 := "observer2"

		regularTx0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndr0,
				"nonce":  101,
				"hash":   "hashRegularTx0",
			},
		}
		rewardsTx0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndr0,
				"nonce":  102,
				"hash":   "hashRewardsTx0",
			},
		}
		scrTx0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndr0,
				"nonce":  103,
				"hash":   "hashSCRTx0",
			},
		}
		regularTx1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndr1,
				"nonce":  111,
				"hash":   "hashRegularTx1",
			},
		}
		rewardsTx1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndr1,
				"nonce":  112,
				"hash":   "hashRewardsTx1",
			},
		}
		scrTx1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": sndr1,
				"nonce":  113,
				"hash":   "hashSCRTx1",
			},
		}

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				if shardId == 0 {
					return []*data.NodeData{
						{Address: addrObs0, ShardId: 0},
						{Address: addrObs1, ShardId: 0},
						{Address: addrObs2, ShardId: 0},
					}, nil
				}

				return nil, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				if address == addrObs0 {
					response := value.(*data.TransactionsPoolApiResponse)
					response.Data.Transactions = data.TransactionsPool{
						RegularTransactions:  []data.WrappedTransaction{regularTx0, regularTx1},
						SmartContractResults: []data.WrappedTransaction{scrTx0, scrTx1},
						Rewards:              []data.WrappedTransaction{rewardsTx0, rewardsTx1},
					}

					return http.StatusOK, nil
				} else if address == addrObs1 || address == addrObs2 {
					response := value.(*data.TransactionsPoolApiResponse)
					response.Data.Transactions = data.TransactionsPool{
						RegularTransactions:  []data.WrappedTransaction{},
						SmartContractResults: []data.WrappedTransaction{},
						Rewards:              []data.WrappedTransaction{},
					}

					return http.StatusOK, nil
				}

				return http.StatusBadGateway, nil
			},
		}, &mock.PubKeyConverterMock{}, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
		require.NotNil(t, tp)

		expectedResponse := &data.TransactionsPool{
			RegularTransactions:  []data.WrappedTransaction{regularTx0, regularTx1},
			SmartContractResults: []data.WrappedTransaction{scrTx0, scrTx1},
			Rewards:              []data.WrappedTransaction{rewardsTx0, rewardsTx1},
		}
		txs, err := tp.GetTransactionsPoolForShard(0, "sender,nonce")
		require.Nil(t, err)
		assert.Equal(t, expectedResponse, txs)
	})

	// GetTransactionsPoolForSender + GetLastPoolNonceForSender + GetTransactionsPoolNonceGapsForSender
	t.Run("no txs in pool", func(t *testing.T) {
		t.Parallel()
		providedPubKeyConverter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, "erd")
		providedShardId := uint32(0)
		providedSenderStr := "erd1kwh72fxl5rwndatsgrvfu235q3pwyng9ax4zxcrg4ss3p6pwuugq3gt3yc"
		addrObs0 := "observer0"

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return providedShardId, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, providedShardId, shardId)
				return []*data.NodeData{
					{Address: addrObs0, ShardId: providedShardId},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				require.True(t, strings.Contains(path, providedSenderStr))
				if strings.Contains(path, "last-nonce") {
					response := value.(*data.TransactionsPoolLastNonceForSenderApiResponse)
					response.Data.Nonce = 0
				} else if strings.Contains(path, "nonce-gaps") {
					response := value.(*data.TransactionsPoolNonceGapsForSenderApiResponse)
					response.Data.NonceGaps.Gaps = make([]data.NonceGap, 0)
				} else {
					response := value.(*data.TransactionsPoolForSenderApiResponse)
					response.Data.TxPool = data.TransactionsPoolForSender{
						Transactions: []data.WrappedTransaction{},
					}
				}

				return http.StatusOK, nil
			},
		}, providedPubKeyConverter, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
		require.NotNil(t, tp)

		txs, err := tp.GetTransactionsPoolForSender(providedSenderStr, "sender,nonce")
		require.NotNil(t, txs)
		assert.NoError(t, err)

		nonce, err := tp.GetLastPoolNonceForSender(providedSenderStr)
		assert.Equal(t, uint64(0), nonce)
		assert.Nil(t, err)

		nonceGaps, err := tp.GetTransactionsPoolNonceGapsForSender(providedSenderStr)
		assert.NotNil(t, nonceGaps)
		assert.NoError(t, err)
	})
	t.Run("txs in pool, with gaps", func(t *testing.T) {
		t.Parallel()

		providedPubKeyConverter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, "erd")
		providedShardId := uint32(0)
		providedSenderStr := "erd1kwh72fxl5rwndatsgrvfu235q3pwyng9ax4zxcrg4ss3p6pwuugq3gt3yc"
		addrObs0 := "observer0"

		lastNonce := uint64(111)
		regularTx0 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": providedSenderStr,
				"nonce":  101,
				"hash":   "hashRegularTx0",
			},
		}
		regularTx1 := data.WrappedTransaction{
			TxFields: map[string]interface{}{
				"sender": providedSenderStr,
				"nonce":  lastNonce,
				"hash":   "hashRegularTx1",
			},
		}
		providedPool := data.TransactionsPoolForSender{
			Transactions: []data.WrappedTransaction{regularTx0, regularTx1},
		}
		providedGaps := []data.NonceGap{
			{
				From: 0,
				To:   101,
			},
			{
				From: lastNonce + 1,
				To:   lastNonce + 2,
			},
		}

		tp, _ := process.NewTransactionProcessor(&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return providedShardId, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, providedShardId, shardId)
				return []*data.NodeData{
					{Address: addrObs0, ShardId: providedShardId},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				require.True(t, strings.Contains(path, providedSenderStr))
				if strings.Contains(path, "last-nonce") {
					response := value.(*data.TransactionsPoolLastNonceForSenderApiResponse)
					response.Data.Nonce = lastNonce
				} else if strings.Contains(path, "nonce-gaps") {
					response := value.(*data.TransactionsPoolNonceGapsForSenderApiResponse)
					response.Data.NonceGaps.Gaps = providedGaps
				} else {
					response := value.(*data.TransactionsPoolForSenderApiResponse)
					response.Data.TxPool = providedPool
				}

				return http.StatusOK, nil
			},
		}, providedPubKeyConverter, hasher, marshalizer, funcNewTxCostHandler, logsMerger, true)
		require.NotNil(t, tp)

		txs, err := tp.GetTransactionsPoolForSender(providedSenderStr, "sender,nonce")
		require.Nil(t, err)
		assert.Equal(t, &providedPool, txs)

		nonce, err := tp.GetLastPoolNonceForSender(providedSenderStr)
		assert.Equal(t, lastNonce, nonce)
		assert.Nil(t, err)

		nonceGaps, err := tp.GetTransactionsPoolNonceGapsForSender(providedSenderStr)
		assert.Nil(t, err)
		assert.Equal(t, providedGaps, nonceGaps.Gaps)
	})
}

func TestTransactionProcessor_computeTransactionStatus(t *testing.T) {
	t.Parallel()

	t.Run("no results should return unknown", func(t *testing.T) {
		t.Parallel()

		testData := loadJsonIntoTxAndScrs(t, "./testdata/pendingNewMoveBalance.json")
		tp := createTestProcessorFromScenarioData(testData)
		status := tp.ComputeTransactionStatus(testData.Transaction, false)
		require.Equal(t, string(data.TxStatusUnknown), status.Status)
	})
	withResults := true
	t.Run("Move balance", func(t *testing.T) {
		t.Run("pending", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/pendingNewMoveBalance.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusPending), status.Status)
		})
		t.Run("executed", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKMoveBalance.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
	})
	t.Run("SC calls", func(t *testing.T) {
		t.Run("pending new", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/pendingNewSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusPending), status.Status)
		})
		t.Run("executing", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/executingSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusPending), status.Status)
		})
		t.Run("tx ok", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
		t.Run("tx ok but with nil logs", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)
			testData.Transaction.Logs = nil
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusPending), status.Status)
		})
		t.Run("tx failed", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
	})
	t.Run("SC deploy", func(t *testing.T) {
		t.Run("ok simple SC deploy", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKSCDeploy.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
		t.Run("ok SC deploy with transfer value", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKSCDeployWithTransfer.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
		t.Run("failed SC deploy with transfer value", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedSCDeployWithTransfer.json")
			tp := createTestProcessorFromScenarioData(testData)
			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
	})
	t.Run("complex scenarios with failed async calls", func(t *testing.T) {
		t.Run("scenario 1: tx failed with ESDTs and SC calls", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedComplexScenario1.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
		t.Run("scenario 2: tx failed with ESDTs and SC calls", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedComplexScenario2.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
		t.Run("scenario 3: tx failed with ESDTs and SC calls", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedComplexScenario3.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
	})
	t.Run("relayed transaction", func(t *testing.T) {
		t.Run("failed relayed transaction un-executable", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedRelayedTxUnexecutable.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
		t.Run("failed relayed transaction with SC call", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedRelayedTxWithSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
		t.Run("failed relayed move balance intra shard transaction", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedFailedRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusFail), status.Status)
		})
		t.Run("ok relayed move balance intra shard transaction", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
		t.Run("ok relayed v2 move balance intra shard transaction", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedV2TxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
		t.Run("ok relayed sc call function balance intra shard transaction still pending", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.SCRs[0].ProcessingTypeOnSource = "SCInvoking"
			testData.SCRs[0].ProcessingTypeOnDestination = "SCInvoking"

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusPending), status.Status)
		})
		t.Run("ok relayed move balance cross shard transaction", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxCrossShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
		t.Run("tx ok", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxWithSCCall.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
		})
	})
	t.Run("reward transaction", func(t *testing.T) {
		t.Parallel()

		testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRewardTx.json")
		tp := createTestProcessorFromScenarioData(testData)

		status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
		require.Equal(t, string(transaction.TxStatusSuccess), status.Status)
	})
	t.Run("invalid transaction", func(t *testing.T) {
		t.Parallel()

		testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedInvalidBuiltinFunction.json")
		tp := createTestProcessorFromScenarioData(testData)

		status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
		require.Equal(t, string(transaction.TxStatusFail), status.Status)
	})
	t.Run("malformed transactions", func(t *testing.T) {
		t.Parallel()

		t.Run("malformed relayed v1 inner transaction - wrong sender", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Sender = "not a sender"

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v1 inner transaction - wrong receiver", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Receiver = "not a sender"

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v1 - relayed v1 marker on wrong position", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Data = append([]byte("A"), testData.Transaction.Data...)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v2 - missing marker", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedV2TxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Data = []byte("aa")

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v2 - not enough arguments", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedV2TxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Data = []byte(process.RelayedTxV2DataMarker)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v1 - not a hex character after the marker", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Data[45] = byte('T')

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v1 - marshaller will fail", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.Transaction.Data = append(testData.Transaction.Data, []byte("aaaaaa")...)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v1 - missing scrs", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/finishedOKRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			testData.SCRs = nil

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v1 - no scr generated", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/malformedRelayedTxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
		t.Run("malformed relayed v2 - no scr generated", func(t *testing.T) {
			t.Parallel()

			testData := loadJsonIntoTxAndScrs(t, "./testdata/malformedRelayedV2TxIntraShard.json")
			tp := createTestProcessorFromScenarioData(testData)

			status := tp.ComputeTransactionStatus(testData.Transaction, withResults)
			require.Equal(t, string(data.TxStatusUnknown), status.Status)
		})
	})
}

func TestTransactionProcessor_GetProcessedTransactionStatus(t *testing.T) {
	t.Parallel()

	hash0 := []byte("hash0")
	providedShardId := uint32(0)
	observerAddress := "observer address"
	tp, _ := process.NewTransactionProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
				return providedShardId, nil
			},
			GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
				require.Equal(t, providedShardId, shardId)
				return []*data.NodeData{
					{
						Address: observerAddress,
						ShardId: providedShardId,
					},
				}, nil
			},
			GetShardIDsCalled: func() []uint32 {
				return []uint32{providedShardId}
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (i int, err error) {
				assert.Contains(t, path, string(hash0))

				txResponse := value.(*data.GetTransactionResponse)
				txResponse.Data.Transaction.Nonce = 0
				txResponse.Data.Transaction.Status = transaction.TxStatusSuccess

				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
		hasher,
		marshalizer,
		funcNewTxCostHandler,
		logsMerger,
		true,
	)

	status, err := tp.GetProcessedTransactionStatus(string(hash0))
	assert.Nil(t, err)
	assert.Equal(t, string(transaction.TxStatusPending), status.Status) // not a move balance tx with missing finish markers
}

func TestTransactionProcessor_GetProcessedStatusIntraShardTxWithPendingSCR(t *testing.T) {
	txWithSCRs := loadJsonIntoTxAndScrs(t, "./testdata/transactionWithScrs.json")

	processorStub := &mock.ProcessorStub{
		GetShardIDsCalled: func() []uint32 {
			return []uint32{0} // force everything intra-shard for test setup simplicity
		},
		ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{
				{
					Address: "test",
					ShardId: 0,
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
			valueC, ok := value.(*data.GetTransactionResponse)
			if !ok {
				return http.StatusOK, nil
			}
			valueC.Data.Transaction = *txWithSCRs.SCRs[0]

			return http.StatusOK, nil
		},
	}
	tp, _ := process.NewTransactionProcessor(
		processorStub,
		testPubkeyConverter,
		hasher,
		marshalizer,
		funcNewTxCostHandler,
		logsMerger,
		false,
	)

	status := tp.ComputeTransactionStatus(txWithSCRs.Transaction, true)
	require.Equal(t, string(transaction.TxStatusPending), status.Status)

}

func TestCheckIfFailed(t *testing.T) {
	t.Parallel()

	logs := `{
        "address": "erd1qqqqqqqqqqqqqpgqzhpcdd8jg77m06zwqmhgw9xdmukn6pfeh2uslry9u8",
        "events": [
          {
            "address": "erd1qqqqqqqqqqqqqpgqzhpcdd8jg77m06zwqmhgw9xdmukn6pfeh2uslry9u8",
            "identifier": "signalError",
            "topics": [
              "Y7snvmIze+8YqIYkhNMYG8zZ7Q8PzaroT/Z7+3rEdCU=",
              "ZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3Q="
            ],
            "data": "QDY1Nzg2NTYzNzU3NDY5NmY2ZTIwNjY2MTY5NmM2NTY0",
            "additionalData": [
              "QDY1Nzg2NTYzNzU3NDY5NmY2ZTIwNjY2MTY5NmM2NTY0"
            ]
          },
          {
            "address": "erd1vwaj00nzxda77x9gscjgf5ccr0xdnmg0plx646z07ealk7kywsjsqf596y",
            "identifier": "internalVMErrors",
            "topics": [
              "AAAAAAAAAAAFABXDhrTyR7236E4G7ocUzd8tPQU5urk=",
              "ZnVsZmlsbA=="
            ],
            "data": "CglydW50aW1lLmdvOjgzMCBbZXhlY3V0aW9uIGZhaWxlZF0gW2Z1bGZpbGxdCglydW50aW1lLmdvOjgzMCBbZXhlY3V0aW9uIGZhaWxlZF0gW2Z1bGZpbGxdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdIFtjbG9zZVRyYWRlTWFya2V0Q2FsbGJhY2tdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdIFtjbG9zZVRyYWRlTWFya2V0Q2FsbGJhY2tdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdIFtjbG9zZVRyYWRlTWFya2V0Q2FsbGJhY2tdCglydW50aW1lLmdvOjgyNyBbc3RvcmFnZSBkZWNvZGUgZXJyb3I6IGlucHV0IHRvbyBzaG9ydF0=",
            "additionalData": [
              "CglydW50aW1lLmdvOjgzMCBbZXhlY3V0aW9uIGZhaWxlZF0gW2Z1bGZpbGxdCglydW50aW1lLmdvOjgzMCBbZXhlY3V0aW9uIGZhaWxlZF0gW2Z1bGZpbGxdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdIFtjbG9zZVRyYWRlTWFya2V0Q2FsbGJhY2tdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdIFtjbG9zZVRyYWRlTWFya2V0Q2FsbGJhY2tdCglydW50aW1lLmdvOjgzMCBbZXJyb3Igc2lnbmFsbGVkIGJ5IHNtYXJ0Y29udHJhY3RdIFtjbG9zZVRyYWRlTWFya2V0Q2FsbGJhY2tdCglydW50aW1lLmdvOjgyNyBbc3RvcmFnZSBkZWNvZGUgZXJyb3I6IGlucHV0IHRvbyBzaG9ydF0="
            ]
          }
        ]
      }`
	var txLogsOnFirstLevel = &transaction.ApiLogs{}
	err := json.Unmarshal([]byte(logs), txLogsOnFirstLevel)
	require.NoError(t, err)

	ok, str := process.CheckIfFailed([]*transaction.ApiLogs{txLogsOnFirstLevel})
	require.True(t, ok)
	require.True(t, strings.Contains(str, "storage decode error: input too short"))
}
