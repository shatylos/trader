package main

import (
	"github.com/shatylos/trader/internal/trading"
	"github.com/shatylos/trader/web"
)

func main() {
	go trading.StartTradingApp()
	web.StartWebApp()
}
