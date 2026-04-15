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

	var lastBuyOrders, lastSellOrders []*entity.Order
	var isBuyFilled, isSellFilled bool
	lastOrder := ""
	for key := len(dealRelation.Orders) - 1; key >= 0; key-- {
		order := dealRelation.Orders[key]
		if order.Side == structs.OrderSideBuy && !isBuyFilled {
			lastOrder = structs.OrderSideBuy
			lastBuyOrders = append(lastBuyOrders, order)
			if len(lastSellOrders) > 0 {
				isSellFilled = true
			}
		}
		if order.Side == structs.OrderSideSell && !isSellFilled {
			lastOrder = structs.OrderSideSell
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

	var BuyOrderConfig *_struct.OrderParams
	var SellOrderConfig *_struct.OrderParams
	switch lastOrder {
	case structs.OrderSideBuy:
		if len(lastBuyOrders) < len(timeFrameItem.Config.BuyOrders) {
			BuyOrderConfig = &timeFrameItem.Config.BuyOrders[len(lastBuyOrders)]
		}
		SellOrderConfig = &timeFrameItem.Config.SellOrders[0]
		break
	case structs.OrderSideSell:
		BuyOrderConfig = &timeFrameItem.Config.BuyOrders[0]
		if len(lastSellOrders) < len(timeFrameItem.Config.SellOrders) {
			SellOrderConfig = &timeFrameItem.Config.SellOrders[len(lastSellOrders)]
		}
		break
	default:
		BuyOrderConfig = &timeFrameItem.Config.BuyOrders[0]
		break
	}

	if BuyOrderConfig == nil {
		dealRelation.QtyToBuy = 0
		dealRelation.PriceToBuy = 0
	} else {
		if len(lastSellOrders) == 0 {
			// if sell orders is empty then percent from full amount
			amount := math.Mul(math.Div(timeFrameFullAmount, 100), BuyOrderConfig.Percentage)
			dealRelation.QtyToBuy = trading.MainCurrencyToTrade(amount, i.State.CurrentPrice)
		} else {
			// if sell orders exists then percent from sum of last sell orders
			var amount float64
			for _, order := range lastSellOrders {
				amount += order.Amount()
			}
			amount = math.Mul(math.Div(amount, 100), BuyOrderConfig.Percentage)
			dealRelation.QtyToBuy = trading.MainCurrencyToTrade(amount, i.State.CurrentPrice)
		}
		_, dealRelation.PriceToBuy = vwap.CalcDeviation(BuyOrderConfig.VwapDeviations)
	}

	if SellOrderConfig == nil {
		dealRelation.QtyToSell = 0
		dealRelation.PriceToSell = 0
	} else {
		// percent from amount in trade
		dealRelation.QtyToSell = math.Mul(math.Div(dealRelation.QtyInTrade, 100), SellOrderConfig.Percentage)
		dealRelation.PriceToSell, _ = vwap.CalcDeviation(SellOrderConfig.VwapDeviations)
	}

	minQty := i.Config.MinQty
	doIncreaseQtyToMinQty := i.Config.DoIncreaseQtyToMinQty
	if dealRelation.QtyToBuy > 0 && dealRelation.QtyToBuy < minQty {
		mainAmountAv := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
		avQtyToBuy := trading.MainCurrencyToTrade(mainAmountAv, i.State.CurrentPrice)
		if doIncreaseQtyToMinQty && avQtyToBuy > minQty {
			dealRelation.QtyToBuy = minQty
		} else {
			dealRelation.QtyToBuy = 0
		}
	}

	if dealRelation.QtyToSell > 0 {
		avQtySell := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
		if avQtySell < dealRelation.QtyInTrade {
			dealRelation.QtyToSell = avQtySell
		}

		if dealRelation.QtyToSell < minQty {
			dealRelation.QtyToSell = 0
		}
	}

	return
}
