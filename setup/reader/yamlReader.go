package reader

import (
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh"
	"gopkg.in/yaml.v2"
	"os"
)

type YamlReader struct{}

type YamlStructure struct {
	Strategy struct {
		BuyCheapSellHigh []buyCheapSellHigh.BuyCheapSellHigh `yaml:"buy_cheap_sell_high"`
	} `yaml:"strategy"`
}

func (r *YamlReader) GetSetupList() ([]*setupStructs.Setup, error) {

	setupList := make([]*setupStructs.Setup, 0)

	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var yamlStructure YamlStructure
	err = yaml.Unmarshal(yamlFile, &yamlStructure)
	if err != nil {
		return nil, err
	}

	for _, item := range yamlStructure.Strategy.BuyCheapSellHigh {
		setupList = append(setupList, &setupStructs.Setup{
			Strategy: &item,
		})
	}

	return setupList, nil
}
