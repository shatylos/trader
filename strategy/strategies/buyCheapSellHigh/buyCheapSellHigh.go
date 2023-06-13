package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain"
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/structs"
	tradeConst "bitbucket.org/shatylos/trader/trading/constant"
	"bitbucket.org/shatylos/trader/trading/services"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"math"
	"strconv"
	"time"
)

type BuyCheapSellHigh struct {
	CoinPare       string
	DomainCode     string
	MainCurrency   string
	TradeCurrency  string
	isInit         bool
	Resolution     string
	TimeoutSeconds time.Duration
	CostRanges     []int64
	PercentRanges  []int64
}

//var candles []structs.DomainCandle
//var orders []structs.DomainOrder
//var positions []structs.DomainPosition
//var ordersToOpen []structs.DomainOrderRequest
//var positionsToOpen []structs.DomainPositionRequest

func (s *BuyCheapSellHigh) IsInit() bool {
	return s.isInit
}

func (s *BuyCheapSellHigh) Initialise() error {
	domainItem, err := domain.GetDomainInterface(s.DomainCode)
	if err != nil {
		return err
	}
	if domainItem.GetType() != constant.DomainTypeSpot {
		return utils.AppError{
			Message: "Strategy buyCheapSellHigh works only with spot domain type",
		}
	}
	s.isInit = true
	return nil
}

func (s *BuyCheapSellHigh) DoAction() error {

	openOrders, err := s.getOpenOrders()
	if err != nil {
		return err
	}

	if len(openOrders) >= 2 {
		return nil
	}

	if len(openOrders) == 1 {
		err = s.cancelAllOrder(openOrders[0])
		if err != nil {
			return err
		}
	}

	historyOrders, err := s.getHistoryOrders()
	if err != nil {
		return err
	}

	baseCurrencyAmount, tradeCurrencyAmount, err := s.getCurrencyAmounts()
	if err != nil {
		return err
	}

	buyPrice, buyQty, sellPrice, sellQty, err := s.getBuySellPriceQty(historyOrders, baseCurrencyAmount, tradeCurrencyAmount)
	if err != nil {
		return err
	}

	if buyPrice > 0 && buyQty > 0 {
		err = s.setLimitOrder(buyPrice, buyQty, tradeConst.SideBuy)
		if err != nil {
			return err
		}
	} else {
		utils.LogInfo(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for buy orders: buyPrice: %f, buyQty: %f", buyPrice, buyQty))
	}
	if sellPrice > 0 && sellQty > 0 {
		err = s.setLimitOrder(sellPrice, sellQty, tradeConst.SideSell)
		if err != nil {
			return err
		}
	} else {
		utils.LogInfo(fmt.Sprintf("[Buy Cheap Sell High] unexpected values for sell orders: sellPrice: %f, sellQty: %f", sellPrice, sellQty))
	}

	return nil
}

func (s *BuyCheapSellHigh) Wait() {
	time.Sleep(time.Second * s.TimeoutSeconds)
}

func (s *BuyCheapSellHigh) getOpenOrders() ([]structs.DomainOrder, error) {
	return services.GetOpenOrderList(s.DomainCode, s.CoinPare)
}

func (s *BuyCheapSellHigh) cancelAllOrder(order structs.DomainOrder) error {
	return services.CancelOrder(s.DomainCode, order.OrderId)
}

func (s *BuyCheapSellHigh) getHistoryOrders() ([]structs.DomainOrder, error) {
	return services.GetHistoryOrders(s.DomainCode, 50)
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

func (s *BuyCheapSellHigh) getBuySellPriceQty(historyOrders []structs.DomainOrder, baseCurrencyAmount float64, tradeCurrencyAmount float64) (float64, float64, float64, float64, error) {

	currentPrice, err := s.getCurrentPrice()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	lastDirection := ""
	countLastDirection := 0
	lastOrderPrice := float64(0)

	for _, historyOrder := range historyOrders {
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

	buyPrice := float64(0)
	buyQty := float64(0)
	sellPrice := float64(0)
	sellQty := float64(0)

	if baseCurrencyAmount == 0 && tradeCurrencyAmount == 0 {
		return 0, 0, 0, 0, nil
	} else if baseCurrencyAmount == 0 {
		sellPrice = currentPrice
		// @TODO: Check currency to sell
		sellQty = tradeCurrencyAmount / 100 * float64(s.PercentRanges[0])
		//sellQty = (tradeCurrencyAmount * currentPrice) / 100 * float64(s.PercentRanges[0])
	} else if tradeCurrencyAmount == 0 {
		buyPrice = currentPrice
		// @TODO: Check currency to sell
		buyQty = (baseCurrencyAmount / currentPrice) / 100 * float64(s.PercentRanges[0])
		//buyQty = baseCurrencyAmount / 100 * float64(s.PercentRanges[0])
	}

	priceToCalcute := currentPrice
	if lastOrderPrice > 0 {
		priceToCalcute = lastOrderPrice
	}

	buyPriceRangeKey := s.getRangeKey("BUY", lastDirection, countLastDirection, s.CostRanges)
	sellPriceRangeKey := s.getRangeKey("SELL", lastDirection, countLastDirection, s.CostRanges)
	buyPercentRangeKey := s.getRangeKey("BUY", lastDirection, countLastDirection, s.PercentRanges)
	sellPercentRangeKey := s.getRangeKey("SELL", lastDirection, countLastDirection, s.PercentRanges)

	buyPrice = priceToCalcute - float64(s.CostRanges[buyPriceRangeKey])
	// @TODO: Check currency to sell
	buyQty = (baseCurrencyAmount / currentPrice) / 100 * float64(s.PercentRanges[buyPercentRangeKey])

	sellPrice = priceToCalcute + float64(s.CostRanges[sellPriceRangeKey])
	// @TODO: Check currency to sell
	sellQty = tradeCurrencyAmount / 100 * float64(s.PercentRanges[sellPercentRangeKey])
	//sellQty = (tradeCurrencyAmount * currentPrice) / 100 * float64(s.PercentRanges[sellPercentRangeKey])

	return math.Round(buyPrice*1e6) / 1e6,
		math.Round(buyQty*1e6) / 1e6,
		math.Round(sellPrice*1e6) / 1e6,
		math.Round(sellQty*1e6) / 1e6,
		nil
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

func (s *BuyCheapSellHigh) getCurrencyAmounts() (float64, float64, error) {
	wallet, err := services.LoadWalletInfo(s.DomainCode)
	if err != nil {
		return 0, 0, err
	}
	mainCurrencyAmount := float64(0)
	tradeCurrencyAmount := float64(0)

	for _, coin := range wallet.Available {
		if coin.Coin == s.MainCurrency {
			mainCurrencyAmount = coin.Amount
		}
		if coin.Coin == s.TradeCurrency {
			tradeCurrencyAmount = coin.Amount
		}
	}

	return mainCurrencyAmount, tradeCurrencyAmount, nil
}

func (s *BuyCheapSellHigh) setLimitOrder(price float64, qty float64, direction string) error {

	request := structs.DomainOrderRequest{
		OrderId:     strconv.FormatInt(time.Now().UnixNano(), 10),
		Price:       price,
		Qty:         qty,
		ReduceOnly:  false,
		Side:        direction,
		Symbol:      s.CoinPare,
		TimeInForce: "GTC",
		Type:        "LIMIT",
	}

	orderId, err := services.OpenOrder(s.DomainCode, request)
	if err != nil {
		return err
	}

	utils.LogInfo(fmt.Sprintf("Created order: %s", orderId))
	return nil
}
