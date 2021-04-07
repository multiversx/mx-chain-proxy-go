package process_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewProofProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	pp, err := process.NewProofProcessor(nil, &mock.PubKeyConverterMock{})

	assert.Nil(t, pp)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewProofProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	pp, err := process.NewProofProcessor(&mock.ProcessorStub{}, nil)

	assert.Nil(t, pp)
	assert.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewProofProcessor(t *testing.T) {
	t.Parallel()

	pp, err := process.NewProofProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})

	assert.NotNil(t, pp)
	assert.Nil(t, err)
}

func TestProofProcessor_GetProofInvalidHexAddressShouldErr(t *testing.T) {
	t.Parallel()

	pp, _ := process.NewProofProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	proof, err := pp.GetProof([]byte("rootHash"), []byte("invalid hex number"))

	assert.Nil(t, proof)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestProofProcessor_GetProofSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := fmt.Errorf("expected error")
	returnedProof := [][]byte{[]byte("valid"), []byte("proof")}

	pp, _ := process.NewProofProcessor(
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
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := value.(*data.GetProofResponse)
				valRespond.Data = returnedProof
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	proof, err := pp.GetProof([]byte("rootHash"), []byte("deadbeef"))
	assert.Nil(t, err)
	assert.NotNil(t, proof)
	assert.Equal(t, returnedProof, proof)
}

func TestProofProcessor_VerifyProofInvalidHexAddressShouldErr(t *testing.T) {
	t.Parallel()

	pp, _ := process.NewProofProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	ok, err := pp.VerifyProof([]byte("rootHash"), []byte("invalid hex number"), [][]byte{})

	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestProofProcessor_VerifyProofSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := fmt.Errorf("expected error")
	proof := [][]byte{[]byte("valid"), []byte("proof")}

	pp, _ := process.NewProofProcessor(
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
			CallPostRestEndPointCalled: func(address string, path string, param interface{}, response interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := response.(*data.VerifyProofResponse)
				valRespond.Data = true
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	ok, err := pp.VerifyProof([]byte("rootHash"), []byte("deadbeef"), proof)
	assert.Nil(t, err)
	assert.True(t, ok)
}
