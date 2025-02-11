package main

import (
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/web"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	go trading.StartTradingApp()
	web.StartWebApp()
}
