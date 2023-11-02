package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain/structs"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"strconv"
	"time"
)

func (s *BuyCheapSellHigh) getOpenOrders() ([]structs.DomainOrder, error) {
	return s.Domain.GetOpenOrderList(s.CoinPare)
}

func (s *BuyCheapSellHigh) getHistoryOrders() ([]structs.DomainOrder, error) {
	return s.Domain.GetHistoryOrders(50)
}

func (s *BuyCheapSellHigh) cancelAllOrder(order structs.DomainOrder) error {
	return s.Domain.CancelOrder(order.OrderId)
}

func (s *BuyCheapSellHigh) cancelOldOrdersWithBigRanges(orders []structs.DomainOrder) error {

	minSellPrice := 0.0
	maxBuyPrice := 0.0

	for _, order := range orders {
		if order.Side == "SELL" && (order.Price < minSellPrice || minSellPrice == 0) {
			minSellPrice = order.Price
		}
		if order.Side == "BUY" && order.Price > maxBuyPrice {
			maxBuyPrice = order.Price
		}
	}

	priceDiff := minSellPrice - maxBuyPrice
	if float64(s.CostRanges[0])*2 >= priceDiff {
		return nil
	}

	currentTime := time.Now().Unix()
	for _, order := range orders {
		orderOldMinutes := (currentTime - order.CreatedTime) / 60
		if orderOldMinutes > s.MinutesToReducePriceRange {
			err := s.Domain.CancelOrder(order.OrderId)
			if err != nil {
				return err
			}
			utils.LogInfo(fmt.Sprintf("Canceled order as it is old with big price range. OrderId: %s, Symbol: %s, Side: %s, Price: %f, Qty: %f", order.OrderId, s.CoinPare, order.Side, order.Price, order.Qty))
		}
	}

	return nil
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

	orderId, err := s.Domain.OpenOrder(request)
	if err != nil {
		return err
	}

	utils.LogSuccess(fmt.Sprintf("Created order: %s. Symbol: %s, Side: %s, Price: %f, Qty: %f", orderId, s.CoinPare, direction, price, qty))
	return nil
}
