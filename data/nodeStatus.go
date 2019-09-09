package data

import "time"

// HeartbeatResponse matches the output structure of an observer's /node/heartbeatstatus endpoint
type HeartbeatResponse struct {
	Heartbeats []PubKeyHeartbeat `json:"message"`
}

// PubKeyHeartbeat represents the heartbeat status struct for one public key
type PubKeyHeartbeat struct {
	HexPublicKey    string    `json:"hexPublicKey"`
	TimeStamp       time.Time `json:"timeStamp"`
	MaxInactiveTime Duration  `json:"maxInactiveTime"`
	IsActive        bool      `json:"isActive"`
	ReceivedShardID uint32    `json:"receivedShardID"`
	ComputedShardID uint32    `json:"computedShardID"`
	TotalUpTime     Duration  `json:"totalUpTime"`
	TotalDownTime   Duration  `json:"totalDownTime"`
	VersionNumber   string    `json:"versionNumber"`
	IsValidator     bool      `json:"isValidator"`
	NodeDisplayName string    `json:"nodeDisplayName"`
}

// StatusResponse represents the status received when trying to find an online node
type StatusResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Running bool   `json:"running"`
}
