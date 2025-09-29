package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/storage"
	"github.com/shatylos/trader/internal/trading/constant"
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

	deal := storage.Deal{}
	err = i.Storage.GetLastDealByTimeframe(ctx, timeFrameItem.Config.Resolution, &deal)
	if err != nil {
		return
	}

	var dealOrders []*storage.Order
	if deal.Id != nil {
		dealOrders, err = i.Storage.GetOrdersByDealId(ctx, *deal.Id)
		if err != nil {
			return
		}
	}

	for _, dealOrder := range dealOrders {
		switch dealOrder.OrderStatus {
		case domainStructs.OrderStatuses.New:
		case domainStructs.OrderStatuses.Open:
		case domainStructs.OrderStatuses.PartiallyFilled:
			err = i.updateOrder(ctx, &deal, dealOrder, timeFrameItem)
			if err != nil {
				return
			}
			break
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
			//if err != nil {
			//	return
			//}
		} else {
			err = i.handlePremium(ctx, &deal, dealOrders, timeFrameItem)
			if err != nil {
				return
			}
		}
	}
	if premiumDiscount < timeFrameItem.Config.SidewaysDiscountCoefficient {
		if timeFrameItem.Config.IsHeap {
			//err = i.handleDiscountHeap(ctx, timeFrameItem)
			//if err != nil {
			//	return
			//}
		} else {
			err = i.handleDiscount(ctx, &deal, dealOrders, timeFrameItem)
			if err != nil {
				return
			}
		}
	}

	return
}

func (i *Investor) loadCandles(timeFrameItem *Timeframe) (err error) {

	var resolutionSeconds int64
	resolutionSeconds, err = constant.ResolutionToSeconds(timeFrameItem.Config.Resolution)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	if timeFrameItem.GetCandleTime+resolutionSeconds < time.Now().Unix() {
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
