package logsevents

import (
	"encoding/hex"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/hashing"
	"github.com/ElrondNetwork/elrond-go/marshal"
)

var log = logger.GetOrCreate("process/logsevents")

type logsMerger struct {
	hasher      hashing.Hasher
	marshalizer marshal.Marshalizer
}

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

func (lm *logsMerger) MergeLogEvents(logSource *transaction.ApiLogs, logDestination *transaction.ApiLogs) *transaction.ApiLogs {
	if logSource == nil {
		return logDestination
	}

	if logDestination == nil {
		return logSource
	}

	mergedEvents := make(map[string]*transaction.Events)
	lm.mergeEvents(mergedEvents, logSource)
	lm.mergeEvents(mergedEvents, logDestination)

	return &transaction.ApiLogs{
		Address: logSource.Address,
		Events:  convertEventsMapInSlice(mergedEvents),
	}
}

func (lm *logsMerger) mergeEvents(mergedEvents map[string]*transaction.Events, apiLog *transaction.ApiLogs) {
	for _, event := range apiLog.Events {
		logHash, err := core.CalculateHash(lm.marshalizer, lm.hasher, event)
		if err != nil {
			log.Warn("logsMerger.addLogsInMap cannot compute event hash", "error", err.Error())
		}

		logHashEncoded := hex.EncodeToString(logHash)
		mergedEvents[logHashEncoded] = event
	}
}

func convertEventsMapInSlice(eventsMap map[string]*transaction.Events) []*transaction.Events {
	events := make([]*transaction.Events, 0, len(eventsMap))
	for _, eventLog := range eventsMap {
		events = append(events, eventLog)
	}

	return events
}
