package investor

import (
	"context"
	"fmt"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"time"
)

type DealRelation struct {
	Deal            *storage.Deal
	Orders          []*storage.Order
	AverageBuyPrice float64
	PriceToSell     float64
	RevenueMainCur  float64
	RevenueTradeCur float64
	RevenueTotal    float64
	QtyInTrade      float64
}

func (d *DealRelation) GetTotalAmountBefore(mainCurrency, tradeCurrency string) (amount float64) {
	if d.Deal.ClosedTime.IsZero() {
		return
	}
	monthKey := fmt.Sprintf("%d-%d", d.Deal.ClosedTime.Year(), d.Deal.ClosedTime.Month())
	var takenOrderTime time.Time
	for _, order := range d.Orders {
		orderMonthKey := fmt.Sprintf("%d-%d", order.CreatedTime.Year(), order.CreatedTime.Month())
		if orderMonthKey != monthKey {
			continue
		}
		if takenOrderTime.IsZero() || takenOrderTime.After(order.CreatedTime) {
			mainCurr := currencyAmountTotal(&order.WalletBefore, mainCurrency)
			tradeCurr := currencyAmountTotal(&order.WalletBefore, tradeCurrency)
			amount = mainCurr + tradeCurrencyToMain(tradeCurr, order.Price)
		}
	}
	return
}

func (i *Investor) GetDealRelationsByPeriod(ctx context.Context, from, to time.Time) (dealRelations []*DealRelation, err error) {
	var deals []*storage.Deal
	deals, err = i.Storage.GetDealsByPeriod(ctx, from, to)
	if err != nil {
		return
	}

	dealRelations = make([]*DealRelation, len(deals))

	for item, deal := range deals {
		var dealRelation *DealRelation
		dealRelation, err = i.GetDealRelation(ctx, deal)
		if err != nil {
			return
		}
		dealRelations[item] = dealRelation
	}
	return
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

	timeFrameItem := i.getTimeframeItemByDeal(deal)
	var revenueMainCur, revenueTradeCur, revenueTotal, lastPrice, boughtQty, spentAmount, soldQty float64

	for _, order := range dealOrders {
		if order.OrderStatus != structs.OrderStatuses.Filled {
			continue
		}

		var mainAmountBefore, tradeAmountBefore, mainAmountAfter, tradeAmountAfter float64

		mainAmountBefore = currencyAmountTotal(&order.WalletBefore, i.config.MainCurrency)
		tradeAmountBefore = currencyAmountTotal(&order.WalletBefore, i.config.TradeCurrency)
		mainAmountAfter = currencyAmountTotal(&order.WalletAfter, i.config.MainCurrency)
		tradeAmountAfter = currencyAmountTotal(&order.WalletAfter, i.config.TradeCurrency)

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

	revenueTotal = revenueMainCur + tradeCurrencyToMain(revenueTradeCur, lastPrice)

	averageBuyPrice := 0.0
	priceToSell := 0.0
	if spentAmount > 0 && boughtQty > 0 {
		averageBuyPrice = math.Div(spentAmount, boughtQty)
		minAmountRange := math.Mul(math.Div(averageBuyPrice, 100), timeFrameItem.Config.MinPercentRangeToSell)
		priceToSell = averageBuyPrice + minAmountRange
	}

	dealRelation = &DealRelation{
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
