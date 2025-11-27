package handlers

import (
	"github.com/shatylos/trader/internal/setup"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/apperrors"
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
		var err error
		data.StrategyStats[i], err = setupItem.Strategy.GetStats()
		if err != nil {
			err = apperrors.Wrap(err, "error get stats %s", setupItem.ID)
			logger.PrintError(err)
			return
		}
	}

	tmpl, err := helper.GetTemplate("web/template/stats.html")
	if err != nil {
		err = apperrors.Wrap(err, "error get template")
		logger.PrintError(err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		err = apperrors.Wrap(err, "error execute template")
		logger.PrintError(err)
		return
	}
}
