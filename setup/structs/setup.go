package structs

import (
	strategyInterface "bitbucket.org/shatylos/trader/strategy/interface"
	"bitbucket.org/shatylos/trader/utils"
)

const StatusReadyForNext = 0
const StatusInProgress = 1
const StatusError = 2

type Setup struct {
	status   int64
	Strategy strategyInterface.StrategyInterface
}

func (s *Setup) GetStatus() int64 {
	return s.status
}

func (s *Setup) SetStatus(status int64) {
	s.status = status
}

func (s *Setup) NextStep() {
	strategy := s.Strategy

	err := strategy.GetData()
	if err != nil {
		utils.LogError(err.Error())
		s.SetStatus(StatusError)
		return
	}
	err = strategy.Analyse()
	if err != nil {
		utils.LogError(err.Error())
		s.SetStatus(StatusError)
		return
	}
	err = strategy.DoAction()
	if err != nil {
		utils.LogError(err.Error())
		s.SetStatus(StatusError)
		return
	}

	strategy.Wait()
	s.SetStatus(StatusReadyForNext)
}
