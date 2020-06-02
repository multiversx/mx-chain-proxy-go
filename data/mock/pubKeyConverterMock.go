package mock

import (
	"encoding/hex"
)

// PubkeyConverterMock -
type PubKeyConverterMock struct {
	len int
}

// Decode -
func (pcm *PubKeyConverterMock) Decode(humanReadable string) ([]byte, error) {
	return hex.DecodeString(humanReadable)
}

// Encode -
func (pcm *PubKeyConverterMock) Encode(pkBytes []byte) string {
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
