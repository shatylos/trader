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

func (i *Investor) handleHeapPremium(ctx context.Context, heapTimeframe *_struct.HeapTimeframe) (err error) {
	// calculate price from last order and average price of heap
	lastOrderPrice := 0.0
	if heapTimeframe.HeapStatus.LastOrderHeap != nil {
		lastOrderPrice = heapTimeframe.HeapStatus.LastOrderHeap.Price
	}
	averagePrice := heapTimeframe.HeapStatus.Price
	currentPrice := i.State.CurrentPrice

	if averagePrice == 0 {
		if i.Config.Verbose {
			logger.Info("Average price of heap is 0. Premium for heap can not be handled")
		}
		return
	}

	priceCheck := false
	priceRange := math.Mul(math.Div(currentPrice, 100), heapTimeframe.Config.MinPercentRangeToSell)
	nextSellPrice := lastOrderPrice + priceRange
	if currentPrice > averagePrice && currentPrice > nextSellPrice {
		priceCheck = true
	}

	qtyRange := math.Mul(math.Div(heapTimeframe.HeapStatus.PurposeQty, 100), heapTimeframe.Config.QtyPercent)

	qty := heapTimeframe.HeapStatus.Qty - heapTimeframe.HeapStatus.PurposeQty + qtyRange

	// do sell
	minQty := math.RoundCell(i.Config.MinQty+math.Mul(math.Div(i.Config.MinQty, 100), i.Config.CommissionSell), i.Config.QtyPrecision)
	if priceCheck && qty > minQty {
		var providerOrderId string
		providerOrderId, err = i.doSell(ctx, heapTimeframe, heapTimeframe.HeapStatus.Deal, qty, currentPrice)
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

func (i *Investor) handleHeapDiscount(ctx context.Context, heapTimeframe *_struct.HeapTimeframe) (err error) {

	// calculate price from last order
	lastOrderPrice := 0.0
	if heapTimeframe.HeapStatus.LastOrderHeap != nil {
		lastOrderPrice = heapTimeframe.HeapStatus.LastOrderHeap.Price
	}
	currentPrice := i.State.CurrentPrice

	priceCheck := false
	priceRange := math.Mul(math.Div(currentPrice, 100), heapTimeframe.Config.MinPercentRangeToSell)
	nextBuyPrice := lastOrderPrice - priceRange

	if currentPrice < nextBuyPrice {
		priceCheck = true
	}

	tradeCurrencyQty := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.TradeCurrency)
	qtyRange := math.Mul(math.Div(heapTimeframe.HeapStatus.PurposeQty, 100), heapTimeframe.Config.QtyPercent)
	qty := heapTimeframe.HeapStatus.PurposeQty - tradeCurrencyQty + qtyRange

	// do buy
	minQty := math.RoundCell(i.Config.MinQty+math.Mul(math.Div(i.Config.MinQty, 100), i.Config.CommissionBuy), i.Config.QtyPrecision)
	if priceCheck && qty > minQty {
		var providerOrderId string
		providerOrderId, err = i.doBuyOnHeap(ctx, heapTimeframe, heapTimeframe.HeapStatus.Deal, qty, currentPrice)
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

func (i *Investor) isTimeToMoveToHeap(timeframeItem *_struct.TimeframeItem, dealRelation *entity.DealRelation) (result bool) {
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

func (i *Investor) moveToHeap(ctx context.Context, dealRelation *entity.DealRelation) (err error) {

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

func (i *Investor) UpdateHeapStatus(ctx context.Context) (err error) {
	var dealRelations []*entity.DealRelation
	dealRelations, err = i.Storage.GetDealRelationsOnHeap(ctx)
	if err != nil {
		return
	}

	buyQty := 0.0
	buyAmount := 0.0
	sellQty := 0.0
	sellAmount := 0.0

	var lastOrderHeap *entity.Order
	var lastOrderMoved *entity.Order

	for _, dealRelation := range dealRelations {
		for _, order := range dealRelation.Orders {
			if order.Timeframe == i.HeapTimeframe.Config.Resolution {
				if lastOrderHeap == nil || order.CreatedTime.After(lastOrderHeap.CreatedTime) {
					lastOrderHeap = order
				}
			} else {
				if lastOrderMoved == nil || order.CreatedTime.After(lastOrderMoved.CreatedTime) {
					lastOrderMoved = order
				}
			}
			if order.Side == domainStructs.OrderSideBuy {
				buyQty += order.Qty
				buyAmount += order.Amount()
			}
			if order.Side == domainStructs.OrderSideSell {
				sellQty += order.Qty
				sellAmount += order.Amount()
			}
		}
	}

	qty := buyQty - sellQty
	price := 0.0
	if buyQty > 0 {
		price = math.Div(buyAmount, buyQty)
	}

	var deal *entity.Deal
	deal, err = i.Storage.GetActiveDealByTimeframe(ctx, &i.HeapTimeframe)
	if err != nil {
		return
	}

	isSideways, sidewaysKLinesAmount := trading.CheckSideways(i.HeapTimeframe.Candles, i.HeapTimeframe.Config.SidewaysMinCandlesAmount, i.HeapTimeframe.Config.SidewaysPercentToPrice)
	i.HeapTimeframe.HeapStatus.IsSidewaysState = isSideways
	i.HeapTimeframe.HeapStatus.Trend, i.HeapTimeframe.HeapStatus.TrendSlope = trading.GetTrendLinearRegression(i.HeapTimeframe.Candles)

	sidewaysKlines := make([]domainStructs.DomainCandle, sidewaysKLinesAmount)
	copy(sidewaysKlines, i.HeapTimeframe.Candles[:sidewaysKLinesAmount])
	i.HeapTimeframe.HeapStatus.PremiumDiscount = trading.PremiumDiscount(sidewaysKlines)

	i.HeapTimeframe.HeapStatus.Qty = qty
	i.HeapTimeframe.HeapStatus.Price = price
	i.HeapTimeframe.HeapStatus.Deal = deal
	i.HeapTimeframe.HeapStatus.LastOrderHeap = lastOrderHeap
	i.HeapTimeframe.HeapStatus.LastOrderMoved = lastOrderMoved

	zone := trading.ZoneNeutral
	if i.HeapTimeframe.HeapStatus.PremiumDiscount > i.HeapTimeframe.Config.SidewaysPremiumCoefficient {
		zone = trading.ZonePremium
	}
	if i.HeapTimeframe.HeapStatus.PremiumDiscount < i.HeapTimeframe.Config.SidewaysDiscountCoefficient {
		zone = trading.ZoneDiscount
	}
	i.HeapTimeframe.HeapStatus.Zone = zone

	historyMinPrice, historyMaxPrice := trading.MinMaxPrice(i.HeapTimeframe.Candles)

	qtyPercentOnMaxPrice := i.HeapTimeframe.Config.QtyPercentOnMaxPrice
	qtyPercentOnMinPrice := i.HeapTimeframe.Config.QtyPercentOnMinPrice

	mainCurrencyAmount := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.MainCurrency)
	tradeCurrencyQty := trading.CurrencyAmountTotal(i.State.Wallet, i.Config.TradeCurrency)
	tradeCurrencyAmount := trading.TradeCurrencyToMain(tradeCurrencyQty, i.State.CurrentPrice)
	totalAmount := mainCurrencyAmount + tradeCurrencyAmount
	qtyPercentForCurrentPrice := math.MapRange(historyMinPrice, historyMaxPrice, qtyPercentOnMinPrice, qtyPercentOnMaxPrice, i.State.CurrentPrice)

	i.HeapTimeframe.HeapStatus.PurposeQtyEqv = totalAmount / 100.0 * qtyPercentForCurrentPrice
	i.HeapTimeframe.HeapStatus.PurposeQty = i.HeapTimeframe.HeapStatus.PurposeQtyEqv / i.State.CurrentPrice

	i.HeapTimeframe.HeapStatus.QtyExcess = i.HeapTimeframe.HeapStatus.Qty - i.HeapTimeframe.HeapStatus.PurposeQty
	i.HeapTimeframe.HeapStatus.QtyExcessEqv = math.Mul(i.HeapTimeframe.HeapStatus.QtyExcess, i.State.CurrentPrice)

	return
}
