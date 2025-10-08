package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

type Timeframe struct {
	Config        TimeframeConfig
	GetCandleTime int64
	Candles       []domainStructs.DomainCandle
}

func (i *Investor) handleTimeframe(ctx context.Context, timeFrameItem *Timeframe) (err error) {
	err = i.loadCandles(timeFrameItem)
	if err != nil {
		return
	}

	var deal *storage.Deal
	deal, err = i.Storage.GetActiveDealByTimeframe(ctx, timeFrameItem.Config.Resolution)
	if err != nil {
		return
	}

	var dealRelation *DealRelation
	dealRelation, err = i.GetDealRelation(ctx, deal)
	if err != nil {
		return
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
				if i.config.Verbose {
					logger.Info(fmt.Sprintf("Wait for fill the order %s", dealOrder.OrderId))
				}
				// @TODO: cancel order if it's open long time
				return
			}
		}
	}

	isSideways, sidewaysKLinesAmount := trading.CheckSideways(timeFrameItem.Candles, timeFrameItem.Config.SidewaysMinCandlesAmount, timeFrameItem.Config.SidewaysPercentToPrice)
	if !isSideways {
		// @TODO: зберегти стан
		//timeFrameItem.TrendState =
		if i.config.Verbose {
			logger.Info(fmt.Sprintf("Timeframe %s is not in sideways", timeFrameItem.Config.Resolution))
		}
		return
	}

	sidewaysKlines := make([]domainStructs.DomainCandle, sidewaysKLinesAmount)
	copy(sidewaysKlines, timeFrameItem.Candles[:sidewaysKLinesAmount])
	premiumDiscount := trading.PremiumDiscount(sidewaysKlines)

	if premiumDiscount > timeFrameItem.Config.SidewaysPremiumCoefficient {
		if timeFrameItem.Config.IsHeap {
			//err = i.handlePremiumHeap(ctx, timeFrameItem)
			err = i.handlePremium(ctx, dealRelation, timeFrameItem)
			if err != nil {
				return
			}
		} else {
			err = i.handlePremium(ctx, dealRelation, timeFrameItem)
			if err != nil {
				return
			}
		}
	}
	if premiumDiscount < timeFrameItem.Config.SidewaysDiscountCoefficient {
		if timeFrameItem.Config.IsHeap {
			//err = i.handleDiscountHeap(ctx, timeFrameItem)
			err = i.handleDiscount(ctx, dealRelation, timeFrameItem)
			if err != nil {
				return
			}
		} else {
			err = i.handleDiscount(ctx, dealRelation, timeFrameItem)
			if err != nil {
				return
			}
		}
	}

	return
}

func (i *Investor) loadCandles(timeFrameItem *Timeframe) (err error) {

	if timeFrameItem.GetCandleTime+timeFrameItem.Config.CandleCacheSeconds < time.Now().Unix() {
		timeFrameItem.Candles, err = i.provider.LoadCandleHistory(i.config.CoinPare, timeFrameItem.Config.Resolution, timeFrameItem.Config.CandleReview)
		if err != nil {
			return
		}
		timeFrameItem.GetCandleTime = time.Now().Unix()
		i.State.CurrentPrice = timeFrameItem.Candles[0].Close
		time.Sleep(time.Second * i.config.RequestDelay)
	}

	return
}
