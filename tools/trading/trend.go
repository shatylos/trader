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

// @TODO: Review it
func GetTrendLinearRegression(candles []domainStructs.DomainCandle) (trend string, slope float64) {
	if len(candles) < 2 {
		trend = TrendUnknown
		return
	}

	reversed := make([]domainStructs.DomainCandle, len(candles))
	for i := range candles {
		reversed[i] = candles[len(candles)-1-i]
	}

	n := float64(len(reversed))
	var sumX, sumY, sumXY, sumXX float64

	for i, c := range reversed {
		x := float64(i)
		y := c.Close

		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	// Slope (m) of best-fit line: m = (N*ΣXY - ΣX*ΣY) / (N*ΣX² - (ΣX)²)
	numerator := n*sumXY - sumX*sumY
	denominator := n*sumXX - sumX*sumX

	if denominator == 0 {
		trend = TrendUnknown
		return
	}

	slope = numerator / denominator

	const threshold = 0.01 // small buffer to ignore noise

	if slope > threshold {
		trend = TrendLong
		return
	} else if slope < -threshold {
		trend = TrendShort
		return
	}

	trend = TrendUnknown
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
	} else {
		return TrendShort
	}
}
