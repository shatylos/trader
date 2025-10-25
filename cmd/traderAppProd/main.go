package main

import (
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web"
	"net/http"
	"os"
	"sync"
)

func main() {
	err := setup.TraderInit()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	mux := http.NewServeMux()
	go trading.StartTradingApp(&wg, mux)

	wg.Wait()
	web.StartWebApp(mux)
}
