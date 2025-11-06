package investor

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	strategyStruct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type ReportTemplateData struct {
	PrevPeriodLink     string
	NextPeriodLink     string
	WsLink             string
	DateFrom           time.Time
	DateTo             time.Time
	DealsRelations     []*entity.DealRelation //@TODO: remove * from the list
	MainCurrency       string
	TradeCurrency      string
	PricePrecision     int
	QtyPrecision       int
	TotalPnl           float64
	TotalPnlPercent    float64
	BalanceMainBefore  float64
	BalanceTradeBefore float64
	BalanceTotalBefore float64
	BalanceTotalAfter  float64
	BalanceMainAfter   float64
	BalanceTradeAfter  float64
	DepoAmount         float64
	WithdrawAmount     float64
	Timeframes         []_struct.TimeframeItem
	HeapTimeframe      _struct.HeapTimeframe
	CurrentPrice       float64
}

func (i *Investor) GetReport(from time.Time, to time.Time) (report strategyStruct.Report, err error) {

	ctx := i.getContext()

	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		return
	}

	var deals, closedDeals, activeDeals []*entity.Deal
	closedDeals, err = i.Storage.GetDealsByPeriod(ctx, from, to)
	if err != nil {
		return
	}

	now := time.Now()
	if from.Year() == now.Year() && from.Month() == now.Month() {
		activeDeals, err = i.Storage.GetActiveDeals(ctx)
		if err != nil {
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
			return
		}
		dealRelations = append(dealRelations, dealRelation)
	}

	var assets []*structs.AssetTransaction
	assets, err = i.Storage.GetAssetTransactions(ctx, from, to)
	totalPnl, totalPnlPercent,
		balanceTotalBefore, balanceMainBefore, balanceTradeBefore,
		balanceTotalAfter, balanceMainAfter, balanceTradeAfter,
		depoAmount, withdrawAmount := i.reportAmounts(from, dealRelations, assets)

	data := ReportTemplateData{
		PrevPeriodLink:     fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink:     fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 1, 0).Format("2006-01")),
		WsLink:             fmt.Sprintf("/%s/ws-report", i.GetId()),
		DateFrom:           from,
		DateTo:             to,
		DealsRelations:     dealRelations,
		MainCurrency:       i.Config.MainCurrency,
		TradeCurrency:      i.Config.TradeCurrency,
		PricePrecision:     int(i.Config.PricePrecision),
		QtyPrecision:       int(i.Config.QtyPrecision),
		TotalPnl:           totalPnl,
		TotalPnlPercent:    totalPnlPercent,
		BalanceMainBefore:  balanceMainBefore,
		BalanceTradeBefore: balanceTradeBefore,
		BalanceTotalBefore: balanceTotalBefore,
		BalanceTotalAfter:  balanceTotalAfter,
		BalanceMainAfter:   balanceMainAfter,
		BalanceTradeAfter:  balanceTradeAfter,
		DepoAmount:         depoAmount,
		WithdrawAmount:     withdrawAmount,
		Timeframes:         i.Timeframes,
		HeapTimeframe:      i.HeapTimeframe,
		CurrentPrice:       i.State.CurrentPrice,
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		return
	}

	htmlStr := resultBuilder.String()

	report = strategyStruct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   i.GetId(),
	}
	return
}

func (i *Investor) reportAmounts(from time.Time, dealRelations []*entity.DealRelation, assets []*structs.AssetTransaction) (
	totalPnl, totalPnlPercent,
	balanceTotalBefore, balanceMainBefore, balanceTradeBefore,
	balanceTotalAfter, balanceMainAfter, balanceTradeAfter,
	depoAmount, withdrawAmount float64) {

	var firstOrder, lastOrder *entity.Order

	for _, dealRelation := range dealRelations {
		if dealRelation.Deal.Status != entity.DealStatusClosed {
			continue
		}

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

	balanceMainBefore = trading.CurrencyAmountTotal(&firstOrder.WalletBefore, i.Config.MainCurrency)
	balanceMainAfter = trading.CurrencyAmountTotal(&lastOrder.WalletAfter, i.Config.MainCurrency)
	balanceTradeBefore = trading.CurrencyAmountTotal(&firstOrder.WalletBefore, i.Config.TradeCurrency)
	balanceTradeAfter = trading.CurrencyAmountTotal(&lastOrder.WalletAfter, i.Config.TradeCurrency)

	balanceTotalBefore = balanceMainBefore + trading.TradeCurrencyToMain(balanceTradeBefore, firstOrder.Price)
	balanceTotalAfter = balanceMainAfter + trading.TradeCurrencyToMain(balanceTradeAfter, lastOrder.Price)

	totalPnl = balanceTotalAfter - balanceTotalBefore
	if balanceTotalBefore != 0 {
		//totalPnlPercent = totalPnl / (balanceBefore / 100)
		totalPnlPercent = math.Div(totalPnl, math.Div(balanceTotalBefore, 100))
	}

	for _, asset := range assets {
		switch asset.TransactionType {
		case structs.TransactionTypeDeposit:
			depoAmount += asset.Amount
			if asset.CreatedTime.Before(firstOrder.CreatedTime) {
				balanceMainBefore -= asset.Amount
				balanceTotalBefore -= asset.Amount
			}
			if asset.CreatedTime.After(lastOrder.CreatedTime) {
				balanceMainAfter += asset.Amount
				balanceTotalAfter += asset.Amount
			}
		case structs.TransactionTypeWithdraw:
			withdrawAmount += asset.Amount
			if asset.CreatedTime.Before(firstOrder.CreatedTime) {
				balanceMainBefore += asset.Amount
				balanceTotalBefore += asset.Amount
			}
			if asset.CreatedTime.After(lastOrder.CreatedTime) {
				balanceMainAfter -= asset.Amount
				balanceTotalAfter -= asset.Amount
			}
		}
	}
	return
}
