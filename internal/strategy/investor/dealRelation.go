package investor

import (
	"context"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
)

type DealRelation struct {
	Deal            *storage.Deal
	Orders          []*storage.Order
	RevenueMainCur  float64
	RevenueTradeCur float64
}

func (i *Investor) GetDealRelation(ctx context.Context, deal *storage.Deal) (dealRelation *DealRelation, err error) {
	if deal.Id == nil {
		msg := "Deal ID is empty. Can not get deal relations"
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}

	var dealOrders []*storage.Order
	dealOrders, err = i.Storage.GetOrdersByDealId(ctx, *deal.Id)
	if err != nil {
		return
	}

	var revenueMainCur, revenueTradeCur float64

	for _, order := range dealOrders {
		var mainAmountBefore, tradeAmountBefore, mainAmountAfter, tradeAmountAfter float64
		for _, coin := range order.WalletBefore.Available {
			if coin.Coin == i.config.MainCurrency {
				mainAmountBefore = coin.Amount
			}
			if coin.Coin == i.config.TradeCurrency {
				tradeAmountBefore = coin.Amount
			}
		}
		for _, coin := range order.WalletAfter.Available {
			if coin.Coin == i.config.MainCurrency {
				mainAmountAfter = coin.Amount
			}
			if coin.Coin == i.config.TradeCurrency {
				tradeAmountAfter = coin.Amount
			}
		}
		revenueMainCur += mainAmountBefore - mainAmountAfter
		revenueTradeCur += tradeAmountBefore - tradeAmountAfter
	}

	dealRelation = &DealRelation{
		Deal:            deal,
		Orders:          dealOrders,
		RevenueMainCur:  revenueMainCur,
		RevenueTradeCur: revenueTradeCur,
	}
	return
}
