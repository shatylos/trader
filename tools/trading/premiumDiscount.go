package trading

import (
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

const (
	ZonePremium  = "PREMIUM"
	ZoneDiscount = "DISCOUNT"
	ZoneNeutral  = "NEUTRAL"
)

// PremiumDiscount calculate if the current price in premium or discount.
// maximum premium: pdKoef 100
// maximum discount: pdKoef -100
func PremiumDiscount(klines []domainStructs.DomainCandle) (pdKoef float64) {
	if len(klines) == 0 {
		return
	}
	price := klines[0].Close
	maxPrice := klines[0].High
	minPrice := klines[0].Low

	for _, kline := range klines {
		if kline.High > maxPrice {
			maxPrice = kline.High
		}
		if kline.Low < minPrice {
			minPrice = kline.Low
		}
	}

	middlePrice := math.Div(maxPrice-minPrice, 2) + minPrice

	if price > middlePrice {
		// premium
		// pdKoef = (price - middlePrice) / ((maxPrice - middlePrice) / 100)
		pdKoef = math.Div(price-middlePrice, math.Div(maxPrice-middlePrice, 100))
	} else if price < middlePrice {
		// discount
		// pdKoef = (middlePrice - price) / ((middlePrice - minPrice) / 100) * -1
		pdKoef = math.Mul(math.Div(middlePrice-price, math.Div(middlePrice-minPrice, 100)), -1)
	}
	// else middle pdKoef = 0

	return
}
