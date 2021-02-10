package process

import (
	"errors"
	"strings"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDnsProcessor_NilPubKeyConverterShouldErr(t *testing.T) {
	t.Parallel()

	dp, err := NewDnsProcessor(nil)
	require.Equal(t, ErrNilPubKeyConverter, err)
	require.Nil(t, dp)
}

func TestNewDnsProcessor(t *testing.T) {
	t.Parallel()

	dp, err := NewDnsProcessor(&mock.PubkeyConverterMock{})
	require.NoError(t, err)
	require.NotNil(t, dp)
}

func TestDnsProcessor_GetDnsAddresses(t *testing.T) {
	t.Parallel()

	dp, _ := NewDnsProcessor(&mock.PubkeyConverterMock{})

	addresses, err := dp.GetDnsAddresses()
	require.NoError(t, err)
	require.Equal(t, 256, len(addresses))
}

func TestDnsProcessor_GetDnsAddressForUsernameInvalidUsernameLength(t *testing.T) {
	t.Parallel()

	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32)
	dp, _ := NewDnsProcessor(converter)

	res, err := dp.GetDnsAddressForUsername("a")
	require.Empty(t, res)
	require.Equal(t, ErrInvalidUsernameLength, err)

	res, err = dp.GetDnsAddressForUsername(strings.Repeat("a", 100))
	require.Empty(t, res)
	require.Equal(t, ErrInvalidUsernameLength, err)
}

func TestDnsProcessor_GetDnsAddressForUsernameInvalidCharacterInUsername(t *testing.T) {
	t.Parallel()

	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32)
	dp, _ := NewDnsProcessor(converter)

	invalidChar := "!"
	username := strings.Repeat(invalidChar, 5)
	res, err := dp.GetDnsAddressForUsername(username)
	require.Empty(t, res)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrInvalidUsername))
}

func TestDnsProcessor_GetDnsAddressForUsernameInvalidSuffix(t *testing.T) {
	t.Parallel()

	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32)
	dp, _ := NewDnsProcessor(converter)

	username := "username.erlond"
	res, err := dp.GetDnsAddressForUsername(username)
	require.Empty(t, res)
	require.Equal(t, ErrInvalidUsername, err)
}

func TestDnsProcessor_GetDnsAddressForUsernameTwoSuffixes(t *testing.T) {
	t.Parallel()

	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32)
	dp, _ := NewDnsProcessor(converter)

	username := "username.elrond.elrond"
	res, err := dp.GetDnsAddressForUsername(username)
	require.Empty(t, res)
	require.Equal(t, ErrInvalidUsername, err)
}

func TestDnsProcessor_GetDnsAddressForUsername(t *testing.T) {
	t.Parallel()

	converter, _ := pubkeyConverter.NewBech32PubkeyConverter(32)
	dp, _ := NewDnsProcessor(converter)

	resWithoutSuffix, err := dp.GetDnsAddressForUsername("test")
	require.NoError(t, err)
	require.Equal(t, "erd1qqqqqqqqqqqqqpgqx4ca3eu4k6w63hl8pjjyq2cp7ul7a4ukqz0skq6fxj", resWithoutSuffix)

	resWithSuffix, err := dp.GetDnsAddressForUsername("test.elrond")
	require.NoError(t, err)
	require.Equal(t, resWithoutSuffix, resWithSuffix)

	resWithoutSuffix, err = dp.GetDnsAddressForUsername("1test1")
	require.NoError(t, err)
	require.Equal(t, "erd1qqqqqqqqqqqqqpgquf7wpxtnln0a8ywf6hwy82pk5644y2c4qpqs66sj90", resWithoutSuffix)

	resWithSuffix, err = dp.GetDnsAddressForUsername("1test1.elrond")
	require.NoError(t, err)
	require.Equal(t, resWithoutSuffix, resWithSuffix)

	resWithoutSuffix, err = dp.GetDnsAddressForUsername("55555")
	require.NoError(t, err)
	require.Equal(t, "erd1qqqqqqqqqqqqqpgqynp5z59pgqjnphfxg00mkpzx54an6038qzwq7apcwt", resWithoutSuffix)

	resWithSuffix, err = dp.GetDnsAddressForUsername("55555.elrond")
	require.NoError(t, err)
	require.Equal(t, resWithoutSuffix, resWithSuffix)
}
