package main

import (
	"bitbucket.org/shatylos/trader/webapi"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	webapi.StartWebApp()
}
