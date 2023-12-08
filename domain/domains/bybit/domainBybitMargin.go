package bybit

import (
	"bitbucket.org/shatylos/trader/domain/constant"
	"bitbucket.org/shatylos/trader/domain/domains/bybit/request"
	bybitStructs "bitbucket.org/shatylos/trader/domain/domains/bybit/structs"
	"bitbucket.org/shatylos/trader/domain/structs"
	"bitbucket.org/shatylos/trader/utils"
)

type DomainBybitMargin struct {
	secrets bybitStructs.Secrets
}

func (d *DomainBybitMargin) GetCode() string {
	//TODO implement me
	panic("implement me")
}

func (d *DomainBybitMargin) SetConfig(m map[interface{}]interface{}) error {
	//TODO implement me
	d.secrets = bybitStructs.Secrets{}
	panic("implement me")
}

func (d *DomainBybitMargin) GetType() int64 {
	return constant.DomainTypeMargin
}

func (d *DomainBybitMargin) GetWallet() (*structs.DomainWallet, error) {
	walletBalances, er := request.GetWalletBalance(d.secrets)
	if er != nil {
		return nil, er
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	for coinCode, walletBalance := range *walletBalances {
		if walletBalance.AvailableBalance == 0 && walletBalance.UsedMargin == 0 {
			continue
		}
		availableCoins = append(availableCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: walletBalance.AvailableBalance,
		})
		reservedCoins = append(reservedCoins, structs.DomainWalletCoinItem{
			Coin:   coinCode,
			Amount: walletBalance.UsedMargin,
		})
	}

	result := structs.DomainWallet{
		Available: availableCoins,
		Reserved:  reservedCoins,
	}

	return &result, nil
}

func (d *DomainBybitMargin) LoadCandleHistory(symbol string, resolution string, from int64, limit int64) ([]structs.DomainCandle, error) {

	candles, err := request.GetKlineList(symbol, resolution, from, limit, d.secrets)
	if err != nil {
		return nil, err
	}

	candlesResult := make([]structs.DomainCandle, len(candles))

	for i, candle := range candles {
		candlesResult[i] = structs.DomainCandle{
			Time:   int64(candle.OpenTime),
			High:   candle.High,
			Low:    candle.Low,
			Open:   candle.Open,
			Close:  candle.Close,
			Volume: candle.Volume,
		}
	}

	return candlesResult, nil
}

func (d *DomainBybitMargin) GetOrder(domainId string) (structs.DomainOrder, error) {
	panic("Not implemented")
}

func (d *DomainBybitMargin) GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error) {
	orders, err := request.GetOrderList(coinPare, "New", d.secrets)
	if err != nil {
		return nil, err
	}

	domainOrders := make([]structs.DomainOrder, len(orders))

	for key, order := range orders {

		createdTime, err := utils.ToInt64(order.CreatedTime)
		if err != nil {
			return nil, err
		}
		panic("Check createdTime and UpdatedTime. It was refactored but not tested")
		domainOrder := structs.DomainOrder{
			CreatedTime: createdTime,
			OrderId:     order.OrderId,
			OrderStatus: order.OrderStatus,
			OrderType:   order.OrderType,
			Price:       order.Price,
			Qty:         order.Qty,
			ReduceOnly:  order.ReduceOnly,
			Side:        order.Side,
			Symbol:      order.Symbol,
			TimeInForce: order.TimeInForce,
			UpdatedTime: 123, // order.UpdatedTime,
		}
		domainOrders[key] = domainOrder
	}

	return domainOrders, nil
}

func (d *DomainBybitMargin) GetPositionList(coinPare string) ([]structs.DomainPosition, error) {

	positions, err := request.GetPositionList(coinPare, d.secrets)
	if err != nil {
		return nil, err
	}
	resultPositions := make([]structs.DomainPosition, len(positions))

	for i, position := range positions {
		resultPosition := structs.DomainPosition{
			Price:         position.EntryPrice,
			Leverage:      int64(position.Leverage),
			Symbol:        position.Symbol,
			Qty:           position.Size,
			RealizedPnl:   position.RealisedPnl,
			StopLoss:      position.StopLoss,
			TakeProfit:    position.TakeProfit,
			Side:          position.Side,
			UnrealizedPnl: position.UnrealisedPnl,
		}
		resultPositions[i] = resultPosition
	}

	return resultPositions, nil
}

func (d *DomainBybitMargin) OpenPosition(positionRequest structs.DomainPositionRequest) (string, error) {

	orderRequest := request.OrderRequest{
		CloseOnTrigger: false,
		OrderLinkId:    positionRequest.PositionId,
		OrderType:      positionRequest.Type,
		Price:          positionRequest.Price,
		Qty:            positionRequest.Qty,
		ReduceOnly:     positionRequest.ReduceOnly, // false,
		Side:           positionRequest.Side,
		StopLoss:       positionRequest.StopLoss,
		Symbol:         positionRequest.Symbol,
		TakeProfit:     positionRequest.TakeProfit,
		TimeInForce:    positionRequest.TimeInForce, // "FillOrKill", // FillOrKill GoodTillCancel ImmediateOrCancel PostOnly
		SlTriggerBy:    "LastPrice",                 // LastPrice IndexPrice MarkPrice
		TpTriggerBy:    "LastPrice",
	}

	order, err := request.CreateOrder(orderRequest, d.secrets)
	if err != nil {
		return "", err
	}

	return order.OrderId, nil
}

func (d *DomainBybitMargin) OpenOrder(orderRequest structs.DomainOrderRequest) (string, error) {

	domainOrderRequest := request.OrderRequest{
		CloseOnTrigger: false,
		OrderType:      orderRequest.Type,
		Price:          orderRequest.Price,
		Qty:            orderRequest.Qty,
		ReduceOnly:     orderRequest.ReduceOnly,
		Side:           orderRequest.Side,
		Symbol:         orderRequest.Symbol,
		TimeInForce:    orderRequest.TimeInForce,
	}

	order, err := request.CreateOrder(domainOrderRequest, d.secrets)
	if err != nil {
		return "", err
	}

	return order.OrderId, nil
}

func (d *DomainBybitMargin) CancelOrder(orderId string) error {
	panic("Not implemented")
}

func (d *DomainBybitMargin) GetHistoryOrders(limit int64) ([]structs.DomainOrder, error) {
	panic("Not implemented")
}
