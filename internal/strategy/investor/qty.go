package investor

import (
	"context"
	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
)

func (i *Investor) calculateQtyToBuy(ctx context.Context, timeFrameItem *_struct.TimeframeItem, deal *entity.Deal) (qty float64, currentPrice float64, err error) {
	var tfCurrencyAvailable float64
	tfCurrencyAvailable, err = i.getTimframeFullAmoun(ctx, timeFrameItem)
	if tfCurrencyAvailable == 0 {
		return
	}

	currentPrice = timeFrameItem.Candles[0].Close

	if timeFrameItem.Config.IsEqualAllOrders && deal.EqualOrdersQty > 0 {
		qty = deal.EqualOrdersQty
		return
	}

	qtyPercent := timeFrameItem.Config.QtyPercent
	minQty := i.Config.MinQty
	doIncreaseQtyToMinQty := i.Config.DoIncreaseQtyToMinQty

	//qty = tfCurrencyAvailable / 100 * qtyPercent / currentPrice
	qty = math.Div(math.Mul(math.Div(tfCurrencyAvailable, 100), qtyPercent), currentPrice)

	if qty < minQty {
		if doIncreaseQtyToMinQty {
			qty = minQty
		} else {
			qty = 0
		}
	}
	qty = i.addCommission(qty, i.Config.CommissionBuy)
	qty = math.RoundCell(qty, i.Config.QtyPrecision)

	return
}

func (i *Investor) calculateQtyToSell(qty float64, isHeap bool) (qtyResult float64) {
	tradeAmountAvailable := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
	remainingBalance := tradeAmountAvailable - qty
	if qty > 0 && !isHeap {
		minCoinReserve := math.Mul(math.Div(qty, 100), i.Config.MinCoinReservePercent)
		if remainingBalance < minCoinReserve {
			qty = i.removeCommission(qty, i.Config.CommissionBuy)
		}
	}
	qtyResult = math.RoundFloor(qty, i.Config.QtyPrecision)
	return
}

func (i *Investor) addCommission(qty, commission float64) (calculatedQty float64) {
	calculatedQty = math.Mul(math.Div(qty, 100), 100+commission)
	return
}

func (i *Investor) removeCommission(qty, commission float64) (calculatedQty float64) {
	calculatedQty = math.Mul(math.Div(qty, 100+commission), 100)
	return
}

func (i *Investor) getTimframeFullAmoun(ctx context.Context, timeFrameItem *_struct.TimeframeItem) (tfFullAmount float64, err error) {
	fullAmount := i.getMainCurrencyAvailable()
	var deals []*entity.Deal
	deals, err = i.Storage.GetActiveDeals(ctx)
	if err != nil {
		err = apperrors.Wrap(err, "error get active deals")
		return
	}

	for _, deal := range deals {
		var dealRelation *entity.DealRelation
		dealRelation, err = i.Storage.GetDealRelation(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error get deal relation")
			return
		}
		for _, order := range dealRelation.Orders {
			if order.OrderStatus == structs.OrderStatuses.PartiallyFilled {
				err = apperrors.New("partially filled order. Wait for fill or cancel the order before calculate full amount")
				return
			}
			if order.OrderStatus == structs.OrderStatuses.Filled {
				if order.Side == structs.OrderSideBuy {
					fullAmount += order.Amount()
				}
				if order.Side == structs.OrderSideSell {
					fullAmount -= order.Amount()
				}
			}
		}
	}

	//tfFullAmount = fullAmount / 100 * timeFrameItem.Config.FullAmountPercent
	tfFullAmount = math.Mul(math.Div(fullAmount, 100), timeFrameItem.Config.FullAmountPercent)
	return
}
