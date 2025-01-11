package main

import (
	"github.com/joho/godotenv"
	"github.com/shatylos/trader/storage"
	"github.com/shatylos/trader/trading"
	"github.com/shatylos/trader/webapi"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

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
