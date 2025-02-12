package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web/helper"
	"regexp"
	"strconv"
	"time"
)

type ReportRequest struct {
	Index  uint   `uri:"index" binding:"required"`
	Period string `uri:"period" binding:"required"`
}

type ReportPage struct {
	StrategyReport _struct.Report
}

type SetupListPage struct {
	SetupItems []SetupItemPage
}

type SetupItemPage struct {
	SeqNum        int
	Title         string
	CurrentPeriod string
}

func SetupListController(c *gin.Context) {

	var setupList []SetupItemPage

	for i, setupItem := range setup.GetSetupList() {
		setupList = append(setupList, SetupItemPage{
			SeqNum:        i + 1,
			Title:         setupItem.Strategy.GetTitle(),
			CurrentPeriod: time.Now().Format("2006-01"),
		})
	}

	data := SetupListPage{
		SetupItems: setupList,
	}

	template, err := helper.GetTemplate("web/template/setupList.html")
	if err != nil {
		logger.Error(err.Error())
		return
	}
	err = template.Execute(c.Writer, data)
	if err != nil {
		logger.Error(err.Error())
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

	period := reportRequest.Period
	matches := regexp.MustCompile("^(\\d{4})-(\\d{2})$").FindStringSubmatch(period)
	if len(matches) != 3 {
		c.String(404, "Not found")
		return
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		c.String(404, "Not found")
		return
	}
	month, err := strconv.Atoi(matches[2])
	if err != nil {
		c.String(404, "Not found")
		return
	}
	now := time.Now()
	firstDayOfNextMonth := time.Date(year, time.Month(month+1), 1, 23, 59, 59, 0, now.Location())
	to := firstDayOfNextMonth.AddDate(0, 0, -1)
	from := time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, now.Location())

	strategyReport, err := strategy.GetReport(from, to)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	data := ReportPage{
		StrategyReport: *strategyReport,
	}

	tmpl, err := helper.GetTemplate("web/template/report.html")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
