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
	Nonce           uint64    `json:"nonce"`
	NumInstances    uint64    `json:"numInstances"`
}

// StatusResponse represents the status received when trying to find an online node
type StatusResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Running bool   `json:"running"`
}

// MaiarReferalApiResponse represents the API response from maiar containing eligible addresses for MEX multiplier
type MaiarReferalApiResponse struct {
	Total     uint32   `json:"total"`
	Addresses []string `json:"addresses"`
}

// DirectStakedValue holds the total staked value for an address
type DirectStakedValue struct {
	Address  string `json:"address"`
	Staked   string `json:"staked"`
	TopUp    string `json:"topUp"`
	Total    string `json:"total"`
	Unstaked string `json:"unstaked"`
}

// DelegationList represents the API response from the DirectStakedValue observer call
type DirectStakedValueList struct {
	List []*DirectStakedValue `json:"list"`
}

// MaiarReferalApiResponse represents the API response from maiar containing eligible addresses for MEX multiplier
type DelegationItem struct {
	DelegationScAddress string `json:"delegationScAddress"`
	UnclaimedRewards    string `json:"unclaimedRewards"`
	UndelegatedValue    string `json:"undelegatedValue"`
	Value               string `json:"value"`
}

// Delegator holds the delegator address and the slice of delegated values
type Delegator struct {
	DelegatorAddress string            `json:"delegatorAddress"`
	DelegatedTo      []*DelegationItem `json:"delegatedTo"`
	Total            string            `json:"total"`
	UnclaimedTotal   string            `json:"unclaimedTotal"`
	UndelegatedTotal string            `json:"undelegatedTotal"`
	WaitingTotal     string            `json:"waitingTotal,omitempty"`
}

// DelegationList represents the API response from the DelegatedInfo observer call
type DelegationList struct {
	List []*Delegator `json:"list"`
}

type AccountBalance struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

// AccountBalanceListResponse defines the list of accounts returned by the account list route
type AccountBalanceListResponse struct {
	List []*AccountBalance `json:"list"`
}

type SnapshotItem struct {
	Address         string `json:"address"`
	Balance         string `json:"balance"`
	Staked          string `json:"staked"`
	Waiting         string `json:"waiting"`
	Unstaked        string `json:"unstaked"`
	Unclaimed       string `json:"unclaimed"`
	IsMaiarEligible bool   `json:"isMaiarEligible"`
}
