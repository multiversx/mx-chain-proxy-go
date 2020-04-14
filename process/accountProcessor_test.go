package process_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(nil, &mock.PubKeyConverterMock{})

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewAccountProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, nil)

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewAccountProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})

	assert.NotNil(t, ap)
	assert.Nil(t, err)
}

//------- GetAccount

func TestAccountProcessor_GetAccountInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{})
	accnt, err := ap.GetAccount("invalid hex number")

	assert.Nil(t, accnt)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestAccountProcessor_GetAccountComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return nil, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return []*data.Observer{
					{Address: "address1", ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
				return errExpected
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestAccountProcessor_GetAccountSendingFailsOnFirstObserverShouldStillSend(t *testing.T) {
	t.Parallel()

	addressFail := "address1"
	errExpected := errors.New("expected error")
	respondedAccount := &data.ResponseAccount{
		AccountData: data.Account{
			Address: "an address",
		},
	}
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
				return []*data.Observer{
					{Address: addressFail, ShardId: 0},
					{Address: "adress2", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
				if address == addressFail {
					return errExpected
				}

				valRespond := value.(*data.ResponseAccount)
				valRespond.AccountData = respondedAccount.AccountData
				return nil
			},
		},
		&mock.PubKeyConverterMock{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Equal(t, &respondedAccount.AccountData, accnt)
	assert.Nil(t, err)
}

func TestAccountProcessor_ValidatorStatisticShouldFailIfNoObserverIsOnline(t *testing.T) {
	t.Parallel()

	processor := &mock.ProcessorStub{
		GetObserversCalled: func(_ uint32) ([]*data.Observer, error) {
			return []*data.Observer{
				{
					ShardId: core.MetachainShardId,
					Address: "address1",
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errors.New("offline")
		},
	}
	ap, _ := process.NewAccountProcessor(processor, &mock.PubKeyConverterMock{})

	res, err := ap.ValidatorStatistics()
	assert.Nil(t, res)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestAccountProcessor_ValidatorStatisticShouldFailIfNoMetachainObserverInList(t *testing.T) {
	t.Parallel()

	processor := &mock.ProcessorStub{
		GetObserversCalled: func(_ uint32) ([]*data.Observer, error) {
			return []*data.Observer{
				{
					ShardId: 0,
					Address: "address1",
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errors.New("offline")
		},
	}
	ap, _ := process.NewAccountProcessor(processor, &mock.PubKeyConverterMock{})

	res, err := ap.ValidatorStatistics()
	assert.Nil(t, res)
	assert.Equal(t, process.ErrSendingRequest, err)
}

func TestAccountProcessor_ValidatorStatisticShouldWork(t *testing.T) {
	t.Parallel()

	mapToRet := make(map[string]*data.ValidatorApiResponse)
	mapToRet["test"] = &data.ValidatorApiResponse{
		NumLeaderSuccess:         4,
		NumLeaderFailure:         5,
		NumValidatorSuccess:      6,
		NumValidatorFailure:      7,
		Rating:                   0.5,
		TempRating:               0.51,
		TotalNumLeaderSuccess:    4,
		TotalNumLeaderFailure:    5,
		TotalNumValidatorSuccess: 6,
		TotalNumValidatorFailure: 7,
	}

	processor := &mock.ProcessorStub{
		GetObserversCalled: func(_ uint32) ([]*data.Observer, error) {
			return []*data.Observer{
				{
					ShardId: core.MetachainShardId,
					Address: "address1",
				},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			val := value.(*process.ValStatsResponse)
			val.Statistics = mapToRet
			return nil
		},
	}
	ap, _ := process.NewAccountProcessor(processor, &mock.PubKeyConverterMock{})

	res, err := ap.ValidatorStatistics()
	assert.Nil(t, err)
	assert.Equal(t, mapToRet, res)
}
