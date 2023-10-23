package trading

import (
	"bitbucket.org/shatylos/trader/setup"
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
	"bitbucket.org/shatylos/trader/utils"
	"os"
)

func StartTradingApp() {
	err := setup.SetupListInit()
	if err != nil {
		utils.LogError(err.Error())
		os.Exit(1)
	}

	bufferChanel := make(chan bool, 1)
	setupChanel := make(chan *setupStructs.Setup)

	go handleSetupNextStep(setupChanel)

	for {
		bufferChanel <- true
		go loadSetupNextStep(bufferChanel, setupChanel)
	}
}

func loadSetupNextStep(bufferChanel chan bool, setupChanelContext chan *setupStructs.Setup) {
	setup.LoadNextSetupStep(setupChanelContext)
	<-bufferChanel
}

func handleSetupNextStep(setupChanel chan *setupStructs.Setup) {
	for setupItem := range setupChanel {
		go setupItem.NextStep()
	}
}
