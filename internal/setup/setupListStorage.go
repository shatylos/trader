package setup

import (
	"github.com/shatylos/trader/internal/config"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/internal/strategy"
	"github.com/shatylos/trader/tools"
	"time"
)

var setupList []*setupStructs.Setup

func SetupListInit() error {
	setupList = make([]*setupStructs.Setup, 0)

	appConfig, err := config.GetConfig()
	if err != nil {
		return err
	}

	setups, err := prepareSetupStrategies(appConfig.StrategyItems, appConfig.DomainItems)
	if err != nil {
		return err
	}

	setupList = setups

	return nil
}

func prepareSetupStrategies(strategyItems []interface{}, domainItems map[string]interface{}) ([]*setupStructs.Setup, error) {

	setupList := make([]*setupStructs.Setup, 0)

	for _, strategyItemConfig := range strategyItems {
		itemMap, ok := strategyItemConfig.(map[interface{}]interface{})
		if !ok {
			return nil, tools.AppError{
				Message: "Can not parse a strategy config.",
			}
		}
		if itemMap["code"] == nil {
			return nil, tools.AppError{
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

	ticker := time.NewTicker(time.Second / 4)
	defer ticker.Stop()

	for range ticker.C {
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

func GetSetupList() []*setupStructs.Setup {
	return setupList
}
