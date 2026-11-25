package common

import (
	"bytes"
	"encoding/json"

	"github.com/multiversx/mx-chain-core-go/core"
)

const (
	UrlParameterNumKeys       = "numKeys"
	UrlParameterIteratorState = "iteratorState"
)

// TrieWalkPage carries an opaque node checkpoint and the client's pagination opt-in.
type TrieWalkPage struct {
	NumKeys       core.OptionalUint64
	IteratorState json.RawMessage
	Paginated     bool
}

// IsContinuation reports whether a checkpoint was supplied.
func (p TrieWalkPage) IsContinuation() bool {
	return len(p.IteratorState) > 0 && !bytes.Equal(bytes.TrimSpace(p.IteratorState), []byte("null"))
}

// HasBlockOptions includes moving block selectors as well as historical coordinates.
func (a AccountQueryOptions) HasBlockOptions() bool {
	return a.OnFinalBlock || a.AreHistoricalCoordinatesSet()
}
