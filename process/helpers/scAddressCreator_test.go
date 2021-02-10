package helpers

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/stretchr/testify/require"
)

func TestCreateScAddress_EmptyOwner(t *testing.T) {
	res, err := CreateScAddress(nil, 0)
	require.Nil(t, res)
	require.Equal(t, ErrEmptyOwnerAddress, err)
}

func TestCreateScAddress_ShouldWork(t *testing.T) {
	bch32, _ := pubkeyConverter.NewBech32PubkeyConverter(32)
	checkResultingAddress(
		t,
		bch32,
		"erd10qaxe763fpshfl37pfzas2fnm7wxhu42t03ht7cuxsqe39szphdst6j434",
		0,
		"erd1qqqqqqqqqqqqqpgqtf3cr084m3pck56wdmjzm6p99dgw3gxsphds4c0faf",
	)

	checkResultingAddress(
		t,
		bch32,
		"erd1d5dm8zwsnyg3juznx3e9t7fd4ez6y5xsy98k82s7l99cfn8eefsqrkvfqn",
		0,
		"erd1qqqqqqqqqqqqqpgqxresaffxsf59k99ck00f8sx3pnwz22tcefsq2wx60c",
	)
}

func checkResultingAddress(t *testing.T, bech32Converter core.PubkeyConverter, ownerBech32 string, nonce uint64, expectedBech32 string) {
	ownerBytes, err := bech32Converter.Decode(ownerBech32)
	require.NoError(t, err)

	res, err := CreateScAddress(ownerBytes, nonce)
	require.NoError(t, err)
	require.Equal(t, expectedBech32, bech32Converter.Encode(res))
}
