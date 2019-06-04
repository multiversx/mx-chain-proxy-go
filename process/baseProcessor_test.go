package process_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	chOutput chan string,
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
			newStr := buf.String()
			chOutput <- newStr
		}
	}))
}

func TestBaseProcessor_CallGetRestEndPoint(t *testing.T) {
	ts := &testStruct{
		Nonce: 10000,
		Name:  "a test struct to be send and received",
	}
	response, _ := json.Marshal(ts)

	server := createTestHttpServer("/some/path", response, nil)
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
		Name:  "a test struct to be send and received",
	}

	chOutput := make(chan string, 10)
	server := createTestHttpServer("/some/path", nil, chOutput)
	fmt.Printf("Server: %s\n", server.URL)
	defer server.Close()

	bp, _ := process.NewBaseProcessor(&mock.AddressConverterStub{})
	err := bp.CallPostRestEndPoint(server.URL, "/some/path", ts)

	assert.Nil(t, err)
	select {
	case data := <-chOutput:
		fmt.Println(data)
	case <-time.After(time.Second * 2):
		assert.Fail(t, "failed to receive data")
	}
}
