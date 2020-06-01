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
	TimeStamp       time.Time `json:"timeStamp"`
	PublicKey       string    `json:"publicKey"`
	VersionNumber   string    `json:"versionNumber"`
	NodeDisplayName string    `json:"nodeDisplayName"`
	Identity        string    `json:"identity"`
	TotalUpTime     int       `json:"totalUpTimeSec"`
	TotalDownTime   int       `json:"totalDownTimeSec"`
	MaxInactiveTime Duration  `json:"maxInactiveTime"`
	ReceivedShardID uint32    `json:"receivedShardID"`
	ComputedShardID uint32    `json:"computedShardID"`
	PeerType        string    `json:"peerType"`
	IsActive        bool      `json:"isActive"`
}

// StatusResponse represents the status received when trying to find an online node
type StatusResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Running bool   `json:"running"`
}
