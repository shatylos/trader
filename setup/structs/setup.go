package structs

import (
	strategyInterface "bitbucket.org/shatylos/trader/strategy/interface"
	"time"
)

const StatusReadyForNext = 0
const StatusInProgress = 1

type Setup struct {
	status     int64
	DomainCode string
	Strategy   strategyInterface.StrategyInterface
}

func (s *Setup) GetStatus() int64 {
	return s.status
}

func (s *Setup) SetStatus(status int64) {
	s.status = status
}

func (s *Setup) NextStep() {

	println("=========================")
	println("==     Handle Step     ==")
	println("=========================")
	time.Sleep(5 * time.Second)
	s.SetStatus(StatusReadyForNext)
}
