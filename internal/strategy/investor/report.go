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
	Orders                map[string][]*entity.Order        //@TODO: remove * from the list
	TimeframeStates       map[string]*entity.TimeframeState //@TODO: remove * from the list
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

	var orders []*entity.Order
	orders, err = i.Storage.GetOrdersByPeriod(ctx, from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get orders by period. From %s, to %s", from, to)
		return
	}

	ordersByTimeframe := make(map[string][]*entity.Order)
	for _, order := range orders {
		ordersByTimeframe[order.Timeframe] = append(ordersByTimeframe[order.Timeframe], order)
	}

	timeframeStates := make(map[string]*entity.TimeframeState)
	for key := range i.Timeframes {
		var state *entity.TimeframeState
		state, err = i.getTimeframeState(ctx, &(i.Timeframes[key]))
		if err != nil {
			err = apperrors.Wrap(err, "error get timeframe state")
			return
		}
		timeframeStates[i.Timeframes[key].Config.Resolution] = state
	}

	var assets []*structs.AssetTransaction
	assets, err = i.Storage.GetAssetTransactions(ctx, from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get asset transactions")
	}

	isCurrentPeriod := false
	if from.Before(time.Now()) && to.After(time.Now()) {
		isCurrentPeriod = true
	}

	amounts := i.reportAmounts(orders, timeframeStates, assets, isCurrentPeriod)

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
	data.Orders = ordersByTimeframe
	data.TimeframeStates = timeframeStates
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

func (i *Investor) reportAmounts(orders []*entity.Order, states map[string]*entity.TimeframeState, assets []*structs.AssetTransaction, isCurrentPeriod bool) (amounts ReportAmountsData) {

	var firstOrder, lastOrder *entity.Order
	lastFilledByTimeframe := make(map[string]*entity.Order)

	// orders are sorted by CreatedTime ascending
	for _, order := range orders {
		if order.OrderStatus != structs.OrderStatuses.Filled {
			continue
		}
		if firstOrder == nil {
			firstOrder = order
		}
		lastOrder = order
		lastFilledByTimeframe[order.Timeframe] = order

		if order.Side == structs.OrderSideSell {
			amounts.RealizedPnl += order.RealizedPNL
		}
	}

	if isCurrentPeriod {
		for _, state := range states {
			if state.QtyInTrade > 0 {
				amounts.UnrealizedPnl += math.Mul(state.QtyInTrade, i.State.CurrentPrice) - math.Mul(state.QtyInTrade, state.AverageBuyPrice)
			}
		}
	} else {
		for _, order := range lastFilledByTimeframe {
			if order.QtyInTrade > 0 {
				amounts.UnrealizedPnl += math.Mul(order.QtyInTrade, order.Price) - math.Mul(order.QtyInTrade, order.AverageBuyPrice)
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
	if isCurrentPeriod {
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
