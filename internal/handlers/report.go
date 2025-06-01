package handlers

import (
	"fmt"
	"github.com/shatylos/trader/internal/setup"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web/helper"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type ReportPage struct {
	StrategyReport _struct.Report
}

type SetupListPage struct {
	SetupItems []SetupItemPage
}

type SetupItemPage struct {
	Title         string
	SetupId       string
	CurrentPeriod string
}

func SetupListController(w http.ResponseWriter, r *http.Request) {

	var setupList []SetupItemPage

	for _, setupItem := range setup.GetSetupList() {
		setupList = append(setupList, SetupItemPage{
			Title:         setupItem.Strategy.GetTitle(),
			SetupId:       setupItem.Strategy.GetId(),
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
	err = template.Execute(w, data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	setupId := r.PathValue("setup_id")
	now := time.Now()
	redirectTo := fmt.Sprintf("/report/%s/%s/", setupId, now.Format("2006-01"))

	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func ReportPeriodHandler(w http.ResponseWriter, r *http.Request) {

	setupId := r.PathValue("setup_id")
	setupList := setup.GetSetupList()

	var strategy _struct.StrategyInterface
	for _, setupItem := range setupList {
		if setupItem.Strategy.GetId() == setupId {
			strategy = setupItem.Strategy
		}
	}
	if strategy == nil {
		http.NotFound(w, r)
		return
	}

	period := r.PathValue("period")
	matches := regexp.MustCompile("^(\\d{4})-(\\d{2})$").FindStringSubmatch(period)
	if len(matches) != 3 {
		http.NotFound(w, r)
		return
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	month, err := strconv.Atoi(matches[2])
	if err != nil {
		http.NotFound(w, r)
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
		StrategyReport: strategyReport,
	}

	tmpl, err := helper.GetTemplate("web/template/report.html")
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
