package structs

import (
	"fmt"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/apperrors"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/tgNotifier"
	"net/http"
	"sync"
	"time"
)

type Setup struct {
	ID         string
	errorCount int64
	Strategy   _struct.StrategyInterface
}

type SetupDelay struct {
	Duration time.Duration
	Setup    *Setup
}

func (s *Setup) Init(mux *http.ServeMux) {
	s.Strategy.Init(mux)
}

func (s *Setup) NextStep(setupDelayChanel chan *SetupDelay, setupWg *sync.WaitGroup) {
	strategy := s.Strategy

	err := strategy.DoAction()
	if err != nil {
		err = apperrors.Wrap(err, "error DoAction for setup: %s", s.ID)
		logger.PrintError(err)
		s.errorCount++
		logger.Info(fmt.Sprintf("error count: %d", s.errorCount))
		if s.errorCount == 10 || s.errorCount%100 == 0 {
			tgNotifier.Notify(fmt.Sprintf("Errors for setup %s. Error count: %d", s.ID, s.errorCount))
		}
		setupDelayChanel <- &SetupDelay{
			Duration: time.Second * time.Duration(s.errorCount) * 5,
			Setup:    s,
		}
		setupWg.Done()
		return
	}

	if s.errorCount > 0 {
		logger.Info(fmt.Sprintf("Success iteration for setup %s after %d errors", s.ID, s.errorCount))
		s.errorCount = 0
	}

	setupDelayChanel <- &SetupDelay{
		Duration: strategy.WaitDuration(),
		Setup:    s,
	}
	setupWg.Done()
}
