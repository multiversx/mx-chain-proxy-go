package observer

import (
	"testing"

	"github.com/multiversx/mx-chain-proxy-go/config"
	"github.com/stretchr/testify/assert"
)

func TestNewObserversProviderFactory_ShouldWork(t *testing.T) {
	t.Parallel()

	opf, err := NewNodesProviderFactory(config.Config{}, "path", 2)
	assert.Nil(t, err)
	assert.NotNil(t, opf)
}

func TestObserversProviderFactory_CreateShouldReturnSimple(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.GeneralSettings.BalancedObservers = false

	opf, _ := NewNodesProviderFactory(cfg, "path", 2)
	op, err := opf.CreateObservers()
	assert.Nil(t, err)
	_, ok := op.(*simpleNodesProvider)
	assert.True(t, ok)
}

func TestObserversProviderFactory_CreateShouldReturnCircularQueue(t *testing.T) {
	t.Parallel()

	cfg := getDummyConfig()
	cfg.GeneralSettings.BalancedObservers = true

	opf, _ := NewNodesProviderFactory(cfg, "path", 2)
	op, err := opf.CreateObservers()
	assert.Nil(t, err)
	_, ok := op.(*circularQueueNodesProvider)
	assert.True(t, ok)
}
