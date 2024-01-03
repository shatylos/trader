package controller

import (
	"bitbucket.org/shatylos/trader/setup"
	"bitbucket.org/shatylos/trader/strategy/struct"
	"bitbucket.org/shatylos/trader/utils"
	"bitbucket.org/shatylos/trader/webapi/helper"
	"github.com/gin-gonic/gin"
	"time"
)

type ReportRequest struct {
	Index uint `uri:"index" binding:"required"`
}

type ReportPage struct {
	StrategyReport _struct.Report
}

type SetupListPage struct {
	SetupItems []SetupItemPage
}

type SetupItemPage struct {
	SeqNum int
	Title  string
}

func SetupListController(c *gin.Context) {

	var setupList []SetupItemPage

	for i, setupItem := range setup.GetSetupList() {
		setupList = append(setupList, SetupItemPage{
			SeqNum: i + 1,
			Title:  setupItem.Strategy.GetTitle(),
		})
	}

	data := SetupListPage{
		SetupItems: setupList,
	}

	template, err := helper.GetTemplate("templates/setupList.html")
	if err != nil {
		utils.LogError(err.Error())
		return
	}
	err = template.Execute(c.Writer, data)
	if err != nil {
		utils.LogError(err.Error())
		return
	}
}

func ReportController(c *gin.Context) {

	var reportRequest ReportRequest
	if err := c.ShouldBindUri(&reportRequest); err != nil {
		c.String(404, "Not found")
		return
	}

	setupList := setup.GetSetupList()
	if uint(len(setupList)) < reportRequest.Index {
		c.String(404, "Not found")
		return
	}
	setupItem := setupList[reportRequest.Index-1]
	strategy := setupItem.Strategy

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

	tmpl, err := helper.GetTemplate("templates/report.html")
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
