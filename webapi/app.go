package webapi

import (
	"bitbucket.org/shatylos/trader/webapi/controller"
	"github.com/gin-gonic/gin"
	"os"
)

const defaultPort = "8080"

var router *gin.Engine

func init() {
	router = gin.Default()
}

func StartWebApp() {
	router.POST("/query", controller.GraphqlHandler())
	router.GET("/", controller.PlaygroundHandler())
	port := os.Getenv("TRADER_PORT")
	if port == "" {
		port = defaultPort
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
