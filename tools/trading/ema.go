package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
)

// CalcEMASeries calculates the exponential moving average of the close prices
// for every candle where enough history is available.
//
// The returned slice is aligned with the candles slice: emas[i] is the EMA value
// at candles[i]. Candles follow the project convention: index 0 is the newest
// candle, index len-1 is the oldest. The oldest period-1 positions have no EMA
// (not enough history to seed it) and are left as 0.
//
// The EMA is seeded with the simple moving average of the oldest `period`
// closes and then rolled forward towards the newest candle with the standard
// smoothing multiplier 2/(period+1).
func CalcEMASeries(candles []domainStructs.DomainCandle, period int64) (emas []float64, err error) {
	if period <= 0 {
		err = apperrors.New("ema period must be positive, got %d", period)
		return
	}
	if int64(len(candles)) < period {
		err = apperrors.New("not enough candles to calculate ema: need %d, got %d", period, len(candles))
		return
	}

	emas = make([]float64, len(candles))

	// Seed: SMA of the oldest `period` closes.
	var sum float64
	for i := int64(0); i < period; i++ {
		sum += candles[int64(len(candles))-1-i].Close
	}
	ema := sum / float64(period)
	seedIdx := int64(len(candles)) - period
	emas[seedIdx] = ema

	// Roll forward oldest→newest; slice index decreases towards the newest candle.
	multiplier := 2.0 / (float64(period) + 1.0)
	for idx := seedIdx - 1; idx >= 0; idx-- {
		ema = (candles[idx].Close-ema)*multiplier + ema
		emas[idx] = ema
	}

	return
}

// CalcEMA returns the exponential moving average of the close prices at the
// newest candle (index 0).
func CalcEMA(candles []domainStructs.DomainCandle, period int64) (ema float64, err error) {
	emas, err := CalcEMASeries(candles, period)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate ema series")
		return
	}
	ema = emas[0]
	return
}
