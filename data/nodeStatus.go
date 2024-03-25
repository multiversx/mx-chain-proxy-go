package data

import "time"

// HeartbeatResponse matches the output structure the data field for an heartbeat response
type HeartbeatResponse struct {
	Heartbeats []PubKeyHeartbeat `json:"heartbeats"`
}

// HeartbeatApiResponse matches the output of an observer's heartbeat endpoint
type HeartbeatApiResponse struct {
	Data  HeartbeatResponse `json:"data"`
	Error string            `json:"error"`
	Code  string            `json:"code"`
}

// PubKeyHeartbeat represents the heartbeat status struct for one public key
type PubKeyHeartbeat struct {
	TimeStamp            time.Time `json:"timeStamp"`
	PublicKey            string    `json:"publicKey"`
	VersionNumber        string    `json:"versionNumber"`
	NodeDisplayName      string    `json:"nodeDisplayName"`
	Identity             string    `json:"identity"`
	ReceivedShardID      uint32    `json:"receivedShardID"`
	ComputedShardID      uint32    `json:"computedShardID"`
	PeerType             string    `json:"peerType"`
	IsActive             bool      `json:"isActive"`
	Nonce                uint64    `json:"nonce"`
	NumInstances         uint64    `json:"numInstances"`
	PeerSubType          uint32    `json:"peerSubType"`
	PidString            string    `json:"pidString"`
	NumTrieNodesReceived uint64    `json:"numTrieNodesReceived"`
}

// StatusResponse represents the status received when trying to find an online node
type StatusResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Running bool   `json:"running"`
}

// NodeStatusResponse holds the metrics returned from the node
type NodeStatusResponse struct {
	Nonce                uint64 `json:"erd_nonce"`
	ProbableHighestNonce uint64 `json:"erd_probable_highest_nonce"`
	AreVmQueriesReady    string `json:"erd_are_vm_queries_ready"`
}

// NodeStatusAPIResponseData holds the mapping of the data field when returning the status of a node
type NodeStatusAPIResponseData struct {
	Metrics NodeStatusResponse `json:"metrics"`
}

// NodeStatusAPIResponse represents the mapping of the response of a node's status
type NodeStatusAPIResponse struct {
	Data  NodeStatusAPIResponseData `json:"data"`
	Error string                    `json:"error"`
	Code  string                    `json:"code"`
}

// TrieStatisticsResponse holds trie statistics metrics
type TrieStatisticsResponse struct {
	AccountsSnapshotNumNodes uint64 `json:"accounts-snapshot-num-nodes"`
}

// TrieStatisticsAPIResponse represents the mapping of the response of a node's trie statistics
type TrieStatisticsAPIResponse struct {
	Data  TrieStatisticsResponse `json:"data"`
	Error string                 `json:"error"`
	Code  string                 `json:"code"`
}

// WaitingEpochsLeftResponse matches the output structure the data field for a waiting epochs left response
type WaitingEpochsLeftResponse struct {
	EpochsLeft uint32 `json:"epochsLeft"`
}

// WaitingEpochsLeftApiResponse matches the output of an observer's waiting epochs left endpoint
type WaitingEpochsLeftApiResponse struct {
	Data  WaitingEpochsLeftResponse `json:"data"`
	Error string                    `json:"error"`
	Code  string                    `json:"code"`
}
