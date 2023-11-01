package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"math"
	"sort"
	"time"
)

func (s *BuyCheapSellHigh) getPricesAndQtysToNewOrders(historyOrders []structs.DomainOrder, baseCurrencyAmount float64, tradeCurrencyAmount float64) (float64, float64, float64, float64, error) {

	currentPrice, err := s.getCurrentPrice()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	lastDirection, countLastDirection, lastOrderPrice, lastOrderCreationTime := s.getLastOrdersInfo(historyOrders)

	qtyLongTermPercentCorrection := s.qtyLongTermPercentCorrection(currentPrice, baseCurrencyAmount, tradeCurrencyAmount)

	buyPrice := s.getBuyPrice(baseCurrencyAmount, tradeCurrencyAmount, currentPrice, lastOrderPrice, lastDirection, countLastDirection, lastOrderCreationTime)
	buyQty := s.getBuyQty(baseCurrencyAmount, tradeCurrencyAmount, currentPrice, lastDirection, countLastDirection, qtyLongTermPercentCorrection, lastOrderCreationTime)
	sellPrice := s.getSellPrice(baseCurrencyAmount, tradeCurrencyAmount, currentPrice, lastOrderPrice, lastDirection, countLastDirection, lastOrderCreationTime)
	sellQty := s.getSellQty(baseCurrencyAmount, tradeCurrencyAmount, lastDirection, countLastDirection, qtyLongTermPercentCorrection, lastOrderCreationTime)

	utils.LogInfo("Calculated order values:")
	utils.LogInfo(fmt.Sprintf("Buy price: %f, Buy qty: %f, Sell price: %f, Sell qty: %f", buyPrice, buyQty, sellPrice, sellQty))

	return buyPrice, buyQty, sellPrice, sellQty, nil
}

func (s *BuyCheapSellHigh) getCurrentPrice() (float64, error) {
	now := time.Now()
	to := now.Unix()
	limit := int64(1)
	from := to - (tradeConst.ResolToSec[s.Resolution] * limit)
	candles, err := s.Domain.LoadCandleHistory(s.CoinPare, s.Resolution, from, limit)
	sort.Slice(candles, func(i, j int) bool {
		return candles[i].Time > candles[j].Time
	})
	if err != nil {
		return 0, err
	}
	return candles[0].Close, nil
}

func (s *BuyCheapSellHigh) getLastOrdersInfo(historyOrders []structs.DomainOrder) (string, int, float64, int64) {
	lastDirection := ""
	countLastDirection := 0
	lastOrderCreationTime := int64(0)
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
		if lastOrderCreationTime == 0 {
			lastOrderCreationTime = historyOrder.CreatedTime
		}

		if lastDirection == historyOrder.Side {
			countLastDirection++
		} else {
			break
		}
	}
	return lastDirection, countLastDirection, lastOrderPrice, lastOrderCreationTime
}

func (s *BuyCheapSellHigh) getBuyPrice(baseCurrencyAmount float64, tradeCurrencyAmount float64, currentPrice float64, lastOrderPrice float64, lastDirection string, countLastDirection int, lastOrderCreationTime int64) float64 {
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

	buyPriceRangeKey := s.getRangeKey("BUY", lastDirection, countLastDirection, s.CostRanges, lastOrderCreationTime)

	buyPrice = priceToCalcute - float64(s.CostRanges[buyPriceRangeKey])
	buyPrice = s.round(buyPrice, float64(s.PurchasePricePrecision))

	return buyPrice
}

func (s *BuyCheapSellHigh) getSellPrice(baseCurrencyAmount float64, tradeCurrencyAmount float64, currentPrice float64, lastOrderPrice float64, lastDirection string, countLastDirection int, lastOrderCreationTime int64) float64 {
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

	sellPriceRangeKey := s.getRangeKey("SELL", lastDirection, countLastDirection, s.CostRanges, lastOrderCreationTime)

	sellPrice = priceToCalcute + float64(s.CostRanges[sellPriceRangeKey])
	sellPrice = s.round(sellPrice, float64(s.PurchasePricePrecision))

	return sellPrice
}

func (s *BuyCheapSellHigh) getBuyQty(baseCurrencyAmount float64, tradeCurrencyAmount float64, currentPrice float64, lastDirection string, countLastDirection int, qtyLongTermPercentCorrection float64, lastOrderCreationTime int64) float64 {
	buyQty := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return buyQty
	}

	if tradeCurrencyAmount == 0 {
		// @TODO: Check currency to sell
		buyQty = (baseCurrencyAmount / currentPrice) / 100 * float64(s.PercentRanges[0])
	}

	buyPercentRangeKey := s.getRangeKey("BUY", lastDirection, countLastDirection, s.PercentRanges, lastOrderCreationTime)

	// @TODO: Check currency to sell
	buyQty = (baseCurrencyAmount / currentPrice) / 100 * float64(s.PercentRanges[buyPercentRangeKey])

	correction := buyQty / 100 * qtyLongTermPercentCorrection
	buyQty += correction
	buyQty = s.round(buyQty, float64(s.PurchaseVolumePrecision))

	return buyQty
}

func (s *BuyCheapSellHigh) getSellQty(baseCurrencyAmount float64, tradeCurrencyAmount float64, lastDirection string, countLastDirection int, qtyLongTermPercentCorrection float64, lastOrderCreationTime int64) float64 {
	sellQty := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return sellQty
	}

	if baseCurrencyAmount == 0 {
		// @TODO: Check currency to sell
		sellQty = tradeCurrencyAmount / 100 * float64(s.PercentRanges[0])
	}

	sellPercentRangeKey := s.getRangeKey("SELL", lastDirection, countLastDirection, s.PercentRanges, lastOrderCreationTime)

	// @TODO: Check currency to sell
	sellQty = tradeCurrencyAmount / 100 * float64(s.PercentRanges[sellPercentRangeKey])

	correction := sellQty / 100 * qtyLongTermPercentCorrection
	sellQty -= correction
	sellQty = s.round(sellQty, float64(s.PurchaseVolumePrecision))

	return sellQty
}

// Percent for long term correction qty. Positive number to increase qty of trade currency amount, negative to decrease
// Maximum value can be 90%, minimum value can be -90%
func (s *BuyCheapSellHigh) qtyLongTermPercentCorrection(currentPrice float64, baseCurrencyAmount float64, tradeCurrencyAmount float64) float64 {

	// linePerc is value of percent you have to have in base currency
	maxPrice := s.LongTermMaxPrice
	minPrice := s.LongTermMinPrice
	percentBuffer := s.LongTermPercentBuffer
	linePerc := s.mathMap(currentPrice, minPrice, maxPrice, percentBuffer, 100-percentBuffer)

	// havePerc is value of percent you have in base currency regarding average price
	avgPrice := minPrice + ((maxPrice - minPrice) / 2)
	baseCurrencyTotal := avgPrice*tradeCurrencyAmount + baseCurrencyAmount
	havePerc := baseCurrencyAmount / (baseCurrencyTotal / 100)

	longTermCorrection := havePerc - linePerc

	if longTermCorrection > 90 {
		longTermCorrection = 90
	}
	if longTermCorrection < -90 {
		longTermCorrection = -90
	}

	return longTermCorrection
}

func (s *BuyCheapSellHigh) mathMap(value float64, fromMin float64, fromMax float64, toMin float64, toMax float64) float64 {
	return (value-fromMin)/(fromMax-fromMin)*(toMax-toMin) + toMin
}

func (s *BuyCheapSellHigh) getRangeKey(orderDirection string, lastDirection string, countLastDirection int, ranges []int64, lastOrderCreationTime int64) int {
	rangeKey := 0

	currentTime := time.Now().Unix()
	orderOldMinutes := (currentTime - lastOrderCreationTime) / 60
	if orderOldMinutes > s.MinutesToReducePriceRange {
		return rangeKey
	}

	if lastDirection == orderDirection {
		if len(ranges) > countLastDirection {
			rangeKey = countLastDirection
		} else {
			rangeKey = len(ranges) - 1
		}
	}
	return rangeKey
}

func (s *BuyCheapSellHigh) round(value float64, precision float64) float64 {
	roundNum := math.Pow(10, float64(precision))
	return math.Round(value*roundNum) / roundNum
}
