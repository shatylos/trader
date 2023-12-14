package controller

import (
	"bitbucket.org/shatylos/trader/setup"
	"bitbucket.org/shatylos/trader/strategy/struct"
	"bitbucket.org/shatylos/trader/utils"
	"bitbucket.org/shatylos/trader/webapi/helper"
	"github.com/gin-gonic/gin"
	"time"
)

type ReportPage struct {
	StrategyReport _struct.Report
}

func ResetController(c *gin.Context) {
	setupList := setup.GetSetupList()
	setupItem := setupList[0]
	strategy := setupItem.Strategy
	strategy.ResetOrderData()
}

func ReportController(c *gin.Context) {

	setupList := setup.GetSetupList()
	//if !slices.Contains(setupList, 70) {
	//	c.Error(utils.AppError{
	//		Message: "Can not find a setup item.",
	//	})
	//	return
	//}
	setupItem := setupList[0]
	strategy := setupItem.Strategy
	println(strategy)

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	strategyReport, err := strategy.GetReport(from, now)
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	data := ReportPage{
		StrategyReport: *strategyReport,
	}

	tmpl, err := helper.GetTemplate("webapi/templates/report.html")
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
}
