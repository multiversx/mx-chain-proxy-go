package testing

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"time"

	"github.com/ElrondNetwork/elrond-go/core/logger"
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

	if strings.Contains(req.URL.Path, "transaction/send") {
		ths.processRequestTransaction(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "transaction/send-user-funds") {
		ths.processRequestSendFunds(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "vm-values") {
		ths.processRequestVmValue(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "/node/status") {
		ths.processRequestOnlineObserver(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "/heartbeat") {
		ths.processRequestGetHeartbeat(rw, req)
		return
	}

	fmt.Printf("Can not serve request: %v\n", req.URL)
}

func (ths *TestHttpServer) processRequestAddress(rw http.ResponseWriter, req *http.Request) {
	_, address := path.Split(req.URL.String())

	responseAccount := &data.ResponseAccount{
		AccountData: data.Account{
			Address:  address,
			Nonce:    45,
			Balance:  "1234",
			CodeHash: []byte(address),
			RootHash: []byte(address),
		},
	}

	responseBuff, _ := json.Marshal(responseAccount)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestTransaction(rw http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(req.Body)
	newStr := buf.String()

	txHash := sha256.Sum256([]byte(newStr))
	txHexHash := hex.EncodeToString(txHash[:])

	fmt.Printf("Got new request: %s, replying with %s\n", newStr, txHexHash)
	response := data.ResponseTransaction{
		TxHash: txHexHash,
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestSendFunds(rw http.ResponseWriter, req *http.Request) {
	response := data.ResponseFunds{
		Message: "ok",
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestVmValue(rw http.ResponseWriter, req *http.Request) {
	response := data.ResponseVmValue{
		HexData: "DEADBEEFDEADBEEFDEADBEEF",
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetHeartbeat(rw http.ResponseWriter, req *http.Request) {
	var heartbeats []data.PubKeyHeartbeat
	heartbeats = append(heartbeats, data.PubKeyHeartbeat{
		HexPublicKey:    "pk1",
		TimeStamp:       time.Time{},
		MaxInactiveTime: data.Duration{},
		IsActive:        false,
		ShardID:         0,
		TotalUpTime:     data.Duration{},
		TotalDownTime:   data.Duration{},
		VersionNumber:   "",
		IsValidator:     false,
		NodeDisplayName: "test1",
	})
	heartbeats = append(heartbeats, data.PubKeyHeartbeat{
		HexPublicKey:    "pk2",
		TimeStamp:       time.Time{},
		MaxInactiveTime: data.Duration{},
		IsActive:        false,
		ShardID:         0,
		TotalUpTime:     data.Duration{},
		TotalDownTime:   data.Duration{},
		VersionNumber:   "",
		IsValidator:     false,
		NodeDisplayName: "test2",
	})
	response := data.HeartbeatResponse{
		Heartbeats: heartbeats,
	}
	responseBuff, _ := json.Marshal(&response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestOnlineObserver(rw http.ResponseWriter, req *http.Request) {
	status := data.StatusResponse{
		Message: "ok",
		Error:   "",
		Running: true,
	}

	responseBuff, _ := json.Marshal(status)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

// Close closes the test http server
func (ths *TestHttpServer) Close() {
	ths.httpServer.Close()
}

// URL returns the connecting url to the http test server
func (ths *TestHttpServer) URL() string {
	return ths.httpServer.URL
}
