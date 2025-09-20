package bybit

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain/domains/bybit/mapping"
	"github.com/shatylos/trader/internal/domain/domains/bybit/request"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/type"
	"sort"
	"strconv"
	"strings"
)

type DomainBybitSpot struct {
	code    string
	secrets bybitStructs.Secrets
}

func (d *DomainBybitSpot) GetCode() string {
	return d.code
}

func (d *DomainBybitSpot) SetConfig(config map[interface{}]interface{}) error {

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

func (d *DomainBybitSpot) GetWallet() (wallet structs.DomainWallet, err error) {
	var walletBalances *map[string]request.SpotWalletBalance
	walletBalances, err = request.GetSpotWalletBalance(d.secrets)
	if err != nil {
		return
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	for coinCode, walletBalance := range *walletBalances {
		var free float64
		free, err = strconv.ParseFloat(walletBalance.Free, 64)
		if err != nil {
			return
		}
		var locked float64
		locked, err = strconv.ParseFloat(walletBalance.Locked, 64)
		if err != nil {
			return
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

	wallet = structs.DomainWallet{
		Available: availableCoins,
		Reserved:  reservedCoins,
	}
	return
}

func (d *DomainBybitSpot) LoadCandleHistory(symbol string, resolution string, limit int64) ([]structs.DomainCandle, error) {
	providerResolution, err := mapping.ToBybitInterval(resolution)
	if err != nil {
		return nil, err
	}
	candles, err := request.GetSpotKlineList(symbol, providerResolution, limit, d.secrets)
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

	sort.Slice(candlesResult, func(i, j int) bool {
		return candlesResult[i].Time > candlesResult[j].Time
	})

	return candlesResult, nil
}

func (d *DomainBybitSpot) GetOrder(domainId string) (structs.DomainOrder, error) {
	order, err := request.GetSpotOrder(domainId, d.secrets)
	if err != nil {
		return structs.DomainOrder{}, err
	}

	if order.Symbol == "" {
		msg := fmt.Sprintf("Order with ID (%s) not found", domainId)
		err = tools.AppError{Message: msg}
		logger.Warning(msg)
		return structs.DomainOrder{}, err
	}

	filledPrice, err := strconv.ParseFloat(order.AvgPrice, 64)
	if err != nil {
		return structs.DomainOrder{}, err
	}
	filledQty, err := strconv.ParseFloat(order.ExecQty, 64)
	if err != nil {
		return structs.DomainOrder{}, err
	}
	createTime, err := strconv.ParseInt(order.CreateTime, 10, 64)
	if err != nil {
		return structs.DomainOrder{}, err
	}
	updateTime, err := strconv.ParseInt(order.UpdateTime, 10, 64)
	if err != nil {
		return structs.DomainOrder{}, err
	}

	domainOrder := structs.DomainOrder{
		CreatedTime: createTime / 1000,
		OrderId:     order.OrderId,
		OrderStatus: order.Status,
		OrderType:   order.OrderType,
		Price:       filledPrice,
		Qty:         filledQty,
		ReduceOnly:  false,
		Side:        order.Side,
		Symbol:      order.Symbol,
		TimeInForce: order.TimeInForce,
		UpdatedTime: updateTime / 1000,
	}

	return domainOrder, nil
}

func (d *DomainBybitSpot) GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error) {
	orders, err := request.GetSpotOpenOrderList(coinPare, d.secrets)
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
		createTime, err := _type.ToInt64(order.CreateTime)
		if err != nil {
			return nil, err
		}
		updateTime, err := _type.ToInt64(order.UpdateTime)
		if err != nil {
			return nil, err
		}

		domainOrder := structs.DomainOrder{
			CreatedTime: createTime / 1000,
			OrderId:     order.OrderId,
			OrderStatus: order.Status,
			OrderType:   order.OrderType,
			Price:       orderPrice,
			Qty:         orderQty,
			ReduceOnly:  false,
			Side:        order.Side,
			Symbol:      order.Symbol,
			TimeInForce: order.TimeInForce,
			UpdatedTime: updateTime / 1000,
		}
		domainOrders[key] = domainOrder
	}

	sort.Slice(domainOrders, func(i, j int) bool {
		return domainOrders[i].CreatedTime > domainOrders[j].CreatedTime
	})

	return domainOrders, nil
}

func (d *DomainBybitSpot) OpenOrder(orderRequest structs.DomainOrderRequest) (orderId string, err error) {
	var orderType string
	orderType, err = mapping.ToBybitOrderType(orderRequest.Type)
	if err != nil {
		return
	}

	orderPriceStr := ""
	if orderRequest.Price > 0 {
		strconv.FormatFloat(orderRequest.Price, 'f', -1, 64)
	}

	domainOrderRequest := request.SpotOrderRequest{
		Symbol:      orderRequest.Symbol,
		OrderQty:    strconv.FormatFloat(orderRequest.Qty, 'f', -1, 64),
		Side:        strings.ToUpper(orderRequest.Side),
		OrderType:   orderType,
		TimeInForce: orderRequest.TimeInForce,
		OrderPrice:  orderPriceStr,
		OrderLinkId: orderRequest.OrderId,
	}

	var order request.SpotOrderResponseTimeStr
	order, err = request.CreateSpotOrder(domainOrderRequest, d.secrets)
	if err != nil {
		return
	}
	orderId = order.OrderId
	return
}

func (d *DomainBybitSpot) CancelOrder(orderId string, coinPare string) error {
	err := request.CancelSpotOrder(orderId, d.secrets, coinPare)
	if err != nil {
		return err
	}
	return nil
}

func (d *DomainBybitSpot) GetHistoryOrders(limit int64, coinPare string) ([]structs.DomainOrder, error) {
	orders, err := request.GetSpotOrderHistory(limit, d.secrets, coinPare)
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
		createTime, err := strconv.ParseInt(order.CreateTime, 10, 64)
		if err != nil {
			return nil, err
		}
		updateTime, err := strconv.ParseInt(order.UpdateTime, 10, 64)
		if err != nil {
			return nil, err
		}

		domainOrder := structs.DomainOrder{
			CreatedTime: createTime / 1000,
			OrderId:     order.OrderId,
			OrderStatus: order.Status,
			OrderType:   order.OrderType,
			Price:       price,
			Qty:         orderQty,
			ReduceOnly:  false,
			Side:        order.Side,
			Symbol:      order.Symbol,
			TimeInForce: order.TimeInForce,
			UpdatedTime: updateTime / 1000,
		}
		domainOrders[key] = domainOrder
	}

	sort.Slice(domainOrders, func(i, j int) bool {
		return domainOrders[i].CreatedTime > domainOrders[j].CreatedTime
	})

	return domainOrders, nil
}
