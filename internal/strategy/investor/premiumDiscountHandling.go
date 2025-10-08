package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
)

func (i *Investor) handlePremium(ctx context.Context, dealRelation *DealRelation, timeFrameItem *Timeframe) (err error) {

	if dealRelation.Deal.Status != storage.DealStatusActive {
		return
	}

	for _, dealOrder := range dealRelation.Orders {
		if i.isActiveOrder(dealOrder) {
			if i.config.Verbose {
				logger.Info(fmt.Sprintf("There is active order to %s. Waiting for the filling of the order", dealOrder.Side))
			}
			return
		}
	}

	if dealRelation.QtyInTrade > 0 && dealRelation.PriceToSell > 0 {
		price := timeFrameItem.Candles[0].Close
		priceToSell := dealRelation.PriceToSell
		if price < priceToSell {
			if i.config.Verbose {
				logger.Info(fmt.Sprintf("Too low price to sell (%g). Expected price: %g", price, priceToSell))
			}
			return
		}

		var providerOrderId string
		providerOrderId, err = i.doSell(ctx, timeFrameItem, dealRelation.Deal, dealRelation.QtyInTrade, price)
		if err != nil {
			if providerOrderId != "" {
				i.config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for sell was created at the provider's side but error happened. Configuration disabled.", i.config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
	}

	return
}

func (i *Investor) handleDiscount(ctx context.Context, dealRelation *DealRelation, timeFrameItem *Timeframe) (err error) {

	if len(dealRelation.Orders) == 0 {
		var providerOrderId string
		providerOrderId, err = i.doBuy(ctx, timeFrameItem, dealRelation.Deal)
		if err != nil {
			if providerOrderId != "" {
				i.config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
	} else if len(dealRelation.Orders) > 0 {
		var minOrderPrice float64
		var countBuyOrders int64
		for _, dealOrder := range dealRelation.Orders {
			if i.isActiveOrder(dealOrder) {
				return
			}
			if dealOrder.Side == structs.OrderSideBuy {
				countBuyOrders++
				if minOrderPrice == 0 || minOrderPrice > dealOrder.Price {
					minOrderPrice = dealOrder.Price
				}
			}
		}

		currentPrice := timeFrameItem.Candles[0].Close
		minAmountRange := math.Mul(math.Div(currentPrice, 100), timeFrameItem.Config.MinPercentRangeToBuyMore)
		currentPriceRange := minOrderPrice - currentPrice
		if currentPriceRange < minAmountRange {
			if i.config.Verbose {
				logger.Info(fmt.Sprintf("Too low price range to handle discount action (%g). Expected range: %g", currentPriceRange, minAmountRange))
			}
			return
		}
		if countBuyOrders < timeFrameItem.Config.MaxNumberOrdersToBuy {
			var providerOrderId string
			providerOrderId, err = i.doBuy(ctx, timeFrameItem, dealRelation.Deal)
			if err != nil {
				if providerOrderId != "" {
					i.config.Enabled = false
					msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.config.Id)
					logger.Warning(msg)
					tgNotifier.Notify(msg)
				}
				return
			}
		} else {
			// @TODO: move to heap and close deal
		}
	}

	return
}

func (i *Investor) handlePremiumHeap(ctx context.Context, timeFrameItem *Timeframe) (err error) {
	return
}

func (i *Investor) handleDiscountHeap(ctx context.Context, timeFrameItem *Timeframe) (err error) {
	return
}

func (i *Investor) isActiveOrder(dealOrder *storage.Order) (result bool) {
	switch dealOrder.OrderStatus {
	case structs.OrderStatuses.New,
		structs.OrderStatuses.Open,
		structs.OrderStatuses.PartiallyFilled:
		if i.config.Verbose {
			logger.Info(fmt.Sprintf("Exists order to %s with status %s. Wait for fill the order.", dealOrder.Side, dealOrder.Side))
		}
		result = true
	}
	return
}
