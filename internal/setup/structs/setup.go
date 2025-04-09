package structs

import (
	"fmt"
	"github.com/shatylos/trader/internal/strategy/struct"
	"github.com/shatylos/trader/tools/logger"
	"time"
)

type Setup struct {
	ID         string
	errorCount int64
	Strategy   _struct.StrategyInterface
}

func (s *Setup) NextStep(setupChanel chan *Setup) {

	defer func() {
		if r := recover(); r != nil {
			logger.Error(fmt.Sprintf("Recovered from panic: %v", r))
		}
	}()

	strategy := s.Strategy

	if !strategy.IsInit() {
		err := strategy.Initialise()
		if err != nil {
			logger.Error(err.Error())
			s.errorCount++
			logger.Info(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
			time.Sleep(time.Second * time.Duration(s.errorCount) * 5)
			setupChanel <- s
			return
		}
	}

	err := strategy.DoAction()
	if err != nil {
		errorMsg := err.Error()
		strategyTitle := strategy.GetTitle()
		logger.Error(fmt.Sprintf("Error DoAction, setup: %s. Error: %s", strategyTitle, errorMsg))
		s.errorCount++
		logger.Info(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
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
