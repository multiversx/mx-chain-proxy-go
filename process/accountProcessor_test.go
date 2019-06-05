package process_test

import (
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountProcessor_NilCoreProcessorShouldErr(t *testing.T) {
	t.Parallel()

	ap, err := process.NewAccountProcessor(nil)

	assert.Nil(t, ap)
	assert.Equal(t, process.ErrNilCoreProcessor, err)
}
