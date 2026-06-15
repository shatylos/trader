package vwapReversion

import (
	"fmt"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/trading"
)

func (v *VwapReversion) actionByPosition(internalPosition structs.Position, currentPrice float64) (err error) {
	if internalPosition.Order.OrderId != "" {
		// Position already has an order; TP/SL live on the exchange. Nothing to do.
		return
	}

	switch internalPosition.Trend {
	case trading.TrendLong:
		err = v.actionBullish(internalPosition, currentPrice)
		if err != nil {
			err = apperrors.Wrap(err, "error action bullish")
		}
	case trading.TrendShort:
		err = v.actionBearish(internalPosition, currentPrice)
		if err != nil {
			err = apperrors.Wrap(err, "error action bearish")
		}
	case trading.TrendUnknown:
		if v.config.Verbose {
			logger.Info(fmt.Sprintf("Trend is %s. Wait for %s or %s.", internalPosition.Trend, trading.TrendLong, trading.TrendShort))
		}
	default:
		err = apperrors.New("unexpected trend \"%s\" for action by position", internalPosition.Trend)
	}
	return
}

func (v *VwapReversion) actionBullish(internalPosition structs.Position, currentPrice float64) (err error) {
	if internalPosition.BalanceBefore <= 0 {
		return apperrors.New("not enough balance (%g) to action", internalPosition.BalanceBefore)
	}

	// Long setup: price stretched below the lower deviation band, in a bullish regime.
	isPriceBelowBand := currentPrice <= internalPosition.Chart.LowerBand
	var confirmed bool
	confirmed, err = v.isReversalConfirmed(domainStructs.PositionSideLong)
	if err != nil {
		err = apperrors.Wrap(err, "error check reversal confirmation")
		return
	}

	if isPriceBelowBand && confirmed {
		_, err = v.openNewPosition(internalPosition, domainStructs.PositionSideLong)
		if err != nil {
			err = apperrors.Wrap(err, "error open new long position")
		}
		return
	}

	return
}

func (v *VwapReversion) actionBearish(internalPosition structs.Position, currentPrice float64) (err error) {
	if internalPosition.BalanceBefore <= 0 {
		logger.Error(fmt.Sprintf("Not enough balance (%g) to action", internalPosition.BalanceBefore))
		return apperrors.New("not enough balance (%g) to action", internalPosition.BalanceBefore)
	}

	// Short setup: price stretched above the upper deviation band, in a bearish regime.
	priceAboveBand := currentPrice >= internalPosition.Chart.UpperBand
	var confirmed bool
	confirmed, err = v.isReversalConfirmed(domainStructs.PositionSideShort)
	if err != nil {
		return apperrors.Wrap(err, "error check reversal confirmation")
	}

	if priceAboveBand && confirmed {
		_, err = v.openNewPosition(internalPosition, domainStructs.PositionSideShort)
		if err != nil {
			err = apperrors.Wrap(err, "error open new short position")
		}
		return
	}

	return
}

// isReversalConfirmed optionally checks that momentum has stopped pushing further
// into the extreme before we fade it: for a long we want the last closed candle
// to be green (close >= open), for a short red (close <= open). When confirmation
// is disabled in config it always returns true.
func (v *VwapReversion) isReversalConfirmed(side string) (confirmed bool, err error) {
	if !v.config.RequireConfirmation {
		confirmed = true
		return
	}

	var candles []domainStructs.DomainCandle
	candles, err = v.provider.LoadCandleHistory(v.config.CoinPare, v.config.Resolution, 2)
	if err != nil {
		err = apperrors.Wrap(err, "error load candle history")
		return
	}
	if len(candles) < 2 {
		err = apperrors.New("no candles to confirm reversal")
		return
	}

	lastClosed := candles[1]
	if side == domainStructs.PositionSideLong {
		confirmed = lastClosed.Close >= lastClosed.Open
	} else {
		confirmed = lastClosed.Close <= lastClosed.Open
	}
	return
}
