package trading

import (
	"github.com/shatylos/trader/internal/setup"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/tools/logger"
	"os"
)

func StartTradingApp() {
	err := setup.SetupListInit()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	setupList := setup.GetSetupList()
	setupChanel := make(chan *setupStructs.Setup, len(setupList))

	for _, setupItem := range setupList {
		setupChanel <- setupItem
	}

	for setupItem := range setupChanel {
		go setupItem.NextStep(setupChanel)
	}
}
