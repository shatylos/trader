package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

func (i *Investor) GetStats() (stats _struct.Stats, err error) {
	stats.SetupId = i.GetId()

	now := time.Now()
	startOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfPrevMonth := startOfCurrentMonth.AddDate(0, -1, 0)
	startOfPrev12Month := startOfCurrentMonth.AddDate(0, -12, 0)
	startTotal := time.Unix(0, 0)

	stats.PNLLastMonth, err = i.calculatePnl(startOfPrevMonth, startOfCurrentMonth)
	stats.PNL12Months, err = i.calculatePnl(startOfPrev12Month, startOfCurrentMonth)
	stats.PNLTotal, err = i.calculatePnl(startTotal, startOfCurrentMonth)

	stats.WithdrawablePrevMonth = 0
	if stats.PNLLastMonth.Amount > 0 {
		stats.WithdrawablePrevMonth = math.Mul(math.Div(stats.PNLLastMonth.Amount, 100), i.config.WithdrawPercent)
	}

	return
}

func (i *Investor) calculatePnl(from, to time.Time) (pnl _struct.Pnl, err error) {
	ctx := context.Background()

	var dealRelations []*entity.DealRelation
	dealRelations, err = i.GetDealRelationsByPeriod(ctx, from, to)

	revPerMonth := make(map[string]float64)
	startAmounts := make(map[string]float64)
	var amount float64
	for _, dealRelation := range dealRelations {
		if len(dealRelation.Orders) == 0 {
			continue
		}
		revenue := dealRelation.RevenueMainCur + trading.TradeCurrencyToMain(dealRelation.RevenueTradeCur, dealRelation.Orders[0].Price)
		amount += revenue

		date := dealRelation.Deal.ClosedTime
		monthKey := fmt.Sprintf("%d-%d", date.Year(), date.Month())
		revPerMonth[monthKey] += revenue
		startAmounts[monthKey] = dealRelation.GetTotalAmountBefore(i.config.MainCurrency, i.config.TradeCurrency)
	}

	revPercAllMonths := 0.0
	for monthKey, revenue := range revPerMonth {
		revPercAllMonths += math.Div(revenue, math.Div(startAmounts[monthKey], 100))
	}
	avPercPerMonth := 0.0
	if revPercAllMonths > 0.0 {
		avPercPerMonth = math.Div(revPercAllMonths, float64(len(revPerMonth)))
	}

	var pnlPercent float64
	if len(dealRelations) > 0 {
		lastDealRelation := dealRelations[len(dealRelations)-1]
		totalAmountBefore := lastDealRelation.GetTotalAmountBefore(i.config.MainCurrency, i.config.TradeCurrency)
		if totalAmountBefore > 0 && amount > 0 {
			onePercent := math.Div(totalAmountBefore, 100)
			pnlPercent = math.Div(amount, onePercent)
		}
	}

	pnl = _struct.Pnl{
		Amount:         amount,
		Percent:        pnlPercent,
		AvPercPerMonth: avPercPerMonth,
		Currency:       i.config.MainCurrency,
	}
	return
}
