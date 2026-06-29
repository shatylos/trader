package vwapReversion

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/vwapReversion/storage"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
	math2 "math"
	"time"
)

// calculateNewPosition derives the long term trend filter and the VWAP deviation
// bands from the provided short-term candles, and produces a fresh internal
// position holding the entry / take-profit / stop-loss levels and the risk sized qty.
func (v *VwapReversion) calculateNewPosition(candles []domainStructs.DomainCandle) (position structs.Position, err error) {
	var ltCandles []domainStructs.DomainCandle
	ltCandles, err = v.LoadLTCandleHistory()
	if err != nil {
		err = apperrors.WrapExcuse(err, "error load long trend candle history")
		return
	}

	// Short term trend is informational only; entry direction is gated by the
	// long term trend filter to avoid fading a strong move.
	stTrend, _ := trading.GetTrendLinearRegression(candles)
	ltTrend, _ := trading.GetTrendLinearRegression(ltCandles)
	if v.config.Verbose {
		logger.Info(fmt.Sprintf("Detected short time trend: %s, long time trend: %s", stTrend, ltTrend))
	}

	chart := v.calculateVwapChart(candles)

	availableBalance := float64(0)
	availableBalance, err = v.getAvailableBalance()
	if err != nil {
		err = apperrors.Wrap(err, "error get available balance")
		return
	}
	if v.config.Verbose {
		logger.Info(fmt.Sprintf("Available balance: %g. VWAP: %g. Lower band: %g. Upper band: %g.",
			availableBalance, chart.Vwap, chart.LowerBand, chart.UpperBand))
	}

	position = structs.Position{
		Id:            nil,
		Chart:         chart,
		Trend:         ltTrend,
		LtTrend:       ltTrend,
		StTrend:       stTrend,
		Order:         domainStructs.DomainOrder{},
		Status:        structs.StatusNew,
		BalanceBefore: availableBalance,
	}

	return
}

// calculateVwapChart builds the VWAP, the upper/lower deviation bands and the
// resulting entry, take-profit (the VWAP line) and stop-loss levels for both a
// long and a short setup. The actual side is decided later by the trend filter.
func (v *VwapReversion) calculateVwapChart(candles []domainStructs.DomainCandle) (chart structs.VwapChart) {
	vwap := trading.CreateVWAP(candles)
	upperEntry, lowerEntry := vwap.CalcDeviation(v.config.EntrySigmaMult)
	// @TODO: maybe remove it
	//upperSl, lowerSl := vwap.CalcDeviation(v.config.SlSigmaMult)

	stdDev := float64(0)
	if v.config.EntrySigmaMult > 0 {
		// Recover sigma from the band offset so the report can show it.
		stdDev = math.Div(upperEntry-vwap.AncVWAP, v.config.EntrySigmaMult)
	}

	chart.Vwap = math.Round(vwap.AncVWAP, v.config.PricePrecision)
	chart.StdDev = math.Round(stdDev, v.config.PricePrecision)
	chart.UpperBand = math.Round(upperEntry, v.config.PricePrecision)
	chart.LowerBand = math.Round(lowerEntry, v.config.PricePrecision)
	chart.TakeProfit = math.Round(vwap.AncVWAP, v.config.PricePrecision)

	// Entry and stop-loss depend on side and are finalised when the order is
	// placed; the chart stores the long-side levels by default and the bearish
	// branch overrides them. We keep both bands available via UpperBand/LowerBand
	// and expose the matching SL through helper selection at order time.
	//_ = upperSl
	//_ = lowerSl

	return
}

// entryTpSlForSide returns the concrete entry, take-profit and stop-loss prices
// for the requested side, using the deviation bands of the chart.
func (v *VwapReversion) entryTpSlForSide(chart structs.VwapChart, side string) (entry, takeProfit, stopLoss float64) {
	takeProfit = chart.TakeProfit
	if side == domainStructs.PositionSideLong {
		entry = chart.LowerBand
		stopLoss = math.Round(chart.Vwap-math.Mul(chart.StdDev, v.config.SlSigmaMult), v.config.PricePrecision)
	} else {
		entry = chart.UpperBand
		stopLoss = math.Round(chart.Vwap+math.Mul(chart.StdDev, v.config.SlSigmaMult), v.config.PricePrecision)
	}
	return
}

// calculateQty sizes the position so that hitting the stop-loss costs a fixed
// RiskPercent of the available balance: qty = (balance * risk%) / |entry - SL|.
func (v *VwapReversion) calculateQty(availableBalance float64, entry float64, stopLoss float64) (qty float64, err error) {
	if availableBalance <= 0 {
		logger.Warning("Not enough deposit")
		err = apperrors.New("not enough deposit")
		return
	}

	priceDistance := math2.Abs(entry - stopLoss)
	if priceDistance <= 0 {
		err = apperrors.New("entry and stop-loss distance must be positive")
		return
	}

	riskAmount := math.Mul(math.Div(availableBalance, 100), v.config.RiskPercent)
	qty = math.Div(riskAmount, priceDistance)
	qty = math.Round(qty, v.config.QtyPrecision)
	if qty < v.config.MinQty && v.config.IncreaseRiskToMinQty {
		logger.Info(fmt.Sprintf("calculated qty %g is less than min qty %g. Increased risk", qty, v.config.MinQty))
		qty = v.config.MinQty
	}
	return
}

func (v *VwapReversion) openNewPosition(internalPosition structs.Position, side string) (isCreated bool, err error) {
	if side != domainStructs.PositionSideShort && side != domainStructs.PositionSideLong {
		err = apperrors.New("unexpected side value: %s", side)
		return
	}

	entry, takeProfit, stopLoss := v.entryTpSlForSide(internalPosition.Chart, side)

	var qty float64
	qty, err = v.calculateQty(internalPosition.BalanceBefore, entry, stopLoss)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate qty")
		return
	}
	if qty < v.config.MinQty {
		logger.Warning(fmt.Sprintf("Order was not created. QTY (%f) less then min qty (%f)", qty, v.config.MinQty))
		return
	}

	internalPosition.Chart.EntryPrice = entry
	internalPosition.Chart.StopLoss = stopLoss
	internalPosition.Chart.TakeProfit = takeProfit
	internalPosition.Chart.Qty = qty

	var orderId string
	// @TODO: review and maybe open position by limit order
	orderId, err = v.provider.OpenPosition(domainStructs.DomainPositionRequest{
		Leverage:   v.config.Leverage,
		Price:      0,
		Qty:        qty,
		ReduceOnly: false,
		Side:       side,
		StopLoss:   stopLoss,
		Symbol:     v.config.CoinPare,
		TakeProfit: takeProfit,
		Type:       domainStructs.PositionTypes.Market,
	})
	if err != nil {
		errMsg := fmt.Sprintf("error opening new position. Leverage: %d. Qty: %g. Side: %s. TakeProfit: %g. StopLoss: %g",
			v.config.Leverage, qty, side, takeProfit, stopLoss)
		tgNotifier.Notify(errMsg)
		err = apperrors.Wrap(err, errMsg)
		return
	}

	var newOrder domainStructs.DomainOrder
	newOrder, err = v.provider.GetOrder(orderId, v.config.CoinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error get order %s", orderId)
		return
	}
	internalPosition.Order = newOrder
	internalPosition.Status = structs.StatusActive

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(v.config.Id)
	if err != nil {
		err = apperrors.Wrap(err, "error get storage")
		return
	}
	internalPosition, err = storage.SaveInternalPosition(internalPosition)
	if err != nil {
		err = apperrors.Wrap(err, "error save internal position")
		return
	}
	if internalPosition.Id == nil {
		err = apperrors.New("internal position was not saved")
		return
	}
	isCreated = true
	logger.Success(fmt.Sprintf("Created new order. PositionId %s. Qty: %g. Price: %g. Side: %s. VWAP: %g. Entry: %g. TP: %g. SL: %g.",
		*internalPosition.Id, newOrder.Qty, newOrder.Price, newOrder.Side,
		internalPosition.Chart.Vwap, entry, takeProfit, stopLoss,
	))
	if v.config.TelegramNotifier {
		tgNotifier.Notify(fmt.Sprintf("Created new position\nSide: %s\nQty: %g\nEntry Price: %g\nTP: %g\nSL: %g",
			newOrder.Side, newOrder.Qty, newOrder.Price, takeProfit, stopLoss))
	}
	return
}

func (v *VwapReversion) closeInternalPosition(internalPosition structs.Position) (err error) {
	if internalPosition.Order.OrderStatus == domainStructs.OrderStatuses.Open {
		var order domainStructs.DomainOrder
		order, err = v.provider.GetOrder(internalPosition.Order.OrderId, v.config.CoinPare)
		if err != nil {
			err = apperrors.Wrap(err, "error get order %s", internalPosition.Order.OrderId)
			return
		}
		if order.OrderStatus != domainStructs.OrderStatuses.Canceled && order.OrderStatus != domainStructs.OrderStatuses.Filled {
			logger.Warning(fmt.Sprintf(
				"Expected order must be cancelled or filled when close the internal position. OrderId: %s, OrderStatus: %s",
				order.OrderId, order.OrderStatus,
			))
		}
		internalPosition.Order = order
	}

	internalPosition.Status = structs.StatusClosed
	internalPosition.ClosedTime = time.Now()

	availableBalance := float64(0)
	availableBalance, err = v.getAvailableBalance()
	if err != nil {
		err = apperrors.Wrap(err, "error get available balance")
		return
	}
	internalPosition.BalanceAfter = availableBalance
	internalPosition.TotalClosePnl = internalPosition.BalanceAfter - internalPosition.BalanceBefore

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(v.config.Id)
	if err != nil {
		err = apperrors.Wrap(err, "error get storage")
		return
	}
	internalPosition, err = storage.SaveInternalPosition(internalPosition)
	if err != nil {
		err = apperrors.Wrap(err, "error save internal position")
		return
	}
	if internalPosition.Id == nil {
		err = apperrors.New("internal position was not closed")
		return
	}
	logger.Info(fmt.Sprintf("Position was closed. Position ID %s.", *internalPosition.Id))
	if v.config.TelegramNotifier {
		tgNotifier.Notify(fmt.Sprintf("Position was closed. PNL: %g.", math.Round(internalPosition.TotalClosePnl, 2)))
	}
	return
}
