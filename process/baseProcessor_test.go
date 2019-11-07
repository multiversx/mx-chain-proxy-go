package process_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/sharding"
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

	bp, err := process.NewBaseProcessor(nil, 5, &mock.ShardCoordinatorMock{})

	assert.Nil(t, bp)
	assert.Equal(t, process.ErrNilAddressConverter, err)
}

func TestNewBaseProcessor_WithInvalidRequestTimeoutShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBaseProcessor(&mock.AddressConverterStub{}, -5, &mock.ShardCoordinatorMock{})

	assert.Nil(t, bp)
	assert.Equal(t, process.ErrInvalidRequestTimeout, err)
}

func TestNewBaseProcessor_WithNilShardCoordinatorShouldErr(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, nil)

	assert.Nil(t, bp)
	assert.Equal(t, process.ErrNilShardCoordinator, err)
}

func TestNewBaseProcessor_WithOkValuesShouldWork(t *testing.T) {
	t.Parallel()

	bp, err := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})

	assert.NotNil(t, bp)
	assert.Nil(t, err)
}

//------- ApplyConfig

func TestBaseProcessor_ApplyConfigNilCfgShouldErr(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
	err := bp.ApplyConfig(nil)

	assert.Equal(t, process.ErrNilConfig, err)
}

func TestBaseProcessor_ApplyConfigNoObserversShouldErr(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
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

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
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

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
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

	msc, _ := sharding.NewMultiShardCoordinator(3, 0)
	bp, _ := process.NewBaseProcessor(
		&mock.AddressConverterStub{
			CreateAddressFromPublicKeyBytesCalled: func(pubKey []byte) (container state.AddressContainer, e error) {
				return &mock.AddressContainerMock{
					BytesField: pubKey,
				}, nil
			},
		},
		5,
		msc,
	)
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
	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
	err := bp.CallGetRestEndPoint(server.URL, "/some/path", tsRecovered)

	assert.Nil(t, err)
	assert.Equal(t, ts, tsRecovered)
}

func TestBaseProcessor_CallGetRestEndPointShouldTimeout(t *testing.T) {
	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be send and received",
	}
	response, _ := json.Marshal(ts)

	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		_, _ = rw.Write(response)
	}))
	fmt.Printf("Server: %s\n", testServer.URL)
	defer testServer.Close()

	tsRecovered := &testStruct{}
	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 1, &mock.ShardCoordinatorMock{})
	err := bp.CallGetRestEndPoint(testServer.URL, "/some/path", tsRecovered)

	assert.NotEqual(t, ts.Name, tsRecovered.Name)
	assert.NotNil(t, err)
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

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
	rc, err := bp.CallPostRestEndPoint(server.URL, "/some/path", ts, tsRecv)

	assert.Nil(t, err)
	assert.Equal(t, ts, tsRecv)
	assert.Equal(t, http.StatusOK, rc)
}

func TestBaseProcessor_CallPostRestEndPointShouldTimeout(t *testing.T) {
	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be send",
	}
	tsRecv := &testStruct{}

	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		tsBytes, _ := json.Marshal(ts)
		_, _ = rw.Write(tsBytes)
	}))

	fmt.Printf("Server: %s\n", testServer.URL)
	defer testServer.Close()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 1, &mock.ShardCoordinatorMock{})
	rc, err := bp.CallPostRestEndPoint(testServer.URL, "/some/path", ts, tsRecv)

	assert.NotEqual(t, tsRecv.Name, ts.Name)
	assert.NotNil(t, err)
	assert.Equal(t, process.HttpNoStatusCode, rc)
}

func TestBaseProcessor_GetAllObserversWithEmptyListShouldFail(t *testing.T) {
	t.Parallel()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
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

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{}, 5, &mock.ShardCoordinatorMock{})
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
