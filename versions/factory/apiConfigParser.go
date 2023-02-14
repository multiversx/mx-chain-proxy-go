package factory

import (
	"fmt"
	"path/filepath"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-proxy-go/data"
)

type apiConfigParser struct {
	baseDir string
}

// NewApiConfigParser returns a new instance of apiConfigParser
func NewApiConfigParser(baseDir string) (*apiConfigParser, error) {
	err := checkDirectoryPath(baseDir)
	if err != nil {
		return nil, err
	}

	return &apiConfigParser{
		baseDir: baseDir,
	}, nil
}

// GetConfigForVersion will open the configuration file and load the api routes config
func (acp *apiConfigParser) GetConfigForVersion(version string) (*data.ApiRoutesConfig, error) {
	filePath := filepath.Join(acp.baseDir, fmt.Sprintf("%s.toml", version))
	return loadApiConfig(filePath)
}

func checkDirectoryPath(baseDirectory string) error {
	file, err := core.OpenFile(baseDirectory)
	if err != nil {
		return err
	}

	fileStats, err := file.Stat()
	if err != nil {
		return err
	}

	if fileStats.IsDir() {
		return nil
	}

	return ErrNoDirectoryAtPath
}

func loadApiConfig(filepath string) (*data.ApiRoutesConfig, error) {
	cfg := &data.ApiRoutesConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
