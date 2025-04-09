package fibonacci

import (
	"fmt"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
)

func (f *Fibonacci) actionByPosition(internalPosition structs.Position, currentPrice float64) (err error) {
	switch internalPosition.Trend {
	case TrendLong:
		err = f.actionByPositionBullish(internalPosition, currentPrice)
		break
	case TrendShort:
		err = f.actionByPositionBearish(internalPosition, currentPrice)
		break
	default:
		err = tools.AppError{
			Message: fmt.Sprintf("Unexpected trend %s for action by position", internalPosition.Trend),
		}
	}
	return
}

func (f *Fibonacci) actionByPositionBullish(internalPosition structs.Position, currentPrice float64) (err error) {
	if internalPosition.Orders.Order1.OrderId == "" &&
		currentPrice < internalPosition.FibonacciChart.EntryPoint1 &&
		currentPrice > internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP1ToFullQtyPercent), f.config.QtyPrecision)
		if epQty > f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 1, "Buy")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 1 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		}
	} else {
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
		currentPrice < internalPosition.FibonacciChart.EntryPoint2 &&
		currentPrice > internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP2ToFullQtyPercent), f.config.QtyPrecision)
		if epQty > f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 2, "Buy")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 2 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		}
	} else {
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
		currentPrice < internalPosition.FibonacciChart.EntryPoint3 &&
		currentPrice > internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP3ToFullQtyPercent), f.config.QtyPrecision)
		if epQty > f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 3, "Buy")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 3 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		}
	} else {
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

	if internalPosition.Orders.Order1.OrderId == "" &&
		currentPrice > internalPosition.FibonacciChart.EntryPoint1 &&
		currentPrice < internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP1ToFullQtyPercent), f.config.QtyPrecision)
		if epQty > f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 1, "Sell")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 1 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		}
	} else {
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
		currentPrice > internalPosition.FibonacciChart.EntryPoint2 &&
		currentPrice < internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP2ToFullQtyPercent), f.config.QtyPrecision)
		if epQty > f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 2, "Sell")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 2 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		}
	} else {
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
		currentPrice > internalPosition.FibonacciChart.EntryPoint3 &&
		currentPrice < internalPosition.FibonacciChart.StopLoss {

		epQty := math.Round(math.Mul(math.Div(internalPosition.FibonacciChart.FullQty, 100), f.config.EP3ToFullQtyPercent), f.config.QtyPrecision)
		if epQty > f.config.MinQty {
			_, err = f.openNewPosition(internalPosition, epQty, 3, "Sell")
			return
		} else if epQty > 0 && epQty < f.config.MinQty {
			logger.Warning(fmt.Sprintf("Order 3 was not created. QTY (%f) less then min qty (%f)", epQty, f.config.MinQty))
		}
	} else {
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
