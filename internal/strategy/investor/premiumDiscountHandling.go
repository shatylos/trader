package investor

import (
	"fmt"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/tgNotifier"
)

func (i *Investor) handlePremium(timeFrameItem *Timeframe) (err error) {

	deal := storage.Deal{}
	err = i.Storage.GetLastDealByTimeframe(timeFrameItem.Config.Resolution, &deal)
	if err != nil {
		return
	}
	if deal.Status != storage.DealStatusActive {
		return
	}

	dealOrders := make([]storage.Order, 0)
	if deal.Id != nil {
		err = i.Storage.GetOrdersByDealId(*deal.Id, &dealOrders)
		if err != nil {
			return
		}
	}

	var qty float64
	for _, dealOrder := range dealOrders {
		qty += dealOrder.DomainOrder.Qty
	}

	if qty > 0 {
		var order *storage.Order
		var providerOrderId string
		order, providerOrderId, err = i.doSell(timeFrameItem, &deal, qty)
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

func (i *Investor) handleDiscount(timeFrameItem *Timeframe) (err error) {
	deal := storage.Deal{}
	err = i.Storage.GetLastDealByTimeframe(timeFrameItem.Config.Resolution, &deal)
	if err != nil {
		return
	}

	if deal.Status == storage.DealStatusClosed {
		deal = storage.Deal{}
	}
	if deal.Id == nil {
		deal.Timeframe = timeFrameItem.Config.Resolution
		deal.Status = storage.DealStatusNew
		err = i.Storage.SaveDeal(&deal)
		if err != nil {
			return
		}
	}

	dealOrders := make([]storage.Order, 0)
	if deal.Id != nil {
		err = i.Storage.GetOrdersByDealId(*deal.Id, &dealOrders)
		if err != nil {
			return
		}
	}

	if len(dealOrders) == 0 {
		var order *storage.Order
		var providerOrderId string
		order, providerOrderId, err = i.doBuy(timeFrameItem, &deal)
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

func (i *Investor) handlePremiumHeap(timeFrameItem *Timeframe) (err error) {
	return
}

func (i *Investor) handleDiscountHeap(timeFrameItem *Timeframe) (err error) {
	return
}
