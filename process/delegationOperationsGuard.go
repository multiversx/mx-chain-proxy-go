package process

import (
	"bytes"

	"github.com/multiversx/mx-chain-proxy-go/data"
)

const (
	delegationSCAddressPrefixLength = 25
	delegationSCAddressSuffixStart  = 29
	maxBlockedFunctionLength        = len("mergeValidatorToDelegationWithWhitelist")
)

var delegationManagerSCAddress = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 255, 255}
var firstDelegationSCAddress = [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 255, 255, 255}

var blockedDelegationFunctions = map[string]struct{}{
	"claimRewards":      {},
	"delegate":          {},
	"reDelegateRewards": {},
	"unDelegate":        {},
	"withdraw":          {},
}

var blockedDelegationManagerFunctions = map[string]struct{}{
	"claimMulti":                              {},
	"mergeValidatorToDelegationSameOwner":     {},
	"mergeValidatorToDelegationWithWhitelist": {},
	"reDelegateMulti":                         {},
}

func (tp *TransactionProcessor) shouldBlockDelegationOperation(tx *data.Transaction) bool {
	blockedDelegation, blockedDelegationManager := getBlockedDelegationFunctionTypes(tx.Data)
	if !blockedDelegation && !blockedDelegationManager {
		return false
	}

	receiver, err := tp.pubKeyConverter.Decode(tx.Receiver)
	if err != nil {
		return false
	}

	if blockedDelegationManager && bytes.Equal(receiver, delegationManagerSCAddress[:]) {
		return true
	}
	if blockedDelegation && isDelegationSCAddress(receiver) {
		return true
	}

	return false
}

func getBlockedDelegationFunctionTypes(txData []byte) (bool, bool) {
	selectorEnd := len(txData)
	if selectorEnd > maxBlockedFunctionLength+1 {
		selectorEnd = maxBlockedFunctionLength + 1
	}

	separatorIndex := bytes.IndexByte(txData[:selectorEnd], '@')
	if separatorIndex >= 0 {
		selectorEnd = separatorIndex
	} else if len(txData) > maxBlockedFunctionLength {
		return false, false
	}

	_, blockedDelegation := blockedDelegationFunctions[string(txData[:selectorEnd])]
	_, blockedDelegationManager := blockedDelegationManagerFunctions[string(txData[:selectorEnd])]

	return blockedDelegation, blockedDelegationManager
}

func isDelegationSCAddress(address []byte) bool {
	if len(address) != len(firstDelegationSCAddress) {
		return false
	}
	if !bytes.Equal(address[:delegationSCAddressPrefixLength], firstDelegationSCAddress[:delegationSCAddressPrefixLength]) {
		return false
	}
	if !bytes.Equal(address[delegationSCAddressSuffixStart:], firstDelegationSCAddress[delegationSCAddressSuffixStart:]) {
		return false
	}

	return bytes.Compare(address, firstDelegationSCAddress[:]) >= 0
}
