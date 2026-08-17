package investor

import (
	"context"
	"fmt"

	"github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
)

func (i *Investor) getQtyAndPriceToBuy(state *entity.TimeframeState) (qty float64, price float64) {

	qty = state.QtyToBuy
	qty = i.addCommission(qty, i.Config.CommissionBuy)
	qty = math.RoundCell(qty, i.Config.QtyPrecision)

	price = math.Round(state.PriceToBuy, i.Config.PricePrecision)

	if price > 0 && price > i.State.CurrentPrice {
		price = i.State.CurrentPrice
	}

	return
}

func (i *Investor) getQtyAndPriceToSell(state *entity.TimeframeState) (qty, price float64) {
	qty = state.QtyToSell

	tradeAmountAvailable := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
	remainingBalance := tradeAmountAvailable - qty
	if qty > 0 {
		minCoinReserve := math.Mul(math.Div(qty, 100), i.Config.MinCoinReservePercent)
		if remainingBalance < minCoinReserve {
			qty = i.removeCommission(qty, i.Config.CommissionBuy)
		}
	}

	price = state.PriceToSell

	if price > 0 && price < i.State.CurrentPrice {
		logger.Info(fmt.Sprintf("Modified sell price. Old price: %g. New price: %g.", price, i.State.CurrentPrice))
		price = i.State.CurrentPrice
	}

	price = math.RoundCell(price, i.Config.PricePrecision)
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

	if timeFrameItem.HasParent() {
		tfFullAmount = timeFrameItem.ChildShareAmount()
		return
	}

	fullAmount := i.getMainCurrencyAvailable()

	for key := range i.Timeframes {
		var state *entity.TimeframeState
		state, err = i.getTimeframeState(ctx, &(i.Timeframes[key]))
		if err != nil {
			err = apperrors.Wrap(err, "error get timeframe state")
			return
		}
		if state.ActiveOrder != nil && state.ActiveOrder.OrderStatus == structs.OrderStatuses.PartiallyFilled {
			err = apperrors.New("partially filled order. Wait for fill or cancel the order before calculate full amount")
			return
		}
		// amount spent on coins that are still in trade minus realized profit of the current cycle
		fullAmount += math.Mul(state.QtyInTrade, state.AverageBuyPrice) - state.RealizedPNL
	}

	tfFullAmount = math.Mul(math.Div(fullAmount, 100), timeFrameItem.Config.FullAmountPercent)
	return
}

func (i *Investor) calculateNextQtyAndPrices(ctx context.Context, timeFrameItem *_struct.TimeframeItem, state *entity.TimeframeState, vwap *trading.VWAP) (err error) {

	var timeFrameFullAmount float64
	timeFrameFullAmount, err = i.getTimeframeFullAmount(ctx, timeFrameItem)
	if err != nil {
		err = apperrors.Wrap(err, "error get timeframe full amount")
		return
	}

	currentPrice := i.State.CurrentPrice

	var buyOrderConfig, sellOrderConfig *_struct.OrderParams
	buyOrderConfig, sellOrderConfig = i.updateOrderConfigsByPrevOrders(timeFrameItem, state.LastSellOrder, state.LastBuyOrder)
	buyOrderConfig, sellOrderConfig = i.updateOrderConfigsByPrice(buyOrderConfig, sellOrderConfig, timeFrameItem, vwap)
	buyOrderConfig, sellOrderConfig = i.updateOrderConfigsByQty(buyOrderConfig, sellOrderConfig, timeFrameFullAmount, timeFrameItem, state.QtyInTrade)

	if buyOrderConfig == nil {
		state.QtyToBuy = 0
		state.PriceToBuy = 0
		state.NumOrderToBuy = 0
	} else {
		state.NumOrderToBuy = buyOrderConfig.ConfigKey + 1
		state.QtyToBuy = 0
		mainAmountAv := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.MainCurrency)
		avQtyToBuy := trading.MainCurrencyToTrade(mainAmountAv, currentPrice)

		purposeAmount := math.Mul(math.Div(timeFrameFullAmount, 100), buyOrderConfig.Percentage)
		purposeQty := trading.MainCurrencyToTrade(purposeAmount, currentPrice)
		qtyToBuy := purposeQty - state.QtyInTrade
		if timeFrameItem.HasChildren() {
			qtyToBuy = purposeQty
		}

		if avQtyToBuy >= qtyToBuy {
			state.QtyToBuy = qtyToBuy
		} else if avQtyToBuy >= i.Config.MinQty {
			state.QtyToBuy = avQtyToBuy - i.Config.MinQty
		}
		if state.QtyToBuy < i.Config.MinQty {
			state.QtyToBuy = 0
		}
		if state.QtyToBuy > 0 {
			_, state.PriceToBuy = vwap.CalcDeviation(buyOrderConfig.VwapDeviations)
		} else {
			state.PriceToBuy = 0
		}
	}

	if sellOrderConfig == nil {
		state.QtyToSell = 0
		state.PriceToSell = 0
		state.NumOrderToSell = 0
	} else {
		state.NumOrderToSell = sellOrderConfig.ConfigKey + 1
		state.QtyToSell = 0
		state.PriceToSell = 0
		avQtySell := trading.CurrencyAmountAvailable(i.State.Wallet, i.Config.TradeCurrency)
		var qtyToSell float64
		if timeFrameItem.HasParent() {
			prevPercentage := sellOrderConfig.Percentage + sellOrderConfig.PercentageDiff
			onePercentQty := math.Div(state.QtyInTrade, prevPercentage)
			purposeQtyInTrade := math.Mul(onePercentQty, sellOrderConfig.Percentage)
			qtyToSell = state.QtyInTrade - purposeQtyInTrade
			if sellOrderConfig.Percentage == 0.0 {
				qtyToSell = state.QtyInTrade
			}
		} else {
			purposeAmountInTrade := math.Mul(math.Div(timeFrameFullAmount, 100), sellOrderConfig.Percentage)
			purposeQtyInTrade := trading.MainCurrencyToTrade(purposeAmountInTrade, currentPrice)
			qtyToSell = state.QtyInTrade - purposeQtyInTrade
		}

		if avQtySell >= qtyToSell {
			state.QtyToSell = qtyToSell
			state.PriceToSell, _ = vwap.CalcDeviation(sellOrderConfig.VwapDeviations)
		} else if avQtySell >= i.Config.MinQty {
			state.QtyToSell = avQtySell - i.Config.MinQty
			state.PriceToSell, _ = vwap.CalcDeviation(sellOrderConfig.VwapDeviations)
		}
		if state.QtyToSell < i.Config.MinQty {
			state.QtyToSell = 0
			state.PriceToSell = 0
		}
	}

	state.UnrealizedPNL = 0
	if state.QtyInTrade > 0 {
		state.UnrealizedPNL = math.Mul(state.QtyInTrade, currentPrice) - math.Mul(state.QtyInTrade, state.AverageBuyPrice)
	}

	return
}
