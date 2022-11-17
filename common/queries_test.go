package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildUrlWithBlockQueryOptions_ShouldWork(t *testing.T) {
	url := BuildUrlWithBlockQueryOptions("/block/by-nonce/15", BlockQueryOptions{})
	require.Equal(t, "/block/by-nonce/15", url)

	url = BuildUrlWithBlockQueryOptions("/block/by-nonce/15", BlockQueryOptions{
		WithTransactions: true,
	})
	require.Equal(t, "/block/by-nonce/15?withTxs=true", url)

	url = BuildUrlWithBlockQueryOptions("/block/by-nonce/15", BlockQueryOptions{
		WithTransactions: true,
		WithLogs:         true,
	})
	require.True(t, url == "/block/by-nonce/15?withTxs=true&withLogs=true" || url == "/block/by-nonce/15?withLogs=true&withTxs=true")
}

func TestBuildUrlWithAlteredAccountsQueryOptions(t *testing.T) {
	url := BuildUrlWithAlteredAccountsQueryOptions("path", GetAlteredAccountsForBlockOptions{})
	require.Equal(t, "path", url)

	url = BuildUrlWithAlteredAccountsQueryOptions("path", GetAlteredAccountsForBlockOptions{
		TokensFilter: "token1,token2,token3",
	})
	// 2C is the ascii hex encoding of (,)
	require.Equal(t, "path?tokens=token1%2Ctoken2%2Ctoken3", url)
}
