package txcost

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestTransactionCostProcessor_RezolveCostRequestWith3LevelsOfAsyncCalls(t *testing.T) {
	t.Parallel()

	sndTx := "0101"
	rcvTx := "0102"
	rcvSCR1 := "0103"
	rcvSCR2 := "0104"
	rcvSCR3 := "0105"

	decodeHexLocal := func(hexStr string) []byte {
		decoded, _ := hex.DecodeString(hexStr)
		return decoded
	}

	gasUsedBigTx := uint64(1500000000) - 1 - 5000
	gasSCR1 := gasUsedBigTx - 4000
	gasSCR2 := gasSCR1 - 2000
	gasSCR3 := uint64(3000)

	count := 0
	coreProc := &mock.ProcessorStub{
		GetObserversCalled: func(shardId uint32, dataAvailability data.ObserverDataAvailabilityType) ([]*data.NodeData, error) {
			return []*data.NodeData{{}}, nil
		},
		ComputeShardIdCalled: func(addressBuff []byte) (uint32, error) {
			switch {
			case bytes.Equal(addressBuff, decodeHexLocal(sndTx)) == true:
				return 0, nil
			case bytes.Equal(addressBuff, decodeHexLocal(rcvTx)) == true:
				return 1, nil
			case bytes.Equal(addressBuff, decodeHexLocal(rcvSCR1)) == true:
				return 2, nil
			case bytes.Equal(addressBuff, decodeHexLocal(rcvSCR2)) == true:
				return 3, nil
			case bytes.Equal(addressBuff, decodeHexLocal(rcvSCR3)) == true:
				return 4, nil
			default:
				return 0, nil
			}
		},
		CallPostRestEndPointCalled: func(address string, path string, req interface{}, response interface{}) (int, error) {
			switch count {
			case 0:
				responseGetTx := response.(*data.ResponseTxCost)
				responseGetTx.Data.TxCost = 1000
			case 1:
				responseGetTx := response.(*data.ResponseTxCost)
				responseGetTx.Data.TxCost = gasUsedBigTx
				responseGetTx.Data.ScResults = map[string]*data.ExtendedApiSmartContractResult{
					"scr1": {
						ApiSmartContractResult: &transaction.ApiSmartContractResult{
							CallType: 1,
							SndAddr:  rcvTx,
							RcvAddr:  rcvSCR1,
							Data:     "scCall2@dummy",
							GasLimit: gasSCR1,
						},
					},
				}
			case 2:
				responseGetTx := response.(*data.ResponseTxCost)
				responseGetTx.Data.TxCost = gasSCR1
				responseGetTx.Data.ScResults = map[string]*data.ExtendedApiSmartContractResult{
					"scr2": {
						ApiSmartContractResult: &transaction.ApiSmartContractResult{
							CallType: 1,
							SndAddr:  rcvSCR1,
							RcvAddr:  rcvSCR2,
							Data:     "scCall3@dummy",
							GasLimit: gasSCR2,
						},
					},
				}
			case 3:
				responseGetTx := response.(*data.ResponseTxCost)
				responseGetTx.Data.TxCost = gasSCR2
				responseGetTx.Data.ScResults = map[string]*data.ExtendedApiSmartContractResult{
					"scr2": {
						ApiSmartContractResult: &transaction.ApiSmartContractResult{
							CallType: 1,
							SndAddr:  rcvSCR2,
							RcvAddr:  rcvSCR3,
							Data:     "scCall4@dummy",
							GasLimit: gasSCR2 - 5000,
						},
					},
				}
			case 4:
				responseGetTx := response.(*data.ResponseTxCost)
				responseGetTx.Data.TxCost = gasSCR3
				responseGetTx.Data.ScResults = map[string]*data.ExtendedApiSmartContractResult{
					"scr3": {
						ApiSmartContractResult: &transaction.ApiSmartContractResult{
							SndAddr:  rcvSCR3,
							RcvAddr:  rcvSCR2,
							CallType: 0,
							Data:     "final@shouldNotBeCall",
						},
					},
				}
			}

			count++
			return http.StatusOK, nil
		},
	}

	newTxCostProcessor, _ := NewTransactionCostProcessor(
		coreProc, &mock.PubKeyConverterMock{})

	tx := &data.Transaction{
		Data:     []byte("scCall1@first"),
		Sender:   sndTx,
		Receiver: rcvTx,
	}

	expectedGas := uint64(14000)
	res, err := newTxCostProcessor.ResolveCostRequest(tx)
	require.Nil(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedGas, res.TxCost)
}
