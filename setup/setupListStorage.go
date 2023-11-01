package setup

import (
	setupReader "bitbucket.org/shatylos/trader/setup/reader"
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	"bitbucket.org/shatylos/trader/strategy"
	"bitbucket.org/shatylos/trader/utils"
	"time"
)

var setupList []*setupStructs.Setup
var reader ReaderInterface

func init() {
	reader = &setupReader.YamlReader{}
}

func SetupListInit() error {
	setupList = make([]*setupStructs.Setup, 0)

	config, err := reader.GetConfig()
	if err != nil {
		return err
	}

	setups, err := prepareSetupStrategies(config.StrategyItems, config.DomainItems)
	if err != nil {
		return err
	}

	setupList = setups

	return nil
}

func prepareSetupStrategies(strategyItems []interface{}, domainItems map[interface{}]interface{}) ([]*setupStructs.Setup, error) {

	setupList := make([]*setupStructs.Setup, 0)

	for _, strategyItemConfig := range strategyItems {
		itemMap, ok := strategyItemConfig.(map[interface{}]interface{})
		if !ok {
			return nil, utils.AppError{
				Message: "Can not parse a strategy config.",
			}
		}
		if itemMap["code"] == nil {
			return nil, utils.AppError{
				Message: "The strategy config must contain a code field",
			}
		}

		strategyItem, err := strategy.GetStrategyByCode(itemMap["code"].(string))
		if err != nil {
			return nil, err
		}
		err = strategyItem.SetConfig(strategyItemConfig, domainItems)
		if err != nil {
			return nil, err
		}

		setupList = append(setupList, &setupStructs.Setup{
			Strategy: strategyItem,
		})
	}

	return setupList, nil
}

func LoadNextSetupStep(setupChanelContext chan *setupStructs.Setup) {
	var setupItemResult *setupStructs.Setup

	for range time.Tick(time.Second / 4) {
		for _, setupListItem := range setupList {
			if setupListItem.GetStatus() == setupStructs.StatusReadyForNext {
				setupListItem.SetStatus(setupStructs.StatusInProgress)
				setupItemResult = setupListItem
				break
			}
		}

		if setupItemResult != nil {
			setupChanelContext <- setupItemResult
			return
		}
	}
}
