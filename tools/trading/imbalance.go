package trading

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/math"
)

type Imbalance struct {
	CandleMax         float64
	CandleMin         float64
	CandleBodyMax     float64
	CandleBodyMin     float64
	ImbalanceMax      float64
	ImbalanceMin      float64
	Trend             string
	ImbalanceTradeMax float64
	ImbalanceTradeMin float64
	Time              int64
	IsClosed          bool
}

// GetOpenImbalancesMinMaxPrice checks for open imbalances in the given candle data
// and determines the minimum and maximum price levels of that imbalances.
//
// Parameters:
//   - candles: The DomainCandle list for analysis.
//   - candlePercent: The threshold percentage used to detect an imbalance.
//
// Returns:
//   - isExistsImbalance: Indicates whether an open imbalance was found.
//   - minPrice: The minimum price of the detected imbalance range.
//   - maxPrice: The maximum price of the detected imbalance range.
func GetOpenImbalancesMinMaxPrice(candles []domainStructs.DomainCandle, candlePercent int64) (isExistsImbalance bool, minPrice float64, maxPrice float64) {
	var imbalances []Imbalance
	imbalances = GetImbalances(candles, candlePercent)

	if len(imbalances) > 0 {
		isExistsImbalance = true
		maxPrice = imbalances[0].ImbalanceMax
		minPrice = imbalances[0].ImbalanceMin
		for _, imbalance := range imbalances {
			if imbalance.ImbalanceMin < minPrice {
				minPrice = imbalance.ImbalanceMin
			}
			if imbalance.ImbalanceMax > maxPrice {
				maxPrice = imbalance.ImbalanceMax
			}
		}
	}

	return
}

// GetImbalances analyzes the provided candles and returns a list of imbalances.
//
// Parameters:
//   - candles: The DomainCandle list for analysis.
//   - candlePercent: The threshold percentage used to determine imbalances.
//
// Returns:
//   - imbalances: A slice of Imbalance values detected based on the provided candles and percentage.
func GetImbalances(candles []domainStructs.DomainCandle, candlePercent int64) (imbalances []Imbalance) {

	var allImbalances []Imbalance
	if len(candles) < 3 {
		return
	}

	var prevCandle, checkCandle, nextCandle *domainStructs.DomainCandle
	for i := len(candles) - 1; i >= 2; i-- {
		prevCandle = &candles[i]
		checkCandle = &candles[i-1]
		nextCandle = &candles[i-2]

		if prevCandle.High < nextCandle.Low {
			// imbalance bullish
			imbSpace := nextCandle.Low - prevCandle.High
			checkBody := checkCandle.Close - checkCandle.Open
			if imbSpace > math.Mul(math.Div(checkBody, 100.0), float64(candlePercent)) {
				allImbalances = append(allImbalances, Imbalance{
					CandleMax:         checkCandle.High,
					CandleMin:         checkCandle.Low,
					CandleBodyMax:     checkCandle.Close,
					CandleBodyMin:     checkCandle.Open,
					ImbalanceMax:      nextCandle.Low,
					ImbalanceMin:      prevCandle.High,
					ImbalanceTradeMax: nextCandle.Low,
					ImbalanceTradeMin: prevCandle.High,
					Trend:             TrendLong,
					Time:              checkCandle.Time,
				})
			}
		} else if prevCandle.Low > nextCandle.High {
			// imbalance bearish
			imbSpace := prevCandle.Low - nextCandle.High
			checkBody := checkCandle.Open - checkCandle.Close
			if imbSpace > math.Mul(math.Div(checkBody, 100.0), float64(candlePercent)) {
				allImbalances = append(allImbalances, Imbalance{
					CandleMax:         checkCandle.High,
					CandleMin:         checkCandle.Low,
					CandleBodyMax:     checkCandle.Open,
					CandleBodyMin:     checkCandle.Close,
					ImbalanceMax:      prevCandle.Low,
					ImbalanceMin:      nextCandle.High,
					ImbalanceTradeMax: prevCandle.Low,
					ImbalanceTradeMin: nextCandle.High,
					Trend:             TrendShort,
					Time:              checkCandle.Time,
				})
			}
		}

		tradeImbalances(allImbalances, nextCandle, candlePercent)
	}

	for _, imb := range allImbalances {
		if !imb.IsClosed {
			imbalances = append(imbalances, imb)
		}
	}

	fmt.Println(len(imbalances))
	fmt.Println(imbalances)

	return
}

func tradeImbalances(imbalances []Imbalance, nextCandle *domainStructs.DomainCandle, candlePercent int64) {
	for i := range imbalances {
		imb := &imbalances[i]
		if imb.IsClosed {
			continue
		}
		if imb.ImbalanceTradeMin > nextCandle.Low && imb.ImbalanceTradeMax < nextCandle.High {
			// candle more than imbalance in both side
			imb.ImbalanceTradeMax = nextCandle.Low
			imb.ImbalanceTradeMin = nextCandle.High
		} else if imb.ImbalanceTradeMin < nextCandle.Low && imb.ImbalanceTradeMax < nextCandle.High && imb.ImbalanceTradeMax > nextCandle.Low {
			// trade imbalance from High
			imb.ImbalanceTradeMax = nextCandle.Low
		} else if imb.ImbalanceTradeMin > nextCandle.Low && imb.ImbalanceTradeMax > nextCandle.High && imb.ImbalanceTradeMin < nextCandle.High {
			// trade imbalance from Low
			imb.ImbalanceTradeMin = nextCandle.High
		} else if imb.ImbalanceTradeMin < nextCandle.Low && imb.ImbalanceTradeMax > nextCandle.High {
			// trade imbalance inside
			if nextCandle.Open < nextCandle.Close {
				// bullish candle
				imb.ImbalanceTradeMin = nextCandle.High
			} else {
				// bearish candle
				imb.ImbalanceTradeMax = nextCandle.Low
			}
		}

		if imb.ImbalanceTradeMin > imb.ImbalanceTradeMax {
			imb.IsClosed = true
			continue
		}
		openImb := imb.ImbalanceTradeMax - imb.ImbalanceTradeMin
		if openImb > 0 {
			//openImb / 100 *candlePercent
		}

	}
}
