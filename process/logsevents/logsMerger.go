package logsevents

import (
	"sort"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/hashing"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("process/logsevents")

type logsMerger struct {
	hasher      hashing.Hasher
	marshalizer marshal.Marshalizer
}

// NewLogsMerger will create a new instance of logsMerger
func NewLogsMerger(hasher hashing.Hasher, marshalizer marshal.Marshalizer) (*logsMerger, error) {
	if check.IfNil(hasher) {
		return nil, ErrNilHasher
	}
	if check.IfNil(marshalizer) {
		return nil, ErrNilMarshalizer
	}

	return &logsMerger{
		hasher:      hasher,
		marshalizer: marshalizer,
	}, nil
}

// MergeLogEvents will merge events from provided logs
func (lm *logsMerger) MergeLogEvents(logSource *transaction.ApiLogs, logDestination *transaction.ApiLogs) *transaction.ApiLogs {
	if logSource == nil {
		return logDestination
	}

	if logDestination == nil {
		return logSource
	}

	eventMap := make(map[string]*transaction.Events)
	allLogs := []*transaction.ApiLogs{logSource, logDestination}
	hashes := make([]string, 0)
	for _, lg := range allLogs {
		for _, ev := range lg.Events {
			hash, _ := core.CalculateHash(lm.marshalizer, lm.hasher, ev)
			_, found := eventMap[string(hash)]
			if found {
				continue
			}

			eventMap[string(hash)] = ev
			hashes = append(hashes, string(hash))
		}
	}
	sort.Strings(hashes)
	mergedEvents := make([]*transaction.Events, 0, len(hashes))
	for _, h := range hashes {
		mergedEvents = append(mergedEvents, eventMap[h])
	}

	return &transaction.ApiLogs{
		Address: logSource.Address,
		Events:  mergedEvents,
	}
}

// IsInterfaceNil returns true if the value under the interface is nil
func (lm *logsMerger) IsInterfaceNil() bool {
	return lm == nil
}
