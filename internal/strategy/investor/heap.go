package investor

import (
	"context"
	"fmt"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
)

func (i Investor) handleHeapPremium(ctx context.Context, timeframeItem *_struct.Timeframe) (err error) {
	// calculate price from last order and average price of heap
	lastOrderPrice := 0.0
	if i.State.Heap.LastOrderHeap != nil {
		lastOrderPrice = i.State.Heap.LastOrderHeap.Price
	}
	averagePrice := i.State.Heap.Price
	currentPrice := i.State.CurrentPrice

	if averagePrice == 0 {
		if i.Config.Verbose {
			logger.Info("Average price of heap is 0. Premium for heap can not be handled")
		}
		return
	}

	priceCheck := false
	priceRange := math.Mul(math.Div(currentPrice, 100), timeframeItem.Config.MinPercentRangeToSell)
	nextSellPrice := lastOrderPrice + priceRange
	if currentPrice > averagePrice && currentPrice > nextSellPrice {
		priceCheck = true
	}

	// calculate qty from qty on heap, available amounts and current price
	historyMinPrice, historyMaxPrice := trading.MinMaxPrice(timeframeItem.Candles)

	qtyPercentOnMaxPrice := timeframeItem.Config.HeapConfig.QtyPercentOnMaxPrice
	qtyPercentOnMinPrice := timeframeItem.Config.HeapConfig.QtyPercentOnMinPrice

	mainCurrencyAmount := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.MainCurrency)
	tradeCurrencyQty := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.TradeCurrency)
	tradeCurrencyAmount := trading.TradeCurrencyToMain(tradeCurrencyQty, currentPrice)
	totalAmount := mainCurrencyAmount + tradeCurrencyAmount
	qtyPercentForCurrentPrice := math.MapRange(historyMinPrice, historyMaxPrice, qtyPercentOnMinPrice, qtyPercentOnMaxPrice, currentPrice)

	purposeAmount := totalAmount / 100.0 * qtyPercentForCurrentPrice
	purposeQty := purposeAmount / currentPrice
	qty := tradeCurrencyQty - purposeQty

	// do sell
	minQty := math.RoundCell(i.Config.MinQty+math.Mul(math.Div(i.Config.MinQty, 100), i.Config.CommissionSell), i.Config.QtyPrecision)
	if priceCheck && qty > minQty {
		var providerOrderId string
		providerOrderId, err = i.doSell(ctx, timeframeItem, i.State.Heap.Deal, qty, currentPrice)
		if err != nil {
			if providerOrderId != "" {
				i.Config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for sell on heap was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
	}

	return
}

func (i Investor) handleHeapDiscount(ctx context.Context, timeframeItem *_struct.Timeframe) (err error) {

	// calculate price from last order

	// calculate qty from qty on heap, available amounts and current price

	// do buy

	return
}
