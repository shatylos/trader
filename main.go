package main

import (
	"github.com/joho/godotenv"
	"github.com/shatylos/trader/storage"
	"github.com/shatylos/trader/trading"
	"github.com/shatylos/trader/webapi"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error load go dot env: " + err.Error())
		return
	}

	err = storage.InitStorage()
	if err != nil {
		panic("Error Init storage: " + err.Error())
		return
	}

	go trading.StartTradingApp()
	webapi.StartWebApp()
}
