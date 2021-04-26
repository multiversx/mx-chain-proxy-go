package process

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func TestSnapshotIndexer(t *testing.T) {
	si, _ := NewSnapshotIndexer()
	_ = si.IndexSnapshot([]*data.SnapshotItem{
		{
			Address: "erwerewdwe",
			Balance: "fsdsd",
		},
	}, "")
}
