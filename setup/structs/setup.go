package structs

import (
	"fmt"
	"github.com/shatylos/trader/strategy/struct"
	"github.com/shatylos/trader/utils"
	"time"
)

const StatusReadyForNext = 0
const StatusInProgress = 1
const StatusError = 2
const MaxErrorsCount = 5

type Setup struct {
	status     int64
	errorCount int64
	Strategy   _struct.StrategyInterface
}

func (s *Setup) GetStatus() int64 {
	return s.status
}

func (s *Setup) SetStatus(status int64) {
	s.status = status
}

func (s *Setup) NextStep() {
	strategy := s.Strategy

	if !strategy.IsInit() {
		err := strategy.Initialise()
		if err != nil {
			utils.LogError(err.Error())
			s.errorCount++
			utils.LogInfo(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
			time.Sleep(time.Second * time.Duration(s.errorCount) * 5)
			s.SetStatus(StatusReadyForNext)
			return
		}
	}

	err := strategy.DoAction()
	if err != nil {
		utils.LogError(err.Error())
		s.errorCount++
		utils.LogInfo(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
		time.Sleep(time.Second * time.Duration(s.errorCount) * 5)
		s.SetStatus(StatusReadyForNext)
		return
	}

	if s.errorCount > 0 {
		utils.LogInfo(fmt.Sprintf("Success iteration after %d errors", s.errorCount))
		s.errorCount = 0
	}

	strategy.Wait()
	s.SetStatus(StatusReadyForNext)
}
