package process

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func TestSnapshotIndexer(t *testing.T) {
	si, _ := NewSnapshotIndexer()
	_ = si.IndexSnapshot(make([]*data.SnapshotItem, 0), "1619440599")
}
