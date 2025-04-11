package main

import (
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web"
	"os"
)

func main() {
	err := setup.TraderInit()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	go trading.StartTradingApp()
	web.StartWebApp()
}
