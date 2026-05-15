package storage

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

func (s *Storage) GetDealRelationsByPeriod(ctx context.Context, from, to time.Time) (dealRelations []*entity.DealRelation, err error) {
	var deals []*entity.Deal
	deals, err = s.GetDealsByPeriod(ctx, from, to)
	if err != nil {
		err = apperrors.Wrap(err, "error get deals by period, from: %s, to: %s", from, to)
		return
	}

	dealRelations = make([]*entity.DealRelation, len(deals))

	for item, deal := range deals {
		var dealRelation *entity.DealRelation
		dealRelation, err = s.GetDealRelation(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error get deal relation")
			return
		}
		dealRelations[item] = dealRelation
	}
	return
}

func (s *Storage) GetDealRelation(ctx context.Context, deal *entity.Deal) (dealRelation *entity.DealRelation, err error) {
	mainCurrency, ok := ctx.Value(_struct.CtxMainCurrencyKey).(string)
	if !ok {
		err = apperrors.New("MainCurrency is not accessible from context")
		return
	}
	tradeCurrency, ok := ctx.Value(_struct.CtxTradeCurrencyKey).(string)
	if !ok {
		err = apperrors.New("TradeCurrency is not accessible from context")
		return
	}

	if deal.Id == nil {
		err = apperrors.New("can not get deal relations, deal.ID is empty")
		return
	}

	var dealOrders []*entity.Order
	dealOrders, err = s.GetOrdersByDealId(ctx, *deal.Id)
	if err != nil {
		err = apperrors.Wrap(err, "error get orders by deal id: %s", *deal.Id)
		return
	}

	var revenueMainCur, revenueTradeCur, realizedPNL, unrealizedPNL, lastPrice, boughtQty, spentAmount, soldQty float64

	for _, order := range dealOrders {
		if order.OrderStatus != structs.OrderStatuses.Filled {
			continue
		}

		var mainAmountBefore, tradeAmountBefore, mainAmountAfter, tradeAmountAfter float64

		mainAmountBefore = trading.CurrencyAmountTotal(&order.WalletBefore, mainCurrency)
		tradeAmountBefore = trading.CurrencyAmountTotal(&order.WalletBefore, tradeCurrency)
		mainAmountAfter = trading.CurrencyAmountTotal(&order.WalletAfter, mainCurrency)
		tradeAmountAfter = trading.CurrencyAmountTotal(&order.WalletAfter, tradeCurrency)

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

	unrealizedPNL = revenueMainCur + trading.TradeCurrencyToMain(revenueTradeCur, lastPrice)

	averageBuyPrice := 0.0
	if spentAmount > 0 && boughtQty > 0 {
		averageBuyPrice = math.Div(spentAmount, boughtQty)
	}

	for _, order := range dealOrders {
		if order.OrderStatus != structs.OrderStatuses.Filled || order.Side != structs.OrderSideSell {
			continue
		}
		avBuyAmount := math.Mul(order.Qty, averageBuyPrice)
		realizedPNL += order.Amount() - avBuyAmount
	}

	if s.activeDealRelations[*deal.Id] != nil {
		dealRelation = s.activeDealRelations[*deal.Id]
	} else {
		dealRelation = &entity.DealRelation{}
	}

	dealRelation.Deal = deal
	dealRelation.Orders = dealOrders
	dealRelation.AverageBuyPrice = averageBuyPrice
	dealRelation.RevenueMainCur = revenueMainCur
	dealRelation.RevenueTradeCur = revenueTradeCur
	dealRelation.RealizedPNL = realizedPNL
	dealRelation.UnrealizedPNL = unrealizedPNL
	dealRelation.QtyInTrade = boughtQty - soldQty
	if deal.Status != entity.DealStatusClosed {
		s.activeDealRelations[*deal.Id] = dealRelation
	}

	return
}
