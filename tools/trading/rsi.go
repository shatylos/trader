package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
)

// CalcRSISeries calculates the Relative Strength Index (Wilder smoothing) of
// the close prices for every candle where enough history is available.
//
// The returned slice is aligned with the candles slice: rsis[i] is the RSI value
// at candles[i]. Candles follow the project convention: index 0 is the newest
// candle, index len-1 is the oldest. The oldest `period` positions have no RSI
// (not enough history to seed it) and are left as 0.
func CalcRSISeries(candles []domainStructs.DomainCandle, period int64) (rsis []float64, err error) {
	if period <= 0 {
		err = apperrors.New("rsi period must be positive, got %d", period)
		return
	}
	// period price changes require period+1 closes
	if int64(len(candles)) < period+1 {
		err = apperrors.New("not enough candles to calculate rsi: need %d, got %d", period+1, len(candles))
		return
	}

	rsis = make([]float64, len(candles))

	// Seed: plain average of gains and losses over the oldest `period` changes.
	var avgGain, avgLoss float64
	for i := int64(0); i < period; i++ {
		olderIdx := int64(len(candles)) - 1 - i
		change := candles[olderIdx-1].Close - candles[olderIdx].Close
		if change > 0 {
			avgGain += change
		} else {
			avgLoss += -change
		}
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	seedIdx := int64(len(candles)) - 1 - period
	rsis[seedIdx] = rsiFromAverages(avgGain, avgLoss)

	// Wilder smoothing towards the newest candle; slice index decreases as time moves forward.
	for idx := seedIdx - 1; idx >= 0; idx-- {
		change := candles[idx].Close - candles[idx+1].Close
		gain := 0.0
		loss := 0.0
		if change > 0 {
			gain = change
		} else {
			loss = -change
		}
		avgGain = (avgGain*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
		rsis[idx] = rsiFromAverages(avgGain, avgLoss)
	}

	return
}

func rsiFromAverages(avgGain, avgLoss float64) float64 {
	if avgLoss == 0 {
		if avgGain == 0 {
			// flat series: no direction, treat as neutral
			return 50
		}
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - 100/(1+rs)
}
