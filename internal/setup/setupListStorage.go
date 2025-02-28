package setup

import (
	"fmt"
	"github.com/shatylos/trader/internal/config"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/internal/strategy"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
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
		if itemMap["id"] == nil {
			return nil, tools.AppError{
				Message: "The strategy config must contain a id field",
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
			ID:       itemMap["id"].(string),
		})
	}

	return setupList, nil
}

func GetSetupList() []*setupStructs.Setup {
	return setupList
}

func GetSetupByID(id string) (*setupStructs.Setup, error) {
	for _, setup := range setupList {
		if setup.ID == id {
			return setup, nil
		}
	}
	msg := fmt.Sprintf("Setup with ID: \"%s\" not found", id)
	logger.Warning(msg)
	return nil, tools.AppError{
		Message:     msg,
		ParentError: nil,
		Code:        404,
	}
}
