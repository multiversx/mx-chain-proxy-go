package data

import "time"

// HeartbeatResponse matches the output structure of an observer's /node/heartbeatstatus endpoint
type HeartbeatResponse struct {
	Heartbeats []PubKeyHeartbeat `json:"message"`
}

// PubKeyHeartbeat represents the heartbeat status struct for one public key
type PubKeyHeartbeat struct {
	HexPublicKey    string
	TimeStamp       time.Time
	MaxInactiveTime Duration
	IsActive        bool
	ShardID         uint32
	TotalUpTime     Duration
	TotalDownTime   Duration
	VersionNumber   string
	IsValidator     bool
	NodeDisplayName string
}

// StatusResponse represents the status received when trying to find an online node
type StatusResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Running bool   `json:"running"`
}
