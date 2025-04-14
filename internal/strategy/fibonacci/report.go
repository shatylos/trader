package fibonacci

import (
	"fmt"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type ReportOrderItem struct {
}

type TemplateData struct {
	PrevPeriodLink  string
	NextPeriodLink  string
	DateFrom        time.Time
	DateTo          time.Time
	Positions       []structs.Position
	MainCurrency    string
	TradeCurrency   string
	PricePrecision  int
	TotalPnl        float64
	TotalPnlPercent float64
	BalanceBefore   float64
	BalanceAfter    float64
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
	totalPnl, totalPnlPercent, balanceBefore, balanceAfter := f.reportAmounts(positions)

	data := TemplateData{
		PrevPeriodLink:  fmt.Sprintf("/report/%s/%s/", f.GetId(), from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink:  fmt.Sprintf("/report/%s/%s/", f.GetId(), from.AddDate(0, 1, 0).Format("2006-01")),
		DateFrom:        from,
		DateTo:          to,
		Positions:       positions,
		MainCurrency:    f.config.MainCurrency,
		TradeCurrency:   f.config.TradeCurrency,
		PricePrecision:  int(f.config.PricePrecision),
		TotalPnl:        totalPnl,
		TotalPnlPercent: totalPnlPercent,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
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

func (f *Fibonacci) reportAmounts(positions []structs.Position) (totalPnl float64, totalPnlPercent float64, balanceBefore float64, balanceAfter float64) {
	filteredPositions := []structs.Position{}
	for _, position := range positions {
		if position.Status == structs.StatusClosed {
			filteredPositions = append(filteredPositions, position)
			totalPnl += position.TotalClosePnl
		}
	}
	if len(filteredPositions) > 0 {
		balanceBefore = filteredPositions[len(filteredPositions)-1].BalanceBefore
		balanceAfter = filteredPositions[0].BalanceAfter
	}
	if totalPnl > 0 && balanceBefore > 0 && balanceAfter > 0 {
		onePercent := math.Div(balanceBefore, 100)
		amountDiff := balanceAfter - balanceBefore
		totalPnlPercent = math.Div(amountDiff, onePercent)
	}
	return
}
