package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/tgNotifier"
)

func (i *Investor) handlePremium(ctx context.Context, deal *storage.Deal, dealOrders []*storage.Order, timeFrameItem *Timeframe) (err error) {

	if deal.Status != storage.DealStatusActive {
		return
	}

	var qty float64
	for _, dealOrder := range dealOrders {
		qty += dealOrder.DomainOrder.Qty
	}

	if qty > 0 {
		var order *storage.Order
		var providerOrderId string
		order, providerOrderId, err = i.doSell(ctx, timeFrameItem, deal, qty)
		if err != nil {
			if providerOrderId != "" {
				i.config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for sell was created at the provider's side but error happened. Configuration disabled.", i.config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
		msg := fmt.Sprintf("[%s] Sold %g %s for timeframe %s", i.config.Id, order.DomainOrder.Qty, i.config.TradeCurrency, timeFrameItem.Config.Resolution)
		logger.Success(msg)
		if i.config.TelegramNotifier {
			tgNotifier.Notify(msg)
		}
	}

	return
}

func (i *Investor) handleDiscount(ctx context.Context, deal *storage.Deal, dealOrders []*storage.Order, timeFrameItem *Timeframe) (err error) {

	if deal.Status == storage.DealStatusClosed {
		deal = &storage.Deal{}
	}
	if deal.Id == nil {
		deal.Timeframe = timeFrameItem.Config.Resolution
		deal.Status = storage.DealStatusNew
		err = i.Storage.SaveDeal(ctx, deal)
		if err != nil {
			return
		}
	}

	if len(dealOrders) == 0 {
		var order *storage.Order
		var providerOrderId string
		order, providerOrderId, err = i.doBuy(ctx, timeFrameItem, deal)
		if err != nil {
			if providerOrderId != "" {
				i.config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for buy was created at the provider's side but error happened. Configuration disabled.", i.config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
		msg := fmt.Sprintf("[%s] Bought %g %s for timeframe %s", i.config.Id, order.DomainOrder.Qty, i.config.TradeCurrency, timeFrameItem.Config.Resolution)
		logger.Success(msg)
		if i.config.TelegramNotifier {
			tgNotifier.Notify(msg)
		}
	} else if len(dealOrders) > 0 {
		// @TODO: check and maybe do buy more
	}

	return
}

func (i *Investor) handlePremiumHeap(ctx context.Context, timeFrameItem *Timeframe) (err error) {
	return
}

func (i *Investor) handleDiscountHeap(ctx context.Context, timeFrameItem *Timeframe) (err error) {
	return
}
