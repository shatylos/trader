package fibonacci

import (
	"fmt"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type ReportOrderItem struct {
}

type TemplateData struct {
	PrevPeriodLink string
	NextPeriodLink string
	DateFrom       time.Time
	DateTo         time.Time
	Positions      []structs.Position
	MainCurrency   string
	TradeCurrency  string
	PricePrecision int
}

func (f *Fibonacci) GetReport(from time.Time, to time.Time) (report _struct.Report, err error) {
	tmpl, err := helper.GetTemplate("web/template/fibonacci/report.html")
	if err != nil {
		return
	}

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	var positions []structs.Position
	positions, err = storage.GetPositions(from, to)

	data := TemplateData{
		PrevPeriodLink: fmt.Sprintf("/report/%s/%s/", f.GetId(), from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink: fmt.Sprintf("/report/%s/%s/", f.GetId(), from.AddDate(0, 1, 0).Format("2006-01")),
		DateFrom:       from,
		DateTo:         to,
		Positions:      positions,
		MainCurrency:   f.config.MainCurrency,
		TradeCurrency:  f.config.TradeCurrency,
		PricePrecision: int(f.config.PricePrecision),
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		return
	}

	htmlStr := resultBuilder.String()

	report = _struct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   f.GetId(),
	}
	return
}
