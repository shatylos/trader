package trading

import domainStructs "github.com/shatylos/trader/internal/domain/structs"

func MinMaxPrice(candles []domainStructs.DomainCandle) (minPrice, maxPrice float64) {

	for _, candle := range candles {
		if candle.Low < minPrice || minPrice == 0.0 {
			minPrice = candle.Low
		}
		if candle.High > maxPrice || maxPrice == 0.0 {
			maxPrice = candle.High
		}
	}

	return
}
