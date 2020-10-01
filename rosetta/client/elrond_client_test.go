package client

import (
	"errors"
	"testing"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client/mock"
	"github.com/stretchr/testify/assert"
)

func TestInitializeElrondClient(t *testing.T) {
	t.Parallel()

	localErr := errors.New("err")
	count := 0
	roundDuration := uint64(4000)
	startTime := uint64(1000)
	elrondProxy := &mock.ElrondProxyClientMock{}
	elrondProxy.GetNetworkConfigMetricsCalled = func() (*data.GenericAPIResponse, error) {
		if count == 2 {
			return &data.GenericAPIResponse{
				Data: map[string]interface{}{
					"config": map[string]interface{}{
						"erd_chain_id":       "1",
						"erd_round_duration": roundDuration,
						"erd_start_time":     startTime,
					},
				},
			}, nil
		}
		count++
		return nil, localErr
	}

	elrondProxyClient, err := NewElrondClient(elrondProxy)
	assert.Nil(t, err)
	assert.Equal(t, roundDuration, elrondProxyClient.roundDurationMilliseconds)
	assert.Equal(t, startTime, elrondProxyClient.blockchainStartTime)
}
