package trading

import (
	"bitbucket.org/shatylos/trader/setup"
	setupStructs "bitbucket.org/shatylos/trader/setup/structs"
)

func StartTradingApp() {
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
