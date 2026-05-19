package trading

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/setup"
	setupStructs "github.com/shatylos/trader/internal/setup/structs"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	"sync"
	"time"
)

func StartTradingApp(shutdownCtx context.Context, initWg, shutDownWg *sync.WaitGroup, mux *http.ServeMux) {
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
				select {
				case <-shutdownCtx.Done():
					logger.Info(fmt.Sprintf("shutdown requested. Will not start next step for setup: %s", setupDelay.Setup.ID))
					return
				default:
					setupChanel <- setupDelay.Setup
				}
			}()
		}
	}()

	go func() {
		for setupItem := range setupChanel {
			shutDownWg.Add(1)
			go setupItem.NextStep(setupDelayChanel, shutDownWg)
		}
	}()
}
