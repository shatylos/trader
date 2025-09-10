package buyCheapSellHigh

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage"
	storageStructs "github.com/shatylos/trader/internal/strategy/buyCheapSellHigh/storage/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"strconv"
	"time"
)

func (s *BuyCheapSellHigh) getOpenOrders() ([]structs.DomainOrder, error) {
	return s.Domain.GetOpenOrderList(s.CoinPare)
}

func (s *BuyCheapSellHigh) getHistoryOrders() ([]structs.DomainOrder, error) {
	return s.Domain.GetHistoryOrders(50, s.CoinPare)
}

func (s *BuyCheapSellHigh) cancelAllOrder(order structs.DomainOrder) error {
	return s.Domain.CancelOrder(order.OrderId, s.CoinPare)
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
			err := s.Domain.CancelOrder(order.OrderId, s.CoinPare)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("Canceled order as it is old with big price range. OrderId: %s, Symbol: %s, Side: %s, Price: %f, Qty: %f", order.OrderId, s.CoinPare, order.Side, order.Price, order.Qty))
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

	logger.Success(fmt.Sprintf("Created order: %s. Symbol: %s, Side: %s, Price: %f, Qty: %f", orderId, s.CoinPare, direction, price, qty))
	return orderId, nil
}

func (s *BuyCheapSellHigh) setOrderToStorage(orderId string, mainCurrencyBalance float64, tradeCurrencyBalance float64) error {

	storage, err := strategyStorage.GetStorage(s.Id)
	if err != nil {
		return err
	}

	_, err = storage.AddDomainOrderOnce(storageStructs.HistoryOrder{
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

	logger.Info("Filling prices")

	storage, err := strategyStorage.GetStorage(s.Id)
	if err != nil {
		return err
	}

	historyOrders, err := storage.GetNotFilledHistoryOrders()
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

			if order.OrderStatus == structs.OrderStatuses.Filled || order.OrderStatus == structs.OrderStatuses.PartiallyFilled {
				historyOrder.FilledPrice = order.Price
				historyOrder.FilledQty = order.Qty
				historyOrder.Side = order.Side
				historyOrder.UpdatedTime = order.UpdatedTime
				err := storage.UpdateOrder(historyOrder)
				if err != nil {
					logger.Error(fmt.Sprintf("Error filling price for the order %s", order.OrderId))
					return err
				}
				logger.Info(fmt.Sprintf("Filled price for the order %s", order.OrderId))
				if order.Side == structs.OrderSideBuy {
					logger.Success(fmt.Sprintf("Bought %g %s for the %g %s", order.Qty, s.TradeCurrency, humanize.CommafWithDigits(order.Price, 2), s.MainCurrency))
				} else if order.Side == structs.OrderSideSell {
					logger.Success(fmt.Sprintf("Sold %g %s for the %g %s", order.Qty, s.TradeCurrency, humanize.CommafWithDigits(order.Price, 2), s.MainCurrency))
				}
				countFilled++
			} else if order.OrderStatus == structs.OrderStatuses.Canceled {
				err := storage.RemoveOrder(historyOrder.DomainOrderId)
				if err != nil {
					logger.Error(fmt.Sprintf("Error removing canceled order %s", order.OrderId))
					return err
				}
				logger.Info(fmt.Sprintf("Removed canceled order %s", order.OrderId))
				countRemoved++
			} else if order.OrderStatus != structs.OrderStatuses.New {
				logger.Warning(fmt.Sprintf("Unexpected order status: %s, for order %s", order.OrderStatus, order.OrderId))
			}
		}
	}
	if countFilled > 0 {
		logger.Info(fmt.Sprintf("Filled prices for %d orders", countFilled))
	}
	if countRemoved > 0 {
		logger.Info(fmt.Sprintf("Removed %d cancelled orders from local storage", countRemoved))
	}
	return nil
}

func (s *BuyCheapSellHigh) calculateHistoryOrderValues() error {

	logger.Info("Calculating average prices")

	storage, err := strategyStorage.GetStorage(s.Id)
	if err != nil {
		return err
	}

	historyOrders, err := storage.GetNotCalculatedHistoryOrders()
	if err != nil {
		return err
	}
	if len(historyOrders) == 0 {
		logger.Info("No orders to calculate average price")
		return nil
	}

	lastCalculatedOrder, err := storage.GetLastCalculatedOrder()
	if err != nil {
		return err
	}

	prevAveragePrice := 0.0
	if lastCalculatedOrder != nil {
		prevAveragePrice = lastCalculatedOrder.AveragePrice
	} else if s.ManualBuyPriceBeforeStart > 0 {
		prevAveragePrice = s.ManualBuyPriceBeforeStart
	} else if len(historyOrders) > 0 {
		prevAveragePrice = historyOrders[0].FilledPrice
	}

	if prevAveragePrice == 0 {
		return tools.AppError{
			Message: "can not get previous average price",
		}
	}

	for _, historyOrder := range historyOrders {
		err := s.calculateAveragePrice(&historyOrder, prevAveragePrice)
		if err != nil {
			return err
		}
		err = s.calculateRevenue(&historyOrder)
		if err != nil {
			return err
		}
		err = s.calculateCommission(&historyOrder)
		if err != nil {
			return err
		}

		err = storage.UpdateOrder(historyOrder)
		if err != nil {
			return err
		}
		prevAveragePrice = historyOrder.AveragePrice
	}

	logger.Info(fmt.Sprintf("Average prices were calculated for %d orders", len(historyOrders)))

	return nil
}

func (s *BuyCheapSellHigh) calculateAveragePrice(historyOrder *storageStructs.HistoryOrder, prevAveragePrice float64) error {
	if historyOrder.FilledPrice == 0 || historyOrder.FilledQty == 0 || historyOrder.Side == "" {
		return tools.AppError{
			Message: fmt.Sprintf("can not calculate average price for order %s", historyOrder.DomainOrderId),
		}
	}

	averagePrice := float64(0)
	if historyOrder.Side == "BUY" {
		averagePrice = math.Div(
			math.Mul(historyOrder.TradeCurrencyAmountBefore, prevAveragePrice)+math.Mul(historyOrder.FilledPrice, historyOrder.FilledQty),
			historyOrder.TradeCurrencyAmountBefore+historyOrder.FilledQty,
		)
	} else if historyOrder.Side == "SELL" {
		averagePrice = prevAveragePrice
	} else {
		return tools.AppError{
			Message: fmt.Sprintf("unexpected order side %s for order %s", historyOrder.Side, historyOrder.DomainOrderId),
		}
	}
	historyOrder.AveragePrice = averagePrice
	return nil
}

func (s *BuyCheapSellHigh) calculateRevenue(historyOrder *storageStructs.HistoryOrder) error {
	if historyOrder.Side == "BUY" {
		return nil
	}

	if historyOrder.Side == "SELL" {
		if historyOrder.FilledPrice == 0 || historyOrder.FilledQty == 0 || historyOrder.Side == "" || historyOrder.AveragePrice == 0 {
			return tools.AppError{
				Message: fmt.Sprintf("can not calculate revenue for order %s", historyOrder.DomainOrderId),
			}
		}

		orderCost := math.Mul(historyOrder.FilledPrice, historyOrder.FilledQty)
		averagePriceCost := math.Mul(historyOrder.AveragePrice, historyOrder.FilledQty)
		historyOrder.Revenue = orderCost - averagePriceCost
		return nil
	}

	return tools.AppError{
		Message: fmt.Sprintf("unexpected order side %s for order %s", historyOrder.Side, historyOrder.DomainOrderId),
	}
}

func (s *BuyCheapSellHigh) calculateCommission(historyOrder *storageStructs.HistoryOrder) error {
	if historyOrder.FilledPrice == 0 || historyOrder.FilledQty == 0 {
		return tools.AppError{
			Message: fmt.Sprintf("can not calculate commission for order %s. Not enough data.", historyOrder.DomainOrderId),
		}
	}
	onePercent := math.Div(historyOrder.FilledQty, 100)
	tradeCurrCommission := math.Mul(onePercent, s.CommissionPercent)
	historyOrder.Comission = math.Mul(tradeCurrCommission, historyOrder.FilledPrice)
	return nil
}
