package main

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	err := setup.TraderInit()
	if err != nil {
		logger.PrintError(err)
		os.Exit(1)
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	var wg, shutDownWg sync.WaitGroup
	wg.Add(1)

	mux := http.NewServeMux()
	go trading.StartTradingApp(shutdownCtx, &wg, &shutDownWg, mux)

	wg.Wait()
	go web.StartWebApp(mux)

	for {
		select {
		case <-shutdownCtx.Done():
			logger.Info(fmt.Sprintf("shutdown requested"))
			shutDownWg.Wait()
			stop()
			logger.Info(fmt.Sprintf("trading app stopped"))
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
