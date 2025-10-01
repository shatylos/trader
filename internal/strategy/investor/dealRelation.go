package investor

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
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

		mainAmountBefore = currencyAmountTotal(&order.WalletBefore, i.config.MainCurrency)
		tradeAmountBefore = currencyAmountTotal(&order.WalletBefore, i.config.TradeCurrency)
		mainAmountAfter = currencyAmountTotal(&order.WalletAfter, i.config.MainCurrency)
		tradeAmountAfter = currencyAmountTotal(&order.WalletAfter, i.config.TradeCurrency)

		if order.Side == structs.OrderSideBuy {
			revenueMainCur -= mainAmountAfter - mainAmountBefore
			revenueTradeCur += tradeAmountAfter - tradeAmountBefore
		} else if order.Side == structs.OrderSideSell {
			revenueMainCur += mainAmountAfter - mainAmountBefore
			revenueTradeCur -= tradeAmountAfter - tradeAmountBefore
		}
	}

	dealRelation = &DealRelation{
		Deal:            deal,
		Orders:          dealOrders,
		RevenueMainCur:  revenueMainCur,
		RevenueTradeCur: revenueTradeCur,
	}
	return
}
