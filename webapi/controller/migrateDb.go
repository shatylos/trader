package controller

import (
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/mongo"
	"bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/sqlite"
	"fmt"
	"github.com/gin-gonic/gin"
)

func MigrateDb(c *gin.Context) {

	sqliteStorage, err := sqlite.New("buyCheapSellHigh.db", "BTC_USDT_1")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	docuementStorage, err := mongo.New("BTC_USDT_1")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	oldOrders, err := sqliteStorage.GetNotCalculatedDomainOrders()
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	for _, order := range oldOrders {

		order.AveragePrice = nil
		order.FilledPrice = nil
		order.FilledQty = nil
		order.Side = nil
		order.UpdatedTime = nil

		isAdded, err := docuementStorage.AddDomainOrderOnce(order)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			continue
		}
		if isAdded {
			fmt.Println("Added: ", order.DomainOrderId)
		}
	}
}
