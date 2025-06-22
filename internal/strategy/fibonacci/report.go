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

type ReportTemplateData struct {
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
	DepoAmount      float64
	WithdrawAmount  float64
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

	var assets []structs.AssetTransaction
	assets, err = storage.GetAssetTransactions(from, to)
	totalPnl, totalPnlPercent, balanceBefore, balanceAfter, depoAmount, withdrawAmount := f.reportAmounts(positions, assets)

	data := ReportTemplateData{
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
		DepoAmount:      depoAmount,
		WithdrawAmount:  withdrawAmount,
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

func (f *Fibonacci) reportAmounts(positions []structs.Position, assets []structs.AssetTransaction) (totalPnl float64, totalPnlPercent float64, balanceBefore float64, balanceAfter float64, depoAmount float64, withdrawAmount float64) {
	filteredPositions := []structs.Position{}
	for _, position := range positions {
		if position.Status == structs.StatusClosed {
			filteredPositions = append(filteredPositions, position)
		}
	}
	var firstPositionCreatedTime int64
	if len(filteredPositions) > 0 {
		firstPositionCreatedTime = filteredPositions[len(filteredPositions)-1].CreatedTime
		balanceBefore = filteredPositions[len(filteredPositions)-1].BalanceBefore
		balanceAfter = filteredPositions[0].BalanceAfter
	}
	totalPnl = balanceAfter - balanceBefore

	var depoAmountBeforeStartTrade float64
	var withdrawAmountBeforeStartTrade float64
	for _, asset := range assets {
		switch asset.TransactionType {
		case "deposit":
			depoAmount += asset.Amount
			if firstPositionCreatedTime < asset.CreatedTime {
				depoAmountBeforeStartTrade += asset.Amount
			}
			break
		case "withdraw":
			withdrawAmount += asset.Amount
			if firstPositionCreatedTime < asset.CreatedTime {
				withdrawAmountBeforeStartTrade += asset.Amount
			}
			break
		}
	}

	if totalPnl > 0 {
		totalPnl += withdrawAmountBeforeStartTrade
		totalPnl -= depoAmountBeforeStartTrade
	}

	if balanceBefore > 0 {
		onePercent := math.Div(balanceBefore, 100)
		if totalPnl != 0.0 {
			totalPnlPercent = math.Div(totalPnl, onePercent)
		}
	}
	return
}
