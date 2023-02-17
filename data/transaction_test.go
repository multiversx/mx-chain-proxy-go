package data

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/data/mock"
	"github.com/stretchr/testify/require"
)

func TestNewTransactionWrapper_NilTransactionShouldErr(t *testing.T) {
	t.Parallel()

	tw, err := NewTransactionWrapper(nil, &mock.PubKeyConverterMock{})
	require.Nil(t, tw)
	require.Equal(t, ErrNilTransaction, err)
}

func TestNewTransactionWrapper_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	tx := Transaction{Nonce: 5}
	tw, err := NewTransactionWrapper(&tx, nil)
	require.Nil(t, tw)
	require.Equal(t, ErrNilPubKeyConverter, err)
}

func TestNewTransactionWrapper_ShouldWork(t *testing.T) {
	t.Parallel()

	tx := Transaction{Nonce: 5}
	tw, err := NewTransactionWrapper(&tx, &mock.PubKeyConverterMock{})
	require.NotNil(t, tw)
	require.NoError(t, err)
}

func TestTransactionWrapper_Getters(t *testing.T) {
	t.Parallel()

	data := "data"
	gasLimit := uint64(37)
	gasPrice := uint64(5)
	rcvr, _ := hex.DecodeString("receiver")

	tx := Transaction{
		Nonce:     0,
		Value:     "",
		Receiver:  hex.EncodeToString(rcvr),
		Sender:    "",
		GasPrice:  gasPrice,
		GasLimit:  gasLimit,
		Data:      []byte(data),
		Signature: "",
	}
	tw, _ := NewTransactionWrapper(&tx, &mock.PubKeyConverterMock{})
	require.NotNil(t, tw)

	require.Equal(t, []byte(data), tw.GetData())
	require.Equal(t, gasLimit, tw.GetGasLimit())
	require.Equal(t, gasPrice, tw.GetGasPrice())
	require.Equal(t, rcvr, tw.GetRcvAddr())
}
