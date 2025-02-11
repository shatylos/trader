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

	bufferChanel := make(chan bool, 1)
	setupChanel := make(chan *setupStructs.Setup)

	go handleSetupNextStep(setupChanel)

	for {
		bufferChanel <- true
		loadSetupNextStep(bufferChanel, setupChanel)
	}
}

func loadSetupNextStep(bufferChanel chan bool, setupChanelContext chan *setupStructs.Setup) {
	setup.LoadNextSetupStep(setupChanelContext)
	<-bufferChanel
}

func handleSetupNextStep(setupChanel chan *setupStructs.Setup) {
	for setupItem := range setupChanel {
		setupItem.NextStep()
	}
}
