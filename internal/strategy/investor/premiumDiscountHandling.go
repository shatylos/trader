package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
)

func (i *Investor) handlePremium(ctx context.Context, dealRelation *entity.DealRelation, timeFrameItem *_struct.TimeframeItem) (err error) {
	currentPrice := timeFrameItem.Candles[0].Close
	priceToSell := dealRelation.PriceToSell

	if dealRelation.Deal.Status != entity.DealStatusActive {
		timeFrameItem.TradeStateMsg = "Handle premium. Deal is not active"
		return
	}

	if dealRelation.PriceToSell == 0 {
		timeFrameItem.TradeStateMsg = "Handle premium. Will not sell"
		return
	}

	for _, dealOrder := range dealRelation.Orders {
		if i.isActiveOrder(dealOrder) {
			timeFrameItem.TradeStateMsg = "There is an active order. Wait for fill the order"
			return
		}
	}

	if currentPrice < priceToSell {
		timeFrameItem.TradeStateMsg = fmt.Sprintf("Handle premium. Expect higher price (%.2f) to sell", priceToSell)
		return
	}

	if dealRelation.QtyInTrade == 0 {
		timeFrameItem.TradeStateMsg = "Handle premium. Qty in trade is 0. Will not sell"
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
			err = apperrors.Wrap(err, "error save deal. DealID: %v", dealRelation.Deal.Id)
			return
		}
		timeFrameItem.TradeStateMsg = "Handle premium. Not enough qty to sell. Deal closed."
		return
	}

	var providerOrderId string
	providerOrderId, err = i.doSell(ctx, timeFrameItem, dealRelation.Deal, qty, currentPrice)
	if err != nil {
		if providerOrderId != "" {
			i.Config.Enabled = false
			msg := fmt.Sprintf("[%s] Order for sell was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
			logger.Warning(msg)
			tgNotifier.Notify(msg)
		}
		err = apperrors.Wrap(err, "error do sell")
		return
	}

	return
}

func (i *Investor) handleDiscount(ctx context.Context, dealRelation *entity.DealRelation, timeFrameItem *_struct.TimeframeItem) (err error) {

	if !timeFrameItem.Config.CanOpenNewOrder {
		timeFrameItem.TradeStateMsg = "Open new order disabled in config"
		return
	}

	currentPrice := timeFrameItem.Candles[0].Close
	if dealRelation.PriceToBuy == 0 {
		timeFrameItem.TradeStateMsg = "Handle discount. Will not buy more"
		return
	}

	if currentPrice > dealRelation.PriceToBuy {
		timeFrameItem.TradeStateMsg = fmt.Sprintf("Handle discount. Expect lower price (%.2f) to buy", dealRelation.PriceToBuy)
		return
	}

	for _, dealOrder := range dealRelation.Orders {
		if i.isActiveOrder(dealOrder) {
			timeFrameItem.TradeStateMsg = "There is an active order"
			return
		}
	}

	var providerOrderId string
	providerOrderId, err = i.doBuy(ctx, timeFrameItem, dealRelation.Deal)
	if err != nil {
		if providerOrderId != "" {
			i.Config.Enabled = false
			msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
			logger.Warning(msg)
			tgNotifier.Notify(msg)
		}
		err = apperrors.Wrap(err, "error do buy")
		return
	}

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
