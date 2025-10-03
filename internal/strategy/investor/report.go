package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/web/helper"
	"html/template"
	"strings"
	"time"
)

type ReportTemplateData struct {
	PrevPeriodLink     string
	NextPeriodLink     string
	DateFrom           time.Time
	DateTo             time.Time
	DealsRelations     []*DealRelation
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
}

func (i *Investor) GetReport(from time.Time, to time.Time) (report _struct.Report, err error) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, CtxSetupKey, i)

	var tmpl *template.Template
	tmpl, err = helper.GetTemplate("web/template/investor/report.html")
	if err != nil {
		return
	}

	var deals, closedDeals, activeDeals []*storage.Deal
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

	var dealRelations []*DealRelation
	for _, deal := range deals {
		var dealRelation *DealRelation
		dealRelation, err = i.GetDealRelation(ctx, deal)
		if err != nil {
			return
		}
		dealRelations = append(dealRelations, dealRelation)
	}

	// @TODO: Add assets
	var assets []*structs.AssetTransaction
	//assets, err = i.Storage.GetAssetTransactions(from, to)
	totalPnl, totalPnlPercent,
		balanceTotalBefore, balanceMainBefore, balanceTradeBefore,
		balanceTotalAfter, balanceMainAfter, balanceTradeAfter,
		depoAmount, withdrawAmount := i.reportAmounts(from, dealRelations, assets)

	data := ReportTemplateData{
		PrevPeriodLink:     fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 0, -1).Format("2006-01")),
		NextPeriodLink:     fmt.Sprintf("/report/%s/%s/", i.GetId(), from.AddDate(0, 1, 0).Format("2006-01")),
		DateFrom:           from,
		DateTo:             to,
		DealsRelations:     dealRelations,
		MainCurrency:       i.config.MainCurrency,
		TradeCurrency:      i.config.TradeCurrency,
		PricePrecision:     int(i.config.PricePrecision),
		QtyPrecision:       int(i.config.QtyPrecision),
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
	}

	var resultBuilder strings.Builder
	err = tmpl.Execute(&resultBuilder, data)
	if err != nil {
		return
	}

	htmlStr := resultBuilder.String()

	report = _struct.Report{
		InnerHtml: template.HTML(htmlStr),
		SetupId:   i.GetId(),
	}
	return
}

func (i *Investor) reportAmounts(from time.Time, dealRelations []*DealRelation, assets []*structs.AssetTransaction) (
	totalPnl, totalPnlPercent,
	balanceTotalBefore, balanceMainBefore, balanceTradeBefore,
	balanceTotalAfter, balanceMainAfter, balanceTradeAfter,
	depoAmount, withdrawAmount float64) {

	var firstOrder, lastOrder *storage.Order

	for _, dealRelation := range dealRelations {
		if dealRelation.Deal.Status != structs.StatusClosed {
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

	balanceMainBefore = currencyAmountTotal(&firstOrder.WalletBefore, i.config.MainCurrency)
	balanceMainAfter = currencyAmountTotal(&lastOrder.WalletAfter, i.config.MainCurrency)
	balanceTradeBefore = currencyAmountTotal(&firstOrder.WalletBefore, i.config.TradeCurrency)
	balanceTradeAfter = currencyAmountTotal(&lastOrder.WalletAfter, i.config.TradeCurrency)

	balanceTotalBefore = balanceMainBefore + tradeCurrencyToMain(balanceTradeBefore, firstOrder.Price)
	balanceTotalAfter = balanceMainAfter + tradeCurrencyToMain(balanceTradeAfter, lastOrder.Price)

	totalPnl = balanceTotalAfter - balanceTotalBefore
	if balanceTotalBefore != 0 {
		//totalPnlPercent = totalPnl / (balanceBefore / 100)
		totalPnlPercent = math.Div(totalPnl, math.Div(balanceTotalBefore, 100))
	}
	return
}
