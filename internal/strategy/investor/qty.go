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

func (i *Investor) getQtyAndPriceToBuy(dealRelation *entity.DealRelation) (qty float64, currentPrice float64, err error) {

	qty = dealRelation.QtyToBuy
	qty = i.addCommission(qty, i.Config.CommissionBuy)
	qty = math.RoundCell(qty, i.Config.QtyPrecision)

	return
}

func (i *Investor) getQtyAndPriceToSell(dealRelation *entity.DealRelation) (qty, price float64) {
	qty = dealRelation.QtyToSell
	price = dealRelation.PriceToSell

	tradeAmountAvailable := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
	remainingBalance := tradeAmountAvailable - qty
	if qty > 0 {
		minCoinReserve := math.Mul(math.Div(qty, 100), i.Config.MinCoinReservePercent)
		if remainingBalance < minCoinReserve {
			qty = i.removeCommission(qty, i.Config.CommissionBuy)
		}
	}
	qty = math.RoundFloor(qty, i.Config.QtyPrecision)
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

func (i *Investor) getTimeframeFullAmount(ctx context.Context, timeFrameItem *_struct.TimeframeItem) (tfFullAmount float64, err error) {
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

func (i *Investor) calculateNextQtyAndPrices(ctx context.Context, timeFrameItem *_struct.TimeframeItem, dealRelation *entity.DealRelation, vwap *trading.VWAP) (err error) {

	var timeFrameFullAmount float64
	timeFrameFullAmount, err = i.getTimeframeFullAmount(ctx, timeFrameItem)
	if err != nil {
		err = apperrors.Wrap(err, "error get timeframe full amount")
		return
	}

	currentPrice := i.State.CurrentPrice
	var lastBuyOrder, lastSellOrder *entity.Order
	var lastBuyOrders, lastSellOrders []*entity.Order
	var isBuyFilled, isSellFilled bool
	for key := len(dealRelation.Orders) - 1; key >= 0; key-- {
		order := dealRelation.Orders[key]
		if order.Side == structs.OrderSideBuy && !isBuyFilled {
			lastBuyOrders = append(lastBuyOrders, order)
			if len(lastSellOrders) > 0 {
				isSellFilled = true
			}
		}
		if order.Side == structs.OrderSideSell && !isSellFilled {
			lastSellOrders = append(lastSellOrders, order)
			if len(lastBuyOrders) > 0 {
				isBuyFilled = true
			}
		}
	}
	for ii, j := 0, len(lastBuyOrders)-1; ii < j; ii, j = ii+1, j-1 {
		lastBuyOrders[ii], lastBuyOrders[j] = lastBuyOrders[j], lastBuyOrders[ii]
	}
	for ii, j := 0, len(lastSellOrders)-1; ii < j; ii, j = ii+1, j-1 {
		lastSellOrders[ii], lastSellOrders[j] = lastSellOrders[j], lastSellOrders[ii]
	}
	if len(lastSellOrders) > 0 {
		lastSellOrder = lastSellOrders[len(lastSellOrders)-1]
	}
	if len(lastBuyOrders) > 0 {
		lastBuyOrder = lastBuyOrders[len(lastBuyOrders)-1]
	}

	var buyOrderConfig, sellOrderConfig *_struct.OrderParams
	buyOrderConfig, sellOrderConfig = i.updateOrderConfigsByPrevOrders(timeFrameItem, lastSellOrder, lastBuyOrder)
	buyOrderConfig, sellOrderConfig = i.updateOrderConfigsByPrice(buyOrderConfig, sellOrderConfig, timeFrameItem, vwap)
	buyOrderConfig, sellOrderConfig = i.updateOrderConfigsByQty(buyOrderConfig, sellOrderConfig, timeFrameFullAmount, timeFrameItem, dealRelation)

	if buyOrderConfig == nil {
		dealRelation.QtyToBuy = 0
		dealRelation.PriceToBuy = 0
		dealRelation.NumOrderToBuy = 0
	} else {
		dealRelation.NumOrderToBuy = buyOrderConfig.ConfigKey + 1
		mainAmountAv := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
		avQtyToBuy := trading.MainCurrencyToTrade(mainAmountAv, currentPrice)

		purposeAmount := math.Mul(math.Div(timeFrameFullAmount, 100), buyOrderConfig.Percentage)
		purposeQty := trading.MainCurrencyToTrade(purposeAmount, currentPrice)
		qtyInTrade := dealRelation.CalcQtyInTrade()
		qtyToBuy := purposeQty - qtyInTrade

		if avQtyToBuy >= qtyToBuy {
			dealRelation.QtyToBuy = qtyToBuy
		} else if avQtyToBuy >= i.Config.MinQty {
			dealRelation.QtyToBuy = i.Config.MinQty
		} else {
			dealRelation.QtyToBuy = 0
		}
		if dealRelation.QtyToBuy > 0 {
			_, dealRelation.PriceToBuy = vwap.CalcDeviation(buyOrderConfig.VwapDeviations)
		} else {
			dealRelation.PriceToBuy = 0
		}
	}

	if sellOrderConfig == nil {
		dealRelation.QtyToSell = 0
		dealRelation.PriceToSell = 0
		dealRelation.NumOrderToSell = 0
	} else {
		dealRelation.NumOrderToSell = sellOrderConfig.ConfigKey + 1
		avQtySell := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
		purposeAmountInTrade := math.Mul(math.Div(timeFrameFullAmount, 100), sellOrderConfig.Percentage)
		purposeQtyInTrade := trading.MainCurrencyToTrade(purposeAmountInTrade, currentPrice)
		qtyInTrade := dealRelation.CalcQtyInTrade()
		qtyToSell := qtyInTrade - purposeQtyInTrade
		if avQtySell >= qtyToSell {
			dealRelation.QtyToSell = qtyToSell
			dealRelation.PriceToSell, _ = vwap.CalcDeviation(sellOrderConfig.VwapDeviations)
		} else if avQtySell >= i.Config.MinQty {
			dealRelation.QtyToSell = i.Config.MinQty
			dealRelation.PriceToSell, _ = vwap.CalcDeviation(sellOrderConfig.VwapDeviations)
		} else {
			dealRelation.QtyToSell = 0
			dealRelation.PriceToSell = 0
		}
	}

	return
}
