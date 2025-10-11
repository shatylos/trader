package investor

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

func (i *Investor) GetDealRelationsByPeriod(ctx context.Context, from, to time.Time) (dealRelations []*entity.DealRelation, err error) {
	var deals []*entity.Deal
	deals, err = i.Storage.GetDealsByPeriod(ctx, from, to)
	if err != nil {
		return
	}

	dealRelations = make([]*entity.DealRelation, len(deals))

	for item, deal := range deals {
		var dealRelation *entity.DealRelation
		dealRelation, err = i.GetDealRelation(ctx, deal)
		if err != nil {
			return
		}
		dealRelations[item] = dealRelation
	}
	return
}

func (i *Investor) GetDealRelation(ctx context.Context, deal *entity.Deal) (dealRelation *entity.DealRelation, err error) {
	if deal.Id == nil {
		msg := "Deal ID is empty. Can not get deal relations"
		logger.Error(msg)
		err = tools.AppError{Message: msg}
		return
	}

	var dealOrders []*entity.Order
	dealOrders, err = i.Storage.GetOrdersByDealId(ctx, *deal.Id)
	if err != nil {
		return
	}

	timeFrameItem := i.getTimeframeItemByDeal(deal)
	var revenueMainCur, revenueTradeCur, revenueTotal, lastPrice, boughtQty, spentAmount, soldQty float64

	for _, order := range dealOrders {
		if order.OrderStatus != structs.OrderStatuses.Filled {
			continue
		}

		var mainAmountBefore, tradeAmountBefore, mainAmountAfter, tradeAmountAfter float64

		mainAmountBefore = trading.CurrencyAmountTotal(&order.WalletBefore, i.config.MainCurrency)
		tradeAmountBefore = trading.CurrencyAmountTotal(&order.WalletBefore, i.config.TradeCurrency)
		mainAmountAfter = trading.CurrencyAmountTotal(&order.WalletAfter, i.config.MainCurrency)
		tradeAmountAfter = trading.CurrencyAmountTotal(&order.WalletAfter, i.config.TradeCurrency)

		revenueMainCur += mainAmountAfter - mainAmountBefore
		revenueTradeCur += tradeAmountAfter - tradeAmountBefore
		lastPrice = order.Price

		if order.Side == structs.OrderSideBuy {
			boughtQty += order.Qty
			spentAmount += order.Amount()
		}
		if order.Side == structs.OrderSideSell {
			soldQty += order.Qty
		}
	}

	revenueTotal = revenueMainCur + trading.TradeCurrencyToMain(revenueTradeCur, lastPrice)

	averageBuyPrice := 0.0
	priceToSell := 0.0
	if spentAmount > 0 && boughtQty > 0 {
		averageBuyPrice = math.Div(spentAmount, boughtQty)
		minAmountRange := math.Mul(math.Div(averageBuyPrice, 100), timeFrameItem.Config.MinPercentRangeToSell)
		priceToSell = averageBuyPrice + minAmountRange
	}

	dealRelation = &entity.DealRelation{
		Deal:            deal,
		Orders:          dealOrders,
		AverageBuyPrice: averageBuyPrice,
		PriceToSell:     priceToSell,
		RevenueMainCur:  revenueMainCur,
		RevenueTradeCur: revenueTradeCur,
		RevenueTotal:    revenueTotal,
		QtyInTrade:      boughtQty - soldQty,
	}
	return
}
