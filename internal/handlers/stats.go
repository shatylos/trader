package handlers

import (
	"github.com/shatylos/trader/internal/setup"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/web/helper"
	"net/http"
)

type StatsPage struct {
	StrategyStats []_struct.Stats
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {

	setupList := setup.GetSetupList()

	strategyStats := make([]_struct.Stats, len(*setupList))
	data := StatsPage{
		StrategyStats: strategyStats,
	}

	for i := range *setupList {
		setupItem := &(*setupList)[i]
		var er error
		data.StrategyStats[i], er = setupItem.Strategy.GetStats()
		if er != nil {
			logger.Error(er.Error())
			return
		}
	}

	tmpl, err := helper.GetTemplate("web/template/stats.html")
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
