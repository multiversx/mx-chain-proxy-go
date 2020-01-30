package observer

import (
	"testing"

	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleObserversProvider_EmptyObserversListShouldErr(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.Observers = make([]*data.Observer, 0)
	sop, err := NewSimpleObserversProvider(cfg)
	assert.Nil(t, sop)
	assert.Equal(t, ErrEmptyObserversList, err)
}

func TestNewSimpleObserversProvider_ShouldWork(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	sop, err := NewSimpleObserversProvider(cfg)
	assert.Nil(t, err)
	assert.False(t, check.IfNil(sop))
}

func TestSimpleObserversProvider_GetObserversByShardIdShouldErrBecauseInvalidShardId(t *testing.T) {
	t.Parallel()

	invalidShardId := uint32(37)
	cfg := getDummyConfig()
	cqop, _ := NewSimpleObserversProvider(cfg)

	res, err := cqop.GetObserversByShardId(invalidShardId)
	assert.Nil(t, res)
	assert.Equal(t, ErrShardNotAvailable, err)
}

func TestSimpleObserversProvider_GetObserversByShardIdShouldWork(t *testing.T) {
	t.Parallel()

	shardId := uint32(0)
	cfg := getDummyConfig()
	cqop, _ := NewSimpleObserversProvider(cfg)

	res, err := cqop.GetObserversByShardId(shardId)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
}

func TestSimpleObserversProvider_GetAllObserversShouldWork(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cqop, _ := NewSimpleObserversProvider(cfg)

	res, err := cqop.GetAllObservers()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))
}
