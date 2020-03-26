package data

import "time"

// HeartbeatResponse matches the output structure of an observer's /node/heartbeatstatus endpoint
type HeartbeatResponse struct {
	Heartbeats []PubKeyHeartbeat `json:"message"`
}

// PubKeyHeartbeat represents the heartbeat status struct for one public key
type PubKeyHeartbeat struct {
	TimeStamp       time.Time `json:"timeStamp"`
	HexPublicKey    string    `json:"hexPublicKey"`
	VersionNumber   string    `json:"versionNumber"`
	NodeDisplayName string    `json:"nodeDisplayName"`
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
