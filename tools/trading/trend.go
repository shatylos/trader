package trading

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/logger"
)

const (
	TrendLong    = "BULLISH"
	TrendShort   = "BEARISH"
	TrendUnknown = "UNKNOWN"
)

func GetFullTrend(candles []domainStructs.DomainCandle, verbose bool) string {
	trendLinearRegression, _ := GetTrendLinearRegression(candles)
	trendSimpleCompare := GetTrendSimpleCompare(candles)
	if verbose {
		logger.Info(fmt.Sprintf("Trend Linear Regression: %s", trendLinearRegression))
		logger.Info(fmt.Sprintf("Trend Simple Compare: %s", trendSimpleCompare))
	}

	if trendLinearRegression == trendSimpleCompare {
		return trendLinearRegression
	}
	return TrendUnknown
}

// GetTrendLinearRegression fits an ordinary least-squares line to the
// candle close prices and classifies the trend.
//
// Two safeguards make the result usable across instruments of any price scale:
//   - the slope is normalized by the mean price, so the threshold is relative
//     (fraction of price per candle) rather than absolute price units;
//   - the coefficient of determination (R²) must clear a minimum, so a series
//     that merely drifts through noise is reported as UNKNOWN rather than a trend.
//
// Candles follow the project convention: index 0 is the newest candle,
// index len-1 is the oldest. We iterate oldest→newest so a rising price
// yields a positive slope.
func GetTrendLinearRegression(candles []domainStructs.DomainCandle) (trend string, slope float64) {
	const (
		// minimum relative slope (per candle) to call a trend, ~0.05% per candle
		slopeThreshold = 0.0005
		// minimum R²: how much of the price variance the line explains
		minRSquared = 0.3
	)

	trend = TrendUnknown
	if len(candles) < 2 {
		return
	}

	candleCount := float64(len(candles))
	var sumTime, sumClose, sumTimeClose, sumTimeSquared float64

	// Walk oldest→newest. Newest is candles[0], oldest is candles[len-1],
	// so timeIndex grows as we move backwards through the slice.
	for idx := 0; idx < len(candles); idx++ {
		timeIndex := float64(idx)
		closePrice := candles[len(candles)-1-idx].Close

		sumTime += timeIndex
		sumClose += closePrice
		sumTimeClose += timeIndex * closePrice
		sumTimeSquared += timeIndex * timeIndex
	}

	// Slope (m) of best-fit line: m = (N*ΣXY - ΣX*ΣY) / (N*ΣX² - (ΣX)²)
	denominator := candleCount*sumTimeSquared - sumTime*sumTime
	if denominator == 0 {
		return
	}

	rawSlope := (candleCount*sumTimeClose - sumTime*sumClose) / denominator
	slope = rawSlope

	meanClose := sumClose / candleCount
	if meanClose == 0 {
		return
	}

	// R² = 1 - SSres/SStot, computed in a second pass.
	intercept := (sumClose - rawSlope*sumTime) / candleCount
	var sumSquaredResiduals, sumSquaredDeviations float64
	for idx := 0; idx < len(candles); idx++ {
		timeIndex := float64(idx)
		closePrice := candles[len(candles)-1-idx].Close
		predicted := intercept + rawSlope*timeIndex
		sumSquaredResiduals += (closePrice - predicted) * (closePrice - predicted)
		sumSquaredDeviations += (closePrice - meanClose) * (closePrice - meanClose)
	}
	if sumSquaredDeviations == 0 {
		return // flat line, no trend
	}
	rSquared := 1 - sumSquaredResiduals/sumSquaredDeviations
	if rSquared < minRSquared {
		return // fit too weak to trust
	}

	// Normalize: fraction of mean price moved per candle.
	relativeSlope := rawSlope / meanClose
	if relativeSlope > slopeThreshold {
		trend = TrendLong
	} else if relativeSlope < -slopeThreshold {
		trend = TrendShort
	}
	return
}

func GetTrendSimpleCompare(candles []domainStructs.DomainCandle) string {
	var minCandle, maxCandle domainStructs.DomainCandle
	if len(candles) == 0 {
		return TrendUnknown
	}
	minCandle = candles[0]
	maxCandle = candles[0]
	for _, candle := range candles {
		if candle.Low < minCandle.Low {
			minCandle = candle
		}
		if candle.High > maxCandle.High {
			maxCandle = candle
		}
	}

	if maxCandle.Time > minCandle.Time {
		return TrendLong
	} else if maxCandle.Time < minCandle.Time {
		return TrendShort
	}
	return TrendUnknown
}
