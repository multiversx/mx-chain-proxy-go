package mock

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core"
)

// PubKeyConverterMock -
type PubKeyConverterMock struct {
	len int
}

// Decode -
func (pcm *PubKeyConverterMock) Decode(humanReadable string) ([]byte, error) {
	return hex.DecodeString(humanReadable)
}

// Encode -
func (pcm *PubKeyConverterMock) Encode(pkBytes []byte) (string, error) {
	return hex.EncodeToString(pkBytes), nil
}

// EncodeSlice -
func (pcm *PubKeyConverterMock) EncodeSlice(pkBytesSlice [][]byte) ([]string, error) {
	results := make([]string, 0)
	for _, pk := range pkBytesSlice {
		results = append(results, hex.EncodeToString(pk))
	}

	return results, nil
}

// SilentEncode -
func (pcm *PubKeyConverterMock) SilentEncode(pkBytes []byte, _ core.Logger) string {
	return hex.EncodeToString(pkBytes)
}

// Len -
func (pcm *PubKeyConverterMock) Len() int {
	return pcm.len
}

// IsInterfaceNil -
func (pcm *PubKeyConverterMock) IsInterfaceNil() bool {
	return pcm == nil
}
