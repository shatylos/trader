package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
)

const (
	TrendLong    = "BULLISH"
	TrendShort   = "BEARISH"
	TrendUnknown = "UNKNOWN"
)

func GetFullTrend(candles []domainStructs.DomainCandle) string {
	trendLinearRegression := GetTrendLinearRegression(candles)
	trendSimpleCompare := GetTrendSimpleCompare(candles)

	if trendLinearRegression == trendSimpleCompare {
		return trendLinearRegression
	}
	return TrendUnknown
}

func GetTrendLinearRegression(candles []domainStructs.DomainCandle) string {
	if len(candles) < 2 {
		return TrendUnknown
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
		return TrendUnknown
	}

	slope := numerator / denominator

	const threshold = 0.01 // small buffer to ignore noise

	if slope > threshold {
		return TrendLong
	} else if slope < -threshold {
		return TrendShort
	}

	return TrendUnknown
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
