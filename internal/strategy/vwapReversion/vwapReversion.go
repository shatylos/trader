package vwapReversion

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/vwapReversion/storage"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/vwapReversion/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	"time"
)

type VwapReversion struct {
	config   Config
	provider domain.FuturesDomainInterface
	state    State
}

type State struct {
	LtTrend           string
	StTrend           string
	Vwap              float64
	UpperBand         float64
	LowerBand         float64
	CurrentPrice      float64
	NewOrderCondition string
	stCandles         []domainStructs.DomainCandle
}

func (v *VwapReversion) Init(mux *http.ServeMux) {
}

func (v *VwapReversion) GetId() string {
	return v.config.Id
}

func (v *VwapReversion) IsEnabled() bool {
	return v.config.Enabled
}

func (v *VwapReversion) GetTitle() string {
	if !v.config.Enabled {
		return fmt.Sprintf("VWAP Reversion: %s (%s) (DISABLED)", v.config.Id, v.config.CoinPare)
	}
	return fmt.Sprintf("VWAP Reversion: %s (%s)", v.config.Id, v.config.CoinPare)
}

func (v *VwapReversion) DoAction() (err error) {

	state := State{}
	if !v.config.Enabled {
		v.state = state
		if v.config.Verbose {
			logger.Info("The setup is disabled. Set enabled:1 in config file to enable it.")
		}
		return
	}

	state.stCandles, err = v.provider.LoadCandleHistory(v.config.CoinPare, v.config.Resolution, v.config.VwapPeriod)
	if err != nil {
		err = apperrors.Wrap(err, "error load short term candle history")
		return
	}

	var internalPosition structs.Position
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(v.config.Id)
	if err != nil {
		err = apperrors.Wrap(err, "error get storage")
		return
	}
	internalPosition, err = storage.GetLastInternalPosition(0)
	if err != nil {
		err = apperrors.Wrap(err, "error get last internal position")
		return
	}

	var providerPosition domainStructs.DomainPosition
	providerPosition, err = v.provider.GetPosition(v.config.CoinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error get provider position")
		return
	}

	currentPrice := state.stCandles[0].Close

	if providerPosition.Size == 0 {
		if internalPosition.Id != nil && internalPosition.Status != structs.StatusClosed {
			err = v.closeInternalPosition(internalPosition)
			if err != nil {
				err = apperrors.Wrap(err, "error close internal position")
				return
			}
		}
		internalPosition, err = v.calculateNewPosition(state.stCandles)
		if err != nil {
			err = apperrors.Wrap(err, "error calculate new position")
			return
		}
		internalPosition.ProviderPosition = providerPosition
	} else if internalPosition.Status == structs.StatusActive {
		internalPosition.ProviderPosition = providerPosition
		internalPosition, err = storage.SaveInternalPosition(internalPosition)
		if err != nil {
			err = apperrors.Wrap(err, "error save internal position")
			return
		}
	} else {
		// @TODO: Move the message to state to show it in report
		//logger.Info("Wait for close current provider position")
		return
	}

	err = v.actionByPosition(internalPosition, currentPrice, state.stCandles)
	if err != nil {
		err = apperrors.Wrap(err, "error action by position")
		return
	}

	state.Vwap = internalPosition.Chart.Vwap
	state.UpperBand = internalPosition.Chart.UpperBand
	state.LowerBand = internalPosition.Chart.LowerBand
	state.CurrentPrice = currentPrice
	state.LtTrend = internalPosition.LtTrend
	state.StTrend = internalPosition.StTrend
	state.NewOrderCondition = v.newOrderCondition(internalPosition)

	v.state = state
	return
}

func (v *VwapReversion) WaitDuration() time.Duration {
	return time.Second * v.config.TimeoutSeconds
}

func (v *VwapReversion) newOrderCondition(internalPosition structs.Position) (condition string) {
	if internalPosition.Order.OrderId != "" {
		return
	}
	switch internalPosition.Trend {
	case "BULLISH":
		condition = fmt.Sprintf("X <= %g (lower band)", internalPosition.Chart.LowerBand)
	case "BEARISH":
		condition = fmt.Sprintf("X >= %g (upper band)", internalPosition.Chart.UpperBand)
	}
	return
}
