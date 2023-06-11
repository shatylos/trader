package bybit

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/bybit/request"
	"bitbucket.org/shatylos/trader/domain/structs"
	"strconv"
	"strings"
)

type DomainBybitSpot struct {
	IsDemo bool
}

func (d *DomainBybitSpot) GetType() int64 {
	return constant.DomainTypeSpot
}

func (d *DomainBybitSpot) GetWallet() (*structs.DomainWallet, error) {
	walletBalances, er := request.GetSpotWalletBalance(d.IsDemoMode())
	if er != nil {
		return nil, er
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	for coinCode, walletBalance := range *walletBalances {
		free, err := strconv.ParseFloat(walletBalance.Free, 64)
		if err != nil {
			return nil, err
		}
		locked, err := strconv.ParseFloat(walletBalance.Locked, 64)
		if err != nil {
			return nil, err
		}
		availableCoins = append(availableCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: free,
		})
		reservedCoins = append(reservedCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: locked,
		})
	}

	result := structs.DomainWallet{
		Available: availableCoins,
		Reserved:  reservedCoins,
	}

	return &result, nil
}

func (d *DomainBybitSpot) IsDemoMode() bool {
	return d.IsDemo
}

func (d *DomainBybitSpot) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error) {
	candles, err := request.GetSpotKlineList(symbol, resolution, from, limit, d.IsDemo)
	if err != nil {
		return nil, err
	}

	candlesResult := make([]structs.DomainCandle, len(candles))

	for i, candle := range candles {
		highPrice, err := strconv.ParseFloat(candle.High, 64)
		if err != nil {
			return nil, err
		}
		lowPrice, err := strconv.ParseFloat(candle.Low, 64)
		if err != nil {
			return nil, err
		}
		openPrice, err := strconv.ParseFloat(candle.Open, 64)
		if err != nil {
			return nil, err
		}
		closePrice, err := strconv.ParseFloat(candle.Close, 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(candle.Volume, 64)
		if err != nil {
			return nil, err
		}

		candlesResult[i] = structs.DomainCandle{
			Time:   candle.Timestamp,
			High:   highPrice,
			Low:    lowPrice,
			Open:   openPrice,
			Close:  closePrice,
			Volume: volume,
		}
	}

	return candlesResult, nil
}

func (d *DomainBybitSpot) GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error) {
	orders, err := request.GetSpotOpenOrderList(coinPare, d.IsDemo)
	if err != nil {
		return nil, err
	}

	domainOrders := make([]structs.DomainOrder, len(orders))

	for key, order := range orders {
		orderPrice, err := strconv.ParseFloat(order.OrderPrice, 64)
		if err != nil {
			return nil, err
		}
		orderQty, err := strconv.ParseFloat(order.OrderQty, 64)
		if err != nil {
			return nil, err
		}

		domainOrder := structs.DomainOrder{
			CreatedTime: strconv.FormatInt(order.CreateTime, 10),
			OrderId:     order.OrderId,
			OrderStatus: order.Status,
			OrderType:   order.OrderType,
			Price:       orderPrice,
			Qty:         orderQty,
			ReduceOnly:  false,
			Side:        order.Side,
			Symbol:      order.Symbol,
			TimeInForce: order.TimeInForce,
			UpdatedTime: strconv.FormatInt(order.UpdateTime, 10),
		}
		domainOrders[key] = domainOrder
	}

	return domainOrders, nil
}

func (d *DomainBybitSpot) GetPositionList(coinPare string) ([]structs.DomainPosition, error) {
	panic("Not implemented")
}

func (d *DomainBybitSpot) OpenPosition(positionRequest structs.DomainPositionRequest) (string, error) {
	panic("Not implemented")
}

func (d *DomainBybitSpot) OpenOrder(orderRequest structs.DomainOrderRequest) (string, error) {

	domainOrderRequest := request.SpotOrderRequest{
		Symbol:      orderRequest.Symbol,
		OrderQty:    strconv.FormatFloat(orderRequest.Qty, 'f', -1, 64),
		Side:        strings.ToUpper(orderRequest.Side),
		OrderType:   orderRequest.Type,
		TimeInForce: orderRequest.TimeInForce,
		OrderPrice:  strconv.FormatFloat(orderRequest.Price, 'f', -1, 64),
		OrderLinkId: orderRequest.OrderId,
	}

	order, err := request.CreateSpotOrder(domainOrderRequest, d.IsDemo)
	if err != nil {
		return "", err
	}

	return order.OrderId, nil

}

func (d *DomainBybitSpot) CancelOrder(orderId string) error {
	err := request.CancelSpotOrder(orderId, d.IsDemo)
	if err != nil {
		return err
	}
	return nil
}

func (d *DomainBybitSpot) GetHistoryOrders(limit int64) ([]structs.DomainOrder, error) {
	orders, err := request.GetSpotOrderHistory(limit, d.IsDemo)
	if err != nil {
		return nil, err
	}

	domainOrders := make([]structs.DomainOrder, len(orders))

	for key, order := range orders {
		orderPrice, err := strconv.ParseFloat(order.OrderPrice, 64)
		if err != nil {
			return nil, err
		}
		avgPrice, err := strconv.ParseFloat(order.AvgPrice, 64)
		if err != nil {
			return nil, err
		}
		price := float64(0)
		if avgPrice > 0 {
			price = avgPrice
		} else {
			price = orderPrice
		}
		orderQty, err := strconv.ParseFloat(order.OrderQty, 64)
		if err != nil {
			return nil, err
		}

		domainOrder := structs.DomainOrder{
			CreatedTime: strconv.FormatInt(order.CreateTime, 10),
			OrderId:     order.OrderId,
			OrderStatus: order.Status,
			OrderType:   order.OrderType,
			Price:       price,
			Qty:         orderQty,
			ReduceOnly:  false,
			Side:        order.Side,
			Symbol:      order.Symbol,
			TimeInForce: order.TimeInForce,
			UpdatedTime: strconv.FormatInt(order.UpdateTime, 10),
		}
		domainOrders[key] = domainOrder
	}

	return domainOrders, nil
}
