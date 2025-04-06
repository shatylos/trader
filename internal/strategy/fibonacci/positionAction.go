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
	}
	return
}
