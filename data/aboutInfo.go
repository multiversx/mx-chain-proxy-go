package data

// AboutInfo defines the structure needed for exposing app info
type AboutInfo struct {
	AppVersion string `json:"appVersion"`
	CommitID   string `json:"commitID"`
}

// NodesVersionProxyResponseData maps the response data for the proxy's nodes version endpoint
type NodesVersionProxyResponseData struct {
	Versions map[uint32][]string `json:"versions"`
}

// NodeVersionAPIResponse maps the format to be used when fetching the node version from API
type NodeVersionAPIResponse struct {
	Data struct {
		Metrics struct {
			Version string `json:"erd_app_version"`
		} `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
