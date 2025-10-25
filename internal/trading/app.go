package trading

import (
	"github.com/shatylos/trader/internal/setup"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	"sync"
)

func StartTradingApp(initWg *sync.WaitGroup, mux *http.ServeMux) {
	logger.Info("Starting trading app")

	setupList := setup.GetSetupList()
	setupChanel := make(chan *setupStructs.Setup, len(*setupList))

	for i := range *setupList {
		setupItem := &(*setupList)[i]
		setupItem.Init(mux)
		setupChanel <- setupItem
	}
	initWg.Done()

	for setupItem := range setupChanel {
		go setupItem.NextStep(setupChanel)
	}
}
