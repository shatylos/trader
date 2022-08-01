package main

import (
	"bitbucket.org/shatylos/trader/trading"
	"bitbucket.org/shatylos/trader/webapi"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	go trading.StartTradingApp()
	webapi.StartWebApp()
}
