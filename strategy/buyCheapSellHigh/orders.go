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
