package process_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/database"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(nil, &mock.PubKeyConverterMock{}, database.NewDisabledElasticSearchConnector())

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewAccountProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, nil, database.NewDisabledElasticSearchConnector())

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilPubKeyConverter, err)
}

func TestNewAccountProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, database.NewDisabledElasticSearchConnector())

	assert.NotNil(t, ap)
	assert.Nil(t, err)
}

//------- GetAccount

func TestAccountProcessor_GetAccountInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.PubKeyConverterMock{}, database.NewDisabledElasticSearchConnector())
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
		database.NewDisabledElasticSearchConnector(),
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
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return nil, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
		database.NewDisabledElasticSearchConnector(),
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
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address1", ShardId: 0},
					{Address: "address2", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				return 0, errExpected
			},
		},
		&mock.PubKeyConverterMock{},
		database.NewDisabledElasticSearchConnector(),
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
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: addressFail, ShardId: 0},
					{Address: "adress2", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				if address == addressFail {
					return 0, errExpected
				}

				valRespond := value.(*data.AccountApiResponse)
				valRespond.Data.AccountData = respondedAccount.AccountData
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
		database.NewDisabledElasticSearchConnector(),
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Equal(t, &respondedAccount.AccountData, accnt)
	assert.Nil(t, err)
}

func TestAccountProcessor_GetValueForAKeyShoudWork(t *testing.T) {
	t.Parallel()

	expectedValue := "dummyValue"
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				valRespond := value.(*data.AccountKeyValueResponse)
				valRespond.Data.Value = expectedValue
				return 0, nil
			},
		},
		&mock.PubKeyConverterMock{},
		database.NewDisabledElasticSearchConnector(),
	)

	key := "key"
	addr1 := "DEADBEEF"
	value, err := ap.GetValueForKey(addr1, key)
	assert.Nil(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestAccountProcessor_GetValueForAKeyShoudError(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("err")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return []*data.NodeData{
					{Address: "address", ShardId: 0},
				}, nil
			},
			CallGetRestEndPointCalled: func(address string, path string, value interface{}) (int, error) {
				return 0, expectedError
			},
		},
		&mock.PubKeyConverterMock{},
		database.NewDisabledElasticSearchConnector(),
	)

	key := "key"
	addr1 := "DEADBEEF"
	value, err := ap.GetValueForKey(addr1, key)
	assert.Equal(t, "", value)
	assert.Equal(t, process.ErrSendingRequest, err)
}
