package config

import (
	"reflect"
	"time"

	"github.com/ElrondNetwork/elrond-go-sandbox/core"
	"github.com/ElrondNetwork/elrond-go-sandbox/core/logger"
)

var log = logger.DefaultLogger()

//ConfigFileMonitor
//TODO - work in progress
type ConfigFileMonitor struct {
	refreshInterval int
	lastConfig      *Config
	filename        string
	readFileHandler func(dest interface{}, relativePath string, log *logger.Logger) error
}

// NewConfigFileMonitor creates a new config file monitor instance
func NewConfigFileMonitor(refreshInterval int, filename string) (*ConfigFileMonitor, error) {
	//TODO checks

	cfm := &ConfigFileMonitor{
		refreshInterval: refreshInterval,
		filename:        filename,
	}
	cfm.readFileHandler = core.LoadTomlFile
	go cfm.readLoop()

	return cfm, nil
}

func (cfm *ConfigFileMonitor) readLoop() {
	for {
		readConfig := &Config{}
		err := cfm.readFileHandler(readConfig, cfm.filename, log)
		if err == nil {
			if !reflect.DeepEqual(cfm.lastConfig, readConfig) {

			}
		} else {
			log.Error(err.Error())
		}

		time.Sleep(time.Duration(cfm.refreshInterval) * time.Second)
	}
}
