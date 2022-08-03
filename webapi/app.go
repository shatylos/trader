package webapi

import (
	"bitbucket.org/shatylos/trader/utils"
	"bitbucket.org/shatylos/trader/webapi/controller"
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartWebApp() {
	router.POST("/query", controller.GraphqlHandler())
	router.GET("/", controller.PlaygroundHandler())
	port := utils.AppConfig("TRADER_WEB_PORT", defaultPort)

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
