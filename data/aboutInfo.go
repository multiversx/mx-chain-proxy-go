package data

// AboutInfo defines the structure needed for exposing app info
type AboutInfo struct {
	AppVersion string `json:"appVersion"`
	CommitID   string `json:"commitID"`
}
