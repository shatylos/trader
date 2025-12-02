package setup

import (
	"github.com/shatylos/trader/internal/config"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/internal/strategy"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/tgNotifier"
)

var setupList []setupStructs.Setup

func TraderInit() error {
	setupList = make([]setupStructs.Setup, 0)

	appConfig, err := config.GetConfig()
	if err != nil {
		err = apperrors.Wrap(err, "error get config")
		return err
	}

	tgNotifierUrl := appConfig.App["tg_notifier_url"]
	tgNotifier.Init(tgNotifierUrl)

	setups, err := prepareSetupStrategies(appConfig.StrategyItems, appConfig.DomainItems)
	if err != nil {
		err = apperrors.Wrap(err, "unable to init setup")
		return err
	}

	setupList = setups

	return nil
}

func prepareSetupStrategies(strategyItems []interface{}, domainItems map[string]interface{}) ([]setupStructs.Setup, error) {

	setups := make([]setupStructs.Setup, 0)

	for _, strategyItemConfig := range strategyItems {
		itemMap, ok := strategyItemConfig.(map[interface{}]interface{})
		if !ok {
			return nil, apperrors.New("Can not parse a strategy config.")
		}
		if itemMap["code"] == nil {
			return nil, apperrors.New("The strategy config must contain a code field")
		}
		if itemMap["id"] == nil {
			return nil, apperrors.New("The strategy config must contain a id field")
		}

		strategyItem, err := strategy.GetStrategyByCode(itemMap["code"].(string))
		if err != nil {
			err = apperrors.Wrap(err, "error get strategy by code %s", itemMap["code"])
			return nil, err
		}
		err = strategyItem.SetConfig(strategyItemConfig, domainItems)
		if err != nil {
			err = apperrors.Wrap(err, "error in setup config \"%s\"", itemMap["id"])
			return nil, err
		}

		setups = append(setups, setupStructs.Setup{
			Strategy: strategyItem,
			ID:       itemMap["id"].(string),
		})
	}

	return setups, nil
}

func GetSetupList() *[]setupStructs.Setup {
	return &setupList
}

func GetSetupByID(id string) (*setupStructs.Setup, error) {
	for i, setup := range setupList {
		if setup.ID == id {
			return &setupList[i], nil
		}
	}
	return nil, apperrors.New("Setup with ID: \"%s\" not found", id)
}
