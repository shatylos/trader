package scalper

import (
	"fmt"
	"github.com/shatylos/trader/internal/domain"
	domainStructs "github.com/shatylos/trader/internal/domain/structs"
	strategyStorage "github.com/shatylos/trader/internal/strategy/scalper/storage"
	"github.com/shatylos/trader/internal/strategy/scalper/storage/mongo"
	"github.com/shatylos/trader/internal/strategy/scalper/structs"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"net/http"
	"time"
)

type Scalper struct {
	config      Config
	provider    domain.FuturesDomainInterface
	state       State
	candleCache map[string]candleCacheEntry
}

type State struct {
	Bias           string
	LtTrend        string
	HtfEmaFast     float64
	HtfEmaSlow     float64
	LtfEma         float64
	Rsi            float64
	Atr            float64
	CurrentPrice   float64
	Signal         string
	SkippedMessage string
}

func (s *Scalper) Init(mux *http.ServeMux) {
}

func (s *Scalper) GetId() string {
	return s.config.Id
}

func (s *Scalper) IsEnabled() bool {
	return s.config.Enabled
}

func (s *Scalper) GetTitle() string {
	if !s.config.Enabled {
		return fmt.Sprintf("Scalper: %s (%s) (DISABLED)", s.config.Id, s.config.CoinPare)
	}
	return fmt.Sprintf("Scalper: %s (%s)", s.config.Id, s.config.CoinPare)
}

func (s *Scalper) DoAction() (err error) {

	if !s.config.Enabled {
		s.state = State{}
		if s.config.Verbose {
			logger.Info("The setup is disabled. Set enabled:1 in config file to enable it.")
		}
		return
	}

	var internalPosition structs.Position
	var storage mongo.MongoStorage
	storage, err = strategyStorage.GetStorage(s.config.Id)
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
	providerPosition, err = s.provider.GetPosition(s.config.CoinPare)
	if err != nil {
		err = apperrors.Wrap(err, "error get provider position")
		return
	}

	if providerPosition.Size != 0 {
		// A position is running; TP/SL live on the exchange, nothing to decide here.
		// Keep the internal snapshot in sync for the report.
		if internalPosition.Id != nil && internalPosition.Status == structs.StatusActive {
			internalPosition.ProviderPosition = providerPosition
			internalPosition, err = storage.SaveInternalPosition(internalPosition)
			if err != nil {
				err = apperrors.Wrap(err, "error save internal position")
				return
			}
		}
		s.state.SkippedMessage = "Wait for close current provider position"
		return
	}

	if internalPosition.Id != nil && internalPosition.Status != structs.StatusClosed {
		err = s.closeInternalPosition(internalPosition)
		if err != nil {
			err = apperrors.Wrap(err, "error close internal position")
			return
		}
	}

	err = s.calculateAndAction()
	if err != nil {
		err = apperrors.Wrap(err, "error calculate and action")
		return
	}

	return
}

// calculateAndAction builds the multi-timeframe indicator snapshot, stores it
// in the state for the report and opens a new position when the entry signal
// fires.
func (s *Scalper) calculateAndAction() (err error) {
	var snapshot signalSnapshot
	snapshot, err = s.calculateSignal()
	if err != nil {
		err = apperrors.Wrap(err, "error calculate signal")
		return
	}

	s.state = State{
		Bias:         snapshot.Bias,
		LtTrend:      snapshot.LtTrend,
		HtfEmaFast:   snapshot.Chart.HtfEmaFast,
		HtfEmaSlow:   snapshot.Chart.HtfEmaSlow,
		LtfEma:       snapshot.Chart.LtfEma,
		Rsi:          snapshot.Chart.Rsi,
		Atr:          snapshot.Chart.Atr,
		CurrentPrice: snapshot.CurrentPrice,
		Signal:       snapshot.Signal,
	}

	if snapshot.Signal == "" {
		s.state.SkippedMessage = snapshot.SkippedMessage
		if s.config.Verbose {
			logger.Info(fmt.Sprintf("No entry signal. Bias: %s. %s", snapshot.Bias, snapshot.SkippedMessage))
		}
		return
	}

	err = s.openNewPosition(snapshot)
	if err != nil {
		err = apperrors.Wrap(err, "error open new position")
		return
	}

	return
}

func (s *Scalper) WaitDuration() time.Duration {
	return time.Second * s.config.TimeoutSeconds
}
