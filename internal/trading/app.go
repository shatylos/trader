package trading

import (
	"github.com/shatylos/trader/internal/setup"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/tools/logger"
)

func StartTradingApp() {
	logger.Info("Starting trading app")

	setupList := setup.GetSetupList()
	setupChanel := make(chan *setupStructs.Setup, len(setupList))

	for _, setupItem := range setupList {
		setupChanel <- setupItem
	}

	for setupItem := range setupChanel {
		go setupItem.NextStep(setupChanel)
	}
}
