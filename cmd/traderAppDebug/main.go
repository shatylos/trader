package main

import (
	"fmt"
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	go func() {
		logger.Info("Starting pprof listen: http://127.0.0.1:6060/debug/pprof/")
		err := http.ListenAndServe(":6060", nil)
		logger.Error(fmt.Sprintf("Error listening pprof %s", err.Error()))
	}()

	err := setup.TraderInit()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	go trading.StartTradingApp()
	web.StartWebApp()
}
