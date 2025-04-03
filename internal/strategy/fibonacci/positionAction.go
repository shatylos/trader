package fibonacci

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/fibonacci/structs"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/math"
)

func (f *Fibonacci) actionByPosition(internalPosition structs.Position, providerPosition domainStructs.DomainPosition) (err error) {
	switch internalPosition.Trend {
	case TrendLong:
		err = f.actionByPositionBullish(internalPosition, providerPosition)
		break
	case TrendShort:
		err = f.actionByPositionBearish(internalPosition, providerPosition)
		break
	default:
		err = tools.AppError{
			Message: fmt.Sprintf("Unexpected trend %s for action by position", internalPosition.Trend),
		}
	}
	return
}

func (f *Fibonacci) actionByPositionBullish(internalPosition structs.Position, providerPosition domainStructs.DomainPosition) (err error) {

	if internalPosition.Orders.Order1.OrderId == "" &&
		providerPosition.MarkPrice < internalPosition.FibonacciChart.EntryPoint1 &&
		providerPosition.MarkPrice > internalPosition.FibonacciChart.StopLoss {

		// calculate position QTY
		var fullQty, ep1Qty float64
		fullQty, err = f.calculateFullQty(providerPosition.MarkPrice)
		if err != nil {
			return
		}
		ep1Qty = math.Round(math.Mul(math.Div(fullQty, 100), f.config.EP1ToFullQtyPercent), f.config.QtyPrecision)
		internalPosition.FullQty = fullQty

		var isCreated bool
		isCreated, err = f.openNewPosition(internalPosition, ep1Qty, 1, "Buy")
		fmt.Println(isCreated)
	}
	return
}

func (f *Fibonacci) actionByPositionBearish(internalPosition structs.Position, providerPosition domainStructs.DomainPosition) (err error) {

	if internalPosition.Orders.Order1.OrderId == "" &&
		providerPosition.MarkPrice > internalPosition.FibonacciChart.EntryPoint1 &&
		providerPosition.MarkPrice < internalPosition.FibonacciChart.StopLoss {

		// calculate position QTY
		var fullQty, ep1Qty float64
		fullQty, err = f.calculateFullQty(providerPosition.MarkPrice)
		if err != nil {
			return
		}
		ep1Qty = math.Round(math.Mul(math.Div(fullQty, 100), f.config.EP1ToFullQtyPercent), f.config.QtyPrecision)
		internalPosition.FullQty = fullQty

		var isCreated bool
		//isCreated, err = f.openNewPosition(ep1Qty, internalPosition.FibonacciChart.TakeProfit1, internalPosition.FibonacciChart.StopLoss, "Sell")
		isCreated, err = f.openNewPosition(internalPosition, ep1Qty, 1, "Sell")
		fmt.Println(isCreated)
	}
	return
}
