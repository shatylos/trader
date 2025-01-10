package reader

import (
	"github.com/shatylos/trader/setup/structs"
	"gopkg.in/yaml.v2"
	"os"
)

type YamlReader struct{}

func (y YamlReader) GetConfig() (*structs.Config, error) {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config structs.Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
