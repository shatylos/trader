package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/shatylos/trader/utils"
	"github.com/shatylos/trader/webapi/controller"
)

const defaultPort = "8080"

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartWebApp() {
	router.GET("/", controller.SetupListController)
	router.GET("/report/:index/:period", controller.ReportController)
	port := utils.AppConfig("TRADER_WEB_PORT", defaultPort)

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
