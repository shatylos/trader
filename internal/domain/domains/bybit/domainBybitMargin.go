package bybit

import (
	"github.com/shatylos/trader/internal/domain/constant"
	"github.com/shatylos/trader/internal/domain/domains/bybit/request"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/type"
)

const orderStatusOpenNew = "New"
const orderStatusOpenPartiallyFilled = "PartiallyFilled"
const orderStatusOpenUntriggered = "Untriggered"
const orderStatusClosedRejected = "Rejected"
const orderStatusClosedPartiallyFilledCanceled = "PartiallyFilledCanceled"
const orderStatusClosedFilled = "Filled"
const orderStatusClosedCancelled = "Cancelled"
const orderStatusClosedTriggered = "Triggered"
const orderStatusClosedDeactivated = "Deactivated"

type DomainBybitMargin struct {
	code     string
	secrets  bybitStructs.Secrets
	leverage int64
}

func (d *DomainBybitMargin) GetCode() string {
	//TODO implement me
	panic("implement me")
}

func (d *DomainBybitMargin) SetConfig(config map[interface{}]interface{}) error {

	secretMap, ok := config["secrets"].(map[interface{}]interface{})
	if !ok {
		return tools.AppError{
			Message: "The field secrets is empty or contains not correct value type. In DomainBybitSpot config. Expects a map with \"key\", \"pass\" and \"endpoint\" keys",
		}
	}

	domainCode, err := _type.ToString(config["code"])
	if err != nil {
		return tools.AppError{
			Message: "The field code is empty or contains not correct value type. In DomainBybitSpot config. Expects a string",
		}
	}
	d.code = domainCode

	secrets := bybitStructs.Secrets{}

	apiEndpoint, err := _type.ToString(secretMap["endpoint"])
	if err != nil {
		return tools.AppError{
			Message: "The field secrets.endpoint is empty or contains not correct value type. In DomainBybitSpot config. Expects a string",
		}
	}
	secrets.ApiEndpoint = apiEndpoint

	key, err := _type.ToString(secretMap["key"])
	if err != nil {
		return tools.AppError{
			Message: "The field secrets.key is empty or contains not correct value type. In DomainBybitSpot config. Expects a string",
		}
	}
	secrets.Key = key

	pass, err := _type.ToString(secretMap["pass"])
	if err != nil {
		return tools.AppError{
			Message: "The field secrets.pass is empty or contains not correct value type. In DomainBybitSpot config. Expects a string",
		}
	}
	secrets.Pass = pass
	d.secrets = secrets

	return nil
}

func (d *DomainBybitMargin) GetType() int64 {
	return constant.DomainTypeMargin
}

func (d *DomainBybitMargin) GetWallet() (*structs.DomainWallet, error) {
	walletBalance, er := request.GetMarginWalletBalance(d.secrets)
	if er != nil {
		return nil, er
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	availableCoins = append(availableCoins, structs.DomainWalletCoinItem{
		Coin:   "USDT",
		Amount: walletBalance.TotalAvailableBalance,
	})
	reservedCoins = append(reservedCoins, structs.DomainWalletCoinItem{
		Coin:   "USDT",
		Amount: walletBalance.TotalInitialMargin,
	})

	result := structs.DomainWallet{
		DomainCode: d.code,
		Available:  availableCoins,
		Reserved:   reservedCoins,
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
		startTime, err := _type.ToInt64(candle.StartTime)
		if err != nil {
			return nil, err
		}
		startTime = startTime / 1000
		highPrice, err := _type.ToFloat64(candle.HighPrice)
		if err != nil {
			return nil, err
		}
		lowPrice, err := _type.ToFloat64(candle.LowPrice)
		if err != nil {
			return nil, err
		}
		openPrice, err := _type.ToFloat64(candle.OpenPrice)
		if err != nil {
			return nil, err
		}
		closePrice, err := _type.ToFloat64(candle.ClosePrice)
		if err != nil {
			return nil, err
		}
		volume, err := _type.ToFloat64(candle.Volume)
		if err != nil {
			return nil, err
		}
		candlesResult[i] = structs.DomainCandle{
			Time:   startTime,
			High:   highPrice,
			Low:    lowPrice,
			Open:   openPrice,
			Close:  closePrice,
			Volume: volume,
		}
	}

	return candlesResult, nil
}

func (d *DomainBybitMargin) GetOrder(orderId string) (order structs.DomainOrder, err error) {
	var orders []*request.OrderResponse
	orders, err = request.GetOrderList("", orderId, d.secrets)
	if err != nil {
		return
	}
	if len(orders) == 0 {
		err = tools.AppError{
			Message: "Order not found",
		}
		return
	}
	order, err = d.mapOrder(orders[0])
	return
}

func (d *DomainBybitMargin) GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error) {
	orders, err := request.GetOrderList(coinPare, "", d.secrets)
	if err != nil {
		return nil, err
	}

	domainOrders := make([]structs.DomainOrder, len(orders))

	for key, order := range orders {
		domainOrder, err := d.mapOrder(order)
		if err != nil {
			return nil, err
		}
		domainOrders[key] = domainOrder
	}

	return domainOrders, nil
}

func (d *DomainBybitMargin) mapOrder(order *request.OrderResponse) (domainOrder structs.DomainOrder, err error) {
	createdTime, er := _type.ToInt64(order.CreatedTime)
	if er != nil {
		err = er
		return
	}
	createdTime = createdTime / 1000
	updatedTime, er := _type.ToInt64(order.UpdatedTime)
	if er != nil {
		err = er
		return
	}
	updatedTime = updatedTime / 1000
	price, er := _type.ToFloat64(order.AvgPrice)
	if er != nil {
		err = er
		return
	}
	qty, er := _type.ToFloat64(order.Qty)
	if er != nil {
		err = er
		return
	}
	takeProfit, er := _type.ToFloat64(order.TakeProfit)
	if er != nil {
		err = er
		return
	}
	stopLoss, er := _type.ToFloat64(order.StopLoss)
	if er != nil {
		err = er
		return
	}
	status := d.mapStatus(order.OrderStatus)
	domainOrder = structs.DomainOrder{
		CreatedTime: createdTime,
		OrderId:     order.OrderId,
		OrderStatus: status,
		OrderType:   order.OrderType,
		Price:       price,
		Qty:         qty,
		ReduceOnly:  order.ReduceOnly,
		Side:        order.Side,
		Symbol:      order.Symbol,
		TimeInForce: order.TimeInForce,
		UpdatedTime: updatedTime,
		TakeProfit:  takeProfit,
		StopLoss:    stopLoss,
	}
	return
}

func (d *DomainBybitMargin) mapStatus(providerStatus string) string {
	switch providerStatus {
	case orderStatusOpenNew:
		return structs.OrderStatusOpen
	case orderStatusOpenPartiallyFilled:
		return structs.OrderStatusOpen
	case orderStatusOpenUntriggered:
		return structs.OrderStatusOpen
	case orderStatusClosedRejected:
		return structs.OrderStatusCancelled
	case orderStatusClosedPartiallyFilledCanceled:
		return structs.OrderStatusFilled
	case orderStatusClosedFilled:
		return structs.OrderStatusFilled
	case orderStatusClosedCancelled:
		return structs.OrderStatusCancelled
	case orderStatusClosedTriggered:
		return structs.OrderStatusCancelled
	case orderStatusClosedDeactivated:
		return structs.OrderStatusCancelled
	}
	return ""
}

func (d *DomainBybitMargin) GetPosition(coinPare string) (resultPosition structs.DomainPosition, err error) {
	var providerPosition request.Position
	providerPosition, err = request.GetPosition(coinPare, d.secrets)
	if err != nil {
		return
	}
	resultPosition.Symbol = providerPosition.Symbol
	resultPosition.Side = providerPosition.Side

	if providerPosition.AvgPrice != "" {
		resultPosition.AvgPrice, err = _type.ToFloat64(providerPosition.AvgPrice)
		if err != nil {
			return
		}
	}
	if providerPosition.MarkPrice != "" {
		resultPosition.MarkPrice, err = _type.ToFloat64(providerPosition.MarkPrice)
		if err != nil {
			return
		}
	}
	resultPosition.Leverage, err = _type.ToInt64(providerPosition.Leverage)
	if err != nil {
		return
	}
	resultPosition.Size, err = _type.ToFloat64(providerPosition.Size)
	if err != nil {
		return
	}
	if providerPosition.StopLoss != "" {
		resultPosition.StopLoss, err = _type.ToFloat64(providerPosition.StopLoss)
		if err != nil {
			return
		}
	}
	if providerPosition.TakeProfit != "" {
		resultPosition.TakeProfit, err = _type.ToFloat64(providerPosition.TakeProfit)
		if err != nil {
			return
		}
	}
	if providerPosition.UnrealizedPnl != "" {
		resultPosition.UnrealizedPnl, err = _type.ToFloat64(providerPosition.UnrealizedPnl)
		if err != nil {
			return
		}
	}
	if providerPosition.CurRealisedPnl != "" {
		resultPosition.RealizedPnl, err = _type.ToFloat64(providerPosition.CurRealisedPnl)
		if err != nil {
			return
		}
	}
	resultPosition.TotalPnl = resultPosition.RealizedPnl + resultPosition.UnrealizedPnl

	return
}

func (d *DomainBybitMargin) setLeverage(symbol string, leverage int64) (err error) {
	leverageRequest := request.LeverageRequest{
		Symbol:       symbol,
		BuyLeverage:  leverage,
		SellLeverage: leverage,
	}

	err = request.SetLeverage(leverageRequest, d.secrets)
	return
}

func (d *DomainBybitMargin) OpenPosition(positionRequest structs.DomainPositionRequest) (orderId string, err error) {

	if d.leverage == 0 && positionRequest.Leverage > 0 {
		err = d.setLeverage(positionRequest.Symbol, positionRequest.Leverage)
		if err != nil {
			return
		}
		d.leverage = positionRequest.Leverage
	}

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

	var order *request.OrderResponse
	order, err = request.CreateOrder(orderRequest, d.secrets)
	if err != nil {
		return
	}

	orderId = order.OrderId
	return
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

func (d *DomainBybitMargin) CancelOrder(orderId string, coinPare string) error {
	panic("Not implemented")
}

func (d *DomainBybitMargin) GetHistoryOrders(limit int64, coinPare string) ([]structs.DomainOrder, error) {
	panic("Not implemented")
}
