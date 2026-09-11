package process

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBlockedDelegationFunctionTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                     string
		data                     []byte
		blockedDelegation        bool
		blockedDelegationManager bool
	}{
		{
			name:              "delegation function",
			data:              []byte("claimRewards"),
			blockedDelegation: true,
		},
		{
			name:              "delegation function with arguments",
			data:              []byte("unDelegate@01"),
			blockedDelegation: true,
		},
		{
			name:                     "delegation manager function",
			data:                     []byte("mergeValidatorToDelegationWithWhitelist@00"),
			blockedDelegationManager: true,
		},
		{
			name: "function prefix",
			data: []byte("claimRewardsExtra"),
		},
		{
			name: "long data",
			data: bytes.Repeat([]byte{'a'}, 1<<20),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			blockedDelegation, blockedDelegationManager := getBlockedDelegationFunctionTypes(testCase.data)

			require.Equal(t, testCase.blockedDelegation, blockedDelegation)
			require.Equal(t, testCase.blockedDelegationManager, blockedDelegationManager)
		})
	}
}

func TestIsDelegationSCAddress(t *testing.T) {
	t.Parallel()

	addressAfterCounterCarry := firstDelegationSCAddress
	addressAfterCounterCarry[27] = 1
	addressAfterCounterCarry[28] = 0

	require.True(t, isDelegationSCAddress(firstDelegationSCAddress[:]))
	require.True(t, isDelegationSCAddress(addressAfterCounterCarry[:]))
	require.False(t, isDelegationSCAddress(delegationManagerSCAddress[:]))
	require.False(t, isDelegationSCAddress(make([]byte, len(firstDelegationSCAddress))))
}
