package observer

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/stretchr/testify/assert"
)

func TestNewObserversProviderFactory_ShouldWork(t *testing.T) {
	t.Parallel()

	opf, err := NewObserversProviderFactory(config.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, opf)
}

func TestObserversProviderFactory_CreateShouldReturnSimple(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.GeneralSettings.BalancedObservers = false

	opf, _ := NewObserversProviderFactory(cfg)
	op, err := opf.Create()
	assert.Nil(t, err)
	_, ok := op.(*SimpleObserversProvider)
	assert.True(t, ok)
}

func TestObserversProviderFactory_CreateShouldReturnCircularQueue(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.GeneralSettings.BalancedObservers = true

	opf, _ := NewObserversProviderFactory(cfg)
	op, err := opf.Create()
	assert.Nil(t, err)
	_, ok := op.(*CircularQueueObserversProvider)
	assert.True(t, ok)
}
