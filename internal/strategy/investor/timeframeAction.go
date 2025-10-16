package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

func (i *Investor) handleTimeframe(ctx context.Context, timeFrameItem *_struct.Timeframe) (err error) {
	err = i.loadCandles(timeFrameItem)
	if err != nil {
		return
	}

	var deal *entity.Deal
	deal, err = i.Storage.GetActiveDealByTimeframe(ctx, timeFrameItem)
	if err != nil {
		return
	}

	var dealRelation *entity.DealRelation
	dealRelation, err = i.Storage.GetDealRelation(ctx, deal)
	if err != nil {
		return
	}

	isSideways, sidewaysKLinesAmount := trading.CheckSideways(timeFrameItem.Candles, timeFrameItem.Config.SidewaysMinCandlesAmount, timeFrameItem.Config.SidewaysPercentToPrice)

	sidewaysKlines := make([]domainStructs.DomainCandle, sidewaysKLinesAmount)
	copy(sidewaysKlines, timeFrameItem.Candles[:sidewaysKLinesAmount])
	premiumDiscount := trading.PremiumDiscount(sidewaysKlines)

	zone := trading.ZoneNeutral
	if premiumDiscount > timeFrameItem.Config.SidewaysPremiumCoefficient {
		zone = trading.ZonePremium
	}
	if premiumDiscount < timeFrameItem.Config.SidewaysDiscountCoefficient {
		zone = trading.ZoneDiscount
	}

	for _, dealOrder := range dealRelation.Orders {
		switch dealOrder.OrderStatus {
		case domainStructs.OrderStatuses.New,
			domainStructs.OrderStatuses.Open,
			domainStructs.OrderStatuses.PartiallyFilled:
			err = i.updateOrder(ctx, dealRelation.Deal, dealOrder, timeFrameItem)
			if err != nil {
				return
			}
			if dealOrder.OrderStatus != domainStructs.OrderStatuses.Filled {

				if (dealOrder.Side == domainStructs.OrderSideBuy && zone == trading.ZonePremium) ||
					(dealOrder.Side == domainStructs.OrderSideSell && zone == trading.ZoneDiscount) {

					err = i.doCancel(ctx, dealOrder)
					if err != nil {
						return
					}
					return
				}

				if i.Config.Verbose {
					logger.Info(fmt.Sprintf("Wait for fill the order %s", dealOrder.OrderId))
				}
				return
			}
		}
	}

	if !isSideways {
		// @TODO: зберегти стан
		//timeFrameItem.TrendState =
		if i.Config.Verbose {
			logger.Info(fmt.Sprintf("Timeframe %s is not in sideways", timeFrameItem.Config.Resolution))
		}
		return
	}

	if zone == trading.ZonePremium {
		err = i.handlePremium(ctx, dealRelation, timeFrameItem)
		if err != nil {
			return
		}
	}
	if zone == trading.ZoneDiscount {
		err = i.handleDiscount(ctx, dealRelation, timeFrameItem)
		if err != nil {
			return
		}
	}

	return
}

func (i *Investor) handleHeapTimeframe(ctx context.Context, timeFrameItem *_struct.Timeframe) (err error) {
	err = i.loadCandles(timeFrameItem)
	if err != nil {
		return
	}

	i.State.Heap, err = i.Storage.UpdateHeap(ctx, i.State.Heap)
	if err != nil {
		return
	}

	_, sidewaysKLinesAmount := trading.CheckSideways(timeFrameItem.Candles, timeFrameItem.Config.SidewaysMinCandlesAmount, timeFrameItem.Config.SidewaysPercentToPrice)
	sidewaysKlines := make([]domainStructs.DomainCandle, sidewaysKLinesAmount)
	copy(sidewaysKlines, timeFrameItem.Candles[:sidewaysKLinesAmount])
	premiumDiscount := trading.PremiumDiscount(sidewaysKlines)

	if premiumDiscount > timeFrameItem.Config.SidewaysPremiumCoefficient {
		// premium
		fmt.Println("premium")
	}
	if premiumDiscount < timeFrameItem.Config.SidewaysDiscountCoefficient {
		// discount
		fmt.Println("discount")
	}

	return
}

func (i *Investor) loadCandles(timeFrameItem *_struct.Timeframe) (err error) {
	if timeFrameItem.GetCandleTime.Before(time.Now().Add(-time.Duration(timeFrameItem.Config.CandleCacheSeconds) * time.Second)) {
		timeFrameItem.Candles, err = i.provider.LoadCandleHistory(i.Config.CoinPare, timeFrameItem.Config.Resolution, timeFrameItem.Config.CandleReview)
		if err != nil {
			return
		}
		timeFrameItem.GetCandleTime = time.Now()
		i.State.CurrentPrice = timeFrameItem.Candles[0].Close
		time.Sleep(time.Second * i.Config.RequestDelay)
	}

	return
}
