package investor

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	strategyStruct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type ReportTemplateData struct {
	AvailableMain         float64
	AvailableTrade        float64
	AvailableTotalMain    float64
	PrevPeriodLink        string
	NextPeriodLink        string
	WsLink                string
	DateFrom              time.Time
	DateTo                time.Time
	DealsRelations        []*entity.DealRelation //@TODO: remove * from the list
	MainCurrency          string
	TradeCurrency         string
	PricePrecision        int
	QtyPrecision          int
	RealizedPNL           float64
	RealizedPNLPercent    float64
	UnrealizedPNL         float64
	UnrealizedPNLPercent  float64
	BalanceMainBefore     float64
	BalanceMainBeforeEqv  float64
	BalanceTradeBefore    float64
	BalanceTradeBeforeEqv float64
	BalanceTotalBefore    float64
	BalanceTotalAfter     float64
	BalanceMainAfter      float64
	BalanceMainAfterEqv   float64
	BalanceTradeAfter     float64
	BalanceTradeAfterEqv  float64
	DepoAmount            float64
	WithdrawAmount        float64
	Timeframes            []_struct.TimeframeItem
	CurrentPrice          float64
	IsCurrentPeriod       bool
}

type ReportAmountsData struct {
	RealizedPnl           float64
	RealizedPnlPercent    float64
	UnrealizedPnl         float64
	UnrealizedPnlPercent  float64
	BalanceTotalBefore    float64
	BalanceMainBefore     float64
	BalanceTradeBefore    float64
	BalanceMainBeforeEqv  float64
	BalanceTradeBeforeEqv float64
	BalanceTotalAfter     float64
	BalanceMainAfter      float64
	BalanceTradeAfter     float64
	BalanceMainAfterEqv   float64
	BalanceTradeAfterEqv  float64
	DepoAmount            float64
	WithdrawAmount        float64
}

func (i *Investor) GetReport(from time.Time, to time.Time) (report strategyStruct.Report, err error) {

	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		err = apperrors.Wrap(err, "error get template")
		return
	}

	var data ReportTemplateData
	data, err = i.GetReportData(from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get report data")
		return
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		err = apperrors.Wrap(err, "error execute template")
		return
	}

	htmlStr := resultBuilder.String()

	report = strategyStruct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   i.GetId(),
	}
	return
}

func (i *Investor) GetReportData(from time.Time, to time.Time) (data ReportTemplateData, err error) {
	ctx := i.getContext()
	var deals, closedDeals, activeDeals []*entity.Deal
	closedDeals, err = i.Storage.GetDealsByPeriod(ctx, from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get deals by period. From %s, to %s", from, to)
		return
	}

	now := time.Now()
	if from.Year() == now.Year() && from.Month() == now.Month() {
		activeDeals, err = i.Storage.GetActiveDeals(ctx)
		if err != nil {
			err = apperrors.Wrap(err, "error get active deals")
			return
		}
		deals = append(deals, activeDeals...)
	}
	deals = append(deals, closedDeals...)

	var dealRelations []*entity.DealRelation
	for _, deal := range deals {
		var dealRelation *entity.DealRelation
		dealRelation, err = i.Storage.GetDealRelation(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error get deal relation")
			return
		}
		dealRelations = append(dealRelations, dealRelation)
	}

	var assets []*structs.AssetTransaction
	assets, err = i.Storage.GetAssetTransactions(ctx, from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get asset transactions")
	}

	amounts := i.reportAmounts(from, dealRelations, assets)

	isCurrentPeriod := false
	if from.Before(time.Now()) && to.After(time.Now()) {
		isCurrentPeriod = true
	}

	availableMain := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
	availableTrade := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)

	data.AvailableMain = availableMain
	data.AvailableTrade = availableTrade
	data.AvailableTotalMain = availableMain + math.Mul(availableTrade, i.State.CurrentPrice)
	data.PrevPeriodLink = fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 0, -1).Format("2006-01"))
	data.NextPeriodLink = fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 1, 0).Format("2006-01"))
	data.WsLink = fmt.Sprintf("/%s/ws-report", i.GetId())
	data.DateFrom = from
	data.DateTo = to
	data.DealsRelations = dealRelations
	data.MainCurrency = i.Config.MainCurrency
	data.TradeCurrency = i.Config.TradeCurrency
	data.PricePrecision = int(i.Config.PricePrecision)
	data.QtyPrecision = int(i.Config.QtyPrecision)
	data.RealizedPNL = amounts.RealizedPnl
	data.RealizedPNLPercent = amounts.RealizedPnlPercent
	data.UnrealizedPNL = amounts.UnrealizedPnl
	data.UnrealizedPNLPercent = amounts.UnrealizedPnlPercent
	data.BalanceMainBefore = amounts.BalanceMainBefore
	data.BalanceMainBeforeEqv = amounts.BalanceMainBeforeEqv
	data.BalanceTradeBefore = amounts.BalanceTradeBefore
	data.BalanceTradeBeforeEqv = amounts.BalanceTradeBeforeEqv
	data.BalanceTotalBefore = amounts.BalanceTotalBefore
	data.BalanceTotalAfter = amounts.BalanceTotalAfter
	data.BalanceMainAfter = amounts.BalanceMainAfter
	data.BalanceMainAfterEqv = amounts.BalanceMainAfterEqv
	data.BalanceTradeAfter = amounts.BalanceTradeAfter
	data.BalanceTradeAfterEqv = amounts.BalanceTradeAfterEqv
	data.DepoAmount = amounts.DepoAmount
	data.WithdrawAmount = amounts.WithdrawAmount
	data.Timeframes = i.Timeframes
	data.CurrentPrice = i.State.CurrentPrice
	data.IsCurrentPeriod = isCurrentPeriod

	return
}

func (i *Investor) reportAmounts(from time.Time, dealRelations []*entity.DealRelation, assets []*structs.AssetTransaction) (amounts ReportAmountsData) {

	var firstOrder, lastOrder *entity.Order

	for _, dealRelation := range dealRelations {
		amounts.RealizedPnl += dealRelation.RealizedPNL
		amounts.UnrealizedPnl += dealRelation.UnrealizedPNL

		for _, order := range dealRelation.Orders {
			if (lastOrder == nil || order.CreatedTime.After(lastOrder.CreatedTime)) &&
				order.CreatedTime.Year() == from.Year() && order.CreatedTime.Month() == from.Month() {

				lastOrder = order
			}
			if (firstOrder == nil || order.CreatedTime.Before(firstOrder.CreatedTime)) &&
				order.CreatedTime.Year() == from.Year() && order.CreatedTime.Month() == from.Month() {
				firstOrder = order
			}
		}
	}
	if firstOrder == nil || lastOrder == nil {
		return
	}

	amounts.BalanceMainBefore = trading.CurrencyAmountTotal(&firstOrder.WalletBefore, i.Config.MainCurrency)
	amounts.BalanceMainAfter = trading.CurrencyAmountTotal(&lastOrder.WalletAfter, i.Config.MainCurrency)
	amounts.BalanceTradeBefore = trading.CurrencyAmountTotal(&firstOrder.WalletBefore, i.Config.TradeCurrency)
	amounts.BalanceTradeAfter = trading.CurrencyAmountTotal(&lastOrder.WalletAfter, i.Config.TradeCurrency)

	amounts.BalanceTotalBefore = amounts.BalanceMainBefore + trading.TradeCurrencyToMain(amounts.BalanceTradeBefore, firstOrder.Price)
	lastPrice := lastOrder.Price
	now := time.Now()
	if lastOrder.CreatedTime.Year() == now.Year() && lastOrder.CreatedTime.Month() == now.Month() {
		lastPrice = i.State.CurrentPrice
	}
	amounts.BalanceTotalAfter = amounts.BalanceMainAfter + trading.TradeCurrencyToMain(amounts.BalanceTradeAfter, lastPrice)

	amounts.BalanceMainBeforeEqv = math.Div(amounts.BalanceMainBefore, firstOrder.Price)
	amounts.BalanceTradeBeforeEqv = math.Mul(amounts.BalanceTradeBefore, firstOrder.Price)
	amounts.BalanceMainAfterEqv = math.Div(amounts.BalanceMainAfter, lastPrice)
	amounts.BalanceTradeAfterEqv = math.Mul(amounts.BalanceTradeAfter, lastPrice)

	if amounts.BalanceTotalBefore != 0 {
		amounts.RealizedPnlPercent = math.Div(amounts.RealizedPnl, math.Div(amounts.BalanceTotalBefore, 100))
		amounts.UnrealizedPnlPercent = math.Div(amounts.UnrealizedPnl, math.Div(amounts.BalanceTotalBefore, 100))
	}

	for _, asset := range assets {
		switch asset.TransactionType {
		case structs.TransactionTypeDeposit:
			amounts.DepoAmount += asset.Amount
			if asset.CreatedTime.Before(firstOrder.CreatedTime) {
				amounts.BalanceMainBefore -= asset.Amount
				amounts.BalanceTotalBefore -= asset.Amount
			}
			if asset.CreatedTime.After(lastOrder.CreatedTime) {
				amounts.BalanceMainAfter += asset.Amount
				amounts.BalanceTotalAfter += asset.Amount
			}
		case structs.TransactionTypeWithdraw:
			amounts.WithdrawAmount += asset.Amount
			if asset.CreatedTime.Before(firstOrder.CreatedTime) {
				amounts.BalanceMainBefore += asset.Amount
				amounts.BalanceTotalBefore += asset.Amount
			}
			if asset.CreatedTime.After(lastOrder.CreatedTime) {
				amounts.BalanceMainAfter -= asset.Amount
				amounts.BalanceTotalAfter -= asset.Amount
			}
		}
	}
	return
}
