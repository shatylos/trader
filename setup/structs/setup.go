package structs

import (
	strategyInterface "bitbucket.org/shatylos/trader/strategy/interface"
	"bitbucket.org/shatylos/trader/utils"
	"fmt"
)

const StatusReadyForNext = 0
const StatusInProgress = 1
const StatusError = 2
const MaxErrorsCount = 5

type Setup struct {
	status     int64
	errorCount int64
	Strategy   strategyInterface.StrategyInterface
}

func (s *Setup) GetStatus() int64 {
	return s.status
}

func (s *Setup) SetStatus(status int64) {
	s.status = status
}

func (s *Setup) NextStep() {
	strategy := s.Strategy

	println(fmt.Sprintf("s.errorCount start: %d", s.errorCount))
	println(fmt.Sprintf("MaxErrorsCount start: %d", MaxErrorsCount))

	if !strategy.IsInit() {
		err := strategy.Initialise()
		if err != nil {
			utils.LogError(err.Error())
			s.errorCount++
			println(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
			if s.errorCount >= MaxErrorsCount {
				s.SetStatus(StatusError)
			} else {
				s.SetStatus(StatusReadyForNext)
			}
			return
		}
	}

	err := strategy.GetData()
	if err != nil {
		utils.LogError(err.Error())
		s.errorCount++
		println(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
		if s.errorCount >= MaxErrorsCount {
			s.SetStatus(StatusError)
		} else {
			s.SetStatus(StatusReadyForNext)
		}
		return
	}
	err = strategy.Analyse()
	if err != nil {
		utils.LogError(err.Error())
		s.errorCount++
		println(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
		if s.errorCount >= MaxErrorsCount {
			s.SetStatus(StatusError)
		} else {
			s.SetStatus(StatusReadyForNext)
		}
		return
	}
	err = strategy.DoAction()
	if err != nil {
		utils.LogError(err.Error())
		s.errorCount++
		println(fmt.Sprintf("s.errorCount after error: %d", s.errorCount))
		if s.errorCount >= MaxErrorsCount {
			s.SetStatus(StatusError)
		} else {
			s.SetStatus(StatusReadyForNext)
		}
		return
	}

	strategy.Wait()
	s.errorCount = 0
	s.SetStatus(StatusReadyForNext)
}
