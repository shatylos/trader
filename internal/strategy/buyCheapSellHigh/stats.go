package buyCheapSellHigh

import (
	strategyStorage "github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage"
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/structs"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"time"
)

func (s *BuyCheapSellHigh) GetStats() (stats _struct.Stats, err error) {
	stats.SetupId = s.GetId()

	now := time.Now()
	startOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfPrevMonth := startOfCurrentMonth.AddDate(0, -1, 0)
	startOfPrev12Month := startOfCurrentMonth.AddDate(0, -12, 0)
	startTotal := time.Unix(0, 0)

	stats.PNLLastMonth, err = s.calculatePnl(startOfPrevMonth, startOfCurrentMonth)
	stats.PNL12Months, err = s.calculatePnl(startOfPrev12Month, startOfCurrentMonth)
	stats.PNLTotal, err = s.calculatePnl(startTotal, startOfCurrentMonth)

	stats.WithdrawablePrevMonth = 0
	if stats.PNLLastMonth.Amount > 0 {
		stats.WithdrawablePrevMonth = math.Mul(math.Div(stats.PNLLastMonth.Amount, 100), s.WithdrawPercent)
	}

	return
}

func (s *BuyCheapSellHigh) calculatePnl(from time.Time, to time.Time) (pnl _struct.Pnl, err error) {
	var storage *mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(s.Id)
	if err != nil {
		return
	}
	var orders []structs.HistoryOrder
	orders, err = storage.GetHistoryOrders(from, to)
	if err != nil || len(orders) == 0 {
		return
	}

	var amount float64
	for _, order := range orders {
		amount += order.Revenue
	}

	firstOrder := orders[len(orders)-1]
	startAmount := math.Mul(firstOrder.TradeCurrencyAmountBefore, firstOrder.FilledPrice)
	startAmount += firstOrder.MainCurrencyAmountBefore
	startOnePercent := math.Div(startAmount, 100)
	pnlPercent := 0.0
	if startOnePercent > 0 {
		pnlPercent = math.Div(amount, startOnePercent)
	}

	pnl = _struct.Pnl{
		Amount:   amount,
		Percent:  pnlPercent,
		Currency: s.MainCurrency,
	}
	return
}
