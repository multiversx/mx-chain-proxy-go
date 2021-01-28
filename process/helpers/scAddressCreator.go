package helpers

import (
	"encoding/binary"

	"github.com/ElrondNetwork/elrond-go/hashing/keccak"
)

var numInitCharsForScAddr = 10
var shardIdLength = 2
var vmTypeLen = 2
var vmType = []byte{5, 0} // WASM
var hasher = keccak.Keccak{}

// CreateScAddress will return a new smart contract address for the given address and its nonce
func CreateScAddress(owner []byte, nonce uint64) ([]byte, error) {
	if len(owner) == 0 {
		return nil, ErrEmptyOwnerAddress
	}

	base := hashFromAddressAndNonce(owner, nonce)
	prefixMask := createPrefixMask(vmType)
	suffixMask := createSuffixMask(owner)

	copy(base[:numInitCharsForScAddr], prefixMask)
	copy(base[len(base)-shardIdLength:], suffixMask)

	return base, nil
}

func hashFromAddressAndNonce(creatorAddress []byte, creatorNonce uint64) []byte {
	buffNonce := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffNonce, creatorNonce)
	adrAndNonce := append(creatorAddress, buffNonce...)
	scAddress := hasher.Compute(string(adrAndNonce))

	return scAddress
}

func createPrefixMask(vmType []byte) []byte {
	prefixMask := make([]byte, numInitCharsForScAddr-vmTypeLen)
	prefixMask = append(prefixMask, vmType...)

	return prefixMask
}

func createSuffixMask(creatorAddress []byte) []byte {
	return creatorAddress[len(creatorAddress)-2:]
}
