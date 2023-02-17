package process

import (
	"github.com/multiversx/mx-chain-proxy-go/common"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

const shortHashSize = 7

type aboutProcessor struct {
	commitID   string
	appVersion string
}

// NewAboutProcessor creates a new instance of about processor
func NewAboutProcessor(appVersion string, commit string) (*aboutProcessor, error) {
	if len(appVersion) == 0 {
		return nil, ErrEmptyAppVersionString
	}
	if len(commit) == 0 {
		return nil, ErrEmptyCommitString
	}

	return &aboutProcessor{
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
