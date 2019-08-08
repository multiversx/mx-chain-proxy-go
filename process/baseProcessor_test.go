package process_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-proxy-go/config"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/process"
	"github.com/ElrondNetwork/elrond-proxy-go/process/mock"
	"github.com/gin-gonic/gin/json"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Nonce int
	Name  string
}

func createTestHttpServer(
	matchingPath string,
	response []byte,
) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			if req.URL.String() == matchingPath {
				_, _ = rw.Write(response)
			}
		}

		if req.Method == "POST" {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(req.Body)
			_, _ = rw.Write(buf.Bytes())
		}
	}))
}

func TestNewBaseProcessor_WithNilAddressConverterShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBaseProcessor(nil)

	assert.Nil(t, bp)
	assert.Equal(t, process.ErrNilAddressConverter, err)
}

func TestNewBaseProcessor_WithValidAddressConverterShouldWork(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBaseProcessor(&mock.AddressConverterStub{})

	assert.NotNil(t, bp)
	assert.Nil(t, err)
}

//------- ApplyConfig

func TestBaseProcessor_ApplyConfigNilCfgShouldErr(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	err := bp.ApplyConfig(nil)

	assert.Equal(t, process.ErrNilConfig, err)
}

func TestBaseProcessor_ApplyConfigNoObserversShouldErr(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	err := bp.ApplyConfig(&config.Config{})

	assert.Equal(t, process.ErrEmptyObserversList, err)
}

func TestBaseProcessor_ApplyConfigShouldProcessConfigAndGetShouldWork(t *testing.T) {
	t.Parallel()

	observersList := []*data.Observer{
		{
			Address: "address1",
			ShardId: 0,
		},
		{
			Address: "address2",
			ShardId: 0,
		},
		{
			Address: "address3",
			ShardId: 1,
		},
	}

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	err := bp.ApplyConfig(&config.Config{
		Observers: observersList,
	})

	assert.Nil(t, err)
	observers, err := bp.GetObservers(0)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(observers))
	assert.Equal(t, observers[0], observersList[0])
	assert.Equal(t, observers[1], observersList[1])

	observers, err = bp.GetObservers(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(observers))
	assert.Equal(t, observers[0], observersList[2])
}

//------- GetObservers

func TestBaseProcessor_GetObserversEmptyListShouldErr(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	observers, err := bp.GetObservers(0)

	assert.Nil(t, observers)
	assert.Equal(t, process.ErrMissingObserver, err)
}

//------- ComputeShardId

func TestBaseProcessor_ComputeShardId(t *testing.T) {
	t.Parallel()

	observersList := []*data.Observer{
		{
			Address: "address1",
			ShardId: 0,
		},
		{
			Address: "address2",
			ShardId: 1,
		},
	}

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{
		CreateAddressFromPublicKeyBytesCalled: func(pubKey []byte) (container state.AddressContainer, e error) {
			return &mock.AddressContainerMock{
				BytesField: pubKey,
			}, nil
		},
	})
	_ = bp.ApplyConfig(&config.Config{
		Observers: observersList,
	})

	//there are 2 shards, compute ID should correctly process
	addressInShard0 := []byte{0}
	shardId, err := bp.ComputeShardId(addressInShard0)
	assert.Nil(t, err)
	assert.Equal(t, uint32(0), shardId)

	addressInShard1 := []byte{1}
	shardId, err = bp.ComputeShardId(addressInShard1)
	assert.Nil(t, err)
	assert.Equal(t, uint32(1), shardId)
}

//------- Calls

func TestBaseProcessor_CallGetRestEndPoint(t *testing.T) {
	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be send and received",
	}
	response, _ := json.Marshal(ts)

	server := createTestHttpServer("/some/path", response)
	fmt.Printf("Server: %s\n", server.URL)
	defer server.Close()

	tsRecovered := &testStruct{}
	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	err := bp.CallGetRestEndPoint(server.URL, "/some/path", tsRecovered)

	assert.Nil(t, err)
	assert.Equal(t, ts, tsRecovered)
}

func TestBaseProcessor_CallPostRestEndPoint(t *testing.T) {
	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be send",
	}
	tsRecv := &testStruct{}

	server := createTestHttpServer("/some/path", nil)
	fmt.Printf("Server: %s\n", server.URL)
	defer server.Close()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	err := bp.CallPostRestEndPoint(server.URL, "/some/path", ts, tsRecv)

	assert.Nil(t, err)
	assert.Equal(t, ts, tsRecv)
}

func TestBaseProcessor_GetAllObserversWithEmptyListShouldFail(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	observer, err := bp.GetAllObservers()
	assert.Equal(t, process.ErrNoObserverConnected, err)
	assert.Nil(t, observer)
}

func TestBaseProcessor_GetAllObserversWithOkValuesShouldPass(t *testing.T) {
	t.Parallel()

	statusResponse := data.StatusResponse{
		Message: "",
		Error:   "",
		Running: true,
	}

	statusResponseBytes, err := json.Marshal(statusResponse)
	assert.Nil(t, err)

	server := createTestHttpServer("/node/status", statusResponseBytes)
	fmt.Printf("Server: %s\n", server.URL)
	defer server.Close()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	var observersList []*data.Observer
	observersList = append(observersList, &data.Observer{
		ShardId: 0,
		Address: server.URL,
	})

	err = bp.ApplyConfig(&config.Config{
		Observers: observersList,
	})
	assert.Nil(t, err)

	observer, err := bp.GetAllObservers()
	assert.Nil(t, err)
	assert.Equal(t, server.URL, observer[0].Address)
}
