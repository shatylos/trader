package main

import (
	"bitbucket.org/shatylos/trader/storage"
	"bitbucket.org/shatylos/trader/trading"
	"bitbucket.org/shatylos/trader/webapi"
	"github.com/joho/godotenv"
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
