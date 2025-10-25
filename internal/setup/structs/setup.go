package structs

import (
	"fmt"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/logger"
	"github.com/shatylos/trader/tools/tgNotifier"
	"net/http"
	"time"
)

type Setup struct {
	ID         string
	errorCount int64
	Strategy   _struct.StrategyInterface
}

func (s *Setup) Init(mux *http.ServeMux) {
	s.Strategy.Init(mux)
}

func (s *Setup) NextStep(setupChanel chan *Setup) {
	strategy := s.Strategy

	err := strategy.DoAction()
	if err != nil {
		errorMsg := err.Error()
		strategyTitle := strategy.GetTitle()
		logger.Error(fmt.Sprintf("Error DoAction, setup: %s. Error: %s", strategyTitle, errorMsg))
		s.errorCount++
		logger.Info(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
		if s.errorCount == 10 || s.errorCount%100 == 0 {
			tgNotifier.Notify(fmt.Sprintf("Errors for setup %s. Error count: %d", s.ID, s.errorCount))
		}
		time.Sleep(time.Second * time.Duration(s.errorCount) * 5)
		setupChanel <- s
		return
	}

	if s.errorCount > 0 {
		logger.Info(fmt.Sprintf("Success iteration after %d errors", s.errorCount))
		s.errorCount = 0
	}

	strategy.Wait()
	setupChanel <- s
}
