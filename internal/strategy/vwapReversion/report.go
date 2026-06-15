package vwapReversion

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	strategyStorage "github.com/shatylos/trader/internal/strategy/vwapReversion/storage"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

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
	State           State
}

func (v *VwapReversion) GetReport(from time.Time, to time.Time) (report _struct.Report, err error) {
	tmpl, err := helper.GetTemplate("web/template/vwapReversion/report.html")
	if err != nil {
		err = apperrors.Wrap(err, "error get template")
		return
	}

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(v.config.Id)
	if err != nil {
		err = apperrors.Wrap(err, "error get storage")
		return
	}

	var positions []structs.Position
	positions, err = storage.GetPositions(from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get positions")
		return
	}

	var assets []structs.AssetTransaction
	assets, err = storage.GetAssetTransactions(from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get asset transactions")
		return
	}
	totalPnl, totalPnlPercent, balanceBefore, balanceAfter, depoAmount, withdrawAmount := v.reportAmounts(positions, assets)

	data := ReportTemplateData{
		PrevPeriodLink:  fmt.Sprintf("/report/%s/%s/", v.GetId(), from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink:  fmt.Sprintf("/report/%s/%s/", v.GetId(), from.AddDate(0, 1, 0).Format("2006-01")),
		DateFrom:        from,
		DateTo:          to,
		Positions:       positions,
		MainCurrency:    v.config.MainCurrency,
		TradeCurrency:   v.config.TradeCurrency,
		PricePrecision:  int(v.config.PricePrecision),
		TotalPnl:        totalPnl,
		TotalPnlPercent: totalPnlPercent,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
		DepoAmount:      depoAmount,
		WithdrawAmount:  withdrawAmount,
		State:           v.state,
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		err = apperrors.Wrap(err, "error execute template")
		return
	}

	htmlStr := resultBuilder.String()

	report = _struct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   v.GetId(),
	}
	return
}

func (v *VwapReversion) reportAmounts(positions []structs.Position, assets []structs.AssetTransaction) (totalPnl float64, totalPnlPercent float64, balanceBefore float64, balanceAfter float64, depoAmount float64, withdrawAmount float64) {
	filteredPositions := []structs.Position{}
	for _, position := range positions {
		if position.Status == structs.StatusClosed {
			filteredPositions = append(filteredPositions, position)
		}
	}
	var firstPositionCreatedTime time.Time
	if len(filteredPositions) > 0 {
		firstPositionCreatedTime = filteredPositions[len(filteredPositions)-1].CreatedTime
		balanceBefore = filteredPositions[len(filteredPositions)-1].BalanceBefore
		balanceAfter = filteredPositions[0].BalanceAfter
	}

	var depoAmountBeforeStartTrade float64
	var withdrawAmountBeforeStartTrade float64
	for _, asset := range assets {
		switch asset.TransactionType {
		case domainStructs.TransactionTypeDeposit:
			depoAmount += asset.Amount
			if firstPositionCreatedTime.After(asset.CreatedTime) {
				depoAmountBeforeStartTrade += asset.Amount
				balanceBefore -= asset.Amount
			}
		case domainStructs.TransactionTypeWithdraw:
			withdrawAmount += asset.Amount
			if firstPositionCreatedTime.After(asset.CreatedTime) {
				withdrawAmountBeforeStartTrade += asset.Amount
				balanceBefore += asset.Amount
			}
		}
	}
	totalPnl = balanceAfter - balanceBefore

	if totalPnl > 0 {
		totalPnl += withdrawAmount
		totalPnl -= depoAmount
	}

	startTradeBalance := balanceBefore - withdrawAmountBeforeStartTrade + depoAmountBeforeStartTrade
	if startTradeBalance > 0 {
		onePercent := math.Div(startTradeBalance, 100)
		if totalPnl != 0.0 {
			totalPnlPercent = math.Div(totalPnl, onePercent)
		}
	}
	return
}
