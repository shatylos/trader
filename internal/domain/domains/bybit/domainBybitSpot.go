package bybit

import (
	"github.com/shatylos/trader/internal/domain/domains"
	"github.com/shatylos/trader/internal/domain/domains/bybit/mapping"
	"github.com/shatylos/trader/internal/domain/domains/bybit/request"
	bybitStructs "github.com/shatylos/trader/internal/domain/domains/bybit/structs"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/type"
	"sort"
	"strconv"
	"time"
)

type DomainBybitSpot struct {
	code    string
	secrets bybitStructs.Secrets
}

func (d *DomainBybitSpot) GetCode() string {
	return d.code
}

func (d *DomainBybitSpot) SetConfig(config map[interface{}]interface{}) (err error) {

	secretMap, ok := config["secrets"].(map[interface{}]interface{})
	if !ok {
		err = apperrors.New("the field secrets is empty or contains not correct value type. In DomainBybitSpot config. Expects a map with \"key\", \"pass\" and \"endpoint\" keys")
		return
	}

	var domainCode string
	domainCode, err = _type.ToString(config["code"])
	if err != nil {
		err = apperrors.New("the field code is empty or contains not correct value type. In DomainBybitSpot config. Expects a string")
		return
	}
	d.code = domainCode

	secrets := bybitStructs.Secrets{}

	var apiEndpoint string
	apiEndpoint, err = _type.ToString(secretMap["endpoint"])
	if err != nil {
		err = apperrors.Wrap(err, "the field secrets.endpoint is empty or contains not correct value type. In DomainBybitSpot config. Expects a string")
		return
	}
	secrets.ApiEndpoint = apiEndpoint

	var key string
	key, err = _type.ToString(secretMap["key"])
	if err != nil {
		err = apperrors.Wrap(err, "the field secrets.key is empty or contains not correct value type. In DomainBybitSpot config. Expects a string")
		return
	}
	secrets.Key = key

	var pass string
	pass, err = _type.ToString(secretMap["pass"])
	if err != nil {
		err = apperrors.Wrap(err, "the field secrets.pass is empty or contains not correct value type. In DomainBybitSpot config. Expects a string")
		return
	}
	secrets.Pass = pass

	var verbose int64
	verbose, err = _type.ToInt64(secretMap["verbose"])
	if err != nil {
		err = apperrors.Wrap(err, "the field verbose in domain config is empty or contains not correct value type. Expects 1 or 0")
		return err
	}
	if verbose == 1 {
		secrets.Verbose = true
	} else {
		secrets.Verbose = false
	}
	d.secrets = secrets

	return nil
}

func (d *DomainBybitSpot) GetWallet() (wallet structs.DomainWallet, err error) {
	var walletBalances *map[string]request.SpotWalletBalance
	throttle()
	walletBalances, err = request.GetSpotWalletBalance(d.secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error get spot wallet balance")
		return
	}

	var availableCoins []structs.DomainWalletCoinItem
	var reservedCoins []structs.DomainWalletCoinItem

	for coinCode, walletBalance := range *walletBalances {
		var free float64
		free, err = strconv.ParseFloat(walletBalance.Free, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float walletBalance.Free: %s", walletBalance.Free)
			return
		}
		var locked float64
		locked, err = strconv.ParseFloat(walletBalance.Locked, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float walletBalance.Locked: %s", walletBalance.Locked)
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
		Available:   availableCoins,
		Reserved:    reservedCoins,
		UpdatedTime: time.Now(),
	}
	return
}

func (d *DomainBybitSpot) LoadCandleHistory(symbol string, resolution string, limit int64) ([]structs.DomainCandle, error) {
	providerResolution, err := mapping.ToBybitInterval(resolution)
	if err != nil {
		err = apperrors.Wrap(err, "error mapping resolution \"%s\" to bybit interval", resolution)
		return nil, err
	}
	throttle()
	candles, err := request.GetSpotKlineList(symbol, providerResolution, limit, d.secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error get spot kline list. Symbol: %s, providerResolution: %s, limit: %d", symbol, providerResolution, limit)
		return nil, err
	}

	candlesResult := make([]structs.DomainCandle, len(candles))

	for i, candle := range candles {
		highPrice, err := strconv.ParseFloat(candle.High, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, candle.High: %s", candle.High)
			return nil, err
		}
		lowPrice, err := strconv.ParseFloat(candle.Low, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, candle.Low: %s", candle.Low)
			return nil, err
		}
		openPrice, err := strconv.ParseFloat(candle.Open, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, candle.Open: %s", candle.Open)
			return nil, err
		}
		closePrice, err := strconv.ParseFloat(candle.Close, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, candle.Close: %s", candle.Close)
			return nil, err
		}
		volume, err := strconv.ParseFloat(candle.Volume, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, candle.Volume: %s", candle.Volume)
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

func (d *DomainBybitSpot) GetOrder(domainId string) (domainOrder structs.DomainOrder, err error) {
	var order request.SpotOrderResponseTimeStr
	throttle()
	order, err = request.GetSpotOrder(domainId, d.secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error get spot order. DomainId: %s", domainId)
		return
	}

	if order.Symbol == "" {
		err = apperrors.Wrap(domains.OrderNotFoundError, "Order with ID (%s) not found", domainId)
		return
	}

	filledPrice, err := strconv.ParseFloat(order.AvgPrice, 64)
	if err != nil {
		err = apperrors.Wrap(err, "error parse float, order.AvgPrice: %s", order.AvgPrice)
		return
	}
	filledQty, err := strconv.ParseFloat(order.ExecQty, 64)
	if err != nil {
		err = apperrors.Wrap(err, "error parse float, order.ExecQty: %s", order.ExecQty)
		return
	}
	createTime, err := strconv.ParseInt(order.CreateTime, 10, 64)
	if err != nil {
		err = apperrors.Wrap(err, "error parse float, order.CreateTime: %s", order.CreateTime)
		return
	}
	updateTime, err := strconv.ParseInt(order.UpdateTime, 10, 64)
	if err != nil {
		err = apperrors.Wrap(err, "error parse float, order.UpdateTime: %s", order.UpdateTime)
		return
	}

	var side string
	side, err = mapping.ToDomainOrderSide(order.Side)
	if err != nil {
		err = apperrors.Wrap(err, "error mapping to domain order side: %s", order.Side)
		return
	}

	domainOrder = structs.DomainOrder{
		CreatedTime: createTime / 1000,
		OrderId:     order.OrderId,
		OrderStatus: order.Status,
		OrderType:   order.OrderType,
		Price:       filledPrice,
		Qty:         filledQty,
		ReduceOnly:  false,
		Side:        side,
		Symbol:      order.Symbol,
		TimeInForce: order.TimeInForce,
		UpdatedTime: updateTime / 1000,
	}

	return
}

func (d *DomainBybitSpot) GetOpenOrderList(coinPare string) ([]structs.DomainOrder, error) {
	throttle()
	orders, err := request.GetSpotOpenOrderList(coinPare, d.secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error get spot open order list. coinPare: %s", coinPare)
		return nil, err
	}

	domainOrders := make([]structs.DomainOrder, len(orders))

	for key, order := range orders {
		orderPrice, err := strconv.ParseFloat(order.OrderPrice, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, order.OrderPrice: %s", order.OrderPrice)
			return nil, err
		}
		orderQty, err := strconv.ParseFloat(order.OrderQty, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, order.OrderQty: %s", order.OrderQty)
			return nil, err
		}
		createTime, err := _type.ToInt64(order.CreateTime)
		if err != nil {
			err = apperrors.Wrap(err, "error parse int64, order.CreateTime: %s", order.CreateTime)
			return nil, err
		}
		updateTime, err := _type.ToInt64(order.UpdateTime)
		if err != nil {
			err = apperrors.Wrap(err, "error parse int64, order.UpdateTime: %s", order.UpdateTime)
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
		err = apperrors.Wrap(err, "error mapping to bybit order type, orderRequest.Type: %s", orderRequest.Type)
		return
	}

	orderPriceStr := ""
	if orderRequest.Price > 0 {
		orderPriceStr = strconv.FormatFloat(orderRequest.Price, 'f', -1, 64)
	}

	var side string
	side, err = mapping.ToBybitOrderSide(orderRequest.Side)
	if err != nil {
		err = apperrors.Wrap(err, "error mapping to bybit order side, orderRequest.Side: %s", orderRequest.Side)
		return
	}

	domainOrderRequest := request.SpotOrderRequest{
		Symbol:      orderRequest.Symbol,
		OrderQty:    strconv.FormatFloat(orderRequest.Qty, 'f', -1, 64),
		Side:        side,
		OrderType:   orderType,
		TimeInForce: orderRequest.TimeInForce,
		OrderPrice:  orderPriceStr,
		OrderLinkId: orderRequest.OrderId,
	}

	var order request.CreateSpotOrderResponse
	throttle()
	order, err = request.CreateSpotOrder(domainOrderRequest, d.secrets)
	if err != nil {
		err = apperrors.Wrap(err, "error create spot order, domainOrderRequest: %s", domainOrderRequest)
		return
	}
	orderId = order.OrderId
	return
}

func (d *DomainBybitSpot) CancelOrder(orderId string, coinPare string) error {
	throttle()
	err := request.CancelSpotOrder(orderId, d.secrets, coinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error cancel spot order, orderId: %s", orderId)
		return err
	}
	return nil
}

func (d *DomainBybitSpot) GetHistoryOrders(limit int64, coinPare string) ([]structs.DomainOrder, error) {
	throttle()
	orders, err := request.GetSpotOrderHistory(limit, d.secrets, coinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error get spot order history, limit: %s, coinPare: %s", limit, coinPare)
		return nil, err
	}

	domainOrders := make([]structs.DomainOrder, len(orders))

	for key, order := range orders {
		orderPrice, err := strconv.ParseFloat(order.OrderPrice, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, order.OrderPrice: %s", order.OrderPrice)
			return nil, err
		}
		avgPrice, err := strconv.ParseFloat(order.AvgPrice, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, order.AvgPrice: %s", order.AvgPrice)
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
			err = apperrors.Wrap(err, "error parse float, order.OrderQty: %s", order.OrderQty)
			return nil, err
		}
		createTime, err := strconv.ParseInt(order.CreateTime, 10, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse float, order.CreateTime: %s", order.CreateTime)
			return nil, err
		}
		updateTime, err := strconv.ParseInt(order.UpdateTime, 10, 64)
		if err != nil {
			err = apperrors.Wrap(err, "error parse int64, order.UpdateTime: %s", order.UpdateTime)
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
