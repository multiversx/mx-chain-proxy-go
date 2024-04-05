package process

import (
	"fmt"
	"net/http"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const shortHashSize = 7

type aboutProcessor struct {
	baseProc   Processor
	commitID   string
	appVersion string
}

// NewAboutProcessor creates a new instance of about processor
func NewAboutProcessor(baseProc Processor, appVersion string, commit string) (*aboutProcessor, error) {
	if check.IfNil(baseProc) {
		return nil, ErrNilCoreProcessor
	}
	if len(appVersion) == 0 {
		return nil, ErrEmptyAppVersionString
	}
	if len(commit) == 0 {
		return nil, ErrEmptyCommitString
	}

	return &aboutProcessor{
		baseProc:   baseProc,
		commitID:   commit,
		appVersion: appVersion,
	}, nil
}

// GetAboutInfo will return the app info parameters
func (ap *aboutProcessor) GetAboutInfo() *data.GenericAPIResponse {
	commit := ap.commitID
	if ap.commitID != common.UndefinedCommitString {
		if len(commit) >= shortHashSize {
			commit = commit[:shortHashSize]
		}
	}

	aboutInfo := &data.AboutInfo{
		AppVersion: ap.appVersion,
		CommitID:   commit,
	}

	resp := &data.GenericAPIResponse{
		Data:  aboutInfo,
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}

	return resp
}

// GetNodesVersions will return the versions of the nodes behind proxy
func (ap *aboutProcessor) GetNodesVersions() (*data.GenericAPIResponse, error) {
	versionsMap := make(map[uint32][]string)
	allObservers, err := ap.baseProc.GetAllObservers(data.AvailabilityRecent)
	if err != nil {
		return nil, err
	}

	for _, observer := range allObservers {
		nodeVersion, err := ap.getNodeAppVersion(observer.Address)
		if err != nil {
			return nil, err
		}

		versionsMap[observer.ShardId] = append(versionsMap[observer.ShardId], nodeVersion)
	}

	return &data.GenericAPIResponse{
		Data: data.NodesVersionProxyResponseData{
			Versions: versionsMap,
		},
		Error: "",
		Code:  data.ReturnCodeSuccess,
	}, nil
}

func (ap *aboutProcessor) getNodeAppVersion(observerAddress string) (string, error) {
	var versionResponse data.NodeVersionAPIResponse
	code, err := ap.baseProc.CallGetRestEndPoint(observerAddress, NodeStatusPath, &versionResponse)
	if code != http.StatusOK {
		return "", fmt.Errorf("invalid return code %d", code)
	}

	if err != nil {
		return "", err
	}

	if len(versionResponse.Error) > 0 {
		return "", fmt.Errorf("%w while extracting the app version", err)
	}

	return versionResponse.Data.Metrics.Version, nil
}
