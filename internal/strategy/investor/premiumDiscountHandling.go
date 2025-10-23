package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
)

func (i *Investor) handlePremium(ctx context.Context, dealRelation *entity.DealRelation, timeFrameItem *_struct.Timeframe) (err error) {

	if dealRelation.Deal.Status != entity.DealStatusActive {
		return
	}

	for _, dealOrder := range dealRelation.Orders {
		if i.isActiveOrder(dealOrder) {
			if i.Config.Verbose {
				logger.Info(fmt.Sprintf("There is active order to %s. Waiting for the filling of the order", dealOrder.Side))
			}
			return
		}
	}

	if dealRelation.QtyInTrade > 0 && dealRelation.PriceToSell > 0 {
		price := timeFrameItem.Candles[0].Close
		priceToSell := dealRelation.PriceToSell
		if price < priceToSell {
			if i.Config.Verbose {
				logger.Info(fmt.Sprintf("Too low price to sell (%g). Expected price: %g", price, priceToSell))
			}
			return
		}

		qty := dealRelation.QtyInTrade
		avQty := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
		if avQty < qty {
			qty = avQty
		}
		if qty < i.Config.MinQty {
			dealRelation.Deal.SetClose()
			err = i.Storage.SaveDeal(ctx, dealRelation.Deal)
			if err != nil {
				return
			}
			return
		}

		var providerOrderId string
		providerOrderId, err = i.doSell(ctx, timeFrameItem, dealRelation.Deal, qty, price)
		if err != nil {
			if providerOrderId != "" {
				i.Config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for sell was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
	}

	return
}

func (i *Investor) handleDiscount(ctx context.Context, dealRelation *entity.DealRelation, timeFrameItem *_struct.Timeframe) (err error) {

	if len(dealRelation.Orders) == 0 {
		var providerOrderId string
		providerOrderId, err = i.doBuy(ctx, timeFrameItem, dealRelation.Deal)
		if err != nil {
			if providerOrderId != "" {
				i.Config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
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
			if dealOrder.Side == structs.OrderSideBuy && dealOrder.OrderStatus != structs.OrderStatuses.Canceled {
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
			if i.Config.Verbose {
				logger.Info(fmt.Sprintf("Too low price range to handle discount action (%g). Expected range: %g", currentPriceRange, minAmountRange))
			}
			return
		}
		if countBuyOrders < timeFrameItem.Config.MaxNumberOrdersToBuy {
			var providerOrderId string
			providerOrderId, err = i.doBuy(ctx, timeFrameItem, dealRelation.Deal)
			if err != nil {
				if providerOrderId != "" {
					i.Config.Enabled = false
					msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
					logger.Warning(msg)
					tgNotifier.Notify(msg)
				}
				return
			}
		} else if i.isTimeToMoveToHeap(timeFrameItem, dealRelation) {
			err = i.moveToHeap(ctx, dealRelation)
			if err != nil {
				return
			}
		}
	}

	return
}

func (i *Investor) handlePremiumHeap(ctx context.Context, dealRelation *entity.DealRelation, timeFrameItem *_struct.Timeframe) (err error) {
	return
}

func (i *Investor) handleDiscountHeap(ctx context.Context, dealRelation *entity.DealRelation, timeFrameItem *_struct.Timeframe) (err error) {
	return
}

func (i *Investor) isActiveOrder(dealOrder *entity.Order) (result bool) {
	switch dealOrder.OrderStatus {
	case structs.OrderStatuses.New,
		structs.OrderStatuses.Open,
		structs.OrderStatuses.PartiallyFilled:
		if i.Config.Verbose {
			logger.Info(fmt.Sprintf("Exists order to %s with status %s. Wait for fill the order.", dealOrder.Side, dealOrder.Side))
		}
		result = true
	}
	return
}
