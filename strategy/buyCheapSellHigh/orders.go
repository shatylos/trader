package buyCheapSellHigh

import (
	"bitbucket.org/shatylos/trader/domain/structs"
	_storage "bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage"
	storageStructs "bitbucket.org/shatylos/trader/strategy/buyCheapSellHigh/storage/structs"
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

func (s *BuyCheapSellHigh) setLimitOrder(price float64, qty float64, direction string) (string, error) {

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
		return "", err
	}

	utils.LogSuccess(fmt.Sprintf("Created order: %s. Symbol: %s, Side: %s, Price: %f, Qty: %f", orderId, s.CoinPare, direction, price, qty))
	return orderId, nil
}

func (s *BuyCheapSellHigh) setOrderToStorage(orderId string, mainCurrencyBalance float64, tradeCurrencyBalance float64) error {

	storage, err := _storage.GetStorage(s.Id)
	if err != nil {
		return err
	}

	_, err = (*storage).AddDomainOrderOnce(storageStructs.HistoryOrder{
		DomainOrderId:             orderId,
		CreatedTime:               time.Now().Unix(),
		MainCurrencyAmountBefore:  mainCurrencyBalance,
		TradeCurrencyAmountBefore: tradeCurrencyBalance,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *BuyCheapSellHigh) fillPrices() error {

	utils.LogInfo("Filling prices")

	storagePointer, err := _storage.GetStorage(s.Id)
	if err != nil {
		return err
	}
	storage := *storagePointer

	historyOrders, err := storage.GetNotCalculatedDomainOrders()
	if err != nil {
		return err
	}

	countFilled := 0
	countRemoved := 0
	for _, historyOrder := range historyOrders {
		if historyOrder.FilledPrice == 0 || historyOrder.FilledQty == 0 || historyOrder.Side == "" {
			order, err := s.Domain.GetOrder(historyOrder.DomainOrderId)
			if err != nil {
				return err
			}

			if order.OrderStatus == "FILLED" || order.OrderStatus == "PARTIALLY_FILLED" {
				historyOrder.FilledPrice = order.Price
				historyOrder.FilledQty = order.Qty
				historyOrder.Side = order.Side
				historyOrder.UpdatedTime = order.UpdatedTime
				err := storage.UpdateOrder(historyOrder)
				if err != nil {
					utils.LogError(fmt.Sprintf("Error filling price for the order %s", order.OrderId))
					return err
				}
				utils.LogInfo(fmt.Sprintf("Filled price for the order %s", order.OrderId))
				countFilled++
			} else if order.OrderStatus == "CANCELED" {
				err := storage.RemoveOrder(historyOrder.DomainOrderId)
				if err != nil {
					utils.LogError(fmt.Sprintf("Error removing canceled order %s", order.OrderId))
					return err
				}
				utils.LogInfo(fmt.Sprintf("Removed canceled order %s", order.OrderId))
				countRemoved++
			} else if order.OrderStatus != "NEW" {
				utils.LogWarning(fmt.Sprintf("Unexpected order status: %s, for order %s", order.OrderStatus, order.OrderId))
			}
		}
	}
	if countFilled > 0 {
		utils.LogInfo(fmt.Sprintf("Filled prices for %d orders", countFilled))
	}
	if countRemoved > 0 {
		utils.LogInfo(fmt.Sprintf("Removed %d cancelled orders from local storage", countRemoved))
	}
	return nil
}
