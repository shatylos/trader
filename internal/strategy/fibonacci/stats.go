package fibonacci

import (
	"fmt"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	_struct "github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/math"
	"time"
)

func (f *Fibonacci) GetStats() (stats _struct.Stats, err error) {
	stats.SetupId = f.GetId()

	now := time.Now()
	startOfCurrentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfPrevMonth := startOfCurrentMonth.AddDate(0, -1, 0)
	startOfPrev12Month := startOfCurrentMonth.AddDate(0, -12, 0)
	startTotal := time.Unix(0, 0)

	stats.PNLLastMonth, err = f.calculatePnl(startOfPrevMonth, startOfCurrentMonth)
	stats.PNL12Months, err = f.calculatePnl(startOfPrev12Month, startOfCurrentMonth)
	stats.PNLTotal, err = f.calculatePnl(startTotal, startOfCurrentMonth)

	stats.WithdrawablePrevMonth = 0
	if stats.PNLLastMonth.Amount > 0 {
		stats.WithdrawablePrevMonth = math.Mul(math.Div(stats.PNLLastMonth.Amount, 100), f.config.WithdrawPercent)
	}

	return
}

func (f *Fibonacci) calculatePnl(from time.Time, to time.Time) (pnl _struct.Pnl, err error) {
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	var positions []structs.Position
	positions, err = storage.GetPositions(from, to)
	if err != nil || len(positions) == 0 {
		return
	}

	revPerMonth := make(map[string]float64)
	startAmounts := make(map[string]float64)
	var amount float64
	for _, position := range positions {
		if position.Status == structs.StatusActive {
			continue
		}
		amount += position.TotalClosePnl
		date := time.Unix(position.CreatedTime, 0)
		monthKey := fmt.Sprintf("%d-%d", date.Year(), date.Month())
		revPerMonth[monthKey] += position.TotalClosePnl
		startAmounts[monthKey] = position.BalanceBefore
	}

	revPercAllMonths := 0.0
	for monthKey, revenue := range revPerMonth {
		revPercAllMonths += math.Div(revenue, math.Div(startAmounts[monthKey], 100))
	}
	avPercPerMonth := 0.0
	if revPercAllMonths > 0.0 {
		avPercPerMonth = math.Div(revPercAllMonths, float64(len(revPerMonth)))
	}

	lastPosition := positions[len(positions)-1]
	onePercent := math.Div(lastPosition.BalanceBefore, 100)
	pnlPercent := math.Div(amount, onePercent)

	pnl = _struct.Pnl{
		Amount:         amount,
		Percent:        pnlPercent,
		AvPercPerMonth: avPercPerMonth,
		Currency:       f.config.MainCurrency,
	}
	return
}
