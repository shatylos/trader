package scalper

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/trading"
)

// The scalper trades pullbacks in the direction of the higher timeframe trend:
//   - the higher (bias) timeframe defines the only allowed side via an EMA
//     regime filter, optionally confirmed by the linear regression trend;
//   - the lower (entry) timeframe waits for price to pull back into the fast
//     EMA and fires when a candle closes back on the trend side of it while
//     the RSI flips through the momentum level, showing the pullback is over;
//   - an ATR floor skips signals in a dead market where the expected move
//     would not cover the trading fees.
//
// All checks run on closed candles only; the forming candle is never used.

type entryParams struct {
	PullbackLookback int64
	RsiMomentumLevel float64
	RsiOverbought    float64
	RsiOversold      float64
}

// detectBias classifies the higher timeframe regime.
//
// Long regime: fast EMA above slow EMA and the last close above the slow EMA.
// Short regime is the mirror. Everything else is UNKNOWN (no trading).
// When requireTrendConfirmation is set, a linear regression trend pointing to
// the opposite side vetoes the EMA regime.
//
// closedCandles[0] must be the last closed candle. emaFast/emaSlow are aligned
// with closedCandles.
func detectBias(closedCandles []domainStructs.DomainCandle, emaFast []float64, emaSlow []float64, requireTrendConfirmation bool) (bias string, lrTrend string) {
	bias = trading.TrendUnknown
	lrTrend, _ = trading.GetTrendLinearRegression(closedCandles)
	if len(closedCandles) == 0 || len(emaFast) == 0 || len(emaSlow) == 0 {
		return
	}

	lastClose := closedCandles[0].Close
	if emaFast[0] > emaSlow[0] && lastClose > emaSlow[0] {
		bias = trading.TrendLong
	} else if emaFast[0] < emaSlow[0] && lastClose < emaSlow[0] {
		bias = trading.TrendShort
	}

	if requireTrendConfirmation && bias != trading.TrendUnknown {
		if (bias == trading.TrendLong && lrTrend == trading.TrendShort) ||
			(bias == trading.TrendShort && lrTrend == trading.TrendLong) {
			bias = trading.TrendUnknown
		}
	}

	return
}

// detectEntrySignal checks the lower timeframe for a finished pullback in the
// direction of the bias and returns the position side to open, or an empty
// string when there is no signal.
//
// Long entry fires when all of the following hold on closed candles
// (index 0 = last closed, the trigger candle):
//  1. bias is bullish;
//  2. the trigger candle is bullish and closes above the fast EMA;
//  3. within the previous PullbackLookback candles price traded at or below
//     the fast EMA (a pullback actually happened);
//  4. RSI on the trigger candle is above the momentum level but below the
//     overbought level (momentum turned up, not chasing an extended move);
//  5. within the previous PullbackLookback candles RSI was at or below the
//     momentum level (the flip is fresh, so the same trend leg is not
//     re-entered on every candle).
//
// Short entry is the exact mirror.
func detectEntrySignal(bias string, closedCandles []domainStructs.DomainCandle, ema []float64, rsi []float64, params entryParams) (side string) {
	lookback := params.PullbackLookback
	if int64(len(closedCandles)) < lookback+1 || int64(len(ema)) < lookback+1 || int64(len(rsi)) < lookback+1 {
		return
	}

	trigger := closedCandles[0]

	switch bias {
	case trading.TrendLong:
		if trigger.Close <= trigger.Open || trigger.Close <= ema[0] {
			return
		}
		if rsi[0] <= params.RsiMomentumLevel || rsi[0] >= params.RsiOverbought {
			return
		}
		pulledBack := false
		rsiFlipped := false
		for i := int64(1); i <= lookback; i++ {
			if closedCandles[i].Low <= ema[i] {
				pulledBack = true
			}
			if rsi[i] <= params.RsiMomentumLevel {
				rsiFlipped = true
			}
		}
		if pulledBack && rsiFlipped {
			side = domainStructs.PositionSideLong
		}
	case trading.TrendShort:
		if trigger.Close >= trigger.Open || trigger.Close >= ema[0] {
			return
		}
		if rsi[0] >= params.RsiMomentumLevel || rsi[0] <= params.RsiOversold {
			return
		}
		pulledBack := false
		rsiFlipped := false
		for i := int64(1); i <= lookback; i++ {
			if closedCandles[i].High >= ema[i] {
				pulledBack = true
			}
			if rsi[i] >= params.RsiMomentumLevel {
				rsiFlipped = true
			}
		}
		if pulledBack && rsiFlipped {
			side = domainStructs.PositionSideShort
		}
	}

	return
}

// isVolatilityEnough checks the ATR floor: the expected move (one ATR) must be
// at least minAtrPercent percent of the current price, otherwise a scalp can
// not cover the fees and spread.
func isVolatilityEnough(atr float64, currentPrice float64, minAtrPercent float64) bool {
	if currentPrice <= 0 {
		return false
	}
	return atr/currentPrice*100 >= minAtrPercent
}
