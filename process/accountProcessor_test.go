package process_test

import (
	"encoding/hex"
	"errors"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/kyber"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(nil, &mock.KeygenStub{})

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}

func TestNewAccountProcessor_NilKeyGenShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, nil)

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilKeyGen, err)
}

func TestNewAccountProcessor_WithCoreProcessorShouldWork(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.KeygenStub{})

	assert.NotNil(t, ap)
	assert.Nil(t, err)
}

//------- GetAccount

func TestAccountProcessor_GetAccountInvalidHexAdressShouldErr(t *testing.T) {
	t.Parallel()

	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{}, &mock.KeygenStub{})
	accnt, err := ap.GetAccount("invalid hex number")

	assert.Nil(t, accnt)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid byte")
}

func TestAccountProcessor_GetAccountComputeShardIdFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, errExpected
		},
	},
		&mock.KeygenStub{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountGetObserversFailsShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return nil, errExpected
		},
	},
		&mock.KeygenStub{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Nil(t, accnt)
	assert.Equal(t, errExpected, err)
}

func TestAccountProcessor_GetAccountSendingFailsOnAllObserversShouldErr(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
		ComputeShardIdCalled: func(addressBuff []byte) (u uint32, e error) {
			return 0, nil
		},
		GetObserversCalled: func(shardId uint32) (observers []*data.Observer, e error) {
			return []*data.Observer{
				{Address: "adress1", ShardId: 0},
				{Address: "adress2", ShardId: 0},
			}, nil
		},
		CallGetRestEndPointCalled: func(address string, path string, value interface{}) error {
			return errExpected
		},
	},
		&mock.KeygenStub{},
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
	ap, _ := process.NewAccountProcessor(&mock.ProcessorStub{
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
		&mock.KeygenStub{},
	)
	address := "DEADBEEF"
	accnt, err := ap.GetAccount(address)

	assert.Equal(t, &respondedAccount.AccountData, accnt)
	assert.Nil(t, err)
}

//-------- PublicKeyFromPrivateKey

func TestAccountProcessor_PublicKeyFromPrivateKeyShouldErrIfCannotConvert(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("error")

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{},
		&mock.KeygenStub{
			PrivateKeyFromByteArrayCalled: func(b []byte) (crypto.PrivateKey, error) {
				return nil, expectedErr
			},
		},
	)

	kg := signing.NewKeyGenerator(kyber.NewBlakeSHA256Ed25519())
	sk, _ := kg.GeneratePair()
	skBytes, _ := sk.ToByteArray()
	skHex := hex.EncodeToString(skBytes)

	_, err := ap.PublicKeyFromPrivateKey(skHex)
	assert.Equal(t, expectedErr, err)
}

func TestAccountProcessor_PublicKeyFromPrivateKeyShouldWork(t *testing.T) {
	t.Parallel()

	kg := signing.NewKeyGenerator(kyber.NewBlakeSHA256Ed25519())
	sk, _ := kg.GeneratePair()
	skBytes, _ := sk.ToByteArray()
	skHex := hex.EncodeToString(skBytes)
	pkFromSk := sk.GeneratePublic()
	pkBytes, _ := pkFromSk.ToByteArray()
	pkHex := hex.EncodeToString(pkBytes)

	ap, _ := process.NewAccountProcessor(
		&mock.ProcessorStub{},
		&mock.KeygenStub{
			PrivateKeyFromByteArrayCalled: func(b []byte) (crypto.PrivateKey, error) {
				return sk, nil
			},
		},
	)

	pk, err := ap.PublicKeyFromPrivateKey(skHex)
	assert.Nil(t, err)
	assert.Equal(t, pkHex, pk)
}
