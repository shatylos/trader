package config

import (
	"github.com/shatylos/trader/tools/apperrors"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type Config struct {
	App           map[string]string      `yaml:"app"`
	StrategyItems []interface{}          `yaml:"strategy"`
	DomainItems   map[string]interface{} `yaml:"domain"`
}

var parsedConfig *Config
var writeMutex sync.Mutex

func GetConfig() (*Config, error) {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	if parsedConfig != nil {
		return parsedConfig, nil
	}

	fileName := "config.yaml"
	yamlFile, err := os.ReadFile(fileName)
	if err != nil {
		err = apperrors.Wrap(err, "error reading config file: %s", fileName)
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &parsedConfig)
	if err != nil {
		err = apperrors.Wrap(err, "error unmarshal config file")
		return nil, err
	}
	return parsedConfig, nil
}
