package fibonacci

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/fibonacci/storage"
	"github.com/shatylos/trader/internal/strategy/fibonacci/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"time"
)

const (
	TrendLong    = "BULLISH"
	TrendShort   = "BEARISH"
	TrendUnknown = "UNKNOWN"
)

func (f *Fibonacci) calculateNewPosition(prevInternalPosition structs.Position, currentPrice float64) (position structs.Position, err error) {
	limit := f.config.MinCandleReview
	if prevInternalPosition.CreatedTime != 0 {
		mins := (time.Now().Unix() - prevInternalPosition.CreatedTime) / 60
		fromPrevLimit := mins / f.config.ResolutionMins
		if limit < fromPrevLimit {
			limit = fromPrevLimit
		}
	}
	if limit > f.config.MaxCandleReview {
		limit = f.config.MaxCandleReview
	}

	var candles []domainStructs.DomainCandle
	candles, err = f.provider.LoadCandleHistory(f.config.CoinPare, f.config.Resolution, time.Now().Unix(), limit)
	if err != nil {
		return
	}

	trend := f.getTrend(candles)

	var fibChart structs.FibonacciChart
	var minPrice, maxPrice float64
	switch trend {
	case TrendLong:
		minPrice, maxPrice, err = f.getMinMaxPriceLong(candles)
		fibChart = f.calculateFibonacciChart(minPrice, maxPrice, true)
		break
	case TrendShort:
		minPrice, maxPrice, err = f.getMinMaxPriceShort(candles)
		fibChart = f.calculateFibonacciChart(minPrice, maxPrice, false)
		break
	default:
		err = tools.AppError{Message: fmt.Sprintf("Unexpected trend value: %s", trend)}
		return
	}
	fibChart.FullQty, err = f.calculateFullQty(currentPrice)
	if err != nil {
		return
	}

	position = structs.Position{
		Id:             nil,
		FibonacciChart: fibChart,
		Trend:          trend,
		Orders:         structs.PositionOrders{},
		Status:         structs.StatusNew,
	}

	return
}

func (f *Fibonacci) getTrend(candles []domainStructs.DomainCandle) string {
	var minCandle, maxCandle domainStructs.DomainCandle
	if len(candles) == 0 {
		return TrendUnknown
	}
	minCandle = candles[0]
	maxCandle = candles[0]
	for _, candle := range candles {
		if candle.Low < minCandle.Low {
			minCandle = candle
		}
		if candle.High > maxCandle.High {
			maxCandle = candle
		}
	}

	if maxCandle.Time > minCandle.Time {
		return TrendLong
	} else {
		return TrendShort
	}
}

func (f *Fibonacci) getMinMaxPriceLong(candles []domainStructs.DomainCandle) (minPrice float64, maxPrice float64, err error) {
	if len(candles) == 0 {
		err = tools.AppError{Message: "candles array are empty"}
		return
	}
	minPrice = candles[0].Low
	maxPrice = candles[0].High
	for _, candle := range candles {
		if candle.Low < minPrice {
			minPrice = candle.Low
		}
	}
	for _, candle := range candles {
		if candle.High > maxPrice {
			maxPrice = candle.High
		}
		if candle.Low == minPrice {
			break
		}
	}
	return
}

func (f *Fibonacci) getMinMaxPriceShort(candles []domainStructs.DomainCandle) (minPrice float64, maxPrice float64, err error) {
	if len(candles) == 0 {
		err = tools.AppError{Message: "candles array are empty"}
		return
	}
	minPrice = candles[0].Low
	maxPrice = candles[0].High
	for _, candle := range candles {
		if candle.High > maxPrice {
			maxPrice = candle.High
		}
	}
	for _, candle := range candles {
		if candle.Low < minPrice {
			minPrice = candle.Low
		}
		if candle.High == maxPrice {
			break
		}
	}
	return
}

func (f *Fibonacci) calculateFibonacciChart(minPrice float64, maxPrice float64, isLong bool) (fibonacciChart structs.FibonacciChart) {
	fibonacciChart.SourceMinPrice = minPrice
	fibonacciChart.SourceMaxPrice = maxPrice

	fibonacciChart.EntryPoint1 = f.valueByCoeff(minPrice, maxPrice, f.config.FibEntryPoint1, isLong)
	fibonacciChart.EntryPoint2 = f.valueByCoeff(minPrice, maxPrice, f.config.FibEntryPoint2, isLong)
	fibonacciChart.EntryPoint3 = f.valueByCoeff(minPrice, maxPrice, f.config.FibEntryPoint3, isLong)
	fibonacciChart.StopLoss = f.valueByCoeff(minPrice, maxPrice, f.config.FibStopLoss, isLong)
	fibonacciChart.TakeProfit1 = f.valueByCoeff(minPrice, maxPrice, f.config.FibTakeProfit1, isLong)
	fibonacciChart.TakeProfit2 = f.valueByCoeff(minPrice, maxPrice, f.config.FibTakeProfit2, isLong)
	fibonacciChart.TakeProfit3 = f.valueByCoeff(minPrice, maxPrice, f.config.FibTakeProfit3, isLong)

	return
}

func (f *Fibonacci) valueByCoeff(minPrice float64, maxPrice float64, coeff float64, isLong bool) (result float64) {
	diff := math.Mul(maxPrice-minPrice, coeff)
	if isLong {
		result = minPrice + diff
	} else {
		result = maxPrice - diff
	}
	result = math.Round(result, f.config.PricePrecision)
	return
}

func (f *Fibonacci) calculateFullQty(price float64) (fullQty float64, err error) {
	depoMainCurrency := float64(0)
	var wallet *domainStructs.DomainWallet
	wallet, err = f.provider.GetWallet()
	if err != nil {
		return
	}

	for _, coin := range wallet.Available {
		if coin.Coin == f.config.MainCurrency {
			depoMainCurrency = coin.Amount
		}
	}

	if depoMainCurrency < 0 {
		logger.Warning("Not enough deposit")
		err = tools.AppError{
			Message: "Not enough deposit",
		}
		return
	}

	fullQty = math.Mul(math.Div(math.Div(depoMainCurrency, price), 100), f.config.FullQtyToDepoPercent)
	fullQty = math.Round(fullQty, f.config.QtyPrecision)
	return
}

func (f *Fibonacci) openNewPosition(internalPosition structs.Position, qty float64, orderNum int, side string) (isCreated bool, err error) {
	if side != "Sell" && side != "Buy" {
		err = tools.AppError{
			Message: fmt.Sprintf("Unexpected side value: %s", side),
		}
		return
	}
	if orderNum < 1 || orderNum > 3 {
		err = tools.AppError{
			Message: fmt.Sprintf("Unexpected orderNum value: %d", orderNum),
		}
		return
	}

	var takeProfit, stopLoss float64
	stopLoss = internalPosition.FibonacciChart.StopLoss
	switch orderNum {
	case 1:
		takeProfit = internalPosition.FibonacciChart.TakeProfit1
		break
	case 2:
		takeProfit = internalPosition.FibonacciChart.TakeProfit2
		break
	case 3:
		takeProfit = internalPosition.FibonacciChart.TakeProfit3
		break
	}

	var orderId string
	orderId, err = f.provider.OpenPosition(domainStructs.DomainPositionRequest{
		Leverage:    f.config.Leverage,
		PositionId:  "",
		Price:       0,
		Qty:         qty,
		ReduceOnly:  false,
		Side:        side,
		StopLoss:    stopLoss,
		Symbol:      f.config.CoinPare,
		TakeProfit:  takeProfit,
		TimeInForce: "",
		Type:        "Market",
	})
	var newOrder domainStructs.DomainOrder
	newOrder, err = f.provider.GetOrder(orderId)
	if err != nil {
		return
	}

	switch orderNum {
	case 1:
		internalPosition.Orders.Order1 = newOrder
		break
	case 2:
		internalPosition.Orders.Order2 = newOrder
		break
	case 3:
		internalPosition.Orders.Order3 = newOrder
		break
	}
	internalPosition.Status = structs.StatusActive

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	internalPosition, err = storage.SaveInternalPosition(internalPosition)
	if err != nil {
		return
	}
	if internalPosition.Id == nil {
		err = tools.AppError{Message: "Internal position was not saved"}
		return
	}
	logger.Success(fmt.Sprintf("Created new order. PositionId %s. Order number: %d. Qty: %f. Price: %f. Side: %s. TakeProfit: %f. StopLoss %f. SourceMinPrice: %f. SourceMaxPrice: %f.",
		*internalPosition.Id,
		orderNum,
		newOrder.Qty,
		newOrder.Price,
		newOrder.Side,
		newOrder.TakeProfit,
		newOrder.StopLoss,
		internalPosition.FibonacciChart.SourceMinPrice,
		internalPosition.FibonacciChart.SourceMaxPrice,
	))
	return
}

func (f *Fibonacci) closeInternalPosition(internalPosition structs.Position) (err error) {
	var order domainStructs.DomainOrder
	if internalPosition.Orders.Order1.OrderStatus == domainStructs.OrderStatusOpen {
		order, err = f.provider.GetOrder(internalPosition.Orders.Order1.OrderId)
		if err != nil {
			return
		}
		if order.OrderStatus != domainStructs.OrderStatusCancelled && order.OrderStatus != domainStructs.OrderStatusFilled {
			logger.Warning(fmt.Sprintf(
				"Expected order must be cancelled or filled when close the internal position. OrderId: %s, OrderStatus: %s",
				order.OrderId,
				order.OrderStatus,
			))
		}
		internalPosition.Orders.Order1 = order
	}
	if internalPosition.Orders.Order2.OrderStatus == domainStructs.OrderStatusOpen {
		order, err = f.provider.GetOrder(internalPosition.Orders.Order2.OrderId)
		if err != nil {
			return
		}
		if order.OrderStatus != domainStructs.OrderStatusCancelled && order.OrderStatus != domainStructs.OrderStatusFilled {
			logger.Warning(fmt.Sprintf(
				"Expected order must be cancelled or filled when close the internal position. OrderId: %s, OrderStatus: %s",
				order.OrderId,
				order.OrderStatus,
			))
		}
		internalPosition.Orders.Order2 = order
	}
	if internalPosition.Orders.Order3.OrderStatus == domainStructs.OrderStatusOpen {
		order, err = f.provider.GetOrder(internalPosition.Orders.Order3.OrderId)
		if err != nil {
			return
		}
		if order.OrderStatus != domainStructs.OrderStatusCancelled && order.OrderStatus != domainStructs.OrderStatusFilled {
			logger.Warning(fmt.Sprintf(
				"Expected order must be cancelled or filled when close the internal position. OrderId: %s, OrderStatus: %s",
				order.OrderId,
				order.OrderStatus,
			))
		}
		internalPosition.Orders.Order3 = order
	}

	internalPosition.Status = structs.StatusClosed

	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	internalPosition, err = storage.SaveInternalPosition(internalPosition)
	if err != nil {
		return
	}
	if internalPosition.Id == nil {
		err = tools.AppError{Message: "Internal position was not closed"}
		return
	}
	logger.Info(fmt.Sprintf("Position was closed. PositionId %s.", *internalPosition.Id))
	return
}
