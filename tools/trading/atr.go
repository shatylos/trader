package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"math"
)

// CalcATR calculates the Average True Range (Wilder smoothing) at the newest
// candle (index 0).
//
// The true range of a candle takes the previous close into account, so
// period+1 candles are required. Candles follow the project convention:
// index 0 is the newest candle, index len-1 is the oldest.
func CalcATR(candles []domainStructs.DomainCandle, period int64) (atr float64, err error) {
	if period <= 0 {
		err = apperrors.New("atr period must be positive, got %d", period)
		return
	}
	if int64(len(candles)) < period+1 {
		err = apperrors.New("not enough candles to calculate atr: need %d, got %d", period+1, len(candles))
		return
	}

	// Seed: plain average of the true range over the oldest `period` candles.
	for i := int64(0); i < period; i++ {
		idx := int64(len(candles)) - 1 - period + i
		atr += trueRange(candles[idx], candles[idx+1])
	}
	atr /= float64(period)

	// Wilder smoothing towards the newest candle.
	for idx := int64(len(candles)) - 1 - period - 1; idx >= 0; idx-- {
		tr := trueRange(candles[idx], candles[idx+1])
		atr = (atr*float64(period-1) + tr) / float64(period)
	}

	return
}

// trueRange returns max(high-low, |high-prevClose|, |low-prevClose|).
func trueRange(candle, prevCandle domainStructs.DomainCandle) float64 {
	tr := candle.High - candle.Low
	highClose := math.Abs(candle.High - prevCandle.Close)
	if highClose > tr {
		tr = highClose
	}
	lowClose := math.Abs(candle.Low - prevCandle.Close)
	if lowClose > tr {
		tr = lowClose
	}
	return tr
}
