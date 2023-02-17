package logsevents

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	hasherFactory "github.com/multiversx/mx-chain-core-go/hashing/factory"
	marshalFactory "github.com/multiversx/mx-chain-core-go/marshal/factory"
	"github.com/stretchr/testify/require"
)

func TestNewLogsMerger(t *testing.T) {
	t.Parallel()

	hasher, _ := hasherFactory.NewHasher("blake2b")
	marshalizer, _ := marshalFactory.NewMarshalizer("json")
	lp, err := NewLogsMerger(nil, marshalizer)
	require.Nil(t, lp)
	require.Equal(t, ErrNilHasher, err)

	lp, err = NewLogsMerger(hasher, nil)
	require.Nil(t, lp)
	require.Equal(t, ErrNilMarshalizer, err)

	lp, err = NewLogsMerger(hasher, marshalizer)
	require.NotNil(t, lp)
	require.Nil(t, err)
}

func TestLogsMerger_MergeLogsNoLogsOnDst(t *testing.T) {
	t.Parallel()

	hasher, _ := hasherFactory.NewHasher("blake2b")
	marshalizer, _ := marshalFactory.NewMarshalizer("json")
	lp, _ := NewLogsMerger(hasher, marshalizer)

	sourceLog := &transaction.ApiLogs{
		Address: "addr1",
		Events: []*transaction.Events{
			{
				Data: []byte("data1"),
			},
		},
	}

	res := lp.MergeLogEvents(sourceLog, nil)
	require.Equal(t, sourceLog, res)
}

func TestLogsMerger_MergeLogsNoLogsOnSource(t *testing.T) {
	t.Parallel()

	hasher, _ := hasherFactory.NewHasher("blake2b")
	marshalizer, _ := marshalFactory.NewMarshalizer("json")
	lp, _ := NewLogsMerger(hasher, marshalizer)

	destinationLog := &transaction.ApiLogs{
		Address: "addr1",
		Events: []*transaction.Events{
			{
				Data: []byte("data1"),
			},
		},
	}

	res := lp.MergeLogEvents(nil, destinationLog)
	require.Equal(t, destinationLog, res)
}

func TestLogsMerger_MergeLogs(t *testing.T) {
	hasher, _ := hasherFactory.NewHasher("blake2b")
	marshalizer, _ := marshalFactory.NewMarshalizer("json")
	lp, _ := NewLogsMerger(hasher, marshalizer)

	sourceLog := &transaction.ApiLogs{
		Address: "addr1",
		Events: []*transaction.Events{
			{
				Data: []byte("data1"),
			},
			{
				Data: []byte("data2"),
			},
		},
	}
	destinationLog := &transaction.ApiLogs{
		Address: "addr1",
		Events: []*transaction.Events{
			{
				Data: []byte("data1"),
			},
			{
				Data: []byte("data3"),
			},
		},
	}

	res := lp.MergeLogEvents(sourceLog, destinationLog)
	require.Len(t, res.Events, 3)
}
