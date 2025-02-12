package main

import (
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/web"
	_ "net/http/pprof"
)

func main() {
	go trading.StartTradingApp()
	web.StartWebApp()
}
