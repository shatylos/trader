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
	"time"
)

func (f *Fibonacci) actionByPosition(internalPosition structs.Position, currentPrice float64) (err error) {
	switch internalPosition.Trend {
	case trading.TrendLong:
		err = f.actionByPositionBullish(internalPosition, currentPrice)
		break
	case trading.TrendShort:
		err = f.actionByPositionBearish(internalPosition, currentPrice)
		break
	case trading.TrendUnknown:
		if f.config.Verbose {
			logger.Info(fmt.Sprintf("Trend is %s. Wait for %s or %s.", internalPosition.Trend, trading.TrendLong, trading.TrendShort))
		}
		break
	default:
		err = tools.AppError{
			Message: fmt.Sprintf("Unexpected trend \"%s\" for action by position", internalPosition.Trend),
		}
	}
	if err != nil {
		return
	}
	err = f.takeProfitCorrection(internalPosition)
	return
}

func (f *Fibonacci) actionByPositionBullish(internalPosition structs.Position, currentPrice float64) (err error) {
	if internalPosition.Status == structs.StatusNew && internalPosition.BalanceBefore <= 0 {
		msg := fmt.Sprintf("Not enough balance (%g) to action", internalPosition.BalanceBefore)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	if internalPosition.Orders.Order1.OrderId == "" &&
		currentPrice < internalPosition.FibonacciChart.EntryPoint1 &&
		currentPrice > internalPosition.FibonacciChart.EntryPoint2 {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP1ToFullQtyPercent), f.config.QtyPrecision)
		if epQty >= f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 1, "Buy")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 1 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		} else {
			logger.Warning(fmt.Sprintf("Order 1 was not created. QTY (%f). Min qty (%f)", epQty, f.config.MinQty))
		}
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Bullish Order 1 was not opened: OrderId: %s. currentPrice: %f. EntryPoint1: %f. StopLoss: %f. SourceMinPrice: %f. SourceMaxPrice: %f. Cond1: %t. Cond2: %t. Cond3: %t.",
			internalPosition.Orders.Order1.OrderId,
			currentPrice,
			internalPosition.FibonacciChart.EntryPoint1,
			internalPosition.FibonacciChart.StopLoss,
			internalPosition.FibonacciChart.SourceMinPrice,
			internalPosition.FibonacciChart.SourceMaxPrice,
			internalPosition.Orders.Order1.OrderId == "",
			currentPrice < internalPosition.FibonacciChart.EntryPoint1,
			currentPrice > internalPosition.FibonacciChart.StopLoss,
		))
	}

	if internalPosition.Orders.Order2.OrderId == "" &&
		internalPosition.Orders.Order1.OrderId != "" &&
		currentPrice < internalPosition.FibonacciChart.EntryPoint2 &&
		currentPrice > internalPosition.FibonacciChart.EntryPoint3 {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP2ToFullQtyPercent), f.config.QtyPrecision)
		if epQty >= f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 2, "Buy")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 2 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		} else {
			logger.Warning(fmt.Sprintf("Order 2 was not created. QTY (%f). Min qty (%f)", epQty, f.config.MinQty))
		}
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Bullish Order 2 was not opened: OrderId: %s. currentPrice: %f. EntryPoint2: %f. StopLoss: %f. SourceMinPrice: %f. SourceMaxPrice: %f. Cond1: %t. Cond2: %t. Cond3: %t.",
			internalPosition.Orders.Order2.OrderId,
			currentPrice,
			internalPosition.FibonacciChart.EntryPoint2,
			internalPosition.FibonacciChart.StopLoss,
			internalPosition.FibonacciChart.SourceMinPrice,
			internalPosition.FibonacciChart.SourceMaxPrice,
			internalPosition.Orders.Order2.OrderId == "",
			currentPrice < internalPosition.FibonacciChart.EntryPoint2,
			currentPrice > internalPosition.FibonacciChart.StopLoss,
		))
	}

	if internalPosition.Orders.Order3.OrderId == "" &&
		internalPosition.Orders.Order2.OrderId != "" &&
		currentPrice < internalPosition.FibonacciChart.EntryPoint3 &&
		currentPrice > internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP3ToFullQtyPercent), f.config.QtyPrecision)
		if epQty >= f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 3, "Buy")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 3 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		} else {
			logger.Warning(fmt.Sprintf("Order 3 was not created. QTY (%f). Min qty (%f)", epQty, f.config.MinQty))
		}
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Bullish Order 3 was not opened: OrderId: %s. currentPrice: %f. EntryPoint3: %f. StopLoss: %f. SourceMinPrice: %f. SourceMaxPrice: %f. Cond1: %t. Cond2: %t. Cond3: %t.",
			internalPosition.Orders.Order3.OrderId,
			currentPrice,
			internalPosition.FibonacciChart.EntryPoint3,
			internalPosition.FibonacciChart.StopLoss,
			internalPosition.FibonacciChart.SourceMinPrice,
			internalPosition.FibonacciChart.SourceMaxPrice,
			internalPosition.Orders.Order3.OrderId == "",
			currentPrice < internalPosition.FibonacciChart.EntryPoint3,
			currentPrice > internalPosition.FibonacciChart.StopLoss,
		))
	}
	return
}

func (f *Fibonacci) actionByPositionBearish(internalPosition structs.Position, currentPrice float64) (err error) {
	if internalPosition.Status == structs.StatusNew && internalPosition.BalanceBefore <= 0 {
		msg := fmt.Sprintf("Not enough balance (%g) to action", internalPosition.BalanceBefore)
		logger.Error(msg)
		err = tools.AppError{
			Message: msg,
		}
		return
	}

	if internalPosition.Orders.Order1.OrderId == "" &&
		currentPrice > internalPosition.FibonacciChart.EntryPoint1 &&
		currentPrice < internalPosition.FibonacciChart.EntryPoint2 {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP1ToFullQtyPercent), f.config.QtyPrecision)
		if epQty >= f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 1, "Sell")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 1 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		} else {
			logger.Warning(fmt.Sprintf("Order 1 was not created. QTY (%f). Min qty (%f)", epQty, f.config.MinQty))
		}
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Bearish Order 1 was not opened: OrderId: %s. currentPrice: %f. EntryPoint1: %f. StopLoss: %f. SourceMinPrice: %f. SourceMaxPrice: %f. Cond1: %t. Cond2: %t. Cond3: %t.",
			internalPosition.Orders.Order1.OrderId,
			currentPrice,
			internalPosition.FibonacciChart.EntryPoint1,
			internalPosition.FibonacciChart.StopLoss,
			internalPosition.FibonacciChart.SourceMinPrice,
			internalPosition.FibonacciChart.SourceMaxPrice,
			internalPosition.Orders.Order1.OrderId == "",
			currentPrice > internalPosition.FibonacciChart.EntryPoint1,
			currentPrice < internalPosition.FibonacciChart.StopLoss,
		))
	}

	if internalPosition.Orders.Order2.OrderId == "" &&
		internalPosition.Orders.Order1.OrderId != "" &&
		currentPrice > internalPosition.FibonacciChart.EntryPoint2 &&
		currentPrice < internalPosition.FibonacciChart.EntryPoint3 {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP2ToFullQtyPercent), f.config.QtyPrecision)
		if epQty >= f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 2, "Sell")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 2 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		} else {
			logger.Warning(fmt.Sprintf("Order 2 was not created. QTY (%f). Min qty (%f)", epQty, f.config.MinQty))
		}
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Bearish Order 2 was not opened: OrderId: %s. currentPrice: %f. EntryPoint2: %f. StopLoss: %f. SourceMinPrice: %f. SourceMaxPrice: %f. Cond1: %t. Cond2: %t. Cond3: %t.",
			internalPosition.Orders.Order2.OrderId,
			currentPrice,
			internalPosition.FibonacciChart.EntryPoint2,
			internalPosition.FibonacciChart.StopLoss,
			internalPosition.FibonacciChart.SourceMinPrice,
			internalPosition.FibonacciChart.SourceMaxPrice,
			internalPosition.Orders.Order2.OrderId == "",
			currentPrice > internalPosition.FibonacciChart.EntryPoint2,
			currentPrice < internalPosition.FibonacciChart.StopLoss,
		))
	}

	if internalPosition.Orders.Order3.OrderId == "" &&
		internalPosition.Orders.Order2.OrderId != "" &&
		currentPrice > internalPosition.FibonacciChart.EntryPoint3 &&
		currentPrice < internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP3ToFullQtyPercent), f.config.QtyPrecision)
		if epQty >= f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 3, "Sell")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 3 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		} else {
			logger.Warning(fmt.Sprintf("Order 3 was not created. QTY (%f). Min qty (%f)", epQty, f.config.MinQty))
		}
	} else if f.config.Verbose {
		logger.Info(fmt.Sprintf("Bearish Order 3 was not opened: OrderId: %s. currentPrice: %f. EntryPoint3: %f. StopLoss: %f. SourceMinPrice: %f. SourceMaxPrice: %f. Cond1: %t. Cond2: %t. Cond3: %t.",
			internalPosition.Orders.Order3.OrderId,
			currentPrice,
			internalPosition.FibonacciChart.EntryPoint3,
			internalPosition.FibonacciChart.StopLoss,
			internalPosition.FibonacciChart.SourceMinPrice,
			internalPosition.FibonacciChart.SourceMaxPrice,
			internalPosition.Orders.Order3.OrderId == "",
			currentPrice > internalPosition.FibonacciChart.EntryPoint3,
			currentPrice < internalPosition.FibonacciChart.StopLoss,
		))
	}
	return
}

func (f *Fibonacci) takeProfitCorrection(internalPosition structs.Position) (err error) {

	if internalPosition.Status != structs.StatusActive {
		return
	}

	newTakeProfit := f.calculateReducedTakeProfit(internalPosition)
	updOrder := 0

	if internalPosition.Orders.Order3.CreatedTime > 0 {
		unixTimeToReduce := time.Now().Add(time.Duration(f.config.HoursToReduceTP3*-1) * time.Hour).Unix()
		if internalPosition.Orders.Order3.CreatedTime < unixTimeToReduce &&
			internalPosition.Orders.Order3.TpModifyTime < unixTimeToReduce {

			updOrder = 3
		}
	} else if internalPosition.Orders.Order2.CreatedTime > 0 {
		unixTimeToReduce := time.Now().Add(time.Duration(f.config.HoursToReduceTP2*-1) * time.Hour).Unix()
		if internalPosition.Orders.Order2.CreatedTime < unixTimeToReduce &&
			internalPosition.Orders.Order2.TpModifyTime < unixTimeToReduce {

			updOrder = 2
		}
	} else if internalPosition.Orders.Order1.CreatedTime > 0 {
		unixTimeToReduce := time.Now().Add(time.Duration(f.config.HoursToReduceTP1*-1) * time.Hour).Unix()
		if internalPosition.Orders.Order1.CreatedTime < unixTimeToReduce &&
			internalPosition.Orders.Order1.TpModifyTime < unixTimeToReduce {

			updOrder = 1
		}
	}

	if updOrder > 0 {
		err = f.provider.ModifyTpSl(domainStructs.TpSlRequest{
			CoinPare:   f.config.CoinPare,
			TakeProfit: newTakeProfit,
			StopLoss:   internalPosition.ProviderPosition.StopLoss,
		})
		if err != nil {
			return
		}
		var storage mongo.MongoStorage
		storage, err = strategyStorage.GetStorage(f.config.Id)
		if err != nil {
			return
		}
		switch updOrder {
		case 3:
			internalPosition.Orders.Order3.TpModifyTime = time.Now().Unix()
			break
		case 2:
			internalPosition.Orders.Order2.TpModifyTime = time.Now().Unix()
			break
		case 1:
			internalPosition.Orders.Order1.TpModifyTime = time.Now().Unix()
			break
		}
		internalPosition, err = storage.SaveInternalPosition(internalPosition)
		if err != nil {
			return
		}
		msg := fmt.Sprintf("Reduced take profit from %g to %g", internalPosition.ProviderPosition.TakeProfit, newTakeProfit)
		logger.Info(msg)
		if f.config.TelegramNotifier {
			tgNotifier.Notify(msg)
		}
	}

	return
}

func (f *Fibonacci) calculateReducedTakeProfit(internalPosition structs.Position) (newTakeProfit float64) {

	oldTakeProfit := internalPosition.ProviderPosition.TakeProfit
	avgPrice := internalPosition.ProviderPosition.AvgPrice
	markPrice := internalPosition.ProviderPosition.MarkPrice

	profDiff := oldTakeProfit - avgPrice
	valueToReduce := math.Mul(math.Div(profDiff, 100), float64(f.config.PercentToReduceTP))
	newTakeProfit = oldTakeProfit - valueToReduce

	if internalPosition.ProviderPosition.Side == trading.SideBuy && newTakeProfit < markPrice {
		newTakeProfit = markPrice
	}
	if internalPosition.ProviderPosition.Side == trading.SideSell && newTakeProfit > markPrice {
		newTakeProfit = markPrice
	}
	newTakeProfit = math.Round(newTakeProfit, f.config.PricePrecision)

	return
}
