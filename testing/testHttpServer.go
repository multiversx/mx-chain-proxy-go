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

	"github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
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

	if strings.Contains(req.URL.Path, "/heartbeatstatus") {
		ths.processRequestGetHeartbeat(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "node/status") {
		ths.processRequestGetNodeStatus(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "validator/statistics") {
		ths.processRequestValidatorStatistics(rw, req)
		return
	}

	if strings.Contains(req.URL.Path, "node/epoch") {
		ths.processRequestGetEpochMetrics(rw, req)
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

	if strings.Contains(req.URL.Path, "/cost") {
		ths.processRequestGetTxCost(rw, req)
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
			Balance:  "100000000000",
			CodeHash: []byte(address),
			RootHash: []byte(address),
		},
	}

	responseBuff, _ := json.Marshal(responseAccount)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetNodeStatus(rw http.ResponseWriter, _ *http.Request) {
	responsStatus := map[string]interface{}{
		"erd_app_version":                      "undefined/go1.13.4/linux-amd64",
		"erd_connected_nodes":                  8,
		"erd_consensus_created_proposed_block": 0,
		"erd_consensus_round_state":            "",
		"erd_consensus_state":                  "participant",
		"erd_count_accepted_blocks":            0,
		"erd_count_consensus":                  13,
		"erd_count_consensus_accepted_blocks":  0,
		"erd_count_leader":                     3,
		"erd_cpu_load_percent":                 0,
		"erd_current_block_hash":               "",
		"erd_current_block_size":               0,
		"erd_current_round":                    114477,
		"erd_current_round_timestamp":          1582631195,
		"erd_fork_choice_count":                3,
		"erd_highest_notarized_block_by_metachain_for_current_shard": 0,
		"erd_is_syncing":                      0,
		"erd_latest_tag_software_version":     "v1.0.88",
		"erd_live_validator_nodes":            8,
		"erd_mem_load_percent":                1,
		"erd_mem_total":                       8219807744,
		"erd_mem_used_golang":                 94572544,
		"erd_mem_used_sys":                    211030264,
		"erd_metric_community_percentage":     "0.100000",
		"erd_metric_consensus_group_size":     4,
		"erd_metric_cross_check_block_height": "0: 1, 1: 1, ",
		"erd_metric_denomination_coefficient": "0.000000000000000001",
		"erd_metric_leader_percentage":        "0.500000",
		"erd_metric_num_validators":           4,
		"erd_mini_blocks_size":                0,
		"erd_network_recv_bps":                6955,
		"erd_network_recv_bps_peak":           188146,
		"erd_network_recv_percent":            3,
		"erd_network_sent_bps":                4254,
		"erd_network_sent_bps_peak":           134154,
		"erd_network_sent_percent":            3,
		"erd_node_display_name":               "",
		"erd_node_type":                       "validator",
		"erd_num_connected_peers":             8,
		"erd_num_mini_blocks":                 0,
		"erd_num_shard_headers_from_pool":     0,
		"erd_num_shard_headers_processed":     0,
		"erd_num_transactions_processed":      0,
		"erd_num_tx_block":                    0,
		"erd_probable_highest_nonce":          0,
		"erd_rewards_value":                   "20000000000000000000",
		"erd_round_time":                      6,
		"erd_shard_id":                        4294967295,
		"erd_min_gas_price":                   100000,
		"erd_chain_id":                        "testChainId",
	}
	responseBuff, _ := json.Marshal(&responsStatus)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

type valStatsResp struct {
	Statistics map[string]*data.ValidatorApiResponse `json:"statistics"`
}

func (ths *TestHttpServer) processRequestValidatorStatistics(rw http.ResponseWriter, _ *http.Request) {
	responseValStats := map[string]*data.ValidatorApiResponse{
		"pubkey1": {
			Rating:                   50,
			TempRating:               70,
			NumLeaderSuccess:         5,
			NumLeaderFailure:         6,
			NumValidatorSuccess:      8,
			NumValidatorFailure:      9,
			TotalNumLeaderFailure:    1,
			TotalNumLeaderSuccess:    2,
			TotalNumValidatorFailure: 5,
			TotalNumValidatorSuccess: 8,
		},
		"pubkey2": {
			Rating:                   90,
			TempRating:               40,
			NumLeaderSuccess:         5,
			NumLeaderFailure:         6,
			NumValidatorSuccess:      2,
			NumValidatorFailure:      9,
			TotalNumLeaderFailure:    12,
			TotalNumLeaderSuccess:    21,
			TotalNumValidatorFailure: 25,
			TotalNumValidatorSuccess: 78,
		},
	}

	valResp := &valStatsResp{Statistics: responseValStats}
	responseBuff, _ := json.Marshal(&valResp)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetEpochMetrics(rw http.ResponseWriter, _ *http.Request) {
	responsStatus := map[string]interface{}{
		"erd_current_round":                  120,
		"erd_epoch_number":                   4,
		"erd_round_at_epoch_start":           90,
		"erd_rounds_passed_in_current_epoch": 30,
		"erd_rounds_per_epoch":               30,
	}
	responseBuff, _ := json.Marshal(&responsStatus)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetNetworkMetrics(rw http.ResponseWriter, _ *http.Request) {
	responsStatus := map[string]interface{}{
		"erd_current_round": 120,
		"erd_epoch_number":  4,
		"erd_nonce":         90,
	}
	responseBuff, _ := json.Marshal(&responsStatus)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetConfigMetrics(rw http.ResponseWriter, _ *http.Request) {
	responsStatus := map[string]interface{}{
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
	}
	responseBuff, _ := json.Marshal(&responsStatus)
	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetTxCost(rw http.ResponseWriter, _ *http.Request) {
	response := data.ResponseTxCost{
		TxCost: 123456,
	}
	responseBuff, _ := json.Marshal(response)

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

func (ths *TestHttpServer) processRequestSendFunds(rw http.ResponseWriter, _ *http.Request) {
	response := data.ResponseFunds{
		Message: "ok",
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestVmValue(rw http.ResponseWriter, _ *http.Request) {
	response := data.ResponseVmValue{
		Data: &vmcommon.VMOutput{},
	}
	responseBuff, _ := json.Marshal(response)

	_, err := rw.Write(responseBuff)
	log.LogIfError(err)
}

func (ths *TestHttpServer) processRequestGetHeartbeat(rw http.ResponseWriter, _ *http.Request) {
	heartbeats := getDummyHeartbeats()
	response := data.HeartbeatResponse{
		Heartbeats: heartbeats,
	}
	responseBuff, _ := json.Marshal(&response)

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
			MaxInactiveTime: data.Duration{Duration: 10 * time.Second},
			IsActive:        getRandomBool(),
			ReceivedShardID: uint32(i % 5),
			ComputedShardID: uint32(i%4) + 1,
			TotalUpTime:     50 + i,
			TotalDownTime:   10 + i,
			VersionNumber:   fmt.Sprintf("v1.0.%d-9e5f4b9a998d/go1.12.7/linux-amd64", i/5),
			PeerType:        peerTypes[randPeerTypeIdx.Int64()],
			NodeDisplayName: fmt.Sprintf("DisplayName%d", i),
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
