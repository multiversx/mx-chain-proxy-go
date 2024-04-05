package process_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-proxy-go/process"
	"github.com/multiversx/mx-chain-proxy-go/process/mock"
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
	proof, err := pp.GetProof("rootHash", "invalid hex number")

	assert.Nil(t, proof)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestProofProcessor_GetProofSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := fmt.Errorf("expected error")
	returnedProof := []string{"valid", "proof"}

	pp, _ := process.NewProofProcessor(
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
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := value.(*data.GenericAPIResponse)
				valRespond.Data = returnedProof
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	response, err := pp.GetProof("rootHash", "deadbeef")
	assert.Nil(t, err)

	proofs, ok := response.Data.([]string)
	assert.True(t, ok)

	assert.Equal(t, returnedProof[0], proofs[0])
	assert.Equal(t, returnedProof[1], proofs[1])
}

func TestProofProcessor_VerifyProofInvalidHexAddressShouldErr(t *testing.T) {
	t.Parallel()

	pp, _ := process.NewProofProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	resp, err := pp.VerifyProof("rootHash", "invalid hex number", []string{})

	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestProofProcessor_VerifyProofSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := fmt.Errorf("expected error")
	proof := []string{"valid", "proof"}

	pp, _ := process.NewProofProcessor(
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
			CallPostRestEndPointCalled: func(address string, path string, param interface{}, response interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := response.(*data.GenericAPIResponse)
				valRespond.Data = true
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	resp, err := pp.VerifyProof("rootHash", "deadbeef", proof)
	assert.Nil(t, err)

	isValid, ok := resp.Data.(bool)
	assert.True(t, ok)
	assert.True(t, isValid)
}

func TestProofProcessor_GetProofDataTrieInvalidHexAddressShouldErr(t *testing.T) {
	t.Parallel()

	pp, _ := process.NewProofProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	proof, err := pp.GetProofDataTrie("abcd", "invalid hex number", "0123")

	assert.Nil(t, proof)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestProofProcessor_GetProofDataTrieSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := fmt.Errorf("expected error")
	returnedProof := []string{"valid", "proof"}

	pp, _ := process.NewProofProcessor(
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
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := value.(*data.GenericAPIResponse)
				valRespond.Data = returnedProof
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	response, err := pp.GetProofDataTrie("rootHash", "deadbeef", "key")
	assert.Nil(t, err)

	proofs, ok := response.Data.([]string)
	assert.True(t, ok)

	assert.Equal(t, returnedProof[0], proofs[0])
	assert.Equal(t, returnedProof[1], proofs[1])
}

func TestProofProcessor_GetProofCurrentRootHashInvalidHexAddressShouldErr(t *testing.T) {
	t.Parallel()

	pp, _ := process.NewProofProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	proof, err := pp.GetProofCurrentRootHash("invalid hex number")

	assert.Nil(t, proof)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestProofProcessor_GetProofCurrentRootHashSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := fmt.Errorf("expected error")
	returnedProof := []string{"valid", "proof"}

	pp, _ := process.NewProofProcessor(
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
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := value.(*data.GenericAPIResponse)
				valRespond.Data = returnedProof
				return http.StatusOK, nil
			},
		},
		&mock.PubKeyConverterMock{},
	)

	response, err := pp.GetProofCurrentRootHash("deadbeef")
	assert.Nil(t, err)

	proofs, ok := response.Data.([]string)
	assert.True(t, ok)

	assert.Equal(t, returnedProof[0], proofs[0])
	assert.Equal(t, returnedProof[1], proofs[1])
}
