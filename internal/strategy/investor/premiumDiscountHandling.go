package investor

import (
	"github.com/shatylos/trader/internal/strategy/investor/storage"
)

func (i *Investor) handlePremium(timeFrameItem *Timeframe) (err error) {

	deal := storage.Deal{}
	err = i.Storage.GetLastDealByTimeframe(timeFrameItem.Config.Resolution, &deal)
	if err != nil {
		return
	}
	if deal.Id == "" {
		return
	}

	// do sell checking last order

	return
}

func (i *Investor) handleDiscount(timeFrameItem *Timeframe) (err error) {

	deal := storage.Deal{}
	err = i.Storage.GetLastDealByTimeframe(timeFrameItem.Config.Resolution, &deal)
	if err != nil {
		return
	}
	if deal.Id == "" {
		deal.Timeframe = timeFrameItem.Config.Resolution
		deal.Status = storage.DealStatusNew
		err = i.Storage.SaveDeal(&deal)
		if err != nil {
			return
		}
	}

	dealOrders := make([]storage.Order, 0)
	err = i.Storage.GetOrdersByDealId(deal.Id, dealOrders)
	if err != nil {
		return
	}

	if len(dealOrders) == 0 {
		// do buy
	} else if len(dealOrders) > 0 {
		
	}

	// do buy checking last order
	//var dealOrders []storage.Order

	return
}

func (i *Investor) handlePremiumHeap(timeFrameItem *Timeframe) (err error) {
	return
}

func (i *Investor) handleDiscountHeap(timeFrameItem *Timeframe) (err error) {
	return
}
