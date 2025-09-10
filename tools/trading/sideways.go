package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

func CheckSideways(klines []domainStructs.DomainCandle, minKlinesAmount int64, pricePercent float64) (isSideways bool, kLinesAmount int64) {

	if len(klines) == 0 {
		return
	}

	price := klines[0].Close
	maxPrice := klines[0].High
	minPrice := klines[0].Low
	maxPriceRange := math.Mul(math.Div(price, 100), pricePercent)

	for i, kline := range klines {
		if kline.High > maxPrice {
			maxPrice = kline.High
		}
		if kline.Low < minPrice {
			minPrice = kline.Low
		}
		if maxPrice-minPrice > maxPriceRange {
			return
		}
		num := int64(i + 1)
		if num >= minKlinesAmount {
			isSideways = true
			kLinesAmount = num
		}
	}

	return
}
