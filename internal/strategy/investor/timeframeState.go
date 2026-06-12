package investor

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
)

// getTimeframeState returns the in-memory state of the timeframe.
// On the first call the state is restored from the order history.
func (i *Investor) getTimeframeState(ctx context.Context, timeFrame _struct.Timeframe) (state *entity.TimeframeState, err error) {
	if i.TimeframeStates == nil {
		i.TimeframeStates = make(map[string]*entity.TimeframeState)
	}
	if state = i.TimeframeStates[timeFrame.Resolution()]; state != nil {
		return
	}

	var orders []*entity.Order
	orders, err = i.Storage.GetOrdersByTimeframe(ctx, timeFrame.Resolution())
	if err != nil {
		err = apperrors.Wrap(err, "error get orders by timeframe %s", timeFrame.Resolution())
		return
	}

	state = &entity.TimeframeState{Timeframe: timeFrame.Resolution()}
	for _, order := range orders {
		if order.Side == structs.OrderSideBuy {
			state.LastBuyOrder = order
		}
		if order.Side == structs.OrderSideSell {
			state.LastSellOrder = order
		}

		switch order.OrderStatus {
		case structs.OrderStatuses.Filled:
			state.ActiveOrder = nil
			if order.StateApplied {
				state.ApplyStoredOrder(order)
			} else {
				// the order was saved before TimeframeState values were stored in orders
				state.ApplyFilledOrder(order, i.Config.MinQty)
				err = i.Storage.SaveOrder(ctx, order)
				if err != nil {
					err = apperrors.Wrap(err, "error save order while restore timeframe state")
					return
				}
			}
		case structs.OrderStatuses.New,
			structs.OrderStatuses.Open,
			structs.OrderStatuses.PartiallyFilled:
			state.ActiveOrder = order
		}
	}

	i.TimeframeStates[timeFrame.Resolution()] = state
	return
}

// resetTimeframeState drops the state of the timeframe. Zero values are saved
// into the last filled order to keep the cycle boundary after restart.
func (i *Investor) resetTimeframeState(ctx context.Context, state *entity.TimeframeState) (err error) {
	lastFilled := state.LastFilledOrder
	if lastFilled != nil && (lastFilled.QtyInTrade != 0 || lastFilled.AverageBuyPrice != 0) {
		lastFilled.QtyInTrade = 0
		lastFilled.AverageBuyPrice = 0
		err = i.Storage.SaveOrder(ctx, lastFilled)
		if err != nil {
			err = apperrors.Wrap(err, "error save order while reset timeframe state")
			return
		}
	}
	state.Reset()
	return
}
