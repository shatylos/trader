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
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
	math2 "math"
	"time"
)

func (f *Fibonacci) calculateNewPosition() (position structs.Position, err error) {
	var calculateFromPosition structs.Position
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(f.config.Id)
	if err != nil {
		return
	}
	calculateFromPosition, err = storage.GetLastInternalPosition(f.config.PrevPositionsReview - 1)
	if err != nil {
		return
	}

	limit := f.config.MinCandleReview
	if calculateFromPosition.CreatedTime != 0 {
		mins := (time.Now().Unix() - calculateFromPosition.CreatedTime) / 60
		fromPrevLimit := mins/f.config.ResolutionMins + 2
		if limit < fromPrevLimit {
			limit = fromPrevLimit
		}
	}
	if limit > f.config.MaxCandleReview {
		limit = f.config.MaxCandleReview
	}

	var candles, ltCandles []domainStructs.DomainCandle
	candles, err = f.provider.LoadCandleHistory(f.config.CoinPare, f.config.Resolution, limit)
	if err != nil {
		return
	}
	ltCandles, err = f.provider.LoadCandleHistory(f.config.CoinPare, f.config.LongTrendResolution, f.config.LongTrendCandleReview)
	if err != nil {
		return
	}

	//shortTrend := trading.GetFullTrend(candles, f.config.Verbose)
	shortTrend := trading.GetTrendLinearRegression(candles)
	if f.config.Verbose {
		logger.Info(fmt.Sprintf("Detected short trend: %s", shortTrend))
	}
	longTrend := trading.GetTrendLinearRegression(ltCandles)
	if f.config.Verbose {
		logger.Info(fmt.Sprintf("Detected long trend: %s", longTrend))
	}
	trend := trading.TrendUnknown
	if shortTrend == longTrend {
		trend = shortTrend
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Long trend (%s) and short (%s) trend are different", longTrend, shortTrend))
	}

	var fibChart structs.FibonacciChart
	var minPrice, maxPrice float64
	switch trend {
	case trading.TrendLong:
		minPrice, maxPrice, err = f.getMinMaxPriceLong(candles)
		fibChart = f.calculateFibonacciChart(minPrice, maxPrice, true)
		break
	case trading.TrendShort:
		minPrice, maxPrice, err = f.getMinMaxPriceShort(candles)
		fibChart = f.calculateFibonacciChart(minPrice, maxPrice, false)
		break
	case trading.TrendUnknown:
		fibChart = f.calculateFibonacciChart(0, 0, true)
		break
	default:
		err = tools.AppError{Message: fmt.Sprintf("Unexpected trend value: %s", trend)}
		return
	}

	availableBalance := float64(0)
	availableBalance, err = f.getAvailableBalance()
	if err != nil {
		return
	}
	if f.config.Verbose {
		logger.Info(fmt.Sprintf("Available balance: %g", availableBalance))
	}

	fibChart.FullQty, err = f.calculateFullQty(availableBalance, fibChart)
	if err != nil {
		return
	}

	position = structs.Position{
		Id:             nil,
		FibonacciChart: fibChart,
		Trend:          trend,
		Orders:         structs.PositionOrders{},
		Status:         structs.StatusNew,
		BalanceBefore:  availableBalance,
	}

	return
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

func (f *Fibonacci) getAvailableBalance() (balance float64, err error) {
	var wallet domainStructs.DomainWallet
	wallet, err = f.provider.GetWallet()
	if err != nil {
		return
	}
	for _, coin := range wallet.Available {
		if coin.Coin == f.config.MainCurrency {
			balance = coin.Amount
		}
	}
	return
}

func (f *Fibonacci) calculateFullQty(availableBalance float64, fibChart structs.FibonacciChart) (fullQty float64, err error) {
	if availableBalance <= 0 {
		logger.Warning("Not enough deposit")
		err = tools.AppError{
			Message: "Not enough deposit",
		}
		return
	}

	tempFullQty := float64(1)                                                      // 1
	tempQty1 := math.Mul(math.Div(tempFullQty, 100), f.config.EP1ToFullQtyPercent) // 1 / 100 * 20 = 0.2
	tempQty2 := math.Mul(math.Div(tempFullQty, 100), f.config.EP2ToFullQtyPercent) // 1 / 100 * 30 = 0.3
	tempQty3 := math.Mul(math.Div(tempFullQty, 100), f.config.EP3ToFullQtyPercent) // 1 / 100 * 50 = 0.5

	var tempLoss1, tempLoss2, tempLoss3 float64
	if fibChart.EntryPoint1 > 0 {
		tempLoss1 = f.calculateLoss(tempQty1, fibChart.EntryPoint1, fibChart.StopLoss) // 2000
	}
	if fibChart.EntryPoint2 > 0 {
		tempLoss2 = f.calculateLoss(tempQty2, fibChart.EntryPoint2, fibChart.StopLoss) // 1800
	}
	if fibChart.EntryPoint3 > 0 {
		tempLoss3 = f.calculateLoss(tempQty3, fibChart.EntryPoint3, fibChart.StopLoss) // 1500
	}

	tempFullLoss := tempLoss1 + tempLoss2 + tempLoss3                           // 2000 + 1800 + 1500 = 5300
	fullRisk := math.Mul(math.Div(availableBalance, 100), f.config.RiskPercent) // 500 / 100 * 10 = 50

	if tempFullLoss > 0 && fullRisk > 0 {
		reduceCoef := math.Div(tempFullLoss, fullRisk) // 5300 / 50 = 106
		fullQty = math.Div(tempFullQty, reduceCoef)    // 1 / 106 = 0,009433962
		fullQty = math.Round(fullQty, f.config.QtyPrecision)
	}
	return
}

func (f *Fibonacci) calculateLoss(qty float64, entryPoint float64, stopLoss float64) (loss float64) {
	amountIn := entryPoint * qty           // 80000 * 0.2 = 16000
	amountOut := stopLoss * qty            // 70000 * 0.2 = 14000
	loss = math2.Abs(amountIn - amountOut) // 16000 - 14000 = 2000
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
	if err != nil {
		err = tools.AppError{
			Message: fmt.Sprintf("Error opening new position. Leverage: %d. Qty: %g. Side: %s, TakeProfit: %g. StopLoss: %g",
				f.config.Leverage,
				qty,
				side,
				takeProfit,
				stopLoss,
			),
			ParentError: err,
		}
		return
	}
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
	logger.Success(fmt.Sprintf("Created new order. PositionId %s. Order number: %d. Qty: %g. Price: %g. Side: %s. TakeProfit: %g. StopLoss %g. SourceMinPrice: %g. SourceMaxPrice: %g.",
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
	if f.config.TelegramNotifier {
		tgNotifier.Notify(fmt.Sprintf("Created new order. Qty: %g. Price: %g. Side: %s.",
			newOrder.Qty,
			newOrder.Price,
			newOrder.Side,
		))
	}
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
	internalPosition.ClosedTime = time.Now().Unix()

	availableBalance := float64(0)
	availableBalance, err = f.getAvailableBalance()
	if err != nil {
		return
	}
	internalPosition.BalanceAfter = availableBalance
	internalPosition.TotalClosePnl = internalPosition.BalanceAfter - internalPosition.BalanceBefore

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
	logger.Info(fmt.Sprintf("Position was closed. Position ID %s.", *internalPosition.Id))
	if f.config.TelegramNotifier {
		tgNotifier.Notify(fmt.Sprintf("Position was closed. PNL: %g.", math.Round(internalPosition.TotalClosePnl, 2)))
	}
	return
}
