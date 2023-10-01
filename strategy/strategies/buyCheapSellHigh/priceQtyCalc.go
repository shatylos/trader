package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/trading/services"
	"math"
	"time"
)

func (s *BuyCheapSellHigh) getPricesAndQtysToNewOrders(historyOrders []structs.DomainOrder, baseCurrencyAmount float64, tradeCurrencyAmount float64) (float64, float64, float64, float64, error) {

	currentPrice, err := s.getCurrentPrice()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	lastDirection, countLastDirection, lastOrderPrice := s.getLastOrdersInfo(historyOrders)

	buyPrice := s.getBuyPrice(baseCurrencyAmount, tradeCurrencyAmount, currentPrice, lastOrderPrice, lastDirection, countLastDirection)
	buyQty := s.getBuyQty(baseCurrencyAmount, tradeCurrencyAmount, currentPrice, lastDirection, countLastDirection)
	sellPrice := s.getSellPrice(baseCurrencyAmount, tradeCurrencyAmount, currentPrice, lastOrderPrice, lastDirection, countLastDirection)
	sellQty := s.getSellQty(baseCurrencyAmount, tradeCurrencyAmount, lastDirection, countLastDirection)

	return s.round(buyPrice), s.round(buyQty), s.round(sellPrice), s.round(sellQty), nil
}

func (s *BuyCheapSellHigh) getCurrentPrice() (float64, error) {
	now := time.Now()
	to := now.Unix()
	limit := int64(1)
	from := to - (tradeConst.ResolToSec[s.Resolution] * limit)
	candles, err := services.GetCandleHistory(s.DomainCode, s.CoinPare, s.Resolution, from, limit)
	if err != nil {
		return 0, err
	}
	return candles[0].Close, nil
}

func (s *BuyCheapSellHigh) getLastOrdersInfo(historyOrders []structs.DomainOrder) (string, int, float64) {
	lastDirection := ""
	countLastDirection := 0
	lastOrderPrice := float64(0)

	for _, historyOrder := range historyOrders {
		// @TODO: move to constant and map values
		if historyOrder.OrderStatus != "FILLED" && historyOrder.OrderStatus != "PARTIALLY_FILLED" {
			continue
		}
		if lastDirection == "" {
			lastDirection = historyOrder.Side
		}
		if lastOrderPrice == 0 {
			lastOrderPrice = historyOrder.Price
		}

		if lastDirection == historyOrder.Side {
			countLastDirection++
		} else {
			break
		}
	}
	return lastDirection, countLastDirection, lastOrderPrice
}

func (s *BuyCheapSellHigh) getBuyPrice(baseCurrencyAmount float64, tradeCurrencyAmount float64, currentPrice float64, lastOrderPrice float64, lastDirection string, countLastDirection int) float64 {
	buyPrice := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return buyPrice
	}

	if tradeCurrencyAmount == 0 {
		buyPrice = currentPrice
	}

	priceToCalcute := currentPrice
	if lastOrderPrice > 0 {
		priceToCalcute = lastOrderPrice
	}

	buyPriceRangeKey := s.getRangeKey("BUY", lastDirection, countLastDirection, s.CostRanges)

	buyPrice = priceToCalcute - float64(s.CostRanges[buyPriceRangeKey])

	return buyPrice
}

func (s *BuyCheapSellHigh) getSellPrice(baseCurrencyAmount float64, tradeCurrencyAmount float64, currentPrice float64, lastOrderPrice float64, lastDirection string, countLastDirection int) float64 {
	sellPrice := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return sellPrice
	}

	if baseCurrencyAmount == 0 {
		sellPrice = currentPrice
	}

	priceToCalcute := currentPrice
	if lastOrderPrice > 0 {
		priceToCalcute = lastOrderPrice
	}

	sellPriceRangeKey := s.getRangeKey("SELL", lastDirection, countLastDirection, s.CostRanges)

	sellPrice = priceToCalcute + float64(s.CostRanges[sellPriceRangeKey])

	return sellPrice
}

func (s *BuyCheapSellHigh) getBuyQty(baseCurrencyAmount float64, tradeCurrencyAmount float64, currentPrice float64, lastDirection string, countLastDirection int) float64 {
	buyQty := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return buyQty
	}

	if tradeCurrencyAmount == 0 {
		// @TODO: Check currency to sell
		buyQty = (baseCurrencyAmount / currentPrice) / 100 * float64(s.PercentRanges[0])
	}

	buyPercentRangeKey := s.getRangeKey("BUY", lastDirection, countLastDirection, s.PercentRanges)

	// @TODO: Check currency to sell
	buyQty = (baseCurrencyAmount / currentPrice) / 100 * float64(s.PercentRanges[buyPercentRangeKey])

	return buyQty
}

func (s *BuyCheapSellHigh) getSellQty(baseCurrencyAmount float64, tradeCurrencyAmount float64, lastDirection string, countLastDirection int) float64 {
	sellQty := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return sellQty
	}

	if baseCurrencyAmount == 0 {
		// @TODO: Check currency to sell
		sellQty = tradeCurrencyAmount / 100 * float64(s.PercentRanges[0])
	}

	sellPercentRangeKey := s.getRangeKey("SELL", lastDirection, countLastDirection, s.PercentRanges)

	// @TODO: Check currency to sell
	sellQty = tradeCurrencyAmount / 100 * float64(s.PercentRanges[sellPercentRangeKey])

	return sellQty
}

func (s *BuyCheapSellHigh) getRangeKey(orderDirection string, lastDirection string, countLastDirection int, ranges []int64) int {
	rangeKey := 0
	if lastDirection == orderDirection {
		if len(ranges) > countLastDirection {
			rangeKey = countLastDirection
		} else {
			rangeKey = len(ranges) - 1
		}
	}
	return rangeKey
}

func (s *BuyCheapSellHigh) round(value float64) float64 {
	return math.Round(value*1e6) / 1e6
}
