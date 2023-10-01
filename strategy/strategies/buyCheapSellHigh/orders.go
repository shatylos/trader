package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain/structs"
	"bitbucket.org/shatylos/trader/trading/services"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
	"strconv"
	"time"
)

func (s *BuyCheapSellHigh) getOpenOrders() ([]structs.DomainOrder, error) {
	return services.GetOpenOrderList(s.DomainCode, s.CoinPare)
}

func (s *BuyCheapSellHigh) getHistoryOrders() ([]structs.DomainOrder, error) {
	return services.GetHistoryOrders(s.DomainCode, 50)
}

func (s *BuyCheapSellHigh) cancelAllOrder(order structs.DomainOrder) error {
	return services.CancelOrder(s.DomainCode, order.OrderId)
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
