package investor

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/apperrors"
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
	if err != nil {
		err = apperrors.Wrap(err, "error calculate PNL from prev month")
		return
	}
	stats.PNL12Months, err = i.calculatePnl(startOfPrev12Month, startOfCurrentMonth)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate PNL from prev 12 months")
		return
	}
	stats.PNLTotal, err = i.calculatePnl(startTotal, startOfCurrentMonth)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate PNL for total period")
		return
	}

	stats.WithdrawablePrevMonth = 0
	if stats.PNLLastMonth.Amount > 0 {
		stats.WithdrawablePrevMonth = math.Mul(math.Div(stats.PNLLastMonth.Amount, 100), i.Config.WithdrawPercent)
	}

	return
}

func (i *Investor) calculatePnl(from, to time.Time) (pnl _struct.Pnl, err error) {
	ctx := i.getContext()

	var orders []*entity.Order
	orders, err = i.Storage.GetOrdersByPeriod(ctx, from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get orders by period from %s to %s", from, to)
		return
	}

	revPerMonth := make(map[string]float64)
	startAmounts := make(map[string]float64)
	var amount, totalAmountBefore float64

	// orders are sorted by CreatedTime ascending
	for _, order := range orders {
		if order.OrderStatus != structs.OrderStatuses.Filled {
			continue
		}

		monthKey := fmt.Sprintf("%d-%d", order.CreatedTime.Year(), order.CreatedTime.Month())
		if _, ok := startAmounts[monthKey]; !ok {
			mainCurr := trading.CurrencyAmountTotal(&order.WalletBefore, i.Config.MainCurrency)
			tradeCurr := trading.CurrencyAmountTotal(&order.WalletBefore, i.Config.TradeCurrency)
			startAmounts[monthKey] = mainCurr + trading.TradeCurrencyToMain(tradeCurr, order.Price)
		}
		if totalAmountBefore == 0 {
			totalAmountBefore = startAmounts[monthKey]
		}

		if order.Side == structs.OrderSideSell {
			amount += order.RealizedPNL
			revPerMonth[monthKey] += order.RealizedPNL
		}
	}

	revPercAllMonths := 0.0
	for monthKey, revenue := range revPerMonth {
		startAmount := math.Div(startAmounts[monthKey], 100)
		if startAmount > 0 {
			revPercAllMonths += math.Div(revenue, startAmount)
		}
	}
	avPercPerMonth := 0.0
	if revPercAllMonths > 0.0 {
		avPercPerMonth = math.Div(revPercAllMonths, float64(len(revPerMonth)))
	}

	var pnlPercent float64
	if totalAmountBefore > 0 && amount > 0 {
		onePercent := math.Div(totalAmountBefore, 100)
		pnlPercent = math.Div(amount, onePercent)
	}

	pnl = _struct.Pnl{
		Amount:         amount,
		Percent:        pnlPercent,
		AvPercPerMonth: avPercPerMonth,
		Currency:       i.Config.MainCurrency,
	}
	return
}
