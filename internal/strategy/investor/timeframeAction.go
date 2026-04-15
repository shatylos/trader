package investor

import (
	"context"
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/investor/entity"
	_struct "github.com/shatylos/trader/internal/strategy/investor/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/math"
	"github.com/shatylos/trader/tools/trading"
	"time"
)

func (i *Investor) handleTimeframe(ctx context.Context, timeFrameItem *_struct.TimeframeItem) (err error) {
	err = i.loadCandles(timeFrameItem)
	if err != nil {
		err = apperrors.Wrap(err, "error load candles")
		return
	}

	var deal *entity.Deal
	deal, err = i.Storage.GetActiveDealByTimeframe(ctx, timeFrameItem)
	if err != nil {
		err = apperrors.Wrap(err, "error get active deal by timeframe %s", timeFrameItem.Config.Resolution)
		return
	}

	var dealRelation *entity.DealRelation
	dealRelation, err = i.Storage.GetDealRelation(ctx, deal)
	if err != nil {
		err = apperrors.Wrap(err, "error get deal relation")
		return
	}

	isSideways, sidewaysKLinesAmount := trading.CheckSideways(timeFrameItem.Candles, timeFrameItem.Config.SidewaysMinCandlesAmount, timeFrameItem.Config.SidewaysPercentToPrice)
	timeFrameItem.IsSidewaysState = isSideways

	timeFrameItem.Trend, timeFrameItem.TrendSlope = trading.GetTrendLinearRegression(timeFrameItem.Candles)

	sidewaysKlines := make([]domainStructs.DomainCandle, sidewaysKLinesAmount)
	copy(sidewaysKlines, timeFrameItem.Candles[:sidewaysKLinesAmount])

	if len(sidewaysKlines) == 0 {
		timeFrameItem.TradeStateMsg = "Not enough sideways candles to trade"
		if i.Config.Verbose {
			logger.Info(fmt.Sprintf("Timeframe %s is not in sideways. Not enough candles", timeFrameItem.Config.Resolution))
		}
		return
	}

	timeFrameItem.SidewaysFrom = time.Unix(sidewaysKlines[len(sidewaysKlines)-1].Time, 0)

	currentPrice := sidewaysKlines[0].Close
	vwap := trading.CreateVWAP(sidewaysKlines)

	err = i.calculateNextQtyAndPrices(ctx, timeFrameItem, dealRelation, &vwap)
	if err != nil {
		err = apperrors.Wrap(err, "error calculate next qty and prices")
		return
	}

	if dealRelation.QtyToBuy == 0 && dealRelation.QtyToSell == 0 {
		logger.Info("qty to buy and sell are 0. Close the deal")
		deal.SetClose()
		err = i.Storage.SaveDeal(ctx, deal)
		if err != nil {
			err = apperrors.Wrap(err, "error saving deal when try to close it")
			return
		}
		return
	}

	zone := trading.ZoneNeutral
	if currentPrice > vwap.AncVWAP {
		zone = trading.ZonePremium
	}
	if currentPrice < vwap.AncVWAP {
		zone = trading.ZoneDiscount
	}
	timeFrameItem.Zone = zone

	for _, dealOrder := range dealRelation.Orders {
		switch dealOrder.OrderStatus {
		case domainStructs.OrderStatuses.New,
			domainStructs.OrderStatuses.Open,
			domainStructs.OrderStatuses.PartiallyFilled:
			timeFrameItem.TradeStateMsg = "Wait for fill the order"
			err = i.updateOrder(ctx, dealRelation, dealOrder, timeFrameItem)
			if err != nil {
				err = apperrors.Wrap(err, "error update order")
				return
			}
			if dealOrder.OrderStatus != domainStructs.OrderStatuses.Filled {

				if (dealOrder.Side == domainStructs.OrderSideBuy && zone == trading.ZonePremium) ||
					(dealOrder.Side == domainStructs.OrderSideSell && zone == trading.ZoneDiscount) {

					err = i.doCancel(ctx, dealOrder)
					if err != nil {
						err = apperrors.Wrap(err, "error do cancel")
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
		timeFrameItem.TradeStateMsg = "Not a sideways market for trading"
		if i.Config.Verbose {
			logger.Info(fmt.Sprintf("Timeframe %s is not in sideways", timeFrameItem.Config.Resolution))
		}
		return
	}

	switch zone {
	case trading.ZonePremium:
		err = i.handlePremium(ctx, dealRelation, timeFrameItem)
		if err != nil {
			err = apperrors.Wrap(err, "error handle premium")
			return
		}
	case trading.ZoneDiscount:
		err = i.handleDiscount(ctx, dealRelation, timeFrameItem)
		if err != nil {
			err = apperrors.Wrap(err, "error handle discount")
			return
		}
	case trading.ZoneNeutral:
		timeFrameItem.TradeStateMsg = "Timeframe zone is neutral, no trade"
	}

	return
}

func (i *Investor) loadCandles(timeFrame _struct.Timeframe) (err error) {
	if timeFrame.GetCandleTime().Before(time.Now().Add(-timeFrame.GetConfig().GetCandleCacheDuration())) {
		var candles []domainStructs.DomainCandle
		candles, err = i.provider.LoadCandleHistory(i.Config.CoinPare, timeFrame.Resolution(), timeFrame.GetConfig().GetCandleReview())
		if err != nil {
			err = apperrors.Wrap(err, "error load candle history")
			return
		}
		timeFrame.SetCandles(candles)
		i.State.CurrentPrice = timeFrame.GetCandles()[0].Close
		time.Sleep(i.Config.RequestDelay)
	}

	return
}

func (i *Investor) modifySidewaysCoeff(origCoeff, percentage float64) (modifiedCoeff float64) {
	if origCoeff == 0 {
		return
	}
	addValue := math.Mul(math.Div(origCoeff, 100), percentage)
	if origCoeff < 0 {
		addValue = math.Mul(addValue, -1)
	}
	modifiedCoeff = origCoeff + addValue
	return
}
