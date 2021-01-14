package factory

import (
	"fmt"
	"path/filepath"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

type apiConfigParser struct {
	baseDir string
}

// NewApiConfigParser returns a new instance of apiConfigParser
func NewApiConfigParser(baseDir string) (*apiConfigParser, error) {
	err := checkPath(baseDir)
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

func checkPath(baseDir string) error {
	file, err := core.OpenFile(baseDir)
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

	return ErrFileIsNotADirectory
}

func loadApiConfig(filepath string) (*data.ApiRoutesConfig, error) {
	cfg := &data.ApiRoutesConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
