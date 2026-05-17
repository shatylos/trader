package trading

import (
	"github.com/shatylos/trader/internal/setup"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	"sync"
	"time"
)

func StartTradingApp(initWg *sync.WaitGroup, mux *http.ServeMux) {
	logger.Info("Starting trading app")

	setupList := setup.GetSetupList()
	setupChanel := make(chan *setupStructs.Setup, len(*setupList))
	setupDelayChanel := make(chan *setupStructs.SetupDelay, len(*setupList))

	for i := range *setupList {
		setupItem := &(*setupList)[i]
		setupItem.Init(mux)
		setupChanel <- setupItem
	}
	initWg.Done()

	go func() {
		for setupDelay := range setupDelayChanel {
			go func() {
				time.Sleep(setupDelay.Duration)
				setupChanel <- setupDelay.Setup
			}()
		}
	}()

	for setupItem := range setupChanel {
		go setupItem.NextStep(setupDelayChanel)
	}
}
