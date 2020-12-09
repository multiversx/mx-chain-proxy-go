package process_test

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/data/state/factory"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/database"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAccountProcessor_GetAccountInvalidHexAddressShouldErr(t *testing.T) {
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

func TestAccountProcessor_GetValueForAKeyShouldWork(t *testing.T) {
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

func TestAccountProcessor_GetValueForAKeyShouldError(t *testing.T) {
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

func TestAccountProcessor_GetShardIForAddressShouldWork(t *testing.T) {
	t.Parallel()

	shardC, err := sharding.NewMultiShardCoordinator(uint32(2), 0)
	require.NoError(t, err)

	bech32C, _ := pubkeyConverter.NewBech32PubkeyConverter(32)

	// this addressShard0 should be in shard 0 for a 2 shards configuration
	addressShard0 := "erd1ffqlrryvwrnfh2523wmzrhvx5d8p2wmxeau64fps4lnqq5qex68q7ax8k5"

	// this addressShard1 should be in shard 1 for a 2 shards configuration
	addressShard1 := "erd1qqe9qll7n66lv4cuuml2wxsv3sd2t0eyajkyjr7rvtqmhha0cgsse4pel3"

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return shardC.ComputeId(addressBuff), nil
			},
			GetObserversCalled: func(shardId uint32) (observers []*data.NodeData, e error) {
				return observers, nil
			},
		},
		bech32C,
		database.NewDisabledElasticSearchConnector(),
	)

	shardID, err := ap.GetShardIDForAddress(addressShard1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), shardID)

	shardID, err = ap.GetShardIDForAddress(addressShard0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), shardID)
}

func TestAccountProcessor_GetShardIDForAddressShouldError(t *testing.T) {
	t.Parallel()

	expectedError := errors.New("err")
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{
			ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
				return 0, expectedError
			},
		},
		&mock.PubKeyConverterMock{},
		database.NewDisabledElasticSearchConnector(),
	)

	shardID, err := ap.GetShardIDForAddress("aaaa")
	assert.Equal(t, uint32(0), shardID)
	assert.Equal(t, expectedError, err)
}

func TestAccountProcessor_GetTransactions(t *testing.T) {
	t.Parallel()

	converter, _ := factory.NewPubkeyConverter(config.PubkeyConfig{
		Length: 32,
		Type:   "bech32",
	})
	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{},
		converter,
		&mock.ElasticSearchConnectorMock{},
	)

	_, err := ap.GetTransactions("invalidAddress")
	assert.True(t, errors.Is(err, process.ErrInvalidAddress))

	_, err = ap.GetTransactions("")
	assert.True(t, errors.Is(err, process.ErrInvalidAddress))

	_, err = ap.GetTransactions("erd1ycega644rvjtgtyd8hfzt6hl5ymaa8ml2nhhs5cv045cz5vxm00q022myr")
	assert.Nil(t, err)
}
