package scalper

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/scalper/storage"
	"github.com/shatylos/trader/internal/strategy/scalper/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/scalper/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
	math2 "math"
	"time"
)

// signalSnapshot is the multi-timeframe indicator state a trading decision is
// made on: the higher timeframe bias, the lower timeframe indicator values and
// the resulting entry signal (empty when no entry).
type signalSnapshot struct {
	Bias           string
	LtTrend        string
	Signal         string
	CurrentPrice   float64
	Chart          structs.ScalpChart
	SkippedMessage string
}

// calculateSignal loads both timeframes, computes the indicators on closed
// candles and evaluates the bias and entry conditions.
func (s *Scalper) calculateSignal() (snapshot signalSnapshot, err error) {
	var entryCandles []domainStructs.DomainCandle
	entryCandles, err = s.LoadEntryCandleHistory()
	if err != nil {
		err = apperrors.WrapExcuse(err, "error load entry candle history")
		return
	}
	if len(entryCandles) < 2 {
		err = apperrors.New("not enough entry candles loaded: %d", len(entryCandles))
		return
	}

	var biasCandles []domainStructs.DomainCandle
	biasCandles, err = s.LoadBiasCandleHistory()
	if err != nil {
		err = apperrors.WrapExcuse(err, "error load bias candle history")
		return
	}
	if len(biasCandles) < 2 {
		err = apperrors.New("not enough bias candles loaded: %d", len(biasCandles))
		return
	}

	// Index 0 is the forming candle on both timeframes; all decisions are made
	// on closed candles only.
	closedEntryCandles := entryCandles[1:]
	closedBiasCandles := biasCandles[1:]
	snapshot.CurrentPrice = entryCandles[0].Close

	var htfEmaFast, htfEmaSlow []float64
	htfEmaFast, err = trading.CalcEMASeries(closedBiasCandles, s.config.HtfEmaFastPeriod)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate htf fast ema")
		return
	}
	htfEmaSlow, err = trading.CalcEMASeries(closedBiasCandles, s.config.HtfEmaSlowPeriod)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate htf slow ema")
		return
	}

	snapshot.Bias, snapshot.LtTrend = detectBias(closedBiasCandles, htfEmaFast, htfEmaSlow, s.config.RequireTrendConfirmation)

	var ltfEma, rsi []float64
	ltfEma, err = trading.CalcEMASeries(closedEntryCandles, s.config.LtfEmaPeriod)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate ltf ema")
		return
	}
	rsi, err = trading.CalcRSISeries(closedEntryCandles, s.config.RsiPeriod)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate rsi")
		return
	}

	var atr float64
	atr, err = trading.CalcATR(closedEntryCandles, s.config.AtrPeriod)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate atr")
		return
	}

	snapshot.Chart = structs.ScalpChart{
		LtfEma:     math.Round(ltfEma[0], s.config.PricePrecision),
		HtfEmaFast: math.Round(htfEmaFast[0], s.config.PricePrecision),
		HtfEmaSlow: math.Round(htfEmaSlow[0], s.config.PricePrecision),
		Rsi:        math.Round(rsi[0], 2),
		Atr:        math.Round(atr, s.config.PricePrecision),
	}

	if snapshot.Bias == trading.TrendUnknown {
		snapshot.SkippedMessage = "No higher timeframe bias, waiting for a clear regime"
		return
	}

	if !isVolatilityEnough(atr, snapshot.CurrentPrice, s.config.MinAtrPercent) {
		snapshot.SkippedMessage = fmt.Sprintf("Volatility too low: ATR %g is less than %g%% of price %g",
			atr, s.config.MinAtrPercent, snapshot.CurrentPrice)
		return
	}

	if !isTpWorthFees(atr, s.config.TpAtrMult, snapshot.CurrentPrice, s.config.FeePercent, s.config.MinTpFeeRatio) {
		snapshot.SkippedMessage = fmt.Sprintf("Take profit %g ATRs (%.4f%%) does not cover round trip fee %g%% x ratio %g",
			s.config.TpAtrMult, math.Mul(atr, s.config.TpAtrMult)/snapshot.CurrentPrice*100, s.config.FeePercent, s.config.MinTpFeeRatio)
		return
	}

	snapshot.Signal = detectEntrySignal(snapshot.Bias, closedEntryCandles, ltfEma, rsi, entryParams{
		PullbackLookback: s.config.PullbackLookback,
		RsiMomentumLevel: s.config.RsiMomentumLevel,
		RsiOverbought:    s.config.RsiOverbought,
		RsiOversold:      s.config.RsiOversold,
	})
	if snapshot.Signal == "" {
		snapshot.SkippedMessage = "Waiting for a pullback entry trigger"
	}

	return
}

// entryTpSl derives the bracket levels from the ATR: the stop-loss sits
// SlAtrMult ATRs away from the entry, the take-profit TpAtrMult ATRs.
func (s *Scalper) entryTpSl(snapshot signalSnapshot) (entry, takeProfit, stopLoss float64) {
	entry = math.Round(snapshot.CurrentPrice, s.config.PricePrecision)
	atr := snapshot.Chart.Atr
	if snapshot.Signal == domainStructs.PositionSideLong {
		takeProfit = math.Round(entry+math.Mul(atr, s.config.TpAtrMult), s.config.PricePrecision)
		stopLoss = math.Round(entry-math.Mul(atr, s.config.SlAtrMult), s.config.PricePrecision)
	} else {
		takeProfit = math.Round(entry-math.Mul(atr, s.config.TpAtrMult), s.config.PricePrecision)
		stopLoss = math.Round(entry+math.Mul(atr, s.config.SlAtrMult), s.config.PricePrecision)
	}
	return
}

// calculateQty sizes the position so that hitting the stop-loss costs a fixed
// RiskPercent of the available balance: qty = (balance * risk%) / |entry - SL|.
func (s *Scalper) calculateQty(availableBalance float64, entry float64, stopLoss float64) (qty float64, err error) {
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

	riskAmount := math.Mul(math.Div(availableBalance, 100), s.config.RiskPercent)
	qty = math.Div(riskAmount, priceDistance)
	qty = math.Round(qty, s.config.QtyPrecision)
	if qty < s.config.MinQty && s.config.IncreaseRiskToMinQty {
		logger.Info(fmt.Sprintf("calculated qty %g is less than min qty %g. Increased risk", qty, s.config.MinQty))
		qty = s.config.MinQty
	}
	return
}

func (s *Scalper) openNewPosition(snapshot signalSnapshot) (err error) {
	if snapshot.Signal != domainStructs.PositionSideShort && snapshot.Signal != domainStructs.PositionSideLong {
		err = apperrors.New("unexpected signal value: %s", snapshot.Signal)
		return
	}

	var availableBalance float64
	availableBalance, err = s.getAvailableBalance()
	if err != nil {
		err = apperrors.Wrap(err, "error get available balance")
		return
	}

	entry, takeProfit, stopLoss := s.entryTpSl(snapshot)

	var qty float64
	qty, err = s.calculateQty(availableBalance, entry, stopLoss)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate qty")
		return
	}
	if qty < s.config.MinQty {
		logger.Warning(fmt.Sprintf("Order was not created. QTY (%f) less then min qty (%f)", qty, s.config.MinQty))
		return
	}

	internalPosition := structs.Position{
		Id:            nil,
		Chart:         snapshot.Chart,
		Bias:          snapshot.Bias,
		LtTrend:       snapshot.LtTrend,
		Signal:        snapshot.Signal,
		Status:        structs.StatusNew,
		BalanceBefore: availableBalance,
	}
	internalPosition.Chart.EntryPrice = entry
	internalPosition.Chart.TakeProfit = takeProfit
	internalPosition.Chart.StopLoss = stopLoss
	internalPosition.Chart.Qty = qty

	var orderId string
	orderId, err = s.provider.OpenPosition(domainStructs.DomainPositionRequest{
		Leverage:    s.config.Leverage,
		Price:       entry,
		Qty:         qty,
		ReduceOnly:  false,
		Side:        snapshot.Signal,
		StopLoss:    stopLoss,
		Symbol:      s.config.CoinPare,
		TakeProfit:  takeProfit,
		Type:        domainStructs.PositionTypes.Limit,
		TpOrderType: domainStructs.PositionTypes.Limit,
	})
	if err != nil {
		errMsg := fmt.Sprintf("error opening new position. Leverage: %d. Qty: %g. Side: %s. TakeProfit: %g. StopLoss: %g",
			s.config.Leverage, qty, snapshot.Signal, takeProfit, stopLoss)
		tgNotifier.Notify(errMsg)
		err = apperrors.Wrap(err, "%s", errMsg)
		return
	}

	var newOrder domainStructs.DomainOrder
	newOrder, err = s.provider.GetOrder(orderId, s.config.CoinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error get order %s", orderId)
		return
	}
	internalPosition.Order = newOrder
	internalPosition.Status = structs.StatusActive

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(s.config.Id)
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
	logger.Success(fmt.Sprintf("Created new order. PositionId %s. Qty: %g. Price: %g. Side: %s. Bias: %s. RSI: %g. ATR: %g. Entry: %g. TP: %g. SL: %g.",
		*internalPosition.Id, newOrder.Qty, newOrder.Price, newOrder.Side,
		snapshot.Bias, snapshot.Chart.Rsi, snapshot.Chart.Atr, entry, takeProfit, stopLoss,
	))
	if s.config.TelegramNotifier {
		tgNotifier.Notify(fmt.Sprintf("Created new position\nSide: %s\nQty: %g\nEntry Price: %g\nTP: %g\nSL: %g",
			newOrder.Side, newOrder.Qty, newOrder.Price, takeProfit, stopLoss))
	}
	return
}

// handlePendingEntryOrder refreshes the entry limit order of an active
// internal position that has no provider position yet. It returns true while
// the entry order still has a chance to be executed (anything but canceled),
// so the caller must not close the internal position or open a new one.
func (s *Scalper) handlePendingEntryOrder(internalPosition *structs.Position) (isPending bool, err error) {
	var order domainStructs.DomainOrder
	order, err = s.provider.GetOrder(internalPosition.Order.OrderId, s.config.CoinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error get order %s", internalPosition.Order.OrderId)
		return
	}

	if order.OrderStatus != internalPosition.Order.OrderStatus {
		internalPosition.Order = order

		var storage mongo.MongoStorage
		storage, err = strategyStorage.GetStorage(s.config.Id)
		if err != nil {
			err = apperrors.Wrap(err, "error get storage")
			return
		}
		*internalPosition, err = storage.SaveInternalPosition(*internalPosition)
		if err != nil {
			err = apperrors.Wrap(err, "error save internal position")
			return
		}
	}

	// A just filled order is still pending here: the provider position was not
	// visible in this tick yet, so let the next tick pick it up.
	isPending = order.OrderStatus != domainStructs.OrderStatuses.Canceled
	return
}

func (s *Scalper) closeInternalPosition(internalPosition structs.Position) (err error) {
	if internalPosition.Order.OrderStatus == domainStructs.OrderStatuses.Open {
		var order domainStructs.DomainOrder
		order, err = s.provider.GetOrder(internalPosition.Order.OrderId, s.config.CoinPare)
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
	availableBalance, err = s.getAvailableBalance()
	if err != nil {
		err = apperrors.Wrap(err, "error get available balance")
		return
	}
	internalPosition.BalanceAfter = availableBalance
	internalPosition.TotalClosePnl = internalPosition.BalanceAfter - internalPosition.BalanceBefore

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(s.config.Id)
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
	if s.config.TelegramNotifier {
		tgNotifier.Notify(fmt.Sprintf("Position was closed. PNL: %g.", math.Round(internalPosition.TotalClosePnl, 2)))
	}
	return
}
