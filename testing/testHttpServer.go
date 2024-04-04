package testing

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-core-go/data/vm"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

var log = logger.GetOrCreate("testing")

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
	if strings.Contains(req.URL.Path, "/esdtnft/") {
		ths.processRequestGetEsdtNftTokenData(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "/esdt/") {
		ths.processRequestGetEsdtTokenData(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "/esdt") {
		ths.processRequestGetAllEsdtTokens(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "address") {
		ths.processRequestAddress(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "block/by-") {
		ths.processFullHistoryBlockRequest(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "transaction/send") {
		ths.processRequestTransaction(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "transaction/simulate") {
		ths.processRequestTransactionSimulation(rw, req)
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

	if strings.Contains(req.URL.Path, "/heartbeatstatus") {
		ths.processRequestGetHeartbeat(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "validator/statistics") {
		ths.processRequestValidatorStatistics(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "network/config") {
		ths.processRequestGetConfigMetrics(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "network/status") {
		ths.processRequestGetNetworkMetrics(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "network/economics") {
		ths.processRequestGetEconomicsMetrics(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "/cost") {
		ths.processRequestGetTxCost(rw, req)
		return
	}

	fmt.Printf("Can not serve request: %v\n", req.URL)
}

func (ths *TestHttpServer) processRequestAddress(rw http.ResponseWriter, req *http.Request) {
	_, address := path.Split(req.URL.String())

	responseAccount := &data.AccountModel{
		Account: data.Account{
			Address:  address,
			Nonce:    45,
			Balance:  "10000000000000000000000000",
			CodeHash: []byte(address),
			RootHash: []byte(address),
		},
	}

	resp := data.GenericAPIResponse{Data: responseAccount, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(resp)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetEsdtTokenData(rw http.ResponseWriter, _ *http.Request) {
	type tkn struct {
		Name       string `json:"tokenName"`
		Balance    string `json:"balance"`
		Properties string `json:"properties"`
	}
	response := data.GenericAPIResponse{
		Data: gin.H{"tokenData": tkn{
			Name:       "testESDTtkn",
			Balance:    "999",
			Properties: "11",
		}},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}

	responseBuff, _ := json.Marshal(response)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetEsdtNftTokenData(rw http.ResponseWriter, _ *http.Request) {
	nftData := struct {
		TokenIdentifier string   `json:"tokenIdentifier"`
		Balance         string   `json:"balance"`
		Properties      string   `json:"properties"`
		Name            string   `json:"name"`
		Creator         string   `json:"creator"`
		Royalties       string   `json:"royalties"`
		Hash            []byte   `json:"hash"`
		URIs            [][]byte `json:"uris"`
		Attributes      []byte   `json:"attributes"`
	}{
		TokenIdentifier: "test identifier",
		Balance:         "1000000",
		Properties:      "1",
		Name:            "name",
		Creator:         "erd1zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zyg3zy0llugszknfpv",
		Royalties:       "10000",
		Hash:            []byte("hash"),
		URIs:            [][]byte{[]byte("uri")},
		Attributes:      []byte("1"),
	}

	response := data.GenericAPIResponse{
		Data:  gin.H{"tokenData": nftData},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}

	responseBuff, _ := json.Marshal(response)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetAllEsdtTokens(rw http.ResponseWriter, _ *http.Request) {
	response := data.GenericAPIResponse{
		Data:  gin.H{"tokens": []string{"testESDTtkn", "testESDTtkn2"}},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}

	responseBuff, _ := json.Marshal(response)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processFullHistoryBlockRequest(rw http.ResponseWriter, _ *http.Request) {
	response := data.GenericAPIResponse{
		Data:  api.Block{Nonce: 10, Round: 11},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}

	responseBuff, _ := json.Marshal(response)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

type valStatsResp struct {
	Statistics map[string]*data.ValidatorApiResponse `json:"statistics"`
}

func (ths *TestHttpServer) processRequestValidatorStatistics(rw http.ResponseWriter, _ *http.Request) {
	responseValStats := map[string]*data.ValidatorApiResponse{
		"pubkey1": {
			TempRating:                         70,
			NumLeaderSuccess:                   5,
			NumLeaderFailure:                   6,
			NumValidatorSuccess:                8,
			NumValidatorFailure:                9,
			NumValidatorIgnoredSignatures:      12,
			Rating:                             50,
			RatingModifier:                     1.1,
			TotalNumLeaderSuccess:              2,
			TotalNumLeaderFailure:              1,
			TotalNumValidatorSuccess:           8,
			TotalNumValidatorFailure:           5,
			TotalNumValidatorIgnoredSignatures: 120,
			ShardId:                            core.MetachainShardId,
			ValidatorStatus:                    "waiting",
		},
		"pubkey2": {
			TempRating:                         40,
			NumLeaderSuccess:                   5,
			NumLeaderFailure:                   6,
			NumValidatorSuccess:                2,
			NumValidatorFailure:                9,
			NumValidatorIgnoredSignatures:      11,
			Rating:                             90,
			RatingModifier:                     1,
			TotalNumLeaderSuccess:              21,
			TotalNumLeaderFailure:              12,
			TotalNumValidatorSuccess:           78,
			TotalNumValidatorFailure:           25,
			TotalNumValidatorIgnoredSignatures: 110,
			ShardId:                            1,
			ValidatorStatus:                    "eligible",
		},
	}

	valResp := &valStatsResp{Statistics: responseValStats}
	resp := data.GenericAPIResponse{Data: valResp, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(&resp)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetNetworkMetrics(rw http.ResponseWriter, _ *http.Request) {
	responseStatus := map[string]interface{}{
		"erd_nonce":                          90,
		"erd_current_round":                  120,
		"erd_epoch_number":                   4,
		"erd_round_at_epoch_start":           90,
		"erd_rounds_passed_in_current_epoch": 30,
		"erd_rounds_per_epoch":               30,
	}
	resp := data.GenericAPIResponse{Data: gin.H{"status": responseStatus}, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(&resp)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetEconomicsMetrics(rw http.ResponseWriter, _ *http.Request) {
	responseStatus := map[string]interface{}{
		"erd_dev_rewards":              "0",
		"erd_inflation":                "120",
		"erd_epoch_number":             4,
		"erd_total_fees":               "3500000000",
		"erd_epoch_for_economics_data": 30,
	}
	type metricsResp struct {
		Metrics map[string]interface{} `json:"metrics"`
	}
	resp := struct {
		Data  metricsResp `json:"data"`
		Error string      `json:"error"`
		Code  string      `json:"code"`
	}{
		Data: metricsResp{Metrics: responseStatus},
		Code: string(data.ReturnCodeSuccess),
	}

	responseBuff, _ := json.Marshal(&resp)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetConfigMetrics(rw http.ResponseWriter, _ *http.Request) {
	responseStatus := map[string]interface{}{
		"erd_chain_id":                   "testnet",
		"erd_gas_per_data_byte":          4,
		"erd_meta_consensus_group_size":  5,
		"erd_min_gas_limit":              5,
		"erd_min_gas_price":              5,
		"erd_num_metachain_nodes":        30,
		"erd_num_nodes_in_shard":         30,
		"erd_num_shards_without_meta":    30,
		"erd_round_duration":             30,
		"erd_shard_consensus_group_size": 30,
		"erd_start_time":                 30,
		"erd_min_transaction_version":    1,
	}
	resp := data.GenericAPIResponse{Data: gin.H{"config": responseStatus}, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(&resp)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetTxCost(rw http.ResponseWriter, _ *http.Request) {
	response := data.ResponseTxCost{
		Data: data.TxCostResponseData{TxCost: 123456},
	}
	resp := data.GenericAPIResponse{Data: response, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(resp)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestTransaction(rw http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(req.Body)
	newStr := buf.String()

	txHash := sha256.Sum256([]byte(newStr))
	txHexHash := hex.EncodeToString(txHash[:])

	response := data.ResponseTransaction{
		Data: data.TransactionResponseData{TxHash: txHexHash},
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestTransactionSimulation(rw http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(req.Body)
	newStr := buf.String()

	txHash := sha256.Sum256([]byte(newStr))
	txHexHash := hex.EncodeToString(txHash[:])

	fmt.Printf("Got new request: %s, replying with %s\n", newStr, txHexHash)
	response := data.ResponseTransactionSimulation{
		Data: data.TransactionSimulationResponseData{
			Result: data.TransactionSimulationResults{
				Status: "executed",
				ScResults: map[string]*transaction.ApiSmartContractResult{
					"scRHash": {
						SndAddr: "erd111",
						RcvAddr: "erd122",
					},
				},
				Receipts: map[string]*transaction.ApiReceipt{
					"rcptHash": {
						SndAddr: "erd111",
						Value:   big.NewInt(10),
					},
				},
				FailReason: "-",
			},
		},
		Error: "",
		Code:  "successful",
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestSendFunds(rw http.ResponseWriter, _ *http.Request) {
	response := data.ResponseFunds{
		Message: "ok",
	}
	resp := data.GenericAPIResponse{Data: response, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(resp)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestVmValue(rw http.ResponseWriter, _ *http.Request) {
	response := data.ResponseVmValue{
		Data: data.VmValuesResponseData{Data: &vm.VMOutputApi{}},
	}
	resp := data.GenericAPIResponse{Data: response, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(resp)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetHeartbeat(rw http.ResponseWriter, _ *http.Request) {
	heartbeats := getDummyHeartbeats()
	response := data.HeartbeatResponse{
		Heartbeats: heartbeats,
	}
	resp := data.GenericAPIResponse{Data: response, Code: data.ReturnCodeSuccess}
	responseBuff, _ := json.Marshal(&resp)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func getDummyHeartbeats() []data.PubKeyHeartbeat {
	noOfHeartbeatsToGenerate := 80
	noOfBytesOfAPubKey := 64
	var heartbeats []data.PubKeyHeartbeat
	peerTypes := []string{"eligible", "waiting", "observer"}
	for i := 0; i < noOfHeartbeatsToGenerate; i++ {
		randPeerTypeIdx, _ := rand.Int(rand.Reader, big.NewInt(3))
		pkBuff := make([]byte, noOfBytesOfAPubKey)
		_, _ = rand.Reader.Read(pkBuff)
		heartbeats = append(heartbeats, data.PubKeyHeartbeat{
			PublicKey:       hex.EncodeToString(pkBuff),
			TimeStamp:       time.Now(),
			IsActive:        getRandomBool(),
			ReceivedShardID: uint32(i % 5),
			ComputedShardID: uint32(i%4) + 1,
			VersionNumber:   fmt.Sprintf("v1.0.%d-9e5f4b9a998d/go1.12.7/linux-amd64", i/5),
			PeerType:        peerTypes[randPeerTypeIdx.Int64()],
			NodeDisplayName: fmt.Sprintf("DisplayName%d", i),
			Identity:        fmt.Sprintf("Identity%d", i),
			Nonce:           uint64(i),
			NumInstances:    1,
		})
	}

	return heartbeats
}

func getRandomBool() bool {
	return mathRand.Int31()%2 == 0
}

// Close closes the test http server
func (ths *TestHttpServer) Close() {
	ths.httpServer.Close()
}

// URL returns the connecting url to the http test server
func (ths *TestHttpServer) URL() string {
	return ths.httpServer.URL
}
