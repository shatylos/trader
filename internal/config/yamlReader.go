package config

import (
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

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &parsedConfig)
	if err != nil {
		return nil, err
	}
	return parsedConfig, nil
}
