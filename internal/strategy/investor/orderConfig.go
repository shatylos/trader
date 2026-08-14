package investor

import (
	"fmt"

	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
)

func (i *Investor) updateOrderConfigsByPrevOrders(timeFrameItem *_struct.TimeframeItem, lastSellOrder, lastBuyOrder *entity.Order) (buyOrderConfig, sellOrderConfig *_struct.OrderParams) {
	lenBuyConf := len(timeFrameItem.Config.BuyOrders)
	lenSellConf := len(timeFrameItem.Config.SellOrders)
	if lenBuyConf == 0 || lenSellConf == 0 {
		return
	}

	if lastBuyOrder == nil {
		buyOrderConfig = &timeFrameItem.Config.BuyOrders[0]
	} else {
		if lastBuyOrder.ConfigKey+1 < len(timeFrameItem.Config.BuyOrders) {
			buyOrderConfig = &timeFrameItem.Config.BuyOrders[lastBuyOrder.ConfigKey+1]
		} else if lastSellOrder != nil && lastSellOrder.CreatedTime.After(lastBuyOrder.CreatedTime) {
			// start to buy again if all buy orders configs are used and there is a sell order after last buy order
			buyOrderConfig = &timeFrameItem.Config.BuyOrders[0]
		}

		if lastSellOrder == nil || lastBuyOrder.CreatedTime.After(lastSellOrder.CreatedTime) {
			sellOrderConfig = &timeFrameItem.Config.SellOrders[0]
		} else if lastSellOrder.ConfigKey+1 < len(timeFrameItem.Config.SellOrders) {
			sellOrderConfig = &timeFrameItem.Config.SellOrders[lastSellOrder.ConfigKey+1]
		}
	}

	return
}

func (i *Investor) updateOrderConfigsByPrice(buyOrderConfig, sellOrderConfig *_struct.OrderParams, timeFrameItem *_struct.TimeframeItem, vwap *trading.VWAP) (resBuyOrderConfig, resSellOrderConfig *_struct.OrderParams) {
	if buyOrderConfig != nil {
		for ii := len(timeFrameItem.Config.BuyOrders) - 1; ii >= 0; ii-- {
			if buyOrderConfig.ConfigKey >= ii {
				resBuyOrderConfig = &timeFrameItem.Config.BuyOrders[ii]
				break
			}
			conf := timeFrameItem.Config.BuyOrders[ii]
			_, lowerBand := vwap.CalcDeviation(conf.VwapDeviations)
			if i.State.CurrentPrice < lowerBand {
				resBuyOrderConfig = &timeFrameItem.Config.BuyOrders[ii]
				break
			}
		}
	}
	if sellOrderConfig != nil {
		for ii := len(timeFrameItem.Config.SellOrders) - 1; ii >= 0; ii-- {
			if sellOrderConfig.ConfigKey >= ii {
				resSellOrderConfig = &timeFrameItem.Config.SellOrders[ii]
				break
			}
			conf := timeFrameItem.Config.SellOrders[ii]
			upperBand, _ := vwap.CalcDeviation(conf.VwapDeviations)
			if i.State.CurrentPrice > upperBand {
				resSellOrderConfig = &timeFrameItem.Config.SellOrders[ii]
				break
			}
		}
	}
	return
}

func (i *Investor) updateOrderConfigsByQty(buyOrderConfig, sellOrderConfig *_struct.OrderParams, timeFrameFullAmount float64, timeFrameItem *_struct.TimeframeItem, qtyInTrade float64) (resBuyOrderConfig, resSellOrderConfig *_struct.OrderParams) {
	currentPrice := i.State.CurrentPrice
	if buyOrderConfig != nil {
		for ii := buyOrderConfig.ConfigKey; ii < len(timeFrameItem.Config.BuyOrders); ii++ {
			tryBuyOrderConfig := timeFrameItem.Config.BuyOrders[ii]

			purposeAmount := math.Mul(math.Div(timeFrameFullAmount, 100), tryBuyOrderConfig.Percentage)
			purposeQty := trading.MainCurrencyToTrade(purposeAmount, currentPrice)
			qtyToBuy := purposeQty - qtyInTrade
			if qtyToBuy >= i.Config.MinQty {
				resBuyOrderConfig = &timeFrameItem.Config.BuyOrders[ii]
				break
			} else if i.Config.Verbose {
				logger.Info(fmt.Sprintf("ID: %s. Timeframe %s. Calculated qty to buy %g. This is less then min qty", i.GetId(), timeFrameItem.Resolution(), qtyToBuy))
			}
		}
	}
	if sellOrderConfig != nil {
		for ii := sellOrderConfig.ConfigKey; ii < len(timeFrameItem.Config.SellOrders); ii++ {
			trySellOrderConfig := timeFrameItem.Config.SellOrders[ii]

			purposeAmountInTrade := math.Mul(math.Div(timeFrameFullAmount, 100), trySellOrderConfig.Percentage)
			purposeQtyInTrade := trading.MainCurrencyToTrade(purposeAmountInTrade, currentPrice)
			qtyToSell := qtyInTrade - purposeQtyInTrade
			if qtyToSell >= i.Config.MinQty {
				resSellOrderConfig = &timeFrameItem.Config.SellOrders[ii]
				break
			} else if i.Config.Verbose {
				logger.Info(fmt.Sprintf("ID: %s. Timeframe %s. Calculated qty to sell %g. This is less then min qty", i.GetId(), timeFrameItem.Resolution(), qtyToSell))
			}
		}
	}
	return
}
