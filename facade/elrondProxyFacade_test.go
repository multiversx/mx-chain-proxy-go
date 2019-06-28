package facade_test

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/facade/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewElrondProxyFacade_NilAccountProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		nil,
		&mock.TransactionProcessorStub{},
		&mock.GetValuesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilAccountProcessor, err)
}

func TestNewElrondProxyFacade_NilTransactionProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		nil,
		&mock.GetValuesProcessorStub{},
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilTransactionProcessor, err)
}

func TestNewElrondProxyFacade_NilGetValuesProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		nil,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilGetValueProcessor, err)
}

func TestNewElrondProxyFacade_ShouldWork(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.GetValuesProcessorStub{},
	)

	assert.NotNil(t, epf)
	assert.Nil(t, err)
}

func TestElrondProxyFacade_GetAccount(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (account *data.Account, e error) {
				wasCalled = true
				return &data.Account{}, nil
			},
		},
		&mock.TransactionProcessorStub{},
		&mock.GetValuesProcessorStub{},
	)

	_, _ = epf.GetAccount("")

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(nonce uint64, sender string,
				receiver string, value *big.Int, data string,
				signature []byte, gasPrice uint64, gasLimit uint64) (s string, e error) {

				wasCalled = true

				return "", nil
			},
		},
		&mock.GetValuesProcessorStub{},
	)

	_, _ = epf.SendTransaction(
		0,
		"",
		"",
		big.NewInt(0),
		"",
		make([]byte, 0),
		0,
		0,
	)

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendUserFunds(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SendUserFundsCalled: func(receiver string) error {
				wasCalled = true

				return nil
			},
		},
		&mock.GetValuesProcessorStub{},
	)

	_ = epf.SendUserFunds("")

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetDataValue(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.GetValuesProcessorStub{
			GetDataValueCalled: func(address string, funcName string, argsBuff ...[]byte) (bytes []byte, e error) {
				wasCalled = true

				return make([]byte, 0), nil
			},
		},
	)

	_, _ = epf.GetDataValue("", "")

	assert.True(t, wasCalled)
}
