package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shatylos/trader/internal/config"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web/controller"
)

const defaultPort = "8080"

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartWebApp() {
	router.GET("/", controller.SetupListController)
	router.GET("/report/:index/:period", controller.ReportController)

	port := defaultPort
	appConfig, err := config.GetConfig()
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting app config %s", err.Error()))
	} else if appConfig.App["web_port"] != "" {
		port = appConfig.App["web_port"]
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
