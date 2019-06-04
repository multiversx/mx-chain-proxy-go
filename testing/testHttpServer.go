package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"

	"github.com/ElrondNetwork/elrond-go-sandbox/core/logger"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

var log = logger.DefaultLogger()

// TestHttpServer is a test http server used for testing the whole binary
type TestHttpServer struct {
	httpServer *httptest.Server
}

// NewTestHttpServer creates a new TestHttpServer instance
func NewTestHttpServer() *TestHttpServer {
	ths := &TestHttpServer{}
	ths.httpServer = httptest.NewServer(
		http.HandlerFunc(ths.processRequest),
	)

	return ths
}

func (ths *TestHttpServer) processRequest(rw http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, "address") {
		ths.processRequestAddress(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "transaction") {
		ths.processRequestTransaction(rw, req)
		return
	}

	fmt.Printf("Can not serve request: %v\n", req.URL)
}

func (ths *TestHttpServer) processRequestAddress(rw http.ResponseWriter, req *http.Request) {
	_, address := path.Split(req.URL.String())

	account := &data.Account{
		Address:  address,
		Nonce:    45,
		Balance:  "1234",
		CodeHash: []byte(address),
		RootHash: []byte(address),
	}

	response, _ := json.Marshal(account)
	_, err := rw.Write(response)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestTransaction(rw http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(req.Body)
	newStr := buf.String()

	fmt.Printf("Got new request: %s\n", newStr)
}

// Close closes the test http server
func (ths *TestHttpServer) Close() {
	ths.httpServer.Close()
}

// URL returns the connecting url to the http test server
func (ths *TestHttpServer) URL() string {
	return ths.httpServer.URL
}
