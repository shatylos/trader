package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/tgNotifier"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

func (i *Investor) handleHeapPremium(ctx context.Context, timeframeItem *_struct.HeapTimeframe) (err error) {
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

	qtyPercentOnMaxPrice := timeframeItem.Config.QtyPercentOnMaxPrice
	qtyPercentOnMinPrice := timeframeItem.Config.QtyPercentOnMinPrice

	mainCurrencyAmount := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.MainCurrency)
	tradeCurrencyQty := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.TradeCurrency)
	tradeCurrencyAmount := trading.TradeCurrencyToMain(tradeCurrencyQty, currentPrice)
	totalAmount := mainCurrencyAmount + tradeCurrencyAmount
	qtyPercentForCurrentPrice := math.MapRange(historyMinPrice, historyMaxPrice, qtyPercentOnMinPrice, qtyPercentOnMaxPrice, currentPrice)

	purposeAmount := totalAmount / 100.0 * qtyPercentForCurrentPrice
	purposeQty := purposeAmount / currentPrice

	qtyRange := math.Mul(math.Div(purposeQty, 100), timeframeItem.Config.QtyPercent)

	qty := tradeCurrencyQty - purposeQty + qtyRange

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

func (i Investor) handleHeapDiscount(ctx context.Context, timeframeItem *_struct.HeapTimeframe) (err error) {

	// calculate price from last order
	lastOrderPrice := 0.0
	if i.State.Heap.LastOrderHeap != nil {
		lastOrderPrice = i.State.Heap.LastOrderHeap.Price
	}
	currentPrice := i.State.CurrentPrice

	priceCheck := false
	priceRange := math.Mul(math.Div(currentPrice, 100), timeframeItem.Config.MinPercentRangeToSell)
	nextBuyPrice := lastOrderPrice - priceRange

	if currentPrice < nextBuyPrice {
		priceCheck = true
	}

	// calculate qty from qty on heap, available amounts and current price
	historyMinPrice, historyMaxPrice := trading.MinMaxPrice(timeframeItem.Candles)

	qtyPercentOnMaxPrice := timeframeItem.Config.QtyPercentOnMaxPrice
	qtyPercentOnMinPrice := timeframeItem.Config.QtyPercentOnMinPrice

	mainCurrencyAmount := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.MainCurrency)
	tradeCurrencyQty := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.TradeCurrency)
	tradeCurrencyAmount := trading.TradeCurrencyToMain(tradeCurrencyQty, currentPrice)
	totalAmount := mainCurrencyAmount + tradeCurrencyAmount
	qtyPercentForCurrentPrice := math.MapRange(historyMinPrice, historyMaxPrice, qtyPercentOnMinPrice, qtyPercentOnMaxPrice, currentPrice)

	purposeAmount := totalAmount / 100.0 * qtyPercentForCurrentPrice
	purposeQty := purposeAmount / currentPrice

	qtyRange := math.Mul(math.Div(purposeQty, 100), timeframeItem.Config.QtyPercent)

	qty := purposeQty - tradeCurrencyQty + qtyRange

	// do buy
	minQty := math.RoundCell(i.Config.MinQty+math.Mul(math.Div(i.Config.MinQty, 100), i.Config.CommissionBuy), i.Config.QtyPrecision)
	if priceCheck && qty > minQty {
		var providerOrderId string
		providerOrderId, err = i.doBuyOnHeap(ctx, timeframeItem, i.State.Heap.Deal, qty, currentPrice)
		if err != nil {
			if providerOrderId != "" {
				i.Config.Enabled = false
				msg := fmt.Sprintf("[%s] Order for buy on heap was created at the provider's side but error happened. Configuration disabled.", i.Config.Id)
				logger.Warning(msg)
				tgNotifier.Notify(msg)
			}
			return
		}
	}

	return
}

func (i Investor) isTimeToMoveToHeap(timeframeItem *_struct.TimeframeItem, dealRelation *entity.DealRelation) (result bool) {
	var lastOrder *entity.Order
	for _, order := range dealRelation.Orders {
		if order.OrderStatus == domainStructs.OrderStatuses.New || order.OrderStatus == domainStructs.OrderStatuses.Open {
			return
		}
		if lastOrder == nil {
			lastOrder = order
		}
		if order.CreatedTime.After(lastOrder.CreatedTime) {
			lastOrder = order
		}
	}

	if lastOrder == nil {
		return
	}

	timeToMove := lastOrder.CreatedTime.Add(timeframeItem.Config.DurationToMoveToHeap)

	if time.Now().After(timeToMove) {
		result = true
	}
	return
}

func (i Investor) moveToHeap(ctx context.Context, dealRelation *entity.DealRelation) (err error) {

	dealRelation.Deal.IsHeap = true
	dealRelation.Deal.SetClose()
	err = i.Storage.SaveDeal(ctx, dealRelation.Deal)
	if err != nil {
		return
	}
	msg := fmt.Sprintf("Deal for timeframe %s moved to heap", dealRelation.Deal.Timeframe)
	logger.Info(msg)
	if i.Config.TelegramNotifier {
		tgNotifier.Notify(msg)
	}
	return
}
